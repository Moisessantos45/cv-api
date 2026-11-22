package post

import (
	"cv_api/internal/features/profile"
	"cv_api/internal/models"
	"cv_api/internal/shared/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	uc             PostService
	profileService profile.ProfileService
}

func NewPostHandlers(service PostService, profileService profile.ProfileService) *PostHandler {
	return &PostHandler{uc: service, profileService: profileService}
}

func (h *PostHandler) GetAllStates(c *gin.Context) {
	states, err := h.uc.GetAllStates(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to retrieve post states"})
		return
	}

	c.JSON(200, gin.H{"data": states, "message": "Post states retrieved successfully"})
}

func (h *PostHandler) GetAll(c *gin.Context) {
	page, pageSize, _, err := utils.ValidateQueryPagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	posts, err := h.uc.GetAll(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch posts: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    posts,
		"message": "Data Success",
	})
}

func (h *PostHandler) GetAllPublic(c *gin.Context) {
	page, pageSize, _, err := utils.ValidateQueryPagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	posts, err := h.uc.GetAllPublic(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch posts: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    posts,
		"message": "Data Success",
	})
}

func (h *PostHandler) GetAllRecents(c *gin.Context) {
	maxSize, err := utils.ValidateParamsQuery[int](c, "max")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	posts, err := h.uc.GetAllRecients(c.Request.Context(), maxSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch posts: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    posts,
		"message": "Data Success",
	})
}

func (h *PostHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "slug es requerido"})
		return
	}

	post, err := h.uc.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch post: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": post,
	})
}

func (h *PostHandler) GetBySlugPrivate(c *gin.Context) {
	_, authID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid ID: " + err.Error()})
		return
	}

	slug := c.Param("slug")
	if slug == "" {
		c.JSON(400, gin.H{"message": "Slug is required"})
		return
	}

	post, err := h.uc.GetBySlugPrivate(c.Request.Context(), slug, authID)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to retrieve post"})
		return
	}

	c.JSON(200, gin.H{"data": post, "message": "Post retrieved successfully"})
}

func (h *PostHandler) GetBySlugPublic(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(400, gin.H{"message": "Slug is required"})
		return
	}

	post, err := h.uc.GetBySlugPublic(c.Request.Context(), slug)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to retrieve post"})
		return
	}

	c.JSON(200, gin.H{"data": post, "message": "Post retrieved successfully"})
}

func (h *PostHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	log.Println("Entrando a post handler")

	_, authID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		log.Printf("Error al bindear JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Datos del post inválidos", "error": err.Error()})
		return
	}

	log.Printf("Post recibida: %+v", post)
	err = h.uc.Create(ctx, authID, &post)
	if err != nil {
		log.Printf("Error al crear el post: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": post, "message": "Post creado correctamente"})
}

func (h *PostHandler) Update(c *gin.Context) {
	_, authID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	log.Println("Entrando a post handler")
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "slug es requerido"})
		return
	}

	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		log.Printf("Error al bindear JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Datos del post inválidos", "error": err.Error()})
		return
	}

	log.Printf("Post recibida: %+v", post)

	err = h.uc.Update(c.Request.Context(), authID, slug, &post)
	if err != nil {
		log.Printf("Error al actualizar el post: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post actualizado correctamente"})
}
