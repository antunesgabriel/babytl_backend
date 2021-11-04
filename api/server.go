package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port   string
	server *gin.Engine
}

func BuildServer() Server {
	return Server{
		"5000",
		gin.Default(),
	}
}

func (s *Server) Run() {
	router := ConfigureRoutes(s.server)
	
	log.Fatalln(router.Run(":" + s.port))
	fmt.Println("Server listen in port: 5000")
}
