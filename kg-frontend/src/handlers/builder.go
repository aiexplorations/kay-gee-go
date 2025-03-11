package handlers

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"kg-frontend/src/models"
	"kg-frontend/src/utils"
)

// StartBuilder starts the knowledge graph builder with the provided parameters
func StartBuilder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params models.BuilderParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid parameters: " + err.Error()})
			return
		}

		// Build command to start the builder
		cmd := exec.Command(
			"/bin/sh",
			utils.GetScriptPath("start-builder.sh"),
			"--seed", params.SeedConcept,
			"--max-nodes", fmt.Sprintf("%d", params.MaxNodes),
			"--timeout", fmt.Sprintf("%d", params.Timeout),
			"--random-relationships", fmt.Sprintf("%d", params.RandomRelationships),
			"--concurrency", fmt.Sprintf("%d", params.Concurrency),
		)

		// Run the command
		output, err := cmd.CombinedOutput()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: fmt.Sprintf("Failed to start builder: %v - %s", err, string(output)),
			})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, models.Response{
			Status:  "success",
			Message: "Builder started successfully",
			Output:  string(output),
		})
	}
}

// StopBuilder stops the knowledge graph builder
func StopBuilder() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Build command to stop the builder
		cmd := exec.Command("/bin/sh", utils.GetScriptPath("stop-builder.sh"))

		// Run the command
		output, err := cmd.CombinedOutput()
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: fmt.Sprintf("Failed to stop builder: %v - %s", err, string(output)),
			})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, models.Response{
			Status:  "success",
			Message: "Builder stopped successfully",
			Output:  string(output),
		})
	}
} 