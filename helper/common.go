package helper

import (
	"regexp"
	"strings"
)

func ValidateNumber(numb string) bool {
	re := regexp.MustCompile(`[^0-9]*1[34578][0-9]{9}[^0-9]*`)
	if len(numb) > 13 {
		return false
	}
	return re.MatchString(numb)
}

func CheckErr(err string) string {
	if strings.Contains(err, "oneof") {
		return "Invalid tag has been detected!"
	}

	return "Something Went's Wrong"
}
