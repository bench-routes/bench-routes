package notify

import(
	"strconv"
	"os/exec"
	"fmt"
)
//NotificationData contains the different data for sending notifications
type NotificationData struct{
	title string
	description string
	icon string
	urgency string
	time int
}
//NotificationManage sends notifications to the user based on different params passed
//All the methods in this interface based on different params
type NotificationManage interface{
	//Notify1 with title and description
	Notify1()
	//Notify2 with title and description and lives for a certain time span
	Notify2()
	//Notify3 with title, description and image icon
	Notify3()
	//Notify4 with title, description and urgency level(low,normal,critical)
	Notify4()
	//Notify5 with title, description, image and lives for a certain time span
	Notify5()
	//Notify6 with title and description, urgency and lives for a certain time span
	Notify6()
}
// Notify1 executes the command to notify with title and description
func (N *NotificationData)Notify1(){
	cmd:=exec.Command("notify-send", N.title, N.description)
	err := cmd.Run()
	if(err!=nil){
		fmt.Println(err)
	}
}
//Notify2 executes the command to notify with title, des and time
func (N *NotificationData)Notify2(){
	cmd:=exec.Command("notify-send", N.title, N.description, "-t", strconv.Itoa(N.time))
	err := cmd.Run()
	if(err!=nil){
		fmt.Println(err)
	}
}
// Notify3 executes the command to notify with title, description, and icon
func (N *NotificationData)Notify3(){
	cmd:=exec.Command("notify-send", N.title, N.description, "-i", N.icon )
	err := cmd.Run()
	if(err!=nil){
		fmt.Println(err)
	}
}
// Notify4 executes the command to notify with title, description and urgency level
func (N *NotificationData)Notify4(){

	cmd:=exec.Command("notify-send", N.title, N.description, "-u", N.urgency)
	err := cmd.Run()
	if(err!=nil){
		fmt.Println(err)
	}
}
// Notify5 executes the command to notify with title, description, icon and time
func (N *NotificationData)Notify5(){
	cmd:=exec.Command("notify-send", N.title, N.description, "-i", "/home/jayashree/go/src/github.com/zairza-cetb/bench-routes/assets/jj.png", "-t", strconv.Itoa(N.time))
	err := cmd.Run()
	if(err!=nil){
		fmt.Println(err)
	}
}
// Notify6 executes the command to notify with title, description, urgency and time
func (N *NotificationData)Notify6(){

	cmd:=exec.Command("notify-send", N.title, N.description, "-u", N.urgency, "-t", strconv.Itoa(N.time))
	err := cmd.Run()
	if(err!=nil){
		fmt.Println(err)
	}
}




