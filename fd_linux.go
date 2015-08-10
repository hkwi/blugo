// +build linux

package blugo

import (
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"
)

type Hci int // fd

func NewHci() (Hci, error) {
	fd, err := syscall.Socket(syscall.AF_BLUETOOTH, syscall.SOCK_RAW|syscall.SOCK_CLOEXEC, BTPROTO_HCI)
	return Hci(fd), err
}

func (self *Hci) Close() {
	syscall.Close(int(*self))
	*self = -1
}

const HCI_MAX_DEV = 32

func (self Hci) GetDevList() ([]HciDevReq, error) {
	buf := make([]byte, SizeofHciDevListReq+HCI_MAX_DEV*SizeofHciDevReq)
	hdr := (*HciDevListReq)(unsafe.Pointer(&buf[0]))
	hdr.Num = uint16(HCI_MAX_DEV)
	if err := ioctl(int(self),
		uintptr(HCIGETDEVLIST),
		uintptr(unsafe.Pointer(&buf[0])),
	); err != nil {
		return nil, err
	}
	var ret []HciDevReq
	for i := 0; i < int(hdr.Num); i++ {
		ret = append(ret, *(*HciDevReq)(unsafe.Pointer(&buf[SizeofHciDevListReq+i*SizeofHciDevReq])))
	}
	return ret, nil
}

func (self Hci) GetDevInfo(devId uint16) (HciDevInfo, error) {
	ret := HciDevInfo{
		Dev_id: devId,
	}
	return ret, ioctl(
		int(self),
		uintptr(HCIGETDEVINFO),
		uintptr(unsafe.Pointer(&ret)),
	)
}

const HCI_MAX_CON = 7 // by specification

func (self Hci) GetConnList(devId uint16) ([]HciConnInfo, error) {
	buf := make([]byte, SizeofHciDevListReq+HCI_MAX_DEV*SizeofHciDevReq)
	hdr := (*HciConnListReq)(unsafe.Pointer(&buf[0]))
	hdr.Dev_id = devId
	hdr.Conn_num = uint16(HCI_MAX_DEV)
	if err := ioctl(int(self),
		uintptr(HCIGETCONNLIST),
		uintptr(unsafe.Pointer(&buf[0])),
	); err != nil {
		return nil, err
	}
	var ret []HciConnInfo
	for i := 0; i < int(hdr.Conn_num); i++ {
		ret = append(ret, *(*HciConnInfo)(unsafe.Pointer(&buf[SizeofHciConnListReq+i*SizeofHciConnInfo])))
	}
	return ret, nil
}

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

func (self HciDev) Request(opcode OpCode, params Parameters) (Parameters, error) {
	req := make([]byte, 4)
	req[0] = HCI_COMMAND_PKT
	binary.LittleEndian.PutUint16(req[1:], uint16(opcode))

	if pbuf, err := params.MarshalBinary(); err != nil {
		return nil, err
	} else {
		req[3] = uint8(len(pbuf))
		req = append(req, pbuf...)
	}

	if filter, err := GetsockoptHciFilter(int(self), SOL_HCI, HCI_FILTER); err != nil {
		return nil, err
	} else {
		defer SetsockoptHciFilter(int(self), SOL_HCI, HCI_FILTER, filter)
	}

	filter := &HciFilter{
		Type_mask:  1 << HCI_EVENT_PKT,
		Event_mask: FilterEventMask(EVT_CMD_STATUS, EVT_CMD_COMPLETE),
		Opcode:     opcode.Native(),
	}
	if opcode.Ogf() == OGF_LE_CTL {
		filter.Event_mask = FilterEventMask(EVT_CMD_STATUS, EVT_CMD_COMPLETE, EVT_LE_META_EVENT)
	}
	if err := SetsockoptHciFilter(int(self), SOL_HCI, HCI_FILTER, filter); err != nil {
		return nil, err
	}

	if n, err := self.Write(req); err != nil {
		return nil, err
	} else if n != len(req) {
		return nil, fmt.Errorf("write incomplete")
	}

	if efd, err := syscall.EpollCreate1(syscall.EPOLL_CLOEXEC); err != nil {
		return nil, err
	} else {
		defer syscall.Close(efd)

		var events [1]syscall.EpollEvent
		syscall.EpollCtl(efd, syscall.EPOLL_CTL_ADD, int(self), &syscall.EpollEvent{
			Events: syscall.EPOLLIN,
			Fd:     int32(self),
		})

		capture := make([]byte, 0, 258)
		buf := make([]byte, 258)
		for {
			if n, err := syscall.EpollWait(efd, events[:], 100); err != nil {
				return nil, err
			} else if n == 0 {
				continue
			}
			// todo: operation timeout

			if n, _, _, _, err := syscall.Recvmsg(int(self), buf, nil, syscall.MSG_DONTWAIT); err != nil {
				if errno, ok := err.(syscall.Errno); ok && errno.Temporary() {
					continue
				}
				return nil, err
			} else if n == 0 {
				continue
			} else {
				capture = append(capture, buf[:n]...)
			}
			if pkt, step := Parse(capture); step == 0 {
				continue
			} else if p, err := pkt.(EventPkt).Parse(); err != nil {
				continue
			} else {
				switch ev := p.(type) {
				case EvtCmdComplete:
					if ev.OpCode == uint16(opcode) {
						return opcode.Response(ev.Params)
					}
				case EvtCmdStatus:
					if ev.OpCode == uint16(opcode) && ev.Status != 0 {
						return nil, HciError(ev.Status)
					}
				case EvtLeMetaEvent:
					// todo: OGF_LE_CTL opcode expects this
				}
				capture = capture[step:]
			}
		}
		return nil, fmt.Errorf("should not reach")
	}
}
