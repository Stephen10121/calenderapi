package functions

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func UintContains(s []uint, str uint) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func Remove(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// var groups []uint
// json.Unmarshal([]byte(user.Groups), &groups)

// groupsJson, _ := json.Marshal([]uint{})
