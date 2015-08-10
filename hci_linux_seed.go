// +build ignore

// +godefs map bdaddr_t Bdaddr /* endian! */

package blugo

// #cgo pkg-config: bluez
// #include <sys/select.h>
// #include <sys/ioctl.h>
// #include <bluetooth/bluetooth.h>
// #include <bluetooth/hci.h>
import "C"

const (
	HCIDEVUP     = C.HCIDEVUP
	HCIDEVDOWN   = C.HCIDEVDOWN
	HCIDEVRESET  = C.HCIDEVRESET
	HCIDEVRESTAT = C.HCIDEVRESTAT

	HCIGETDEVLIST  = C.HCIGETDEVLIST
	HCIGETDEVINFO  = C.HCIGETDEVINFO
	HCIGETCONNLIST = C.HCIGETCONNLIST
	HCIGETCONNINFO = C.HCIGETCONNINFO
	HCIGETAUTHINFO = C.HCIGETAUTHINFO

	HCISETRAW      = C.HCISETRAW
	HCISETSCAN     = C.HCISETSCAN
	HCISETAUTH     = C.HCISETAUTH
	HCISETENCRYPT  = C.HCISETENCRYPT
	HCISETPTYPE    = C.HCISETPTYPE
	HCISETLINKPOL  = C.HCISETLINKPOL
	HCISETLINKMODE = C.HCISETLINKMODE
	HCISETACLMTU   = C.HCISETACLMTU
	HCISETSCOMTU   = C.HCISETSCOMTU

	HCIBLOCKADDR   = C.HCIBLOCKADDR
	HCIUNBLOCKADDR = C.HCIUNBLOCKADDR

	HCIINQUIRY = C.HCIINQUIRY
)

type SockaddrHci C.struct_sockaddr_hci

// OpCode is little endian
type HciFilter C.struct_hci_filter

type HciDevStats C.struct_hci_dev_stats

type HciDevInfo C.struct_hci_dev_info

type HciConnInfo C.struct_hci_conn_info

type HciDevReq C.struct_hci_dev_req

type HciDevListReq C.struct_hci_dev_list_req

type HciConnListReq C.struct_hci_conn_list_req

type HciConnInfoReq C.struct_hci_conn_info_req

type HciAuthInfoReq C.struct_hci_auth_info_req

type HciInquiryReq C.struct_hci_inquiry_req

const (
	SizeofSockaddrHci = C.sizeof_struct_sockaddr_hci
	SizeofHciFilter   = C.sizeof_struct_hci_filter

	SizeofHciDevStats    = C.sizeof_struct_hci_dev_stats
	SizeofHciDevInfo     = C.sizeof_struct_hci_dev_info
	SizeofHciConnInfo    = C.sizeof_struct_hci_conn_info
	SizeofHciDevReq      = C.sizeof_struct_hci_dev_req
	SizeofHciDevListReq  = C.sizeof_struct_hci_dev_list_req
	SizeofHciConnListReq = C.sizeof_struct_hci_conn_list_req
	SizeofHciConnInfoReq = C.sizeof_struct_hci_conn_info_req
	SizeofHciAuthInfoReq = C.sizeof_struct_hci_auth_info_req
	SizeofHciInquiryReq  = C.sizeof_struct_hci_inquiry_req
)
