package v1

import (
    "errors"
    "math/rand"
    "encoding/binary"
    "time"
)

type Packet []byte
type Version byte
type Format byte
type Size byte
type OpCode byte
type XId []byte
type Payload []byte

func (p Packet) Version() Version { return Version(p[0])  }
func (p Packet) OpCode()  OpCode  { return OpCode(p[1])   }
func (p Packet) Format()  Format  { return Format(p[2])   }
func (p Packet) Size()    Size    { return Size(p[3])     }
func (p Packet) XId()     XId     { return XId(p[4:8])    }
func (p Packet) Payload() Payload { return Payload(p[8:]) }

func (p Packet) SetVersion(v Version)    error { p[0] = byte(v);      return nil }
func (p Packet) SetOpCode(opcode OpCode) error { p[1] = byte(opcode); return nil }
func (p Packet) SetFormat(format Format) error { p[2] = byte(format); return nil }
func (p Packet) SetSize(size Size)       error { p[3] = byte(size);   return nil }
func (p Packet) SetXId(xid XId)          error { copy(p.XId(), xid);  return nil }
func (p Packet) SetPayload(pl Payload)   error {
    if len(pl) > 254 {
        return errors.New("payload max size exeeded")
    }
    copy(p.Payload(), pl)
    p.SetSize(Size(len(pl)))
    return nil
}

const (
    WrongVersion OpCode = 128
    WrongFormat  OpCode = 129
    MalFormed    OpCode = 130
    Request      OpCode = 0
    Answer       OpCode = 1
)

const (
    Clear       Format = 0
    NaCl        Format = 1
)

const (
    V1          Version = 1
)

func GenXid() XId {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    rand := r.Uint32()
    bs := make([]byte, 4)
    binary.LittleEndian.PutUint32(bs, rand)
    return XId(bs)
}

func NewPacket(format Format, opcode OpCode, xid XId, pl Payload) (Packet, error) {
    p := make(Packet, 270)
    p.SetVersion(V1)
    p.SetOpCode(opcode)
    p.SetFormat(format)
    p.SetXId(xid)
    //rpl := make(Payload, 254)
    //switch format {
    //case Clear:
    //    copy(pl, rpl)
    //case NaCl:
    //    //TODO
    //    copy(pl, rpl)
    //default:
    //    return nil, errors.New("wrong payload format")
    //}
    p.SetPayload(pl)
    return p[0:8 + len(pl)], nil
}

func NewRequest(format Format) (Packet, XId, error) {
    xid := GenXid()
    p, err := NewPacket(format, Request, xid, nil)
    p.SetSize(0)
    return p[0:8], xid, err
}

func NewAnswer(in Packet, pl Payload) (Packet, XId, error) {
    if in.Version() != V1 {
        return nil, nil, errors.New("bad version")
    }

    if in.OpCode() != Request {
        return nil, nil, errors.New("not a request")
    }
    xid := in.XId()
    format := in.Format()
    p, err := NewPacket(format, Answer, xid, pl)
    return p, xid, err
}
