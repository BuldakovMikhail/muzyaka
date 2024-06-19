package main

import "src/cmd/muzyaka"

// @title Muzyaka API
// @version 1.0
// @description API Server for musical service

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	muzyaka.App()
}
