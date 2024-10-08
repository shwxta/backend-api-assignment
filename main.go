package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Input struct {
	Numbers []int `json:"numbers" binding:"required"`
	Target  int   `json:"target" binding:"required"`
}

func findPairs(c *gin.Context) {
	var input Input

	// Input validation
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	numbers := input.Numbers
	target := input.Target

	// Hash map to store indices of elements
	indexMap := make(map[int]int)
	solutions := [][]int{}

	// Iterate through the array
	for i, num := range numbers {
		complement := target - num
		if j, ok := indexMap[complement]; ok {
			// Pair found, add the indices
			solutions = append(solutions, []int{j, i})
		}
		indexMap[num] = i
	}

	// Return the result
	c.JSON(http.StatusOK, gin.H{"solutions": solutions})
}

func main() {
	r := gin.Default()

	// Define POST route
	r.POST("/find-pairs", findPairs)

	// Start server
	r.Run(":8080")
}
