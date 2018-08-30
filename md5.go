package main

import (
	"crypto/md5"
	"fmt"
	_ "strings"
)

/*
* takes string trims and lowercases it, converts to md5 if not md5
* @param line string
* @return match bool
 */
func loopForceMd5(emailChan, hashedLineChan *chan string) {

	for line := range *emailChan {

		hashedTrimmed := forceMd5(line)

		*hashedLineChan <- hashedTrimmed
	}
	//close( *hashedLineChan )
	//defer wg.Done()

}

func forceMd5(line string) (hashedTrimmed string) {

	bytes := []byte(hashedTrimmed)
	hashedBytes := md5.Sum(bytes)
	hashedTrimmed = fmt.Sprintf("%x", hashedBytes)

	return
}
