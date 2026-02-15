package routes

import (
	"cv_api/config"
	"cv_api/internal/handlers"
	"cv_api/internal/repository"
	"cv_api/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

func ProjectRoutes(api *gin.RouterGroup) {
	client := config.Client
	clientRd := config.Rdb
	projectRepo := repository.NewProjectRepository(client, clientRd, time.Minute*15)
	projectSvc := services.NewProjectService(projectRepo)
	projectHandler := handlers.NewProjectHandlers(projectSvc)

	project := api.Group("/project")
	{
		project.GET("", projectHandler.GetProjects)
		project.GET("/:id", projectHandler.GetProjectById)
		project.PUT("/likes/:id", projectHandler.UpdateCounter)

		// auth.POST("/register", registerHandler)
	}
}
