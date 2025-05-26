package server

import (
	"github.com/gin-gonic/gin"
	"github.com/umekikazuya/logleaf/internal/interface/handler"
)

// ルーティングを設定
func NewRouter(leafHandler *handler.LeafHandler) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.GET("/leaves", leafHandler.ListLeaves)
		api.POST("/leaves", leafHandler.AddLeaf)
		api.GET("/leaves/:id", leafHandler.GetLeaf)
		api.PATCH("/leaves/:id", leafHandler.UpdateLeaf)
		api.PATCH("/leaves/:id/read", leafHandler.ReadLeaf)
		api.DELETE("/leaves/:id", leafHandler.DeleteLeaf)
	}
	return r
}
