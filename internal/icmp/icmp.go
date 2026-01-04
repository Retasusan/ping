package icmp

import "encoding/binary"

type ICMPEcho struct {
	Type       uint8
	Code       uint8
	Checksum   uint16
	Identifier uint16
	Sequence   uint16
	Data       []byte
}

func (e *ICMPEcho) Marshal() []byte {
	b := make([]byte, 8+len(e.Data))

	b[0] = e.Type
	b[1] = e.Code
	binary.BigEndian.PutUint16(b[2:4], e.Checksum)
	binary.BigEndian.PutUint16(b[4:6], e.Identifier)
	binary.BigEndian.PutUint16(b[6:8], e.Sequence)
	copy(b[8:], e.Data)

	return b
}

func (e *ICMPEcho) MarshalWithChecksum() []byte {
	b := e.Marshal()
	binary.BigEndian.PutUint16(b[2:4], 0)
	csum := icmpChecksum(b)
	binary.BigEndian.PutUint16(b[2:4], csum)
	return b
}
