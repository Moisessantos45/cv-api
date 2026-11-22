package experience

import (
	"cv_api/internal/models"
	"cv_api/internal/shared/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ExperienceHandler struct {
	service ExperienceService
}

func NewExperienceHandler(service ExperienceService) *ExperienceHandler {
	return &ExperienceHandler{service: service}
}

func (h *ExperienceHandler) GetExperiences(c *gin.Context) {
	page, pageSize, _, err := utils.ValidateQueryPagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	experiences, err := h.service.GetAllExperiences(c.Request.Context(), page, pageSize, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch experiences: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    experiences,
		"message": "Data Success",
	})
}

func (h *ExperienceHandler) GetAllExperiences(c *gin.Context) {
	page, pageSize, _, err := utils.ValidateQueryPagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	experiences, err := h.service.GetAllExperiences(c.Request.Context(), page, pageSize, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch experiences: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    experiences,
		"message": "Data Success",
	})
}

func (h *ExperienceHandler) GetAllExperiencesRecents(c *gin.Context) {
	maxSizeStr := c.DefaultQuery("max", "3")

	max, err := strconv.Atoi(maxSizeStr)
	if err != nil || max < 1 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Parámetro 'page' inválido: debe ser un número entero positivo",
		})
		return
	}

	experiences, err := h.service.GetAllExperiencesRecients(c.Request.Context(), max)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch experiences: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    experiences,
		"message": "Data Success",
	})
}

func (h *ExperienceHandler) GetExperienceById(c *gin.Context) {

	id, err := utils.ValidateParamsId(c, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	experience, err := h.service.GetExperienceById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch experience: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": experience,
	})
}

func (h *ExperienceHandler) AddExperience(c *gin.Context) {
	ctx := c.Request.Context()
	log.Println("Entrando a experience handler")

	var experience models.Experience
	if err := c.ShouldBindJSON(&experience); err != nil {
		log.Printf("Error al bindear JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Datos de la experiencia inválidos", "error": err.Error()})
		return
	}

	log.Printf("Experience recibida: %+v", experience)

	err := h.service.AddExperience(ctx, &experience)
	if err != nil {
		log.Printf("Error al crear la experiencia: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": experience, "message": "Experiencia creada correctamente"})
}

func (h *ExperienceHandler) UpdateExperience(c *gin.Context) {

	log.Println("Entrando a experience handler")
	id, err := utils.ValidateParamsId(c, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var experience models.Experience
	if err := c.ShouldBindJSON(&experience); err != nil {
		log.Printf("Error al bindear JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Datos de experiencia inválidos", "error": err.Error()})
		return
	}

	log.Printf("Experience recibida: %+v", experience)

	err = h.service.Update(c.Request.Context(), id, &experience)
	if err != nil {
		log.Printf("Error al actualizar experiencia: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": experience, "message": "Experiencia actualizada correctamente"})
}
