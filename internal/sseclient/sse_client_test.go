package sseclient

import (
	"io"
	"upfcc/internal/types"

	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSocialPost_GetValue(t *testing.T) {
	tests := []struct {
		name      string
		post      SocialPost
		dimension types.Dimension
		expected  int
	}{
		{
			name:      "Get Likes",
			post:      SocialPost{Likes: 10},
			dimension: types.Likes,
			expected:  10,
		},
		{
			name:      "Get Comments",
			post:      SocialPost{Comments: 5},
			dimension: types.Comments,
			expected:  5,
		},
		{
			name:      "Get Favorites",
			post:      SocialPost{Favorites: 7},
			dimension: types.Favorites,
			expected:  7,
		},
		{
			name:      "Get Retweets",
			post:      SocialPost{Retweets: 3},
			dimension: types.Retweets,
			expected:  3,
		},
		{
			name:      "Invalid Dimension",
			post:      SocialPost{Likes: 10},
			dimension: types.Dimension("invalid"),
			expected:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.post.GetValue(tt.dimension); got != tt.expected {
				t.Errorf("GetValue() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestNew(t *testing.T) {
	url := "http://example.com"
	client := New(url)
	if client == nil {
		t.Error("New() returned nil")
	} else {
		if client.url != url {
			t.Errorf("New() url = %s, want %s", client.url, url)
		}
	}
}

func TestSSEClient_ReadStream(t *testing.T) {
	tests := []struct {
		name     string
		server   *httptest.Server
		duration time.Duration
		expected []Post
	}{
		{
			name: "Valid Response",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/event-stream")
				w.Write([]byte("data: {\"post\":{\"timestamp\":1234567890,\"likes\":10}}\n\n"))
			})),
			duration: 5 * time.Second,
			expected: []Post{
				{
					Type: "post",
					Data: SocialPost{Timestamp: 1234567890, Likes: 10},
				},
			},
		},
		{
			name: "Invalid Response",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/event-stream")
				w.Write([]byte("invalid data\n\n"))
			})),
			duration: 5 * time.Second,
			expected: []Post{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()

			client := New(tt.server.URL)
			posts := client.ReadStream(tt.duration)

			var got []Post
			for post := range posts {
				got = append(got, post)
			}

			if len(got) != len(tt.expected) {
				t.Errorf("ReadStream() got %d posts, want %d", len(got), len(tt.expected))
			}

			for i, post := range got {
				if post.Type != tt.expected[i].Type || post.Data.Timestamp != tt.expected[i].Data.Timestamp || post.Data.Likes != tt.expected[i].Data.Likes {
					t.Errorf("ReadStream() got %v, want %v", post, tt.expected[i])
				}
			}
		})
	}
}

func TestSSEClient_scanResponse(t *testing.T) {
	tests := []struct {
		name     string
		response string
		expected []Post
	}{
		{
			name:     "Valid Response",
			response: "data: {\"post\":{\"timestamp\":1234567890,\"likes\":10}}\n\n",
			expected: []Post{
				{
					Type: "post",
					Data: SocialPost{Timestamp: 1234567890, Likes: 10},
				},
			},
		},
		{
			name:     "Invalid Response",
			response: "invalid data\n\n",
			expected: []Post{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := New("")
			resp := &http.Response{
				Body: io.NopCloser(bytes.NewReader([]byte(tt.response))),
			}

			postChan := make(chan Post)
			go client.scanResponse(resp, postChan)

			var got []Post
			for post := range postChan {
				got = append(got, post)
			}

			if len(got) != len(tt.expected) {
				t.Errorf("scanResponse() got %d posts, want %d", len(got), len(tt.expected))
			}

			for i, post := range got {
				if post.Type != tt.expected[i].Type || post.Data.Timestamp != tt.expected[i].Data.Timestamp || post.Data.Likes != tt.expected[i].Data.Likes {
					t.Errorf("scanResponse() got %v, want %v", post, tt.expected[i])
				}
			}
		})
	}
}

func TestSSEClient_processDataLine(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected []Post
	}{
		{
			name: "Valid Data Line",
			data: "{\"post\":{\"timestamp\":1234567890,\"likes\":10}}",
			expected: []Post{
				{
					Type: "post",
					Data: SocialPost{Timestamp: 1234567890, Likes: 10},
				},
			},
		},
		{
			name:     "Invalid Data Line",
			data:     "invalid data",
			expected: []Post{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := New("")
			postChan := make(chan Post)
			go func() {
				client.processDataLine(tt.data, postChan)
				close(postChan)
			}()

			var got []Post
			for post := range postChan {
				got = append(got, post)
			}

			if len(got) != len(tt.expected) {
				t.Errorf("processDataLine() got %d posts, want %d", len(got), len(tt.expected))
			}

			for i, post := range got {
				if post.Type != tt.expected[i].Type || post.Data.Timestamp != tt.expected[i].Data.Timestamp || post.Data.Likes != tt.expected[i].Data.Likes {
					t.Errorf("processDataLine() got %v, want %v", post, tt.expected[i])
				}
			}
		})
	}
}
