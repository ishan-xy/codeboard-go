package utility

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const PROFILE_QUERY = `query userPublicProfile($username: String!) {
	matchedUser(username: $username) {
		profile {
			ranking
			userAvatar
			realName
		}
	}
}`

const NUM_QUESTION_QUERY = `query userSessionProgress($username: String!) {
  allQuestionsCount {
    count
  }
  matchedUser(username: $username) {
    submitStats {
      acSubmissionNum {
        count
      }
    }
  }
}
`

type GraphQLQuery struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type GraphQLResponse struct {
	Data   map[string]interface{} `json:"data"`
	Errors []interface{}          `json:"errors"`
}

const baseURL = "https://leetcode.com/graphql"

func SendQuery(query string, variables map[string]interface{}) (GraphQLResponse, error) {
	// Prepare the GraphQL request
	graphqlRequest := GraphQLQuery{
		Query:     query,
		Variables: variables,
	}

	// Serialize the request to JSON
	requestBody, err := json.Marshal(graphqlRequest)
	if err != nil {
		return GraphQLResponse{}, fmt.Errorf("failed to serialize request: %v", err)
	}

	// Create a new HTTP client
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Send the HTTP POST request
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return GraphQLResponse{}, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return GraphQLResponse{}, fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GraphQLResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read and parse the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GraphQLResponse{}, fmt.Errorf("failed to read response: %v", err)
	}

	var graphqlResponse GraphQLResponse
	if err := json.Unmarshal(responseBody, &graphqlResponse); err != nil {
		return GraphQLResponse{}, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	// Check for GraphQL errors
	if len(graphqlResponse.Errors) > 0 {
		return graphqlResponse, fmt.Errorf("graphql errors: %v", graphqlResponse.Errors)
	}

	return graphqlResponse, nil
}
