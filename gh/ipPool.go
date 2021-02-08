package gh

import (
	"encoding/binary"
	"net"
	"sync"

	"github.com/lucheng0127/groundhog/common"
)

type ipPool struct {
	sync.Mutex
	subnet *net.IPNet
	pool   []net.IP
}

func getIPPool(ipNet *net.IPNet) *ipPool {
	// Convert IPNet mask and address to unit32
	ipAddr := net.ParseIP(ipNet.IP.String())
	ipAddr = ipAddr.To4()
	mask := binary.BigEndian.Uint32(ipNet.Mask)
	start := binary.BigEndian.Uint32(ipAddr)
	// Get the last address
	end := (start & mask) | (mask ^ 0xffffffff)

	pool := new(ipPool)
	pool.subnet = ipNet
	for i := start + 1; i <= end; i++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		pool.pool = append(pool.pool, ip)
	}
	return pool
}

func (pool *ipPool) pop() (net.IP, error) {
	s := pool.pool
	if len(pool.pool) < 1 {
		return nil, common.IPPoolEmptyError("No enough ip")
	}
	ip, s := s[0], s[1:]
	pool.Lock()
	pool.pool = s
	pool.Unlock()
	return ip, nil
}

func (pool *ipPool) release(ip net.IP) {
	pool.Lock()
	pool.pool = append(pool.pool, ip)
	pool.Unlock()
}
