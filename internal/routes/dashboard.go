package routes

import (
	"cv_api/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDashboardMetrics(c *gin.Context) {
	ctx := c.Request.Context()

	projectStats, err := projectUC.GetStats(ctx)
	if err != nil {
		log.Printf("Error getting project stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error getting project stats"})
		return
	}

	postStats, err := postUC.GetStats(ctx)
	if err != nil {
		log.Printf("Error getting post stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error getting post stats"})
		return
	}

	videoStats, err := videoUC.GetStats(ctx)
	if err != nil {
		log.Printf("Error getting video stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error getting video stats"})
		return
	}

	metrics := models.DashboardMetrics{
		Projects: *projectStats,
		Posts:    *postStats,
		Videos:   *videoStats,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    metrics,
		"message": "Success",
	})
}
