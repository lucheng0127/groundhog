package gh

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
	"golang.org/x/net/ipv4"

	"github.com/lucheng0127/groundhog/common"
)

// Server struct of server
type Server struct {
	sync.Mutex
	iface         *water.Interface
	ipAddr        *netlink.Addr
	ipPool        *ipPool
	port          int
	ghMTU         int
	clientConnMap map[string]Connection
}

func handleUDPPkg(buf []byte, addr *net.UDPAddr) {
	common.Logger.Info(*addr)
}

func forwardIfacePackageToSocket(conn *net.UDPConn, server *Server) {
	buf := make([]byte, server.ghMTU)
	for {
		len, err := server.iface.Read(buf)
		if err != nil {
			err = common.IfaceSetupError(
				fmt.Sprintf("Read data from interface error\n%+v\n", err))
			common.Logger.Error(err)
		} else {
			header, _ := ipv4.ParseHeader(buf[:len])
			if clientConn, keyInMap := server.clientConnMap[header.Dst.String()]; keyInMap {
				_, err = conn.WriteToUDP(buf[:len], clientConn.remoteAddr)
				if err != nil {
					err = common.UDPConnectionError(
						fmt.Sprintf("Write data to udp connection error\n%+v\n", err))
					common.Logger.Error(err)
				}
			}
		}
	}
}

// Launch : Launch server
func (server *Server) Launch() error {
	// Launch UDP server, get package from socket and forward to iface
	addr := fmt.Sprintf(
		"%s:%d", strings.Split(server.ipAddr.String(), "/")[0], server.port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return common.LaunchServerError(
			fmt.Sprintf(
				"Failed to resolve udp addr\n%+v\n", err))
	}
	ServerConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return common.LaunchServerError(
			fmt.Sprintf(
				"Failed to listen udp\n%+v\n", err))
	}

	// Forward packages from tun interface to socket
	go forwardIfacePackageToSocket(ServerConn, server)

	for {
		readBuf := make([]byte, server.ghMTU)
		len, addr, err := ServerConn.ReadFromUDP(readBuf)
		if err != nil {
			err = common.UDPConnectionError(
				fmt.Sprintf("Read data from udp connection error\n%+v\n", err))
			common.Logger.Errorf("%+v", err)
			continue
		}
		connSubnetIP, err := server.ipPool.pop()
		if err != nil {
			common.Logger.Error(err)
		} else {
			clientConn := Connection{subnetIP: connSubnetIP, remoteAddr: addr, status: "ACTIVE"}
			common.Logger.Debugf("Assign %s for connection: %+v", connSubnetIP.String(), clientConn)
			server.Lock()
			server.clientConnMap[connSubnetIP.String()] = clientConn
			server.Unlock()
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
	server.ipPool = getIPPool(ipnet)
	server.port = cfg.Port
	server.ghMTU = cfg.MTU - IPHeaderLength - UDPHeaderLength - GHHeaderLength
	server.clientConnMap = make(map[string]Connection)
	// Create and setup tun interface
	iface, err := setupIface(ipAddr, cfg.MTU)
	if err != nil {
		return nil, err
	}
	server.iface = iface

	return server, nil
}
