package system

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	disk "github.com/mackerelio/go-osstat/memory"
	gostats "github.com/shirou/gopsutil/cpu"
	gomem "github.com/shirou/gopsutil/mem"
)

const (
	tcpSegmentSizeInBytes = 1200
)

// CPU interface for CPU related details.
type CPU interface {
	GetTotalCPUUsage() string
	GetTotalCPUStats() gostats.TimesStat
}

// Memory interface for Memory related details.
type Memory interface {
	GetVirtualMemoryStats() MemoryStats
}

// Disk interface for Disk related details.
type Disk interface {
	GetDiskIOStats() DiskStats
}

// SystemMetrics contains system level metrics.
// These are very basic as of now and can be
// made to collect advanced informations with
// the same implementation.
type SystemMetrics struct{}

// New returns a new system-metrics reader.
func New() *SystemMetrics {
	return &SystemMetrics{}
}

// GetTotalCPUUsage returns the total CPU usage by all the
// available cores over the previous second.
func (s *SystemMetrics) GetTotalCPUUsage(c chan *string) {
	usage, err := gostats.Percent(time.Duration(time.Second), false)
	if err != nil {
		panic(err)
	}
	str := fmt.Sprintf("%.2f", usage[0])
	c <- &str
}

// GetTotalCPUStats returns the stats related to the CPU
// as a whole.
func (s *SystemMetrics) GetTotalCPUStats() gostats.TimesStat {
	stats, err := gostats.Times(false)
	if err != nil {
		panic(err)
	}

	return stats[0]
}

// MemoryStats for memory based statistics.
type MemoryStats struct {
	Total       uint64  `json:"totalBytes"`
	Available   uint64  `json:"availableBytes"`
	Used        uint64  `json:"usedBytes"`
	UsedPercent float64 `json:"usedPercent"`
	Free        uint64  `json:"freeBytes"`
}

// MemoryStatsStringified for http response.
type MemoryStatsStringified struct {
	Total       string `json:"totalBytes"`
	Available   string `json:"availableBytes"`
	Used        string `json:"usedBytes"`
	UsedPercent string `json:"usedPercent"`
	Free        string `json:"freeBytes"`
}

// GetVirtualMemoryStats returns the memory statistics of the host machine.
func (s *SystemMetrics) GetVirtualMemoryStats(c chan *MemoryStats) {
	stats, err := gomem.VirtualMemory()
	if err != nil {
		panic(err)
	}

	c <- &MemoryStats{
		// default is always in bytes. hence, convert into the required format.
		Total:       stats.Total / 1000000,
		Available:   stats.Available / 1000000,
		Used:        stats.Used / 1000000,
		UsedPercent: stats.UsedPercent,
		Free:        stats.Free / 1000000,
	}
}

// DiskStats statistics for information related to the disk.
type DiskStats struct {
	DiskIO int `json:"diskIO"`
	Cached int `json:"cached"`
}

// DiskStatsStringified for http response.
type DiskStatsStringified struct {
	DiskIO string `json:"diskIO"`
	Cached string `json:"cached"`
}

// GetDiskIOStats returns the disk stats: IO per sec and cached volume.
func (s *SystemMetrics) GetDiskIOStats(c chan *DiskStats) {
	before, err := disk.Get()
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Duration(time.Second * 3))

	now, err := disk.Get()
	if err != nil {
		panic(err)
	}

	c <- &DiskStats{
		// default is always in bytes. hence, convert into the required format.
		DiskIO: (int(now.Used) - int(before.Used)) / 1000, // in kilo-bytes
		Cached: int(now.Cached / 1000000),                 // mega-bytes
	}
}

// NetworkStats for the network stats.
type NetworkStats struct {
	// P stands for packets.
	// S stands for size in kilo-bytes.
	PtcpIncoming int `json:"PtcpIncoming"`
	PtcpOutgoing int `json:"PtcpOutgoing"`
	StcpIncoming int `json:"StcpIncoming"`
	StcpOutgoing int `json:"StcpOutgoing"`
	PudpIncoming int `json:"PudpIncoming"`
	PudpOutgoing int `json:"PudpOutgoing"`
}

