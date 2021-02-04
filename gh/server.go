package gh

import (
	"fmt"

	"github.com/lucheng0127/groundhog/common"
)

// Server struct of server
type Server struct {
}

// Launch server
func (server *Server) Launch() error {
	fmt.Println("Launch server")
	return nil
}

// NewServer create server and return
func NewServer(cfg *common.ServerConfig) (*Server, error) {
	common.Logger.Debugf("Create server with config: %v\n", cfg)
	return new(Server), nil
}
