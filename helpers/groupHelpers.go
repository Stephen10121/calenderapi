package helpers

import (
	"encoding/json"
	"errors"
)

func UnmarshalPendingGroups(pendingGroups string) ([]uint, error) {
	var userPendingGroups []uint
	err := json.Unmarshal([]byte(pendingGroups), &userPendingGroups)

	if err != nil {
		return nil, errors.New("Failed to unmarshal array of pending groups")
	} else {
		return userPendingGroups, nil
	}
}

func UnmarshalGroupParticapants(groupParticapants string) ([]uint, error) {
	var unmarshalledGroupParticapants []uint
	err := json.Unmarshal([]byte(groupParticapants), &unmarshalledGroupParticapants)

	if err != nil {
		return nil, errors.New("Failed to unmarshal array of group particapants.")
	} else {
		return unmarshalledGroupParticapants, nil
	}
}
