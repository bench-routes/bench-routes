package journal

import (
	"os/exec"
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
	points
	warn, err string
}

type cmd struct{}
type points struct {
	cwarn, cerr, ckwarn, ckerr, cevents, ckevents int
}

// New returns a Journal.
func New() *Journal {
	return &Journal{}
}

// Run runs and fetches the journal based information.
func (j *Journal) Run() {
	logs := make([]chan string, 2)
	klogs := make([]chan int, 4)
	go j.cmd.warnLog(logs[0])
	go j.cmd.errLog(logs[1])
	go j.cmd.logs(klogs[2], klogs[3])
	j.warn = <-logs[0]
	j.err = <-logs[1]
	j.points.cevents = <-klogs[2]
	j.points.ckevents = <-klogs[3]
	go j.cmd.kernels(j.warn, klogs[0])
	go j.cmd.kernels(j.err, klogs[1])
	j.points.cwarn = len(strings.Split(j.warn, "\n"))
	j.points.cerr = len(strings.Split(j.err, "\n"))
	j.points.ckwarn = <-klogs[0]
	j.points.ckerr = <-klogs[1]
}

func (j *cmd) warnLog(c chan string) {
	o, err := exec.Command("journalctl", "-p", "warning", "-b", "-o", "json-seq", "-b", "|", "cat").Output()
	if err != nil {
		panic(err)
	}
	c <- string(o)
}

func (j *cmd) errLog(c chan string) {
	o, err := exec.Command("journalctl", "-p", "err", "-b", "-o", "json-seq", "-b", "|", "cat").Output()
	if err != nil {
		panic(err)
	}
	c <- string(o)
}

func (j *cmd) logs(c, ck chan int) {
	o, err := exec.Command("journalctl", "-b", "|", "cat").Output()
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
