package services

import (
	"autohost-cli/internal/plugins/agent/types"
	"fmt"
	"net"
	"os/exec"
)

func GetAgentData() *types.NodeData {

	nodeData := &types.NodeData{}

	if output, err := exec.Command("hostname").Output(); err == nil {
		nodeData.HostName = string(output)
	}
	// if output, err := exec.Command("hostname", "-I").Output(); err == nil {
	// 	nodeData.IPLocal = string(output)
	// }
	if output, err := exec.Command("uname", "-o").Output(); err == nil {
		nodeData.OS = string(output)
	}
	if output, err := exec.Command("uname", "-m").Output(); err == nil {
		nodeData.Arch = string(output)
	}
	// if output, err := exec.Command("hostname", "-I").Output(); err == nil {
	nodeData.IPLocal = getLocalIP()
	fmt.Println("Local IP:", nodeData.IPLocal)
	// }
	nodeData.VersionAgent = "0.1.0"

	return nodeData
}
func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "unknown"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
