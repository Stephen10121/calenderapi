package realtime

import "fmt"

func NotifyPeople(groupId uint, title string, message string) {
	fmt.Print("Notification: ")
	fmt.Print(title, " ")
	fmt.Println(groupId, message)
}

func NotifyGroupOwner(groupId uint, title string, message string) {
	fmt.Print("To: Group ")
	fmt.Print(groupId)
	fmt.Println(" owner.")
	fmt.Print(title)
	fmt.Print(": ")
	fmt.Println(message)
}
