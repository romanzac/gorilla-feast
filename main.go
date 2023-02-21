package main

import (
	"fmt"
	"github.com/romanzac/gorilla-feast/cmd"
)

func main() {
	err := cmd.GorillaFeastCmd.Execute()
	if err != nil && err.Error() != "" {
		fmt.Println(err)
	}
}
