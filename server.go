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
	r.GET("/storage/:site/:bucket", BucketHandler)
	r.GET("/storage/:site/:bucket/:file", FileHandler)

	// POST routes
	r.POST("/storage/:site/:bucket", BucketPostHandler)
	r.POST("/storage/:site/:bucket/:file", FilePostHandler)

	// PUT routes
	r.PUT("/storage/:site/:bucket", BucketPutHandler)
	r.PUT("/storage/:site/:bucket/:file", FilePutHandler)

	// DELETE routes
	r.DELETE("/storage/:site/:bucket", BucketDeleteHandler)
	r.DELETE("/storage/:site/:bucket/:file", FileDeleteHandler)
	return r
}

func Server(configFile string) {
	r := setupRouter()
	sport := fmt.Sprintf(":%d", Config.Port)
	log.Printf("Start HTTP server %s", sport)
	r.Run(sport)
}
