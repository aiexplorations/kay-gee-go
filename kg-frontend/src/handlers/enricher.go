package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"kg-frontend/src/models"
	"kg-frontend/src/utils"
)

// StartEnricher starts the knowledge graph enricher with the provided parameters
func StartEnricher(runner utils.CommandRunnerInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params models.EnricherParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid parameters: " + err.Error()})
			return
		}

		// Get the script path
		scriptPath := utils.GetScriptPath("start-enricher.sh")

		// Build command to start the enricher
		args := []string{
			scriptPath,
			"--batch-size", fmt.Sprintf("%d", params.BatchSize),
			"--interval", fmt.Sprintf("%d", params.Interval),
			"--max-relationships", fmt.Sprintf("%d", params.MaxRelationships),
			"--concurrency", fmt.Sprintf("%d", params.Concurrency),
		}

		// Run the command
		output, err := runner.RunCommand("bash", args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: fmt.Sprintf("Failed to start enricher: %v - %s", err, string(output)),
			})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, models.Response{
			Status:  "success",
			Message: "Enricher started successfully",
			Output:  string(output),
		})
	}
}

// StopEnricher stops the knowledge graph enricher
func StopEnricher(runner utils.CommandRunnerInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the script path
		scriptPath := utils.GetScriptPath("stop-enricher.sh")

		// Run the command
		output, err := runner.RunCommand("bash", scriptPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: fmt.Sprintf("Failed to stop enricher: %v - %s", err, string(output)),
			})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, models.Response{
			Status:  "success",
			Message: "Enricher stopped successfully",
			Output:  string(output),
		})
	}
} 