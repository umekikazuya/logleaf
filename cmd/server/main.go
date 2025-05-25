package main

import (
	"github.com/umekikazuya/logleaf/internal/server"
)

func main() {
	leafHandler, port := server.InitializeDependencies()
	r := server.NewRouter(leafHandler)
	if err := r.Run(":" + port); err != nil {
		panic("failed to start server: " + err.Error())
	}
}
