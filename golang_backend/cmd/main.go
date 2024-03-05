package main

import "clean_arch/internal/app"

// @title Ufaelectrto API
// @version 1.3
// @description API Server for Ufaelectrto market

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

// @contact.name   API Creator
// @contact.url    https://t.me/amirich18
// @contact.email    adamstradvers@gmail.com
func main() {
	app.Run()
}
