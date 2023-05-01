package main

import (
	"fmt"
	"time"
)

func main() {
	m := Monitor{}

	err := m.Init("config.ini", "database", "aws")
	if err != nil {
		fmt.Println(err)
	}

	err = m.Run(time.Second * 10)
	if err != nil {
		fmt.Println(err)
	}
}
