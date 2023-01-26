package ip_regex

import "regexp"

var IPorCIDRregex *regexp.Regexp
var IPonlyRegex *regexp.Regexp

func GetIPorCIDRregex() *regexp.Regexp {
	if IPorCIDRregex == nil {
		numBlock := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
		subnet := "(\\/3[0-2]|\\/[1-2][0-9]|\\/[0-9])?"
		regexPattern := numBlock + "\\." + numBlock + "\\." + numBlock + "\\." + numBlock + subnet
		IPorCIDRregex = regexp.MustCompile(regexPattern)
	}
	return IPorCIDRregex
}

func GetIPonlyRegex() *regexp.Regexp {
	if IPonlyRegex == nil {
		numBlock := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
		regexPattern := numBlock + "\\." + numBlock + "\\." + numBlock + "\\." + numBlock
		IPonlyRegex = regexp.MustCompile(regexPattern)
	}
	return IPonlyRegex
}
