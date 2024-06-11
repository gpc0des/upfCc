// Package aggregator provides functionality to aggregate social media posts
// data over a specified duration and calculate statistical results such as
// total posts, minimum timestamp, maximum timestamp, and average value for a
// given dimension.
package aggregator

import (
	"upfcc/internal/sseclient"
	"upfcc/internal/types"

	"time"
)

type SSEClientInterface interface {
	ReadStream(duration time.Duration) chan sseclient.Post
}

// Aggregator is responsible for aggregating data from an SSE client.
type Aggregator struct {
	sseClient SSEClientInterface
}

// New creates a new Aggregator with the provided SSE client.
//
// Parameters:
//   - sseClient: An instance of SSEClientInterface to read the stream of posts.
//
// Returns:
//   - A pointer to the newly created Aggregator.
func New(sseClient SSEClientInterface) *Aggregator {
	return &Aggregator{sseClient: sseClient}
}

// AnalysisResult holds the results of the aggregation process.
type AnalysisResult struct {
	TotalPosts   int     `json:"total_posts"`       // Total number of posts analyzed
	MinTimestamp int64   `json:"minimum_timestamp"` // The timestamp of the first post analyzed
	MaxTimestamp int64   `json:"maximum_timestamp"` // The timestamp of the last post analyzed
	AvgValue     float64 `json:"avg_value"`         // Average value of the specified dimension
}

// AggregateData reads social media posts for a specified duration and calculates
// the total number of posts, minimum timestamp, maximum timestamp, and average value
// for the specified dimension.
//
// Parameters:
//   - duration: The duration for which to read and aggregate posts.
//   - dimension: The dimension for which to calculate the average value (e.g., likes, comments).
//   - resultChan: A channel to send the result of the aggregation.
func (a *Aggregator) AggregateData(duration time.Duration, dimension types.Dimension, resultChan chan AnalysisResult) {
	var analysisResult AnalysisResult
	totalValue := 0
	postChan := a.sseClient.ReadStream(duration) // ReadStream will close the channel after the duration has elapsed

	for post := range postChan {
		if analysisResult.TotalPosts == 0 {
			analysisResult.MinTimestamp = post.Data.Timestamp
		}
		analysisResult.MaxTimestamp = post.Data.Timestamp
		totalValue += post.Data.GetValue(dimension)
		analysisResult.TotalPosts++
	}

	// Calculate and send the result after the postChan is closed
	if analysisResult.TotalPosts > 0 {
		analysisResult.AvgValue = float64(totalValue) / float64(analysisResult.TotalPosts)
	} else {
		analysisResult.AvgValue = 0
	}

	resultChan <- analysisResult
}
