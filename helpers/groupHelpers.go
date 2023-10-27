package helpers

import (
	"encoding/json"
	"errors"
)

func UnmarshalPendingGroups(pendingGroups string) ([]uint, error) {
	var userPendingGroups []uint
	err := json.Unmarshal([]byte(pendingGroups), &userPendingGroups)

	if err != nil {
		return nil, errors.New("Failed to unmarshal array")
	} else {
		return userPendingGroups, nil
	}
}
