package blugo

import (
	"encoding/binary"
	"fmt"
	"log"
	"syscall"
	"testing"
	"unsafe"
)

type HciDev int // fd

func NewHciDev(devId uint16) (HciDev, error) {
	fd, err := syscall.Socket(syscall.AF_BLUETOOTH, syscall.SOCK_RAW|syscall.SOCK_CLOEXEC, BTPROTO_HCI)
	if err != nil {
		return -1, err
	}
	if _, _, errno := syscall.Syscall(syscall.SYS_BIND,
		uintptr(fd),
		uintptr(unsafe.Pointer(&SockaddrHci{
			Family: syscall.AF_BLUETOOTH,
			Dev:    devId,
		})),
		uintptr(SizeofSockaddrHci),
	); errno != 0 {
		syscall.Close(fd)
		return -1, errno
	}
	return HciDev(fd), nil
}

func (self HciDev) Write(data []byte) (int, error) {
	return syscall.Write(int(self), data)
}

func (self HciDev) Read(data []byte) (int, error) {
	return syscall.Read(int(self), data)
}

func (self *HciDev) Close() {
	syscall.Close(int(*self))
	*self = -1
}

func (self HciDev) getConnInfo(addr Bdaddr) (*HciConnInfo, error) {
	buf := make([]byte, SizeofHciConnInfoReq+SizeofHciConnInfo)
	hdr := (*HciConnInfoReq)(unsafe.Pointer(&buf[0]))
	hdr.Bdaddr = addr
	hdr.Type = uint8(ACL_LINK)

	if err := ioctl(int(self),
		uintptr(HCIGETCONNINFO),
		uintptr(unsafe.Pointer(&buf[0])),
	); err != nil {
		return nil, err
	}

	info := *(*HciConnInfo)(unsafe.Pointer(&buf[SizeofHciConnInfoReq]))
	return &info, nil
}

func GetsockoptHciFilter(fd, level, opt int) (*HciFilter, error) {
	filter := &HciFilter{}
	sz := SizeofHciFilter
	_, _, errno := syscall.Syscall6(syscall.SYS_GETSOCKOPT,
		uintptr(fd),
		uintptr(level),
		uintptr(opt),
		uintptr(unsafe.Pointer(filter)),
		uintptr(unsafe.Pointer(&sz)),
		0)
	if errno != 0 {
		return nil, errno
	}
	return filter, nil
}

func SetsockoptHciFilter(fd, level, opt int, filter *HciFilter) error {
	return syscall.SetsockoptString(
		fd,
		level,
		opt,
		string((*[SizeofHciFilter]byte)(unsafe.Pointer(filter))[:]),
	)
}

func (self HciDev) Request(ogf, ocf uint16, params Parameters) (Parameters, error) {
	req := make([]byte, 4)
	req[0] = HCI_COMMAND_PKT
	binary.LittleEndian.PutUint16(req[1:], uint16((ogf<<10)|ocf))

	if pbuf, err := params.MarshalBinary(); err != nil {
		return nil, err
	} else {
		req[3] = uint8(len(pbuf))
		req = append(req, pbuf...)
	}

	var save *HciFilter
	if filter, err := GetsockoptHciFilter(int(self), SOL_HCI, HCI_FILTER); err != nil {
		return nil, err
	} else {
		save = filter
	}
	_ = save

	if err := SetsockoptHciFilter(
		int(self),
		SOL_HCI,
		HCI_FILTER,
		&HciFilter{
			Type_mask:  1 << HCI_EVENT_PKT,
			Event_mask: FilterEventMask(EVT_CMD_STATUS, EVT_CMD_COMPLETE, EVT_LE_META_EVENT),
			Opcode:     uint16(MakeOpCode(ogf, ocf)),
		},
	); err != nil {
		return nil, err
	}

	if n, err := self.Write(req); err != nil {
		return nil, err
	} else if n != len(req) {
		return nil, fmt.Errorf("write incomplete")
	}
	for {
		// xxx: add timeout
		buf := make([]byte, 1500)
		if n, err := self.Read(buf); err != nil {
			return nil, err
		} else {
			// xxx: parse the response
			log.Print(buf[:n])
		}
	}

	return nil, nil
}

func TestHoge(t *testing.T) {
	hci, err := NewHci()
	if err != nil {
		t.Error(err)
		return
	}
	if devs, err := hci.GetDevList(); err != nil {
		t.Error(err)
		return
	} else {
		for _, dev := range devs {
			if devio, err := NewHciDev(dev.Id); err != nil {
				t.Error(err)
			} else {
				if conns, err := hci.GetConnList(dev.Id); err != nil {
					t.Error(err)
				} else {
					for _, con := range conns {
						t.Logf("connection %v", con)
						devio.Request(0x05, 0x0005, Parameters{
							Handle(con.Handle),
						})
					}
				}
			}
		}
	}
}
