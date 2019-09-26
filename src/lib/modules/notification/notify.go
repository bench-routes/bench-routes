package notify

import (
	"fmt"
	"os/exec"
	"strconv"
)

const (
	//CmdNotify is an alias for cli command notify-send
	CmdNotify = "notify-send"
)

//NotificationData contains the different data for sending notifications
type NotificationData struct {
	title       string
	description string
	icon        string
	urgency     string
	time        int
}

//Notification sends notifications to the user based on different params passed
//All the methods in this interface based on different params
type Notification interface {
	//NotifyBasic with title and description
	NotifyBasic()
	//NotifyWithTimeSpan with title and description and lives for a certain time span
	NotifyWithTimeSpan()
	//NotifyWithImageIcon with title, description and image icon
	NotifyWithImageIcon()
	//NotifyWithUrgencyLevel with title, description and urgency level(low,normal,critical)
	NotifyWithUrgencyLevel()
	//NotifyWithImageAndTimeSpan with title, description, image and lives for a certain time span
	NotifyWithImageAndTimeSpan()
	//NotifyWithUrgencyLevelAndTimeSpan with title and description, urgency and lives for a certain time span
	NotifyWithUrgencyLevelAndTimeSpan()
}

// NotifyBasic executes the command to notify with title and description
func (N *NotificationData) NotifyBasic() {
	cmd := exec.Command(CmdNotify, N.title, N.description)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

//NotifyWithTimeSpan executes the command to notify with title, des and time
func (N *NotificationData) NotifyWithTimeSpan() {
	cmd := exec.Command(CmdNotify, N.title, N.description, "-t", strconv.Itoa(N.time))
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

// NotifyWithImageIcon executes the command to notify with title, description, and icon
func (N *NotificationData) NotifyWithImageIcon() {
	cmd := exec.Command(CmdNotify, N.title, N.description, "-i", N.icon)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

// NotifyWithUrgencyLevel executes the command to notify with title, description and urgency level
func (N *NotificationData) NotifyWithUrgencyLevel() {

	cmd := exec.Command(CmdNotify, N.title, N.description, "-u", N.urgency)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

// NotifyWithImageAndTimeSpan executes the command to notify with title, description, icon and time
func (N *NotificationData) NotifyWithImageAndTimeSpan() {
	cmd := exec.Command(CmdNotify, N.title, N.description, "-i", N.icon, "-t", strconv.Itoa(N.time))
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

// NotifyWithUrgencyLevelAndTimeSpan executes the command to notify with title, description, urgency and time
func (N *NotificationData) NotifyWithUrgencyLevelAndTimeSpan() {

	cmd := exec.Command(CmdNotify, N.title, N.description, "-u", N.urgency, "-t", strconv.Itoa(N.time))
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
