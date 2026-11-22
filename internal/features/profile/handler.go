package profile

import (
	"cv_api/internal/models"
	"cv_api/internal/shared/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	service ProfileService
}

func NewProfileHandler(service ProfileService) *ProfileHandler {
	return &ProfileHandler{service: service}
}

func (h *ProfileHandler) GetById(c *gin.Context) {
	_, authID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	profile, err := h.service.GetByID(c.Request.Context(), authID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch profile: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    profile,
		"message": "Data Success",
	})
}

func (h *ProfileHandler) Create(c *gin.Context) {
	_, authID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	ctx := c.Request.Context()
	log.Println("Entrando a profile handler")

	var profile models.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		log.Printf("Error al bindear JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Datos del perfil inválidos", "error": err.Error()})
		return
	}

	log.Printf("Profile recibida: %+v", profile)

	err = h.service.Create(ctx, authID, &profile)
	if err != nil {
		log.Printf("Error al crear el perfil: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": profile, "message": "Perfil creado correctamente"})
}

func (h *ProfileHandler) Update(c *gin.Context) {

	log.Println("Entrando a profile handler")
	_, authID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	var profile models.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		log.Printf("Error al bindear JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Datos de perfil inválidos", "error": err.Error()})
		return
	}

	log.Printf("Profile recibida: %+v", profile)

	err = h.service.Update(c.Request.Context(), authID, &profile)
	if err != nil {
		log.Printf("Error al actualizar perfil: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": profile, "message": "Perfil actualizado correctamente"})
}

func (h *ProfileHandler) Contact(c *gin.Context) {
	var input struct {
		Name    string `json:"name" binding:"required"`
		Email   string `json:"email" binding:"required"`
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Datos de contacto inválidos", "error": err.Error()})
		return
	}

	err := h.service.Contact(c.Request.Context(), input.Name, input.Email, input.Message)
	if err != nil {
		log.Printf("Error al enviar mensaje de contacto: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mensaje enviado correctamente"})
}
