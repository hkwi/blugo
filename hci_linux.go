// +build linux

//go:generate sh gen.sh hci linux_seed $GOOS $GOARCH

package blugo

import (
	"net"
	"strings"
)

type Bdaddr [6]byte

func (self Bdaddr) String() string {
	return net.HardwareAddr([]byte{
		self[5],
		self[4],
		self[3],
		self[2],
		self[1],
		self[0],
	}).String()
}

func ParseMAC(s string) (Bdaddr, error) {
	var self Bdaddr
	addr, err := net.ParseMAC(s)
	if err == nil {
		copy(self[:], addr)
		self[0], self[1], self[2], self[3], self[4], self[5] = self[5], self[4], self[3], self[2], self[1], self[0]
	}
	return self, err
}

type HciBus uint8

const (
	HCI_VIRTUAL HciBus = iota
	HCI_USB
	HCI_PCCARD
	HCI_UART
	HCI_RS232
	HCI_PCI
	HCI_SDIO
)

func (self HciBus) String() string {
	switch self {
	case HCI_VIRTUAL:
		return "VIRTUAL"
	case HCI_USB:
		return "USB"
	case HCI_PCCARD:
		return "PCCARD"
	case HCI_UART:
		return "UART"
	case HCI_RS232:
		return "RS232"
	case HCI_PCI:
		return "PCI"
	case HCI_SDIO:
		return "SDIO"
	default:
		return "UNKNOWN"
	}
}

type HciType uint8

const (
	HCI_BREDR HciType = iota
	HCI_AMP
)

func (self HciType) String() string {
	switch self {
	case HCI_BREDR:
		return "BR/EDR"
	case HCI_AMP:
		return "AMP"
	default:
		return "UNKNOWN"
	}
}

type LinkMode uint32

const (
	HCI_LM_MASTER = 1 << iota
	HCI_LM_AUTH
	HCI_LM_ENCRYPT
	HCI_LM_TRUSTED
	HCI_LM_RELIABLE
	HCI_LM_SECURE
	HCI_LM_ACCEPT = 0x8000
)

func (self LinkMode) String() string {
	var comps []string
	if self&HCI_LM_ACCEPT != 0 {
		comps = append(comps, "ACCEPT")
	}
	if self&HCI_LM_MASTER != 0 {
		comps = append(comps, "MASTER")
	} else {
		comps = append(comps, "SLAVE")
	}
	if self&HCI_LM_AUTH != 0 {
		comps = append(comps, "AUTH")
	}
	if self&HCI_LM_ENCRYPT != 0 {
		comps = append(comps, "ENCRYPT")
	}
	if self&HCI_LM_TRUSTED != 0 {
		comps = append(comps, "TRUSTED")
	}
	if self&HCI_LM_RELIABLE != 0 {
		comps = append(comps, "RELIABLE")
	}
	if self&HCI_LM_SECURE != 0 {
		comps = append(comps, "SECURE")
	}
	if len(comps) == 0 {
		return "NONE"
	} else {
		return strings.Join(comps, " ")
	}
}

type LinkType uint8 // Baseband link

const (
	SCO_LINK  LinkType = 0x00
	ACL_LINK  LinkType = 0x01
	ESCO_LINK LinkType = 0x02
	LE_LINK   LinkType = 0x80
)

func (self LinkType) String() string {
	switch self {
	case SCO_LINK:
		return "SCO"
	case ACL_LINK:
		return "ACL"
	case ESCO_LINK:
		return "eSCO"
	case LE_LINK:
		return "LE"
	default:
		return "Unknown"
	}
}

// socket option
const (
	_ = iota
	HCI_DATA_DIR
	HCI_FILTER
	HCI_TIME_STAMP
)

func FilterEventMask(index ...int) [2]uint32 {
	var ret [2]uint32
	for _, i := range index {
		ret[i/32] |= uint32(1 << uint8(i%32))
	}
	return ret
}
