package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	str := struct {
		Time time.Time
	}{
		Time: time.Now(),
	}

	bytes, err := json.Marshal(str)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(bytes))
	fmt.Println(str.Time.Nanosecond())
}
