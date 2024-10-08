package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Helper function to sort pairs (to handle order-insensitive comparison)
func sortPairs(pairs [][]int) [][]int {
	for i := range pairs {
		sort.Ints(pairs[i])
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i][0] < pairs[j][0]
	})
	return pairs
}

func TestFindPairs(t *testing.T) {
	r := gin.Default()
	r.POST("/find-pairs", findPairs)

	tests := []struct {
		name         string
		body         map[string]interface{}
		expectedCode int
		expectedResp map[string]interface{}
	}{
		{
			name: "Valid input with pairs",
			body: map[string]interface{}{
				"numbers": []int{1, 2, 3, 4, 5},
				"target":  6,
			},
			expectedCode: http.StatusOK,
			expectedResp: map[string]interface{}{
				"solutions": [][]int{{0, 4}, {1, 3}},
			},
		},
		{
			name: "Valid input with no pairs",
			body: map[string]interface{}{
				"numbers": []int{1, 2, 3},
				"target":  10,
			},
			expectedCode: http.StatusOK,
			expectedResp: map[string]interface{}{
				"solutions": [][]int{},
			},
		},
		{
			name: "Empty array",
			body: map[string]interface{}{
				"numbers": []int{},
				"target":  6,
			},
			expectedCode: http.StatusOK,
			expectedResp: map[string]interface{}{
				"solutions": [][]int{},
			},
		},
		{
			name: "Invalid input",
			body: map[string]interface{}{
				"numbers": "invalid",
				"target":  6,
			},
			expectedCode: http.StatusBadRequest,
			expectedResp: map[string]interface{}{
				"error": "Invalid input",
			},
		},
		{
			name: "Missing numbers key",
			body: map[string]interface{}{
				"target": 6,
			},
			expectedCode: http.StatusBadRequest,
			expectedResp: map[string]interface{}{
				"error": "Invalid input",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare request
			jsonValue, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/find-pairs", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			// Record response
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedCode, w.Code)

			// Unmarshal actual response
			var actualResp map[string]interface{}
			_ = json.Unmarshal(w.Body.Bytes(), &actualResp)

			// Handle the case where "solutions" contains unordered pairs
			if solutions, ok := actualResp["solutions"].([]interface{}); ok {
				// Convert to [][]int for sorting
				actualPairs := make([][]int, len(solutions))
				for i, solution := range solutions {
					pair := solution.([]interface{})
					actualPairs[i] = []int{int(pair[0].(float64)), int(pair[1].(float64))}
				}

				expectedPairs := tt.expectedResp["solutions"].([][]int)

				// Sort both actual and expected pairs
				actualPairs = sortPairs(actualPairs)
				expectedPairs = sortPairs(expectedPairs)

				// Replace the "solutions" key for comparison
				actualResp["solutions"] = actualPairs
				tt.expectedResp["solutions"] = expectedPairs
			}

			// Compare responses
			expectedJSON, _ := json.Marshal(tt.expectedResp)
			actualJSON, _ := json.Marshal(actualResp)
			assert.JSONEq(t, string(expectedJSON), string(actualJSON))
		})
	}
}
