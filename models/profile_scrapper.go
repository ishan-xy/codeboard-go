package models

import (
	"log"
	"sync"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"gorm.io/gorm"
)

var (
	List QuestionList
	db   *gorm.DB // Assume initialized elsewhere
)

type recentlyAnsweredResponse struct {
	TitleSlug string 
	Timestamp string 
}

var UserList = cmap.New[*UserPublicProfile]()

type user struct {
	*UserPublicProfile
	mu sync.Mutex
}

// sync updates the user's public profile in the database.
func (u *user) sync() error {
	err := db.Model(u.UserPublicProfile).Updates(u.UserPublicProfile).Error
	if err != nil {
		log.Printf("Error syncing user '%s': %v", u.LeetcodeUsername, err)
	} else {
		log.Printf("Successfully synced user '%s'", u.LeetcodeUsername)
	}
	return err
}

// mergeRecent fetches the recently answered questions and updates the user's solved list.
func (u *user) mergeRecent() (bool, error) {
	log.Printf("Merging recent submissions for user '%s'", u.LeetcodeUsername)
	solved, err := RecentlyAnswered(u.LeetcodeUsername)
	if err != nil {
		log.Printf("Error fetching recent submissions for user '%s': %v", u.LeetcodeUsername, err)
		return false, err
	}

	altered := false
	for _, q := range solved {
		dest, ok := List.SlugMap[q.TitleSlug]
		if !ok {
			// Skip if question is not found in the master list.
			continue
		}

		u.mu.Lock()
		if _, exists := u.Solved[dest.Title]; !exists {
			u.Solved[dest.Title] = q.Timestamp
			altered = true
			log.Printf("User '%s': added solved question '%s'", u.LeetcodeUsername, dest.Title)
		}
		u.mu.Unlock()
	}
	if !altered {
		log.Printf("User '%s' has no new solved questions", u.LeetcodeUsername)
	}
	return altered, nil
}

// RecentlyAnswered queries the API to retrieve recent submissions.
func RecentlyAnswered(username string) ([]recentlyAnsweredResponse, error) {
	type J map[string]any
	var result struct {
		Data struct {
			List []recentlyAnsweredResponse `json:"recentAcSubmissionList"`
		} `json:"data"`
	}

	err := (&query{
		Query: `query recentAcSubmissions($username: String!, $limit: Int!) {
			recentAcSubmissionList(username: $username, limit: $limit) {
				s:titleSlug
				p:timestamp
			}
		}`,
		Variables: J{
			"username": username,
			"limit":    1024,
		},
		OperationName: "recentAcSubmissions",
	}).jsonResponse(&result)

	if err != nil {
		log.Printf("Error in RecentlyAnswered for user '%s': %v", username, err)
		return nil, err
	}
	return result.Data.List, nil
}

const maxConcurrent = 100

// syncAll iterates through all users, merges their recent submissions,
// and updates their profiles if any new solved questions are found.
func syncAll() {
	log.Println("Starting syncAll")
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, maxConcurrent)

	UserList.IterCb(func(k string, upp *UserPublicProfile) {
		wg.Add(1)
		go func(u *user) {
			// Acquire semaphore for limiting concurrency.
			sem <- struct{}{}
			defer func() {
				<-sem
				wg.Done()
			}()

			altered, err := u.mergeRecent()
			if err != nil {
				log.Printf("Skipping sync for user '%s' due to merge error: %v", u.LeetcodeUsername, err)
				return
			}

			if altered {
				if err := u.sync(); err != nil {
					log.Printf("Failed to sync user '%s': %v", u.LeetcodeUsername, err)
				}
			} else {
				log.Printf("No changes detected for user '%s'", u.LeetcodeUsername)
			}
		}(&user{UserPublicProfile: upp})
	})

	wg.Wait()
	log.Println("Completed syncAll")
}

const syncInterval = 15 * time.Second

// SyncLoop runs syncAll at a regular interval.
func SyncLoop() {
	log.Println("Starting SyncLoop")
	for {
		beginning := time.Now()
		syncAll()
		remaining := time.Until(beginning.Add(syncInterval))
		if remaining > 0 {
			log.Printf("Sleeping for %v until next sync", remaining)
			time.Sleep(remaining)
		}
	}
}
