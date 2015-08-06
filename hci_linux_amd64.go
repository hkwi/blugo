// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs -- -funsigned-char hci_linux_seed.go

package blugo

const (
	HCIDEVUP     = 0x400448c9
	HCIDEVDOWN   = 0x400448ca
	HCIDEVRESET  = 0x400448cb
	HCIDEVRESTAT = 0x400448cc

	HCIGETDEVLIST  = 0x800448d2
	HCIGETDEVINFO  = 0x800448d3
	HCIGETCONNLIST = 0x800448d4
	HCIGETCONNINFO = 0x800448d5
	HCIGETAUTHINFO = 0x800448d7

	HCISETRAW      = 0x400448dc
	HCISETSCAN     = 0x400448dd
	HCISETAUTH     = 0x400448de
	HCISETENCRYPT  = 0x400448df
	HCISETPTYPE    = 0x400448e0
	HCISETLINKPOL  = 0x400448e1
	HCISETLINKMODE = 0x400448e2
	HCISETACLMTU   = 0x400448e3
	HCISETSCOMTU   = 0x400448e4

	HCIBLOCKADDR   = 0x400448e6
	HCIUNBLOCKADDR = 0x400448e7

	HCIINQUIRY = 0x800448f0
)

type SockaddrHci struct {
	Family  uint16
	Dev     uint16
	Channel uint16
}

type HciFilter struct {
	Type_mask  uint32
	Event_mask [2]uint32
	Opcode     uint16
	Pad_cgo_0  [2]byte
}

type HciDevStats struct {
	Err_rx  uint32
	Err_tx  uint32
	Cmd_tx  uint32
	Evt_rx  uint32
	Acl_tx  uint32
	Acl_rx  uint32
	Sco_tx  uint32
	Sco_rx  uint32
	Byte_rx uint32
	Byte_tx uint32
}

type HciDevInfo struct {
	Dev_id      uint16
	Name        [8]uint8
	Bdaddr      Bdaddr /* endian! */
	Flags       uint32
	Type        uint8
	Features    [8]uint8
	Pad_cgo_0   [3]byte
	Pkt_type    uint32
	Link_policy uint32
	Link_mode   uint32
	Acl_mtu     uint16
	Acl_pkts    uint16
	Sco_mtu     uint16
	Sco_pkts    uint16
	Stat        HciDevStats
}

type HciConnInfo struct {
	Handle uint16
	Bdaddr Bdaddr /* endian! */
	Type   uint8
	Out    uint8
	State  uint16
	Mode   uint32
}

type HciDevReq struct {
	Id        uint16
	Pad_cgo_0 [2]byte
	Opt       uint32
}

type HciDevListReq struct {
	Num       uint16
	Pad_cgo_0 [2]byte
	Req       [0]HciDevReq
}

type HciConnListReq struct {
	Dev_id    uint16
	Conn_num  uint16
	Conn_info [0]HciConnInfo
}

type HciConnInfoReq struct {
	Bdaddr    Bdaddr /* endian! */
	Type      uint8
	Pad_cgo_0 [1]byte
	Info      [0]HciConnInfo
}

type HciAuthInfoReq struct {
	Bdaddr Bdaddr /* endian! */
	Type   uint8
}

type HciInquiryReq struct {
	Dev_id    uint16
	Flags     uint16
	Lap       [3]uint8
	Length    uint8
	Num_rsp   uint8
	Pad_cgo_0 [1]byte
}

type _Socklen uint32

const (
	SizeofSockaddrHci = 0x6
	SizeofHciFilter   = 0x10

	SizeofHciDevStats    = 0x28
	SizeofHciDevInfo     = 0x5c
	SizeofHciConnInfo    = 0x10
	SizeofHciDevReq      = 0x8
	SizeofHciDevListReq  = 0x4
	SizeofHciConnListReq = 0x4
	SizeofHciConnInfoReq = 0x8
	SizeofHciAuthInfoReq = 0x7
	SizeofHciInquiryReq  = 0xa
)
