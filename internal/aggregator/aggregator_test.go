package aggregator

import (
	"upfcc/internal/sseclient"
	"upfcc/internal/testingTools"
	"upfcc/internal/types"

	"testing"
	"time"
)

func TestAggregateData(t *testing.T) {
	tests := []struct {
		name      string
		posts     []sseclient.Post
		duration  time.Duration
		dimension types.Dimension
		wantPosts int
		wantAvg   float64
		wantMinTs int64
		wantMaxTS int64
	}{
		{
			name:      "NoPosts",
			posts:     []sseclient.Post{},
			duration:  5 * time.Second,
			dimension: types.Likes,
			wantPosts: 0,
			wantAvg:   0,
			wantMinTs: 0,
			wantMaxTS: 0,
		},
		{
			name: "NoMatchingDimension",
			posts: []sseclient.Post{
				{
					Type: "instagram_media",
					Data: sseclient.SocialPost{
						Timestamp: testingTools.FakeTimestamp,
						Likes:     10,
					},
				},
			},
			duration:  5 * time.Second,
			dimension: types.Dimension("comments"),
			wantPosts: 1,
			wantAvg:   0,
			wantMinTs: testingTools.FakeTimestamp,
			wantMaxTS: testingTools.FakeTimestamp,
		},
		{
			name: "ValidDimension",
			posts: []sseclient.Post{
				{
					Type: "instagram_media",
					Data: sseclient.SocialPost{
						Timestamp: testingTools.FakeTimestamp,
						Likes:     10,
					},
				},
			},
			duration:  5 * time.Second,
			dimension: types.Likes,
			wantPosts: 1,
			wantAvg:   10,
			wantMinTs: testingTools.FakeTimestamp,
			wantMaxTS: testingTools.FakeTimestamp,
		},
		{
			name: "MultiplePosts",
			posts: []sseclient.Post{
				{
					Type: "instagram_media",
					Data: sseclient.SocialPost{
						Timestamp: testingTools.FakeTimestamp,
						Likes:     10,
						Comments:  5,
					},
				},
				{
					Type: "instagram_media",
					Data: sseclient.SocialPost{
						Timestamp: testingTools.FakeTimestamp2,
						Likes:     20,
						Comments:  10,
					},
				},
			},
			duration:  5 * time.Second,
			dimension: types.Comments,
			wantPosts: 2,
			wantAvg:   7.5,
			wantMinTs: testingTools.FakeTimestamp,
			wantMaxTS: testingTools.FakeTimestamp2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockSSEClient{posts: tt.posts}
			aggregator := New(mockClient)
			resultChan := make(chan AnalysisResult)

			go aggregator.AggregateData(tt.duration, tt.dimension, resultChan)
			result := <-resultChan

			if result.TotalPosts != tt.wantPosts {
				t.Errorf("Expected total posts to be %d, got %d", tt.wantPosts, result.TotalPosts)
			}
			if result.AvgValue != tt.wantAvg {
				t.Errorf("Expected average value to be %f, got %f", tt.wantAvg, result.AvgValue)
			}
			if result.MinTimestamp != tt.wantMinTs {
				t.Errorf("Expected MinTimestamp to be %d, got %d", tt.wantMinTs, result.MinTimestamp)
			}
			if result.MaxTimestamp != tt.wantMaxTS {
				t.Errorf("Expected MaxTimestamp to be %d, got %d", tt.wantMaxTS, result.MaxTimestamp)
			}
		})
	}
}

/////// Helpers

// MockSSEClient simulates an SSE client for testing purposes.
type MockSSEClient struct {
	posts []sseclient.Post
}

// ReadStream simulates reading a stream of posts for the specified duration.
func (m *MockSSEClient) ReadStream(duration time.Duration) chan sseclient.Post {
	postChan := make(chan sseclient.Post)
	go func() {
		defer close(postChan)
		for _, post := range m.posts {
			postChan <- post
		}
	}()
	return postChan
}
