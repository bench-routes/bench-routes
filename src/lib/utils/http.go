package utils

import (
	"fmt"
	"github.com/zairza-cetb/bench-routes/src/lib/filters"
	"net/http"
	"os/exec"
	"strconv"
)

// cli command base
const (
	CmdPingBasedOnPacketsNumber = "ping"
	CmdAdministrator            = "sudo"
	CmdEcho                     = "echo"
)

//will retrieve this from a user settings file at later stages
var password = ""

// CLIPing works as an *independent subroutine*, for ping operations with the external networks
// Takes in a pointer channel in the last params inorder to implement subroutines since the ping
// operations might take time thereby avoiding delay in the other operations
// Use of pointers necessary since the data received is of large size, thereby slowing the traditional
// method of variables, as using variables require the time involved in loading into and out from cpu registers.
// Specifying addresses directly speeds the entire process manyfolds.
func CLIPing(url *string, packets int, cliPingChannel chan *string) {
	url = filters.HTTPPingFilter(url)
	cmd, err := exec.Command(CmdPingBasedOnPacketsNumber, "-c", strconv.Itoa(packets), *url).Output()
	if err != nil {
		panic(err)
	}
	cmdStr := string(cmd)
	cliPingChannel <- &cmdStr
}

// CLIFloodPing in another subroutine, for ping operation with -f flag
// which sends multiple ping request at once i.e. floods the url with requests.
func CLIFloodPing(url *string, packets int, cliPingChannel chan *string) {
	url = filters.HTTPPingFilter(url)

	cmd := fmt.Sprintf("%s -e \"%s\n\" | %s -S %s -f -c %s %s", CmdEcho, password, CmdAdministrator, CmdPingBasedOnPacketsNumber, strconv.Itoa(packets), *url)

	cmdPing, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		panic(err)
	}
	cmdStr := string(cmdPing)
	cliPingChannel <- &cmdStr
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
