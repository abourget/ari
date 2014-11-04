package ari

import (
	"fmt"
	"strings"
	"time"
)

type Time time.Time

func (j *Time) UnmarshalJSON(input []byte) error {
	// ARI stamps in this format: "2014-10-30T06:04:39.113+0000"
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)
	newTime, err := time.Parse("2006-01-02T15:04:05.999-0700", strInput)
	if err != nil {
		//fmt.Printf(" - ERROR PARSING ARITIME: %s - ", err)
		return fmt.Errorf("Error parsing Time: %s", err)
	}
	*j = Time(newTime)
	return nil
}


// FIXME: This doesn't work to improve "pretty.Formatter"
func (j *Time) MarshalText() ([]byte, error) {
	t := time.Time(*j)
	return []byte(t.Format("2006-01-02T15:04:05.999-0700")), nil
}
