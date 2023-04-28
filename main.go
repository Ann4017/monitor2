package main

import "fmt"

func main() {
	m := Monitor{}

	err := m.Init("config.ini", "database")
	if err != nil {
		fmt.Println(err)
	}

	err = m.Run()
	if err != nil {
		fmt.Println(err)
	}
}
