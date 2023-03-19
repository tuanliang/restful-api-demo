package main

import (
	"fmt"

	"github.com/tuanliang/restful-api-demo/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
