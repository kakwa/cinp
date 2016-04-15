package proto

import (
    "net"
//    "github.com/kakwa/cinp/proto/v1"
    "golang.org/x/net/ipv4"
)

const (
    V1          Version = 1
)

type Version byte

type Packet []byte

func (p Packet) Version() Version { return Version(p[0]) }

type Handler interface {
    ServeCinp(req Packet) Packet
}

type serveIfConn struct {
    ifIndex int
    conn    *ipv4.PacketConn
    cm      *ipv4.ControlMessage
}

type ServeConn interface {
    ReadFrom(b []byte) (n int, addr net.Addr, err error)
    WriteTo(b []byte, addr net.Addr) (n int, err error)
}

func (s *serveIfConn) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
    n, s.cm, addr, err = s.conn.ReadFrom(b)
    if s.cm != nil && s.cm.IfIndex != s.ifIndex { // Filter all other interfaces
        n = 0 // Packets < 240 are filtered in Serve().
    }
    return
}

func (s *serveIfConn) WriteTo(b []byte, addr net.Addr) (n int, err error) {

    // ipv4 docs state that Src is "specify only", however testing by tfheen
    // shows that Src IS populated.  Therefore, to reuse the control message,
    // we set Src to nil to avoid the error "write udp4: invalid argument"
    s.cm.Src = nil

    return s.conn.WriteTo(b, s.cm, addr)
}

func ServeIf(ifIndex int, conn net.PacketConn, handler Handler) error {
    p := ipv4.NewPacketConn(conn)
    if err := p.SetControlMessage(ipv4.FlagInterface, true); err != nil {
        return err
    }
    return Serve(&serveIfConn{ifIndex: ifIndex, conn: p}, handler)
}

// ListenAndServe listens on the UDP network address addr and then calls
// Serve with handler to handle requests on incoming packets.
// i.e. ListenAndServeIf("eth0",handler, port)
func ListenAndServeIf(interfaceName string, handler Handler, port string) error {
    iface, err := net.InterfaceByName(interfaceName)
    if err != nil {
        return err
    }
    addr, err := iface.Addrs()
    ip, _, err := net.ParseCIDR(addr[0].String())
    l, err := net.ListenPacket("udp", ip.String() + ":" + port)
    if err != nil {
        return err
    }
    defer l.Close()
    return ServeIf(iface.Index, l, handler)
}

func ListenAndClientIf(interfaceName string, handler Handler, port string, timeout int) error {
    iface, err := net.InterfaceByName(interfaceName)
    if err != nil {
        return err
    }
    l, err := net.ListenPacket("udp", port)
    if err != nil {
        return err
    }
    defer l.Close()
    return ServeIf(iface.Index, l, handler)
}

func Serve(conn ServeConn, handler Handler) error {
    buffer := make([]byte, 1500)
    for {
        n, _, err := conn.ReadFrom(buffer)
        if err != nil {
            return err
        }
        req := Packet(buffer[:n])
        switch req.Version() {
            case V1: {
            }
            default: {
            }
        }
//        options := req.ParseOptions()
//        var reqType MessageType
//        if t := options[OptionDHCPMessageType]; len(t) != 1 {
//            continue
//        } else {
//            reqType = MessageType(t[0])
//            if reqType < Discover || reqType > Inform {
//                continue
//            }
//        }
//        if res := handler.ServeDHCP(req, reqType, options); res != nil {
//            // If IP not available, broadcast
//            ipStr, portStr, err := net.SplitHostPort(addr.String())
//            if err != nil {
//                return err
//            }
//
//            if net.ParseIP(ipStr).Equal(net.IPv4zero) || req.Broadcast() {
//                port, _ := strconv.Atoi(portStr)
//                addr = &net.UDPAddr{IP: net.IPv4bcast, Port: port}
//            }
//            if _, e := conn.WriteTo(res, addr); e != nil {
//                return e
//            }
//        }
    }
}
