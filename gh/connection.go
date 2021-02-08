package gh

import "net"

// Connection : Connection of client
type Connection struct {
	subnetIP   net.IP       // Subnet ip get from ip pool
	remoteAddr *net.UDPAddr // Remote ip maybe a public ip and port
	status     string       // Connection status
}
