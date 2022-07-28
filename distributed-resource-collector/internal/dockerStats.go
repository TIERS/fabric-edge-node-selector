package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/docker"
)

// -- DOCKER INTERNAL STRUCTS
type DockerPort struct {
	IP          string `json:"IP"`
	PrivatePort int16  `json:"PrivatePort"`
	PublicPort  int16  `json:"PublicPort"`
	Type        string `json:"Type"`
}

type DockerNetworkConfig struct {
	IPAMConfig          []string `json:"IPAMConfig"`
	Links               []string `json:"Links"`
	Aliases             []string `json:"Aliases"`
	NetworkID           string   `json:"NetworkID"`
	EndpointID          string   `json:"EndpointID"`
	Gateway             string   `json:"Gateway"`
	IPAddress           string   `json:"IPAddress"`
	IPPrefixLen         int16    `json:"IPPrefixLen"`
	IPv6Gateway         string   `json:"IPv6Gateway"`
	GlobalIPv6Address   string   `json:"GlobalIPv6Address"`
	GlobalIPv6PrefixLen string   `json:"GlobalIPv6PrefixLen"`
	MacAddress          string   `json:"MacAddress"`
	DriverOpts          []string `json:"DriverOpts"`
}

type DockerNetwork struct {
	Networks map[string]DockerNetworkConfig `json:"Networks"`
}

type DockerMount struct {
	Type        string `json:"Type"`
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
	Mode        string `json:"Mode"`
	RW          bool   `json:"RW"`
	Propagation string `json:"Propagation"`
}

type DrcDockerSocketStats struct {
	ID              string            `json:"ID"`
	Names           []string          `json:"Names"`
	Image           string            `json:"Image"`
	ImageID         string            `json:"ImageID"`
	Command         string            `json:"Command"`
	Created         int64             `json:"Created"`
	Ports           []DockerPort      `json:"Ports"`
	Labels          map[string]string `json:"Labels"`
	State           string            `json:"State"`
	Status          string            `json:"Status"`
	HostConfig      map[string]string `json:"HostConfig"`
	NetworkSettings DockerNetwork     `json:"NetworkSettings"`
	Mounts          []DockerMount     `json:"Mounts"`
}

func DrcSocketJsonToStruct(v string) (socketStats []DrcDockerSocketStats) {
	json.Unmarshal([]byte(v), &socketStats)
	return socketStats
}

func (d DrcDockerSocketStats) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

func GetDockerSocketStats() (result string) {
	conn, _ := net.Dial("unix", "/var/run/docker.sock")
	fmt.Fprintf(conn, "GET /containers/json HTTP/1.0\r\n\r\n")
	status, _ := bufio.NewReader(conn).ReadString('\t')
	//fmt.Println("Message from server: " + status)
	splitStatus := strings.Split(status, "\n")
	result = splitStatus[len(splitStatus)-2]
	return result
}

// There are basically two ways to get Docker stats. If this program is running inside a container, we must connect to the Docker Unix Socket.
// The Unix Socket has much more information than psUtil. At this point in time, I've configured so that we get the same information from both implementations.
// In the future it might be worth it to remove the usage of psUtil's docker function, but it's not guaranteed that all palces where this application
// is running will be using docker at all.

func GetDockerStats() (dockerStats []DrcDockerStats) {
	if InDockerContainer() {
		tmpStats := DrcSocketJsonToStruct(GetDockerSocketStats())
		for _, container := range tmpStats {
			//fmt.Println(container.String())
			tmp := DrcDockerStats{
				ContainerID: container.ID,
				Name:        container.Names[0],
				Image:       container.Image,
				Status:      container.Status,
				State:       container.State,
			}
			dockerStats = append(dockerStats, tmp)
		}
	} else {
		tmpStats, _ := docker.GetDockerStat()
		for _, docker := range tmpStats {
			tmp := DrcDockerStats{
				ContainerID: docker.ContainerID,
				Name:        docker.Name,
				Image:       docker.Image,
				Status:      docker.Status,
				State:       strconv.FormatBool(docker.Running),
			}
			dockerStats = append(dockerStats, tmp)
		}

	}
	return dockerStats
}
