package models

import (
	"fmt"
	"log"
	// "regexp"
	// "strings"

	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
	utils "github.com/ItsMeSamey/go_utils"
)

type query struct {
	Query         string         `json:"query"`
	Variables     map[string]any `json:"variables"`
	OperationName string         `json:"operationName"`
}

// var whitespaceRegex = regexp.MustCompile(`\s+`)

// marshal cleans up the query string and marshals the query struct to JSON.
func (q *query) marshal() ([]byte, error) {
	// q.Query = whitespaceRegex.ReplaceAllString(strings.TrimSpace(q.Query), " ")
	data, err := json.Marshal(q)
	if err != nil {
		log.Printf("Error marshalling query: %v", err)
		return nil, err
	}
	log.Printf("Marshalled query: %s", data)
	return data, nil
}

// getResponse sends the GraphQL query to the LeetCode API and returns the response.
func (q *query) getResponse() (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI("https://leetcode.com/graphql")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")

	body, err := q.marshal()
	if err != nil {
		log.Printf("Failed to marshal query: %v", err)
		return nil, utils.WithStack(err)
	}
	req.SetBody(body)

	log.Printf("Sending request to LeetCode API")
	resp := fasthttp.AcquireResponse()
	err = fasthttp.Do(req, resp)
	if err != nil {
		log.Printf("Error performing request: %v", err)
		fasthttp.ReleaseResponse(resp)
		return nil, utils.WithStack(err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		errMsg := fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		log.Printf("Request failed: %v", errMsg)
		fasthttp.ReleaseResponse(resp)
		return nil, errMsg
	}

	log.Printf("Received response with status code: %d", resp.StatusCode())
	log.Printf("Response body: %s", resp.Body())
	return resp, nil
}

// jsonResponse unmarshals the JSON response into the provided result interface.
func (q *query) jsonResponse(result any) error {
	resp, err := q.getResponse()
	if err != nil {
		return err
	}
	defer fasthttp.ReleaseResponse(resp)

	err = json.Unmarshal(resp.Body(), result)
	if err != nil {
		log.Printf("Error unmarshalling response body: %v", err)
		return err
	}
	log.Printf("Successfully unmarshalled response")
	return nil
}
