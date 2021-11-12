package api

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port   string
	server *gin.Engine
}

func BuildServer() Server {
	port := os.Getenv("BTL_PORT")

	if port == "" {
		log.Fatalln("BTL_PORT NOT DEFINED")
	}

	return Server{
		port,
		gin.Default(),
	}
}

func (s *Server) Run() {
	router := ConfigureRoutes(s.server)

	log.Fatalln(router.Run(":" + s.port))
	fmt.Println("Server listen in port: 5000")
}
