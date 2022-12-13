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

func UserKickedOut(groupID uint, userId uint) {
	fmt.Print(userId)
	fmt.Print(" got kicked out of ")
	fmt.Println(groupID)
}

func UserGotAccepted(groupID uint, userId uint) {
	fmt.Print(userId)
	fmt.Print(" joined ")
	fmt.Println(groupID)
}

func UserGotRejected(groupID uint, userId uint) {
	fmt.Print(userId)
	fmt.Print(" got rejected from ")
	fmt.Println(groupID)
}

func UserLeft(groupID uint, userId uint) {
	fmt.Print(userId)
	fmt.Print(" left ")
	fmt.Println(groupID)
}

func UserLeftTransfered(groupID uint, userId uint, newOwner uint) {
	fmt.Print(userId)
	fmt.Print(" left ")
	fmt.Print(groupID)
	fmt.Print(" and transfered ownership to ")
	fmt.Println(newOwner)
}

func UserLeftWhilePending(groupID uint, userId uint) {
	fmt.Print(userId)
	fmt.Print(" left while pending ")
	fmt.Println(groupID)
}

func GroupDeleted(groupID uint) {
	fmt.Print("Deleted ")
	fmt.Println(groupID)
}
