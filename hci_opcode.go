package blugo

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

type OpCode uint16

const (
	_ = iota
	OGF_LINK_CTL
	OGF_LINK_POLICY
	OGF_HOST_CTL
	OGF_INFO_PARAM
	OGF_STATUS_PARAM
	OGF_TEST
	_
	OGF_LE_CTL
)

const (
	_ = iota | (OGF_STATUS_PARAM << 10)
	HCI_Read_Failed_Contact_Counter
	HCI_Reset_Failed_Contact_Counter
	HCI_Read_Link_Quality
	_
	HCI_Read_RSSI
	HCI_Read_AFH_Channel_Map
	HCI_Read_Clock
	HCI_Read_Encryption_Key_Size
	HCI_Read_Local_AMP_Info
	HCI_Read_Local_AMP_ASSOC
	HCI_Write_Remote_AMP_ASSOC
	HCI_Get_MWS_Transport_Layer_Configuration
	HCI_Set_Triggered_Clock_Capture
)

func (self OpCode) Ogf() uint8 {
	return uint8(self >> 10)
}

func (self OpCode) Ocf() uint16 {
	return uint16(self) & 0x03ff
}

func (self OpCode) String() string {
	return fmt.Sprintf("ogf=0x%02x,ocf=0x%04x", self.Ogf(), self.Ocf())
}

func (self OpCode) Native() uint16 {
	var ret [2]byte
	binary.LittleEndian.PutUint16(ret[:], uint16(self))
	return *(*uint16)(unsafe.Pointer(&ret))
}

func MakeOpCode(ogf uint8, ocf uint16) OpCode {
	return OpCode(uint16(ogf)<<10 | (ocf & 0x03ff))
}

func (self OpCode) Response(data []byte) (Parameters, error) {
	switch self {
	case HCI_Read_RSSI:
		if len(data) < 4 {
			return nil, fmt.Errorf("too short")
		}
		return Parameters{
			data[0],
			binary.LittleEndian.Uint16(data[1:]),
			int8(data[3]),
		}, nil
	// XXX: add more opcodes
	default:
		return nil, fmt.Errorf("unknown opcode")
	}
}
