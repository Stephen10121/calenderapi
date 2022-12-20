package realtime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

var endpointOfficial string = "https://socket.gruzservices.com"

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

func UserKickedOut(groupID string, userId uint) {
	fmt.Print(userId)
	fmt.Print(" got kicked out of ")
	fmt.Println(groupID)
	fmt.Print(userId)
	fmt.Print(" joined ")
	fmt.Println(groupID)
	userIdString := string(strconv.FormatUint(uint64(userId), 10))
	endpoint := endpointOfficial + "/particapantDeleted?groupId=" + url.QueryEscape(groupID) + "&userId=" + url.QueryEscape(userIdString)
	fmt.Println(endpoint)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Add("Authorization", "Bearer secretKey")

	if err != nil {
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return
	}
	fmt.Println(resp)
	return
}

func UserJoiningGroup(groupId string, newUser string, ownerId uint) {
	ownerToString := string(strconv.FormatUint(uint64(ownerId), 10))
	endpoint := endpointOfficial + "/newPendingUser?groupId=" + url.QueryEscape(groupId) + "&newUser=" + url.QueryEscape(newUser) + "&ownerId=" + url.QueryEscape(string(ownerToString))
	fmt.Println(endpoint)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Add("Authorization", "Bearer secretKey")

	if err != nil {
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return
	}
	fmt.Println(resp)
	return
}

func UserGotAccepted(groupID string, userId uint, owner string, othersCanAdd bool) {
	fmt.Print(userId)
	fmt.Print(" joined ")
	fmt.Println(groupID)
	var othersCanAdd2 string
	if othersCanAdd {
		othersCanAdd2 = "1"
	} else {
		othersCanAdd2 = "0"
	}
	userIdString := string(strconv.FormatUint(uint64(userId), 10))
	endpoint := endpointOfficial + "/groupAccepted?groupId=" + url.QueryEscape(groupID) + "&userId=" + url.QueryEscape(userIdString) + "&owner=" + url.QueryEscape(owner) + "&othersCanAdd=" + othersCanAdd2
	fmt.Println(endpoint)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Add("Authorization", "Bearer secretKey")

	if err != nil {
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return
	}
	fmt.Println(resp)
	return
}

func UserGotRejected(groupID string, userId uint) {
	fmt.Print(userId)
	fmt.Print(" got rejected from ")
	fmt.Println(groupID)
	userIdString := string(strconv.FormatUint(uint64(userId), 10))
	endpoint := endpointOfficial + "/particapantDeletedPending?groupId=" + url.QueryEscape(groupID) + "&userId=" + url.QueryEscape(userIdString)
	fmt.Println(endpoint)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Add("Authorization", "Bearer secretKey")

	if err != nil {
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return
	}
	fmt.Println(resp)
	return
}

func UserLeft(groupID uint, userId uint) {
	fmt.Print(userId)
	fmt.Print(" left ")
	fmt.Println(groupID)
}

func UserLeftTransfered(groupID string, particapants string, newOwner string) {
	fmt.Print(" left ")
	fmt.Print(groupID)
	fmt.Print(" and transfered ownership to ")
	fmt.Println(newOwner)
	endpoint := endpointOfficial + "/userLeftTransfered?groupId=" + url.QueryEscape(groupID) + "&particapants=" + url.QueryEscape(particapants) + "&newOwner=" + url.QueryEscape(newOwner)
	fmt.Println(endpoint)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Add("Authorization", "Bearer secretKey")

	if err != nil {
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return
	}
	fmt.Println(resp)
	return
}

func UserLeftWhilePending(groupID uint, userId uint) {
	fmt.Print(userId)
	fmt.Print(" left while pending ")
	fmt.Println(groupID)
}

func GroupDeleted(groupID string, particapants []uint, pendingParticapants []uint) {
	fmt.Print("Deleted ")
	fmt.Println(groupID)
	particapantsJson, _ := json.Marshal(particapants)
	pendingParticapantsJson, _ := json.Marshal(pendingParticapants)
	req, err := http.NewRequest(http.MethodGet, endpointOfficial+"/groupDeleted?groupId="+url.QueryEscape(groupID)+"&particapants="+string(particapantsJson)+"&pendingParticapants="+string(pendingParticapantsJson), nil)
	req.Header.Add("Authorization", "Bearer secretKey")

	if err != nil {
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return
	}
	fmt.Println(resp)
	return
}

func UserLoggedIn(token string, userId uint) {
	fmt.Print("User Logged In ")
	fmt.Print(token)
	fmt.Print(" id: ")
	fmt.Println(userId)
	req, err := http.NewRequest(http.MethodGet, endpointOfficial+"/newConnection?token="+token+"&id="+string(strconv.FormatUint(uint64(userId), 10)), nil)
	req.Header.Add("Authorization", "Bearer secretKey")

	if err != nil {
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return
	}
	fmt.Println(resp)
	return
}
