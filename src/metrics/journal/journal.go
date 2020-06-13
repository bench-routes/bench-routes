package journal

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const (
	indentifier      = `"SYSLOG_IDENTIFIER":"kernel"`
	logsIdentifier   = "localhost.localdomain"
	kernelIndetifier = "kernel: "
)

// Journal implements for collecting journal/systemd information.
type Journal struct {
	cmd
	Points
	warn, err string
}

type cmd struct{}

// Points are data-points for journal information counters.
// C stands for count.
// K stands for kernel.
// CK for count kernel.
type Points struct {
	Cwarn    int `json:"cwarn"`
	Cerr     int `json:"cerr"`
	Ckwarn   int `json:"ckwarn"`
	Ckerr    int `json:"ckerr"`
	Cevents  int `json:"cevents"`
	Ckevents int `json:"ckevents"`
}

// New returns a Journal.
func New() *Journal {
	return &Journal{}
}

// Run runs and fetches the journal based information.
func (j *Journal) Run() *Points {
	logs1 := make(chan string)
	logs2 := make(chan string)
	klogs1 := make(chan int)
	klogs2 := make(chan int)
	klogs3 := make(chan int)
	klogs4 := make(chan int)

	go j.cmd.warnLog(logs1)
	go j.cmd.errLog(logs2)
	go j.cmd.logs(klogs1, klogs2)

	j.warn = <-logs1
	j.err = <-logs2
	j.Points.Cevents = <-klogs1
	j.Points.Ckevents = <-klogs2

	go j.cmd.kernels(j.warn, klogs3)
	go j.cmd.kernels(j.err, klogs4)

	j.Points.Cwarn = len(strings.Split(j.warn, "\n"))
	j.Points.Cerr = len(strings.Split(j.err, "\n"))
	j.Points.Ckwarn = <-klogs3
	j.Points.Ckerr = <-klogs4

	return &j.Points
}

// Get returns the data-points.
func (p *Points) Get() *Points {
	return p
}

// Encode encodes the data-points into br's tsdb compatible
// blocks.
func (p *Points) Encode() *string {
	enc := fmt.Sprintf("%d|%d|%d|%d|%d|%d", p.Cerr, p.Cwarn, p.Ckerr, p.Ckwarn, p.Cevents, p.Ckevents)
	return &enc
}

func (j *cmd) warnLog(c chan string) {
	o, err := exec.Command("journalctl", "-p", "warning", "-b", "-o", "short-full", "--no-pager").Output()
	if err != nil {
		panic(err)
	}
	c <- string(o)
}

func (j *cmd) errLog(c chan string) {
	o, err := exec.Command("journalctl", "-p", "err", "-b", "-o", "short-full", "--no-pager").Output()
	if err != nil {
		panic(err)
	}
	c <- string(o)
}

func (j *cmd) logs(c, ck chan int) {
	o, err := exec.Command("journalctl", "-b", "-o", "short-full", "--no-pager").Output()
	if err != nil {
		panic(err)
	}
	s := string(o)
	lines := strings.Split(s, "\n")
	cl := len(lines)
	count := 0
	kcount := 0
	curr := 0
	for {
		if curr == cl {
			break
		}
		if strings.ContainsAny(lines[curr], logsIdentifier) {
			count++
			if strings.ContainsAny(lines[curr], kernelIndetifier) {
				kcount++
			}
		}
		curr++
	}
	c <- count
	ck <- kcount
}

func (j *cmd) kernels(l string, count chan int) {
	lines := strings.Split(l, "\n")
	c := 0
	curr := 0
	cl := len(lines)
	for {
		if curr == cl {
			break
		}
		if strings.ContainsAny(lines[curr], indentifier) {
			c++
		}
		curr++
	}
	count <- c
}

// Decode converts the marshalled data-point to valid data-points.
func Decode(datapointArr []string) Points {
	var p Points
	for i, b := range datapointArr {
		switch i {
		case 0:
			p.Cerr = sToi(b)
		case 1:
			p.Cwarn = sToi(b)
		case 2:
			p.Ckerr = sToi(b)
		case 3:
			p.Ckwarn = sToi(b)
		case 4:
			p.Cevents = sToi(b)
		case 5:
			p.Ckevents = sToi(b)
		default:
			panic(fmt.Sprintf("invalid decoding with index: %d", i))
		}
	}
	return p
}

func sToi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
