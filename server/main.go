package main

import (
	"ForgettiServer/routes"
	"ForgettiServer/services"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	serviceContainer := services.CreateServiceContainer()
	routes.AddEncRoutes(r, serviceContainer)

	r.Run(":8080")
}
