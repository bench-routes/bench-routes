package utils

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/zairza-cetb/bench-routes/src/lib/filters"
)

// cli command base
const (
	CmdPingBasedOnPacketsNumber = "ping"
	CmdAdministrator            = "sudo"
	CmdEcho                     = "echo"
)

// CLIPing works as an *independent subroutine*, for ping operations with the external networks
// Takes in a pointer channel in the last params inorder to implement subroutines since the ping
// operations might take time thereby avoiding delay in the other operations
// Use of pointers necessary since the data received is of large size, thereby slowing the traditional
// method of variables, as using variables require the time involved in loading into and out from cpu registers.
// Specifying addresses directly speeds the entire process manyfolds.
func CLIPing(url string, packets int) (*string, error) {
	url = *filters.HTTPPingFilter(&url)
	cmd, err := exec.Command(CmdPingBasedOnPacketsNumber, "-c", strconv.Itoa(packets), url).Output()
	if err != nil {
		return nil, fmt.Errorf("err: %s, url: %s", err.Error(), url)
	}
	cmdStr := string(cmd)
	return &cmdStr, err
}

// CLIFloodPing in another subroutine, for ping operation with -f flag
// which sends multiple ping request at once i.e. floods the url with requests.
func CLIFloodPing(url string, packets int, password string) (*string, error) {
	url = *filters.HTTPPingFilter(&url)
	cmd := fmt.Sprintf("%s -e \"%s\n\" | %s -S %s -f -c %s %s", CmdEcho, password, CmdAdministrator, CmdPingBasedOnPacketsNumber, strconv.Itoa(packets), url)
	cmdPing := exec.Command("bash", "-c", cmd)
	cmdOuput, err := cmdPing.CombinedOutput()
	if err != nil {
		// There was an issue
		// executing the command.
		panic(err)
	}
	cmdStr := string(cmdOuput)
	return &cmdStr, err
}

//SendGETRequest sends http GET request to the specified url(both resp_delay and monitor_response_status module use it)
func SendGETRequest(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		// Prone to alerting, printing for now
		fmt.Println(err)
	}
	return resp
}
