package main

// Site object
type Site struct {
	Name         string `json:"name" binding:"required"`
	URL          string `json:"url" binding:"required"`
	Endpoint     string `json:"endpoint" binding:"required"`
	AccessKey    string `json:"access_key" binding:"required"`
	AccessSecret string `json:"access_secret" binding:"required"`
	UseSSL       bool   `json:"use_ssl"`
	Description  string `json:"description"`
}

// global site list we keep in server
var _sites []Site
