package main

import (
	"fmt"
	"gitpack-core/utils"
)

func main() {
	err := utils.RepoToNew("/Users/zhoucheng/Desktop/Develop/DAV-Server", "/Users/zhoucheng/Downloads")
	if err != nil {
		fmt.Println(err.Error())
	}
}
