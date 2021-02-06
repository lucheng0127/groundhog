package gh

import (
	"fmt"
	"net"
	"strings"

	"github.com/vishvananda/netlink"

	"github.com/songgao/water"

	"github.com/lucheng0127/groundhog/common"
)

// Server struct of server
type Server struct {
	iface  *water.Interface
	ipAddr *netlink.Addr
	ipPool *net.IPNet
	port   int
	ghMTU  int
}

func handleUDPPkg(buf []byte, addr *net.UDPAddr) {
	common.Logger.Info(*addr)
}

// Launch : Launch server
func (server *Server) Launch() error {
	// Launch UDP server, get package from socket and forward to iface
	addr := fmt.Sprintf("%s:%d", strings.Split(server.ipAddr.String(), "/")[0], server.port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		common.Logger.Error(err)
		return err
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		common.Logger.Error(err)
	}

	for {
		readBuf := make([]byte, server.ghMTU)
		len, addr, err := udpConn.ReadFromUDP(readBuf)
		if err != nil {
			common.Logger.Error(err)
		} else {
			go handleUDPPkg(readBuf[:len], addr)
		}
	}
}

// NewServer : Create server and return it
func NewServer(cfg *common.ServerConfig) (*Server, error) {
	common.Logger.Debugf("Create server with config: %v\n", cfg)
	ipnet, err := netlink.ParseIPNet(cfg.IPAddr)
	if err != nil {
		return nil, err
	}
	ipAddr, err := netlink.ParseAddr(cfg.IPAddr)
	if err != nil {
		return nil, err
	}

	server := new(Server)
	server.ipAddr = ipAddr
	server.ipPool = ipnet
	server.port = cfg.Port
	server.ghMTU = cfg.MTU - IPHeaderLength - UDPHeaderLength - GHHeaderLength
	// Create and setup tun interface
	iface, err := setupIface(ipAddr, cfg.MTU)
	if err != nil {
		return nil, err
	}
	server.iface = iface

	return server, nil
}
