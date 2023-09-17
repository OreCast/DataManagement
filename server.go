package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// examples: https://go.dev/doc/tutorial/web-service-gin

// helper function to setup our server router
func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// GET routes
	r.GET("/storage", StorageHandler)
	r.GET("/storage/:site", SiteHandler)
	r.GET("/storage/:site/:dataset", DatasetHandler)
	r.GET("/storage/:site/:dataset/:file", FileHandler)

	// POST routes
	r.POST("/storage/:site/:dataset", DatasetPostHandler)
	r.POST("/storage/:site/:dataset/:file", FilePostHandler)

	return r
}

func Server(configFile string) {
	r := setupRouter()
	sport := fmt.Sprintf(":%d", Config.Port)
	log.Printf("Start HTTP server %s", sport)
	r.Run(sport)
}
