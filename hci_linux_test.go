package blugo

import (
	"testing"
)

func TestDevList(t *testing.T) {
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
			if di, err := hci.GetDevInfo(dev.Id); err != nil {
				t.Error(err)
				return
			} else {
				t.Logf("%v Type=%v Bus=%v",
					string(di.Name[:]),
					HciType((di.Type&0x30)>>4),
					HciBus(di.Type&0x0f))
			}
		}
	}
}

func TestConnList(t *testing.T) {
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
			if conns, err := hci.GetConnList(dev.Id); err != nil {
				t.Error(err)
				return
			} else {
				for _, con := range conns {
					t.Logf("%v %v %v handle=%v state=%v %v",
						[]string{">", "<"}[con.Out],
						LinkType(con.Type),
						con.Bdaddr,
						con.Handle,
						con.State,
						LinkMode(con.Mode),
					)
				}
			}
		}
	}
}
