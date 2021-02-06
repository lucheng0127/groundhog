package gh

import (
	"fmt"

	"github.com/songgao/water"
	"github.com/vishvananda/netlink"

	"github.com/lucheng0127/groundhog/common"
)

func doSetupIface(ifaceName string, ip *netlink.Addr, mtu int) error {
	iface, err := netlink.LinkByName(ifaceName)
	if err != nil {
		return common.IfaceSetupError(fmt.Sprintf("Get link failed: %+v\n", err))
	}

	err = netlink.AddrAdd(iface, ip)
	if err != nil {
		return common.IfaceSetupError(fmt.Sprintf("Tun add ip failed: %+v\n", err))
	}

	err = netlink.LinkSetMTU(iface, mtu)
	if err != nil {
		return common.IfaceSetupError(fmt.Sprintf("Tun set mtu failed: %+v\n", err))
	}

	err = netlink.LinkSetUp(iface)
	if err != nil {
		return common.IfaceSetupError(fmt.Sprintf("Set tun up failed: %+v\n", err))
	}

	return nil
}

func setupIface(ip *netlink.Addr, mtu int) (*water.Interface, error) {
	config := water.Config{DeviceType: water.TUN}
	iface, err := water.New(config)
	if err != nil {
		return nil, common.IfaceSetupError(fmt.Sprintf("Create tun failed: %+v\n", err))
	}

	err = doSetupIface(iface.Name(), ip, mtu)
	if err != nil {
		return nil, err
	}
	common.Logger.Infof("Setup %s with ip %s mtu %d\n", iface.Name(), ip.String(), mtu)

	return iface, nil
}
