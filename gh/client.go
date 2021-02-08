package gh

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	"github.com/lucheng0127/groundhog/common"
)

// Client : struct of client
type Client struct {
	sync.Mutex
	port      int
	ghMTU     int
	subnetIP  string
	secretKey string
	conn      *net.UDPConn
}

func doAuth(cfg *common.ClientConfig, client *Client) error {
	// Build groundhog auth package and do request
	ghPkg := new(groundhogPackage)
	ghHeader := new(groundhogHeader)
	ghHeader.flag = GHFlagAuth
	ghHeader.index = GHPKGEND
	ghPkg.payload = []byte(cfg.SecretKey)
	ghHeader.length = uint16(len(ghPkg.payload))
	ghPkg.header = ghHeader
	buf := bytes.NewBuffer(make([]byte, 0, 4+ghHeader.length))
	binary.Write(buf, binary.BigEndian, ghPkg.header)
	buf.Write(ghPkg.payload)
	_, err := client.conn.Write(buf.Bytes())
	if err != nil {
		return common.ClientError(
			fmt.Sprintf("Send data to udp connection error\n%+v\n", err))
	}
	return nil
}

func setupClient(cfg *common.ClientConfig, client *Client) (*net.UDPConn, error) {
	// Connect to udp server
	localAddr, err := net.ResolveUDPAddr("udp", cfg.LocalAddr)
	if err != nil {
		return nil, common.InitailizeError(
			fmt.Sprintf("Failed to resolve local address\n%+v\n", err))
	}
	remoteAddr, err := net.ResolveUDPAddr("udp", cfg.ServerAddr)
	if err != nil {
		return nil, common.InitailizeError(
			fmt.Sprintf("Failed to resolve server address\n%+v\n", err))
	}

	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	if err != nil {
		return nil, common.ClientError(
			fmt.Sprintf("Dial error\n%+V\n", err))
	}
	client.conn = conn

	// Block until finish setup
	// Auth
	err = doAuth(cfg, client)
	if err != nil {
		return nil, common.AuthError(
			fmt.Sprintf("Authorise error\n%+v\n", err))
	}

	// Request subnet ip and setup tun

	// Request and setup routings
	return conn, nil
}

// NewClientAndRun : Create client and launch it
func NewClientAndRun(cfg *common.ClientConfig) error {
	common.Logger.Debugf("Create client witch config: %+v", cfg)
	client := Client{}

	// Setup client, it will connect to udp server
	// and do auth and get subnet ip and routings first
	// then setup tun interface and routing
	udpConn, err := setupClient(cfg, &client)
	if err != nil {
		return err
	}

	fmt.Println(udpConn)
	return nil
}
