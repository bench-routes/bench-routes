package utils

import (
	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"os/exec"
	"strconv"
	"fmt"
)

// cli command base
const (
	CmdPingBasedOnPacketsNumber = "ping"
)

// CLIPing works as an *independent subroutine*, for ping operations with the external networks
// Takes in a pointer channel in the last params inorder to implement subroutines since the ping
// operations might take time thereby avoiding delay in the other operations
// Use of pointers necessary since the data received is of large size, thereby slowing the traditional
// method of variables, as using variables require the time involved in loading into and out from cpu registers.
// Specifying addresses directly speeds the entire process manyfolds.
func CLIPing(url *string, packets int, cliPingChannel chan *string) {
	url = filters.HTTPPingFilter(url)
	fmt.Println("url:::::::::", url)
	cmd, err := exec.Command(CmdPingBasedOnPacketsNumber, "-c", strconv.Itoa(packets), *url).Output()
	if err != nil {
		panic(err)
	}
	cmdStr := string(cmd)
	cliPingChannel <- &cmdStr
}
