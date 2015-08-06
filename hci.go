package blugo

import (
	"encoding/binary"
)

// Bluetooth Core speicification, Vol 4, Part A, Section 2
// HCI packet indicator shall be sent immediately before HCI packet.

const (
	HCI_COMMAND_PKT = 0x01
	HCI_ACLDATA_PKT = 0x02
	HCI_SCODATA_PKT = 0x03
	HCI_EVENT_PKT   = 0x04
	HCI_VENDOR_PKT  = 0xff
)

// Bluetooth Core specification, Vol 2, Part E, Section 5.2

const (
	EVT_CMD_COMPLETE  = 0x0E
	EVT_CMD_STATUS    = 0x0F
	EVT_LE_META_EVENT = 0x3E
)

type Parameter interface {
	MarshalBinary() ([]byte, error)
}

type Parameters []Parameter

func (self Parameters) MarshalBinary() ([]byte, error) {
	var ret []byte
	for _, p := range []Parameter(self) {
		if b, err := p.MarshalBinary(); err != nil {
			return nil, err
		} else {
			ret = append(ret, b...)
		}
	}
	return ret, nil
}

type Handle uint16

func (self Handle) MarshalBinary() ([]byte, error) {
	var ret [2]byte
	binary.LittleEndian.PutUint16(ret[:], uint16(self))
	return ret[:], nil
}
