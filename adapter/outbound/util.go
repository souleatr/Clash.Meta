package outbound

import (
	"bytes"
	"crypto/tls"
	"github.com/Dreamacro/clash/component/resolver"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/transport/socks5"
	xtls "github.com/xtls/go"
	"net"
	"strconv"
	"sync"
)

var (
	globalClientSessionCache  tls.ClientSessionCache
	globalClientXSessionCache xtls.ClientSessionCache
	once                      sync.Once
)

func tcpKeepAlive(c net.Conn) {
	//if tcp, ok := c.(*net.TCPConn); ok {
	//	_ = tcp.SetKeepAlive(true)
	//	_ = tcp.SetKeepAlivePeriod(30 * time.Second)
	//}
}

func getClientSessionCache() tls.ClientSessionCache {
	once.Do(func() {
		globalClientSessionCache = tls.NewLRUClientSessionCache(128)
	})
	return globalClientSessionCache
}

func getClientXSessionCache() xtls.ClientSessionCache {
	once.Do(func() {
		globalClientXSessionCache = xtls.NewLRUClientSessionCache(128)
	})
	return globalClientXSessionCache
}

func serializesSocksAddr(metadata *C.Metadata) []byte {
	var buf [][]byte
	aType := uint8(metadata.AddrType)
	p, _ := strconv.ParseUint(metadata.DstPort, 10, 16)
	port := []byte{uint8(p >> 8), uint8(p & 0xff)}
	switch metadata.AddrType {
	case socks5.AtypDomainName:
		lenM := uint8(len(metadata.Host))
		host := []byte(metadata.Host)
		buf = [][]byte{{aType, lenM}, host, port}
	case socks5.AtypIPv4:
		host := metadata.DstIP.AsSlice()
		buf = [][]byte{{aType}, host, port}
	case socks5.AtypIPv6:
		host := metadata.DstIP.AsSlice()
		buf = [][]byte{{aType}, host, port}
	}
	return bytes.Join(buf, nil)
}

func resolveUDPAddr(network, address string) (*net.UDPAddr, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	ip, err := resolver.ResolveProxyServerHost(host)
	if err != nil {
		return nil, err
	}
	return net.ResolveUDPAddr(network, net.JoinHostPort(ip.String(), port))
}

func safeConnClose(c net.Conn, err error) {
	if err != nil {
		_ = c.Close()
	}
}
