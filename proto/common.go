package proto

type Version byte

type Packet []byte

func (p Packet) Version() Version { return Version(p[0]) }
