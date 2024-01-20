package api

import (
	db "avancedGo/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store   `json:"store,omitempty"`
	router *gin.Engine `json:"router,omitempty"`
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)

	server.router = router
	return server
}
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}