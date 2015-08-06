package blugo

import (
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
