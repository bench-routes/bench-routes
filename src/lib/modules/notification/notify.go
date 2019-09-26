package notify

import (
	"fmt"
	"os/exec"
	"strconv"
)

//NotificationData contains the different data for sending notifications
type NotificationData struct {
	title       string
	description string
	icon        string
	urgency     string
	time        int
}

//NotificationManage sends notifications to the user based on different params passed
//All the methods in this interface based on different params
type NotificationManage interface {
	//SendNotification1 with title and description
	SendNotification1()
	//SendNotification2 with title and description and lives for a certain time span
	SendNotification2()
	//SendNotification3 with title, description and image icon
	SendNotification3()
	//SendNotification4 with title, description and urgency level(low,normal,critical)
	SendNotification4()
	//SendNotification5 with title, description, image and lives for a certain time span
	SendNotification5()
	//SendNotification6 with title and description, urgency and lives for a certain time span
	SendNotification6()
}

// SendNotification1 executes the command to notify with title and description
func (N *NotificationData) SendNotification1() {
	cmd := exec.Command("notify-send", N.title, N.description)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

//SendNotification2 executes the command to notify with title, des and time
func (N *NotificationData) SendNotification2() {
	cmd := exec.Command("notify-send", N.title, N.description, "-t", strconv.Itoa(N.time))
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

// SendNotification3 executes the command to notify with title, description, and icon
func (N *NotificationData) SendNotification3() {
	cmd := exec.Command("notify-send", N.title, N.description, "-i", N.icon)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

// SendNotification4 executes the command to notify with title, description and urgency level
func (N *NotificationData) SendNotification4() {

	cmd := exec.Command("notify-send", N.title, N.description, "-u", N.urgency)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

// SendNotification5 executes the command to notify with title, description, icon and time
func (N *NotificationData) SendNotification5() {
	cmd := exec.Command("notify-send", N.title, N.description, "-i", N.icon, "-t", strconv.Itoa(N.time))
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

// SendNotification6 executes the command to notify with title, description, urgency and time
func (N *NotificationData) SendNotification6() {

	cmd := exec.Command("notify-send", N.title, N.description, "-u", N.urgency, "-t", strconv.Itoa(N.time))
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
