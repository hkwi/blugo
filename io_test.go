package blugo

import (
	"testing"
)

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
						if ret, err := devio.Request(HCI_Read_RSSI, Parameters{
							U16(con.Handle),
						}); err != nil {
							t.Error(err)
						} else {
							t.Logf("rssi=%v", ret[2])
						}
					}
				}
			}
		}
	}
}
