package system

import (
	"os/exec"
)

// SysDetails serves as a type for system details.
type SysDetails struct {
	CPU float32
	Memory float32
	DiskIO float32
}

// Metrics implements the metrics interface.
type Metrics struct {}

// New returns a new Metrics type that helps to implement
// the system metrics interface.
func New() *Metrics {
	return &Metrics{}
}

func (m *Metrics) GetTotalCPUUsage() float32 {
	ps, err := exec.Command("awk", )
}

