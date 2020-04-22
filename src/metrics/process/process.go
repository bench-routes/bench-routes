package process

import (
	"fmt"
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
	State           string
	StartTime       string
	Time            string
	Command         string
	FilteredCommand string
	ThreadCount     int
}

// PBuffer corresponds to the in-memory store for current running-processes.
type PBuffer struct {
	ProcessesDetails      *[]PDetails
	TotalRunningProcesses int
}

// DecodeType generally used as decoding value for responding to the querier.
type DecodeType struct {
	CPUUtilization    string `json:"CPUUtilization"`
	MemoryUtilization string `json:"MemoryUtilization"`
	VMS               string `json:"VMS"`
	RSS               string `json:"RSS"`
	ThreadCount       string `json:"ThreadCount"`
}

// New returns a reader that reads over the running processes in a system.
func New() *PBuffer {
	return &PBuffer{}
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
	prc.TotalRunningProcesses = len(*tmp)

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
		pDetail.FilterCommandToUseableAddress()

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
	opt, err := exec.Command("ps", "-o", "nlwp", strconv.Itoa(pid)).Output()
	if err != nil {
		return -1
	}
	tmp := strings.TrimSpace(strings.Split(string(opt), "\n")[1])
	if i, err := strconv.Atoi(tmp); err != nil {
		panic(err)
	} else {
		return i
	}
}

// FilterCommandToUseableAddress filters illegal symbols from COMMAND
// so as to use as a normal path
func (p *PDetails) FilterCommandToUseableAddress() (s *string) {
	s = &p.Command
	sreplace(s, " ", "_")
	sreplace(s, "/", "@")
	p.FilteredCommand = *s
	return
}

// UnFilterCommandToUseableCommand converts the filtered string back to
// the original command.
func (p *PDetails) UnFilterCommandToUseableCommand() (s *string) {
	s = &p.FilteredCommand
	sreplace(s, "_", " ")
	sreplace(s, "@", "/")
	p.Command = *s
	return
}

// Encode encodes the process type block for inserting into the tsdb.
func (p *PDetails) Encode() string {
	return fmt.Sprintf("%f|%f|%f|%f|%d", p.CPUUtilization, p.MemoryUtilization, p.VMS, p.RSS, p.ThreadCount)
}

// Decode decodes the blocks from tsdb to be sent as monitor
// to the calling querier.
func (p *PDetails) Decode(s string) DecodeType {
	tmp := strings.Split(s, "|")
	return DecodeType{
		CPUUtilization:    tmp[0],
		MemoryUtilization: tmp[1],
		VMS:               tmp[2],
		RSS:               tmp[3],
		ThreadCount:       tmp[4],
	}
}

func sreplace(s *string, prev, new string) {
	*s = strings.ReplaceAll(*s, prev, new)
}