var (
	// buff stores the previous values that are required
	// for calculating the network speed.
	buffIncoming, buffOutgoing int
)

// GetNetworkStats returns the network stats.
func (s *SystemMetrics) GetNetworkStats(c chan *NetworkStats) {
	o, err := exec.Command("netstat", "-s").Output()
	if err != nil {
		panic(err)
	}
	var netStats NetworkStats
	str := string(o)
	lines := strings.Split(str, "\n")
	l := len(lines)
	curr := -1
	for {
		curr++
		if curr == l {
			break
		}
		line := lines[curr]
		if strings.TrimSpace(line) == "Tcp:" {
			for {
				curr++
				if strings.Contains(lines[curr], "segments received") {
					arr := strings.Split(strings.TrimSpace(lines[curr]), " ")
					netStats.PtcpIncoming, err = strconv.Atoi(arr[0])
					if err != nil {
						panic(err)
					}
					if buffIncoming == 0 {
						buffIncoming = netStats.PtcpIncoming
					} else {
						netStats.StcpIncoming = ((netStats.PtcpIncoming - buffIncoming) * tcpSegmentSizeInBytes) / 1024
						buffIncoming = netStats.PtcpIncoming
					}
					continue
				}
				if strings.Contains(lines[curr], "segments sent out") {
					arr := strings.Split(strings.TrimSpace(lines[curr]), " ")
					netStats.PtcpOutgoing, err = strconv.Atoi(arr[0])
					if err != nil {
						panic(err)
					}
					if buffOutgoing == 0 {
						buffOutgoing = netStats.PtcpOutgoing
					} else {
						netStats.StcpOutgoing = ((netStats.PtcpOutgoing - buffOutgoing) * tcpSegmentSizeInBytes) / 1024
						buffOutgoing = netStats.PtcpOutgoing
					}
					break
				}
			}
		}
		if strings.TrimSpace(line) == "Udp:" {
			for {
				curr++
				if strings.Contains(lines[curr], "packets received") {
					arr := strings.Split(strings.TrimSpace(lines[curr]), " ")
					netStats.PudpIncoming, err = strconv.Atoi(arr[0])
					if err != nil {
						panic(err)
					}
					continue
				}
				if strings.Contains(lines[curr], "packets sent") {
					arr := strings.Split(strings.TrimSpace(lines[curr]), " ")
					netStats.PudpOutgoing, err = strconv.Atoi(arr[0])
					if err != nil {
						panic(err)
					}
					break
				}
			}
			break
		}
	}
	c <- &netStats
}

// Encode encodes the blocks into format that can be consumed
// by the tsdb module.
func (s *SystemMetrics) Encode(block interface{}) string {
	switch node := block.(type) {
	case DiskStats:
		return fmt.Sprintf("%d|%d", node.DiskIO, node.Cached)
	case MemoryStats:
		return fmt.Sprintf("%d|%d|%d|%f|%d",
			node.Total, node.Available, node.Used, node.UsedPercent, node.Free,
		)
	case NetworkStats:
		return fmt.Sprintf("%d|%d|%d|%d|%d|%d",
			node.PtcpIncoming, node.PtcpOutgoing, node.StcpIncoming, node.StcpOutgoing, node.PudpIncoming, node.PudpOutgoing,
		)
	case string:
		return node
	}

	data, ok := block.(string)
	if !ok {
		panic(fmt.Sprintf("Invalid block type: %v (block)", block))
	}

	return data
}

// Response is used to decode the tsdb blocks to data points that supports JSON encoding.
type Response struct {
	CPUTotalUsage string                 `json:"cpuTotalUsage"`
	Memory        MemoryStatsStringified `json:"memory"`
	Disk          DiskStatsStringified   `json:"disk"`
	Network       NetworkStats           `json:"network"`
}

// Combine combines segments into a single data point for a block.
func (s *SystemMetrics) Combine(cpu, memory, disk, net string) string {
	return cpu + "|" + memory + "|" + disk + "|" + net
}
