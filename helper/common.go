package helper

import "regexp"

func ValidateNumber(numb string) bool {
	re := regexp.MustCompile(`[^0-9]*1[34578][0-9]{9}[^0-9]*`)
	if len(numb) > 13 {
		return false
	}
	return re.MatchString(numb)
}
