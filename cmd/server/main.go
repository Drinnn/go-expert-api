package main

import (
	"fmt"

	"github.com/Drinnn/go-expert-api/configs"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	fmt.Println(config)
}
