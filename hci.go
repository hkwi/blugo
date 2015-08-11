// bluetooth library in pure golang

package blugo

import (
	"encoding/binary"
	"fmt"
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

type Pkt interface {
	Indicator() uint8
}

type CommandPkt struct {
	OpCode OpCode
	Params []byte
}

func (self CommandPkt) Indicator() uint8 {
	return HCI_COMMAND_PKT
}

type AcldataPkt struct {
	Handle uint16
	PB     uint8
	BC     uint8
	Data   []byte
}

func (self AcldataPkt) Indicator() uint8 {
	return HCI_ACLDATA_PKT
}

type ScodataPkt struct {
	ConnectionHandle uint16
	PacketStatusFlag uint8
	Data             []byte
}

func (self ScodataPkt) Indicator() uint8 {
	return HCI_SCODATA_PKT
}

type EventPkt struct {
	Code   uint8
	Params []byte
}

func (self EventPkt) Indicator() uint8 {
	return HCI_EVENT_PKT
}

// Bluetooth Core specification, Vol 2, Part E, Section 5.2

const (
	EVT_REMOTE_NAME_REQ_COMPLETE = 0x07
	EVT_CMD_COMPLETE             = 0x0E
	EVT_CMD_STATUS               = 0x0F
	EVT_LE_META_EVENT            = 0x3E
)

type EvtRemoteNameReqComplete struct {
	Status uint8
	Bdaddr Bdaddr
	Name   string
}

func (self *EvtRemoteNameReqComplete) UnmarshalBinary(data []byte) error {
	if len(data) < 255 {
		return fmt.Errorf("too short")
	}
	self.Status = data[0]
	copy(self.Bdaddr[:], data[1:])
	self.Name = string(data[7:]) // 248 octets
	return nil
}

type EvtCmdComplete struct {
	Ncmd   uint8
	OpCode uint16
	Params []byte
}

func (self *EvtCmdComplete) UnmarshalBinary(data []byte) error {
	self.Ncmd = data[0]
	self.OpCode = binary.LittleEndian.Uint16(data[1:])
	self.Params = data[3:]
	return nil
}

type EvtLeMetaEvent struct {
	Subevent uint8
	Data     []byte
}

func (self *EvtLeMetaEvent) UnmarshalBinary(data []byte) error {
	self.Subevent = data[0]
	self.Data = data[1:]
	return nil
}

type EvtCmdStatus struct {
	Status uint8
	Ncmd   uint8
	OpCode uint16
}

func (self *EvtCmdStatus) UnmarshalBinary(data []byte) error {
	self.Status = data[0]
	self.Ncmd = data[1]
	self.OpCode = binary.LittleEndian.Uint16(data[2:])
	return nil
}

type EventPktParams interface{}

func (self EventPkt) Parse() (EventPktParams, error) {
	switch self.Code {
	case EVT_REMOTE_NAME_REQ_COMPLETE:
		params := EvtRemoteNameReqComplete{}
		if err := params.UnmarshalBinary(self.Params); err != nil {
			return nil, err
		} else {
			return params, nil
		}
	case EVT_CMD_COMPLETE:
		params := EvtCmdComplete{}
		if err := params.UnmarshalBinary(self.Params); err != nil {
			return nil, err
		} else {
			return params, nil
		}
	case EVT_CMD_STATUS:
		params := EvtCmdStatus{}
		if err := params.UnmarshalBinary(self.Params); err != nil {
			return nil, err
		} else {
			return params, nil
		}
	case EVT_LE_META_EVENT:
		params := EvtRemoteNameReqComplete{}
		if err := params.UnmarshalBinary(self.Params); err != nil {
			return nil, err
		} else {
			return params, nil
		}
	default:
		return nil, fmt.Errorf("unknown EVT_CMD_ %02x", self.Code)
	}
}

// Parse extracts an HCI packet from binary sequence
func Parse(buf []byte) (Pkt, int) {
	if len(buf) == 0 {
		return nil, 0
	}
	switch buf[0] {
	case HCI_COMMAND_PKT:
		if len(buf) < 4 {
			return nil, 0
		}
		params_length := int(buf[3])
		if len(buf) < 4+params_length {
			return nil, 0
		}
		return CommandPkt{
			OpCode: OpCode(binary.LittleEndian.Uint16(buf[1:])),
			Params: buf[4 : 4+params_length],
		}, 4 + params_length

	case HCI_ACLDATA_PKT:
		if len(buf) < 5 {
			return nil, 0
		}
		data_length := int(binary.LittleEndian.Uint16(buf[3:]))
		if len(buf) < 5+data_length {
			return nil, 0
		}
		hdr := binary.LittleEndian.Uint16(buf[1:])
		return AcldataPkt{
			Handle: hdr & 0x0FFF,
			PB:     uint8(hdr>>12) & 0x3,
			BC:     uint8(hdr>>14) & 0x3,
			Data:   buf[5 : 5+data_length],
		}, 5 + data_length

	case HCI_SCODATA_PKT:
		if len(buf) < 3 {
			return nil, 0
		}
		data_length := int(buf[3])
		if len(buf) < 4+data_length {
			return nil, 0
		}
		hdr := binary.LittleEndian.Uint16(buf[1:])
		return ScodataPkt{
			ConnectionHandle: hdr & 0x0FFF,
			PacketStatusFlag: uint8((hdr >> 12) & 0x3),
			Data:             buf[4 : 4+data_length],
		}, 4 + data_length

	case HCI_EVENT_PKT:
		if len(buf) < 3 {
			return nil, 0
		}
		params_length := int(buf[2])
		if len(buf) < 3+params_length {
			return nil, 0
		}
		return EventPkt{
			Code:   buf[1],
			Params: buf[3 : 3+params_length],
		}, 3 + params_length
	default:
		return nil, 0
	}
}

type Parameter interface{}

type Parameters []Parameter

func (self Parameters) MarshalBinary() ([]byte, error) {
	var buf [8]byte
	var ret []byte
	for _, p := range []Parameter(self) {
		switch v := p.(type) {
		case int8:
			ret = append(ret, uint8(v))
		case uint8:
			ret = append(ret, v)
		case int16:
			binary.LittleEndian.PutUint16(buf[:], uint16(v))
			ret = append(ret, buf[:2]...)
		case uint16:
			binary.LittleEndian.PutUint16(buf[:], v)
			ret = append(ret, buf[:2]...)
		case int32:
			binary.LittleEndian.PutUint32(buf[:], uint32(v))
			ret = append(ret, buf[:4]...)
		case uint32:
			binary.LittleEndian.PutUint32(buf[:], v)
			ret = append(ret, buf[:4]...)
		case int64:
			binary.LittleEndian.PutUint64(buf[:], uint64(v))
			ret = append(ret, buf[:8]...)
		case uint64:
			binary.LittleEndian.PutUint64(buf[:], v)
			ret = append(ret, buf[:8]...)
		case []byte:
			ret = append(ret, v...)
		default:
			return nil, fmt.Errorf("unknown type")
		}
	}
	return ret, nil
}

type HciError uint8

const (
	_ HciError = iota
	HCI_UNKNOWN_COMMAND
	HCI_NO_CONNECTION
	HCI_HARDWARE_FAILURE
	HCI_PAGE_TIMEOUT
	HCI_AUTHENTICATION_FAILURE
	HCI_PIN_OR_KEY_MISSING
	HCI_MEMORY_FULL
	HCI_CONNECTION_TIMEOUT
	HCI_MAX_NUMBER_OF_CONNECTIONS
	HCI_MAX_NUMBER_OF_SCO_CONNECTIONS
	HCI_ACL_CONNECTION_EXISTS
	HCI_COMMAND_DISALLOWED
	HCI_REJECTED_LIMITED_RESOURCES
	HCI_REJECTED_SECURITY
	HCI_REJECTED_PERSONAL
	HCI_HOST_TIMEOUT
	HCI_UNSUPPORTED_FEATURE
	HCI_INVALID_PARAMETERS
	HCI_OE_USER_ENDED_CONNECTION
	HCI_OE_LOW_RESOURCES
	HCI_OE_POWER_OFF
	HCI_CONNECTION_TERMINATED
	HCI_REPEATED_ATTEMPTS
	HCI_PAIRING_NOT_ALLOWED
	HCI_UNKNOWN_LMP_PDU
	HCI_UNSUPPORTED_REMOTE_FEATURE
	HCI_SCO_OFFSET_REJECTED
	HCI_SCO_INTERVAL_REJECTED
	HCI_AIR_MODE_REJECTED
	HCI_INVALID_LMP_PARAMETERS
	HCI_UNSPECIFIED_ERROR
	HCI_UNSUPPORTED_LMP_PARAMETER_VALUE
	HCI_ROLE_CHANGE_NOT_ALLOWED
	HCI_LMP_RESPONSE_TIMEOUT
	HCI_LMP_ERROR_TRANSACTION_COLLISION
	HCI_LMP_PDU_NOT_ALLOWED
	HCI_ENCRYPTION_MODE_NOT_ACCEPTED
	HCI_UNIT_LINK_KEY_USED
	HCI_QOS_NOT_SUPPORTED
	HCI_INSTANT_PASSED
	HCI_PAIRING_NOT_SUPPORTED
	HCI_TRANSACTION_COLLISION
	_
	HCI_QOS_UNACCEPTABLE_PARAMETER
	HCI_QOS_REJECTED
	HCI_CLASSIFICATION_NOT_SUPPORTED
	HCI_INSUFFICIENT_SECURITY
	HCI_PARAMETER_OUT_OF_RANGE
	_
	HCI_ROLE_SWITCH_PENDING
	_
	HCI_SLOT_VIOLATION
	HCI_ROLE_SWITCH_FAILED
	HCI_EIR_TOO_LARGE
	HCI_SIMPLE_PAIRING_NOT_SUPPORTED
	HCI_HOST_BUSY_PAIRING
)

func (self HciError) String() string {
	switch self {
	case HCI_UNKNOWN_COMMAND:
		return "HCI_UNKNOWN_COMMAND"
	case HCI_NO_CONNECTION:
		return "HCI_NO_CONNECTION"
	case HCI_HARDWARE_FAILURE:
		return "HCI_HARDWARE_FAILURE"
	case HCI_PAGE_TIMEOUT:
		return "HCI_PAGE_TIMEOUT"
	case HCI_AUTHENTICATION_FAILURE:
		return "HCI_AUTHENTICATION_FAILURE"
	case HCI_PIN_OR_KEY_MISSING:
		return "HCI_PIN_OR_KEY_MISSING"
	case HCI_MEMORY_FULL:
		return "HCI_MEMORY_FULL"
	case HCI_CONNECTION_TIMEOUT:
		return "HCI_CONNECTION_TIMEOUT"
	case HCI_MAX_NUMBER_OF_CONNECTIONS:
		return "HCI_MAX_NUMBER_OF_CONNECTIONS"
	case HCI_MAX_NUMBER_OF_SCO_CONNECTIONS:
		return "HCI_MAX_NUMBER_OF_SCO_CONNECTIONS"
	case HCI_ACL_CONNECTION_EXISTS:
		return "HCI_ACL_CONNECTION_EXISTS"
	case HCI_COMMAND_DISALLOWED:
		return "HCI_COMMAND_DISALLOWED"
	case HCI_REJECTED_LIMITED_RESOURCES:
		return "HCI_REJECTED_LIMITED_RESOURCES"
	case HCI_REJECTED_SECURITY:
		return "HCI_REJECTED_SECURITY"
	case HCI_REJECTED_PERSONAL:
		return "HCI_REJECTED_PERSONAL"
	case HCI_HOST_TIMEOUT:
		return "HCI_HOST_TIMEOUT"
	case HCI_UNSUPPORTED_FEATURE:
		return "HCI_UNSUPPORTED_FEATURE"
	case HCI_INVALID_PARAMETERS:
		return "HCI_INVALID_PARAMETERS"
	case HCI_OE_USER_ENDED_CONNECTION:
		return "HCI_OE_USER_ENDED_CONNECTION"
	case HCI_OE_LOW_RESOURCES:
		return "HCI_OE_LOW_RESOURCES"
	case HCI_OE_POWER_OFF:
		return "HCI_OE_POWER_OFF"
	case HCI_CONNECTION_TERMINATED:
		return "HCI_CONNECTION_TERMINATED"
	case HCI_REPEATED_ATTEMPTS:
		return "HCI_REPEATED_ATTEMPTS"
	case HCI_PAIRING_NOT_ALLOWED:
		return "HCI_PAIRING_NOT_ALLOWED"
	case HCI_UNKNOWN_LMP_PDU:
		return "HCI_UNKNOWN_LMP_PDU"
	case HCI_UNSUPPORTED_REMOTE_FEATURE:
		return "HCI_UNSUPPORTED_REMOTE_FEATURE"
	case HCI_SCO_OFFSET_REJECTED:
		return "HCI_SCO_OFFSET_REJECTED"
	case HCI_SCO_INTERVAL_REJECTED:
		return "HCI_SCO_INTERVAL_REJECTED"
	case HCI_AIR_MODE_REJECTED:
		return "HCI_AIR_MODE_REJECTED"
	case HCI_INVALID_LMP_PARAMETERS:
		return "HCI_INVALID_LMP_PARAMETERS"
	case HCI_UNSPECIFIED_ERROR:
		return "HCI_UNSPECIFIED_ERROR"
	case HCI_UNSUPPORTED_LMP_PARAMETER_VALUE:
		return "HCI_UNSUPPORTED_LMP_PARAMETER_VALUE"
	case HCI_ROLE_CHANGE_NOT_ALLOWED:
		return "HCI_ROLE_CHANGE_NOT_ALLOWED"
	case HCI_LMP_RESPONSE_TIMEOUT:
		return "HCI_LMP_RESPONSE_TIMEOUT"
	case HCI_LMP_ERROR_TRANSACTION_COLLISION:
		return "HCI_LMP_ERROR_TRANSACTION_COLLISION"
	case HCI_LMP_PDU_NOT_ALLOWED:
		return "HCI_LMP_PDU_NOT_ALLOWED"
	case HCI_ENCRYPTION_MODE_NOT_ACCEPTED:
		return "HCI_ENCRYPTION_MODE_NOT_ACCEPTED"
	case HCI_UNIT_LINK_KEY_USED:
		return "HCI_UNIT_LINK_KEY_USED"
	case HCI_QOS_NOT_SUPPORTED:
		return "HCI_QOS_NOT_SUPPORTED"
	case HCI_INSTANT_PASSED:
		return "HCI_INSTANT_PASSED"
	case HCI_PAIRING_NOT_SUPPORTED:
		return "HCI_PAIRING_NOT_SUPPORTED"
	case HCI_TRANSACTION_COLLISION:
		return "HCI_TRANSACTION_COLLISION"
	case HCI_QOS_UNACCEPTABLE_PARAMETER:
		return "HCI_QOS_UNACCEPTABLE_PARAMETER"
	case HCI_QOS_REJECTED:
		return "HCI_QOS_REJECTED"
	case HCI_CLASSIFICATION_NOT_SUPPORTED:
		return "HCI_CLASSIFICATION_NOT_SUPPORTED"
	case HCI_INSUFFICIENT_SECURITY:
		return "HCI_INSUFFICIENT_SECURITY"
	case HCI_PARAMETER_OUT_OF_RANGE:
		return "HCI_PARAMETER_OUT_OF_RANGE"
	case HCI_ROLE_SWITCH_PENDING:
		return "HCI_ROLE_SWITCH_PENDING"
	case HCI_SLOT_VIOLATION:
		return "HCI_SLOT_VIOLATION"
	case HCI_ROLE_SWITCH_FAILED:
		return "HCI_ROLE_SWITCH_FAILED"
	case HCI_EIR_TOO_LARGE:
		return "HCI_EIR_TOO_LARGE"
	case HCI_SIMPLE_PAIRING_NOT_SUPPORTED:
		return "HCI_SIMPLE_PAIRING_NOT_SUPPORTED"
	case HCI_HOST_BUSY_PAIRING:
		return "HCI_HOST_BUSY_PAIRING"
	default:
		return "HCI_?"
	}
}

func (self HciError) Error() string {
	return self.String()
}
