package realtime

import "fmt"

func NotifyPeople(groupId string, title string, message string) {
	fmt.Print("Notification: ")
	fmt.Print(title, " ")
	fmt.Println(groupId, message)
}
