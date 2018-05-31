package mystatus

import (
	"strings"

	"github.com/shirou/gopsutil/net"
)

type ipAddressBlock struct {
	interfaceName string
}

func (cb *ipAddressBlock) Render() barBlockData {
	var text string
	interfaces, err := net.Interfaces()
	if err != nil {
		text = "IP error: " + err.Error()
	} else {
		for _, iface := range interfaces {
			if iface.Name == cb.interfaceName {
				if len(iface.Addrs) == 0 {
					text = "D"
					break
				}
				text = strings.Split(iface.Addrs[0].Addr, "/")[0]
				text = strings.Replace(text, "192.168", `<span size="6000">192.168</span>`, -1)
				break
			}
		}
	}
	return barBlockData{
		FullText: text,
		Markup:   "pango",
	}
}
