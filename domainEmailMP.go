package main

import (
	"fmt"
	"regexp"
	"strings"
)

//Regex
var md5Regex, _ = regexp.Compile("^[a-f0-9]{32}$")
var domainRegex, _ = regexp.Compile("^((|\\*)@)[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,})$")
var emailRegex, _ = regexp.Compile("^((|\\*|([a-z0-9_!#$%&\\'+\\/=?^`{|}~-]+\\.?)+)@)[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,})$")

func domainEmailMultiPlex(recs *int, errorChan, lineChan, domainChan, hashedLineChan *chan string) {
	var line, trimmed string
	for line = range *lineChan {

		// clean the input
		trimmed = strings.TrimSpace(line)
		trimmed = strings.ToLower(trimmed)

		switch {
		case domainRegex.MatchString(trimmed):
			//			fmt.Println("domain", trimmed)
			*domainChan <- trimmed
			*recs++
		case emailRegex.MatchString(trimmed):
			//			fmt.Println("email", trimmed)
			trimmed = forceMd5(trimmed)
			*hashedLineChan <- trimmed
			*recs++
		case md5Regex.MatchString(trimmed):
			//			fmt.Println("md5", trimmed)
			*hashedLineChan <- trimmed
			*recs++
		default:
			//			fmt.Println("no match", trimmed)
			//save as errors maybe?
			*errorChan <- trimmed
		}

	}
	fmt.Println("Total Records: ", *recs)

	defer close(*domainChan)
	defer close(*hashedLineChan)
	defer close(*errorChan)

}
