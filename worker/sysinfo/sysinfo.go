package sysinfo

import (
	"bytes"
	"net"
	"os/exec"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/tzY15368/lazarus/config"
)

func GetCPUPercent() int32 {
	p, err := cpu.Percent(50*time.Millisecond, false)
	if err != nil {
		return -1
	}
	return int32(p[0])
}

func GetMemPercent() int32 {
	m, err := mem.VirtualMemory()
	if err != nil {
		return -1
	}
	return int32(m.UsedPercent)
}

func GetActiveConn() int32 {
	cmd := exec.Command("bash", "-c", "netstat -nat|grep -i “80”|wc -l")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return -1
	}
	s := string(out)
	result, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return int32(result)
}

func GetCurrentData() int32 {
	return -1
}

func GetTotalData() int32 {
	return config.Cfg.Worker.TotalDataMB
}

func GetMacAddr() (addr string) {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				// Don't use random as we have a real address
				addr = i.HardwareAddr.String()
				break
			}
		}
	}
	return
}
