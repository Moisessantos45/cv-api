package stream

import (
	"cv_api/internal/models"
	"cv_api/internal/shared/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VideoHandlers struct {
	service VideoService
}

func NewVideoHandlers(service VideoService) *VideoHandlers {
	return &VideoHandlers{service: service}
}

func (h *VideoHandlers) GetAllData(c *gin.Context) {
	projects, err := h.service.GetAllData(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=projects.json")
	c.Header("Content-Type", "application/json")

	c.JSON(http.StatusOK, projects)
}

func (h *VideoHandlers) GetAll(c *gin.Context) {
	page, pageSize, _, err := utils.ValidateQueryPagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	results, err := h.service.GetAll(c.Request.Context(), page, pageSize, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch videos: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    results,
		"message": "Data Success",
	})
}

func (h *VideoHandlers) GetAllRecents(c *gin.Context) {
	maxSize, err := utils.ValidateParamsQuery[int](c, "max")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	projects, err := h.service.GetAllRecients(c.Request.Context(), maxSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch videos: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    projects,
		"message": "Data Success",
	})
}

func (h *VideoHandlers) GetByID(c *gin.Context) {

	id, err := utils.ValidateParamsId(c, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	project, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch videos: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": project,
	})
}

func (h *VideoHandlers) Create(c *gin.Context) {
	ctx := c.Request.Context()
	log.Println("Entrando a videos handler")

	var project models.Video
	if err := c.ShouldBindJSON(&project); err != nil {
		log.Printf("Error al bindear JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Datos del videos inválidos", "error": err.Error()})
		return
	}

	log.Printf("Video recibida: %+v", project)

	err := h.service.Create(ctx, &project)
	if err != nil {
		log.Printf("Error al crear el video: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": project, "message": "Video creada correctamente"})
}

func (h *VideoHandlers) CreateAll(c *gin.Context) {
	ctx := c.Request.Context()

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no se envió el archivo json",
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no se pudo abrir el archivo",
		})
		return
	}
	defer file.Close()

	var videos []models.Video
	if err := json.NewDecoder(file).Decode(&videos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "json inválido",
		})
		return
	}

	if len(videos) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no se proporcionaron videos",
		})
		return
	}

	if err := h.service.CreateAll(ctx, videos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "videos creados correctamente",
	})
}

func (h *VideoHandlers) Update(c *gin.Context) {

	log.Println("Entrando a video handler")
	id, err := utils.ValidateParamsId(c, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var video models.Video
	if err := c.ShouldBindJSON(&video); err != nil {
		log.Printf("Error al bindear JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Datos de video inválidos", "error": err.Error()})
		return
	}

	log.Printf("Video recibida: %+v", video)

	err = h.service.Update(c.Request.Context(), id, &video)
	if err != nil {
		log.Printf("Error al crear proyecto: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": video, "message": "Video actulizado correctamente"})
}