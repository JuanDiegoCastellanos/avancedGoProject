package api

import (
	db "avancedGo/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store    `json:"store"`
	router *gin.Engine `json:"router"`
}

// NewServer function to create a new server passing by parameter the store interface where there are all
// the functions to persist data or execute transactions
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccountById)
	router.GET("/accounts/", server.listAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)
	router.PUT("/accounts/", server.updateAccount)

	server.router = router
	return server
}
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
