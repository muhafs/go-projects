package validators

import "regexp"

func IsEmail(email string) bool {
	pattern := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if len(email) < 3 {
		return false
	}

	if len(email) > 254 {
		return false
	}

	if !pattern.MatchString(email) {
		return false
	}

	return true
}
