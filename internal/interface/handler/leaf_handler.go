package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/umekikazuya/logleaf/internal/domain"
	"github.com/umekikazuya/logleaf/internal/usecase"
)

type LeafHandler struct {
	Usecase *usecase.LeafUsecase
}

func NewLeafHandler(u *usecase.LeafUsecase) *LeafHandler {
	return &LeafHandler{Usecase: u}
}

// Index /api/leaves
func (h *LeafHandler) ListLeaves(c *gin.Context) {
	opts := domain.ListOptions{}
	leaves, err := h.Usecase.ListLeaves(c.Request.Context(), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch leaves"})
		return
	}
	c.JSON(http.StatusOK, leaves)
}

// GET /api/leaves/:id
func (h *LeafHandler) GetLeaf(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "leaf ID is required"})
		return
	}

	leaf, err := h.Usecase.GetLeaf(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch leaf"})
		return
	}
	if leaf == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "leaf not found"})
		return
	}
	c.JSON(http.StatusOK, leaf)
}

// POST /api/leaves
func (h *LeafHandler) AddLeaf(c *gin.Context) {
	input := struct {
		ID       string `json:"id"`
		Title    string `json:"title" binding:"required"`
		URL      string `json:"url" binding:"required"`
		Platform string `json:"platform" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	var leafID string = input.ID
	if input.ID == "" {
		leafID = fmt.Sprintf("leaf-%s", uuid.New().String())
	}
	l := domain.NewLeaf(
		leafID,
		input.Title,
		input.URL,
		input.Platform,
	)

	app_err := h.Usecase.AddLeaf(c.Request.Context(), l)
	if app_err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add leaf"})
		return
	}
	c.JSON(http.StatusCreated, l)
}

// PATCH /api/leaves/:id
func (h *LeafHandler) UpdateLeaf(c *gin.Context) {
	id := c.Param("id")

	var input domain.Leaf
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if err := h.Usecase.UpdateLeaf(c.Request.Context(), id, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// DELETE /api/leaves/:id
func (h *LeafHandler) DeleteLeaf(c *gin.Context) {
	id := c.Param("id")
	if err := h.Usecase.DeleteLeaf(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
