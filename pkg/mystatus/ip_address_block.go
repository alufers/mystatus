package mystatus

import (
	"strings"

	"github.com/shirou/gopsutil/net"
)

type ipAddressBlock struct {
	interfaceName  string
	interfaceNames []string
}

func (cb *ipAddressBlock) hadInterface(name string) bool {
	for _, i := range cb.interfaceNames {
		if i == name {
			return true
		}
	}
	return false
}

func (cb *ipAddressBlock) Render() barBlockData {
	fullText := ""
	interfaces, err := net.Interfaces()
	if err != nil {
		fullText = "IP error: " + err.Error()
	} else {
		for _, iface := range interfaces {
			text := ""
			if cb.hadInterface((iface.Name)) {
				if len(iface.Addrs) == 0 {
					text = "Disconnected"
				} else {
					text = strings.Split(iface.Addrs[0].Addr, "/")[0]
					text = strings.Replace(text, "192.168", `<span size="6000">192.168</span>`, -1)
				}

				fullText += text + " "
			}
		}
	}

	return barBlockData{
		FullText: fullText,
		Markup:   "pango",
	}
}
