// Package sseclient provides a client for consuming Server-Sent Events (SSE)
// streams and processing social media posts from the stream.
package sseclient

import (
	"upfcc/internal/types"

	"bufio"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// SocialPost represents a social media post with various metrics such as likes,
// comments, favorites, and retweets.
type SocialPost struct {
	Timestamp int64 `json:"timestamp"`
	Likes     int   `json:"likes,omitempty"`
	Comments  int   `json:"comments,omitempty"`
	Favorites int   `json:"favorites,omitempty"`
	Retweets  int   `json:"retweets,omitempty"`
}

// Post represents a structured event containing a type and associated social post data.
type Post struct {
	Type string     `json:"type"`
	Data SocialPost `json:"data"`
}

// GetValue returns the value of a specific dimension (Likes, Comments, Favorites, Retweets)
// from the SocialPost. If the dimension is not recognized, it returns 0.
func (p *SocialPost) GetValue(dimension types.Dimension) int {
	switch dimension {
	case types.Likes:
		return p.Likes
	case types.Comments:
		return p.Comments
	case types.Favorites:
		return p.Favorites
	case types.Retweets:
		return p.Retweets
	default:
		return 0
	}
}

// SSEClient represents a client that connects to an SSE stream and reads events.
type SSEClient struct {
	url string // url is the endpoint of the SSE stream.
}

// New creates a new instance of SSEClient with the specified URL.
func New(url string) *SSEClient {
	return &SSEClient{url: url}
}

// ReadStream starts reading the SSE stream from the specified URL for a given duration.
// It returns a channel of Post structs that can be consumed by the caller.
func (c *SSEClient) ReadStream(duration time.Duration) chan Post {
	postChan := make(chan Post)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", c.url, nil)
		if err != nil {
			close(postChan)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			close(postChan)
			return
		}
		defer resp.Body.Close()

		c.scanResponse(resp, postChan)
	}()
	return postChan
}

// scanResponse reads the response body line by line and sends parsed events to the channel.
func (c *SSEClient) scanResponse(resp *http.Response, postChan chan<- Post) {
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			c.processDataLine(data, postChan)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading SSE stream: %v", err)
	}
	close(postChan)
}

// processDataLine parses a data line and sends the resulting events to the channel.
func (c *SSEClient) processDataLine(data string, postChan chan<- Post) {
	var event map[string]SocialPost
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		return
	}
	for eventType, post := range event {
		postChan <- Post{Type: eventType, Data: post}
	}
}
