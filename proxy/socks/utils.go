package socks

import (
	"errors"
	"net"

	"github.com/brobird/clash/common/pool"
	"github.com/brobird/clash/component/socks5"
)

type fakeConn struct {
	net.PacketConn
	rAddr   net.Addr
	payload []byte
	bufRef  []byte
}

func (c *fakeConn) Data() []byte {
	return c.payload
}

// WriteBack wirtes UDP packet with source(ip, port) = `addr`
func (c *fakeConn) WriteBack(b []byte, addr net.Addr) (n int, err error) {
	if addr == nil {
		err = errors.New("Invalid udp packet")
		return
	}

	udpaddr, ok := addr.(*net.UDPAddr)
	if !ok || udpaddr == nil {
		err = errors.New("Invalid udp packet")
		return
	}

	packet, err := socks5.EncodeUDPPacket(socks5.ParseAddrToSocksAddr(addr), b)
	if err != nil {
		return
	}
	return c.PacketConn.WriteTo(packet, c.rAddr)
}

// LocalAddr returns the source IP/Port of UDP Packet
func (c *fakeConn) LocalAddr() net.Addr {
	return c.rAddr
}

func (c *fakeConn) Close() error {
	err := c.PacketConn.Close()
	pool.BufPool.Put(c.bufRef[:cap(c.bufRef)])
	return err
}
