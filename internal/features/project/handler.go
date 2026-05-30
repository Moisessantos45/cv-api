package project

import (
	"cv_api/internal/models"
	"cv_api/internal/shared/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	service ProjectService
}

func NewProjectHandlers(service ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) GetAllData(c *gin.Context) {
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

func (h *ProjectHandler) GetAll(c *gin.Context) {
	page, pageSize, _, err := utils.ValidateQueryPagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	projects, err := h.service.GetAll(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch projects: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    projects,
		"message": "Data Success",
	})
}

func (h *ProjectHandler) GetAllPublic(c *gin.Context) {
	page, pageSize, _, err := utils.ValidateQueryPagination(c)
	if err != nil {
		log.Printf("Error validating pagination query parameters: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	projects, err := h.service.GetAllPublic(c.Request.Context(), page, pageSize)
	if err != nil {
		log.Printf("Error fetching public projects: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch projects: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    projects,
		"message": "Data Success",
	})
}

func (h *ProjectHandler) GetALLRecents(c *gin.Context) {
	maxSize, err := utils.ValidateParamsQuery[int](c, "max")
	if err != nil {
		log.Printf("Error validating 'max' query parameter: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	projects, err := h.service.GetAllRecents(c.Request.Context(), maxSize)
	if err != nil {
		log.Printf("Error fetching recent projects: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch projects: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    projects,
		"message": "Data Success",
	})
}

func (h *ProjectHandler) GetByID(c *gin.Context) {

	id, err := utils.ValidateParamsId(c, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	project, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch project: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": project,
	})
}

func (h *ProjectHandler) GetBySlugPublic(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(400, gin.H{"message": "Slug is required"})
		return
	}

	project, err := h.service.GetBySlugPublic(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch project: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    project,
		"message": "Data Success",
	})
}

func (h *ProjectHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(400, gin.H{"message": "Slug is required"})
		return
	}

	project, err := h.service.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch project: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    project,
		"message": "Data Success",
	})
}

func (h *ProjectHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	log.Println("Entrando a project handler")

	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		log.Printf("Error al bindear JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Datos del proyectto inválidos", "error": err.Error()})
		return
	}

	log.Printf("Proyecto recibida: %+v", project)

	err := h.service.Create(ctx, &project)
	if err != nil {
		log.Printf("Error al crear el proyecto: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": project, "message": "Proyecto creada correctamente"})
}

func (h *ProjectHandler) CreateAll(c *gin.Context) {
	ctx := c.Request.Context()

	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error al obtener el archivo: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "no se envió el archivo json",
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("Error al abrir el archivo: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "no se pudo abrir el archivo",
		})
		return
	}
	defer file.Close()

	var projects []models.Project
	if err := json.NewDecoder(file).Decode(&projects); err != nil {
		log.Printf("Error al decodificar el archivo JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "json inválido",
		})
		return
	}

	if len(projects) == 0 {
		log.Printf("No se proporcionaron proyectos en el archivo JSON")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "no se proporcionaron proyectos",
		})
		return
	}

	if err := h.service.CreateAll(ctx, projects); err != nil {
		log.Printf("Error al crear los proyectos: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "proyectos creados correctamente",
	})
}

func (h *ProjectHandler) Update(c *gin.Context) {
	log.Println("Entrando a project handler")
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(400, gin.H{"message": "Slug is required"})
		return
	}

	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		log.Printf("Error al bindear JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Datos de project inválidos", "error": err.Error()})
		return
	}

	log.Printf("Project recibida: %+v", project)

	err := h.service.Update(c.Request.Context(), slug, &project)
	if err != nil {
		log.Printf("Error al crear proyecto: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Projecto actulizado correctamente"})
}

func (h *ProjectHandler) UpdateState(c *gin.Context) {
	id, err := utils.ValidateParamsId(c, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var req struct {
		StateID uint64 `json:"state_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Datos de estado inválidos", "error": err.Error()})
		return
	}

	err = h.service.UpdateState(c.Request.Context(), id, req.StateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estado actualizado correctamente"})
}

func (h *ProjectHandler) UpdateCounter(c *gin.Context) {
	id, err := utils.ValidateParamsId(c, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = h.service.UpdateCounter(c.Request.Context(), id, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update counter: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    "ok",
		"message": "Counter updated successfully",
	})
}
