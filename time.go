package ari

import (
	"fmt"
	"strings"
	"time"
)

type AriTime time.Time

func (j *AriTime) UnmarshalJSON(input []byte) error {
	// ARI stamps in this format: "2014-10-30T06:04:39.113+0000"
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)
	newTime, err := time.Parse("2006-01-02T15:04:05.999-0700", strInput)
	if err != nil {
		//fmt.Printf(" - ERROR PARSING ARITIME: %s - ", err)
		return fmt.Errorf("Error parsing AriTime: %s", err)
	}
	*j = AriTime(newTime)
	return nil
}
