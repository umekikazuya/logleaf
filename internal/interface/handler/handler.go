package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/umekikazuya/logleaf/internal/application"
	"github.com/umekikazuya/logleaf/internal/domain"
)

type LeafHandler struct {
	Usecase *application.LeafUsecase
}

func NewLeafHandler(u *application.LeafUsecase) *LeafHandler {
	return &LeafHandler{Usecase: u}
}

// Index /api/leaves
func (h *LeafHandler) ListLeaves(c *gin.Context) {
	// Parse query parameters for filtering options
	opts := domain.ListOptions{
		Limit: 100, // default limit
	}
	leaves, err := h.Usecase.ListLeaves(c.Request.Context(), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch leaves"})
		return
	}
	// Convert to output DTOs
	outputDTOs := make([]*application.LeafOutputDTO, len(leaves))
	for i, leaf := range leaves {
		outputDTOs[i] = application.LeafDomainToOutputDTO(&leaf)
	}
	c.JSON(http.StatusOK, outputDTOs)
}

// GET /api/leaves/:id
func (h *LeafHandler) GetLeaf(c *gin.Context) {
	leaf, err := h.Usecase.GetLeaf(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if leaf == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "leaf not found"})
		return
	}
	c.JSON(http.StatusOK, application.LeafDomainToOutputDTO(leaf))
}

// POST /api/leaves
func (h *LeafHandler) AddLeaf(c *gin.Context) {
	// Request
	var req CreateLeafRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	// Convert to Dto
	inputDto := application.LeafInputDTO{
		Title:    req.Title,
		URL:      req.URL,
		Platform: req.Platform,
		Tags:     req.Tags,
	}
	// Add Leaf
	leaf, err := h.Usecase.AddLeaf(c.Request.Context(), &inputDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Response
	c.JSON(http.StatusCreated, application.LeafDomainToOutputDTO(leaf))
}

// PATCH /api/leaves/:id
func (h *LeafHandler) UpdateLeaf(c *gin.Context) {
	id := c.Param("id")

	var req UpdateLeafRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Convert to Dto
	inputDto := application.LeafInputDTO{
		ID:       id,
		Title:    req.Title,
		URL:      req.URL,
		Platform: req.Platform,
		Tags:     req.Tags,
	}

	err := h.Usecase.UpdateLeaf(c.Request.Context(), &inputDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// PATCH /api/leaves/:id/read
func (h *LeafHandler) ReadLeaf(c *gin.Context) {
	id := c.Param("id")

	err := h.Usecase.ReadLeaf(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

// DELETE /api/leaves/:id
func (h *LeafHandler) DeleteLeaf(c *gin.Context) {
	id := c.Param("id")
	if err := h.Usecase.DeleteLeaf(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
}
