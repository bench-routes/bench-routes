package process

import (
	"os/exec"
	"strconv"
	"strings"
)

// PDetails contains the details of all the processes
type PDetails struct {
	User              string
	Pid               int
	CPUUtilization    float32
	MemoryUtilization float32
	// Virtual Memory Size
	VMS float32
	// Resident Set Size
	RSS float32
	// System State. Reference link for more information: https://askubuntu.com/a/360253
	State       string
	StartTime   string
	Time        string
	Command     string
	ThreadCount int
}

// PBuffer type
type PBuffer struct {
	ProcessesDetails *[]PDetails
}

// NewProcessReader returns a reader that reads over the running processes in a system.
func NewProcessReader() *PBuffer {
	return &PBuffer{
		ProcessesDetails: nil,
	}
}

// UpdateCurrentProcesses updates the process details list
func (prc *PBuffer) UpdateCurrentProcesses() (*[]PDetails, error) {
	cmd, err := exec.Command("ps", "aux").Output()
	if err != nil {
		return nil, err
	}
	rawProcessesTable := string(cmd)
	tmp := processDetailedRawTableFromAUX(rawProcessesTable)
	prc.ProcessesDetails = tmp

	return tmp, nil
}

func processDetailedRawTableFromAUX(table string) *[]PDetails {
	processlines := strings.Split(table, "\n")[1:] // discard first line to avoid headings.
	var pDetails []PDetails
	for _, processInstance := range processlines {
		var pDetail PDetails
		for i, w := range removeEmptyStringsFromArray(strings.Split(processInstance, " ")) {
			switch i {
			case 0:
				pDetail.User = w
			case 1:
				n, err := strconv.Atoi(w)
				if err != nil {
					panic(err)
				}
				pDetail.Pid = n
			case 2:
				pDetail.CPUUtilization = parseToFloat32(w)
			case 3:
				pDetail.MemoryUtilization = parseToFloat32(w)
			case 4:
				pDetail.VMS = parseToFloat32(w)
			case 5:
				pDetail.RSS = parseToFloat32(w)
			case 7:
				pDetail.State = w
			case 8:
				pDetail.StartTime = w
			case 9:
				pDetail.Time = w
			case 10:
				pDetail.Command = w
			}
		}
		pDetail.ThreadCount = getProcessThreadCount(pDetail.Pid)
		pDetails = append(pDetails, pDetail)
	}

	return &pDetails
}

func removeEmptyStringsFromArray(s []string) []string {
	var a []string
	for _, inst := range s {
		if len(inst) > 0 {
			a = append(a, inst)
		}
	}

	return a
}

func parseToFloat32(s string) float32 {
	f, err := strconv.ParseFloat(s, 32)
	f32 := float32(f)
	if err != nil {
		panic(err)
	}

	return f32
}

// getProcessThreadCount returns the count of threads involved in a particular process.
func getProcessThreadCount(pid int) int {
	// number of light-weight processes that corresponds to the threads.
	opt, err := exec.Command("ps", "-o", "nlwp", string(pid)).Output()
	if err != nil {
		panic(err)
	}

	if i, err := strconv.Atoi(string(opt)); err != nil {
		panic(err)
	} else {
		return i
	}
}
