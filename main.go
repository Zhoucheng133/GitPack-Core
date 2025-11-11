package main

import (
	"fmt"
	"gitpack-core/utils"

	"C"
)

//export RepoToNew
func RepoToNew(repoPath *C.char, outputPath *C.char) C.int {
	err := utils.RepoToNew(C.GoString(repoPath), C.GoString(outputPath))

	if err != nil {
		return 1
	}
	return 0
}

func main() {
	err := utils.RepoToNew("/Users/zhoucheng/Desktop/Develop/DAV-Server", "/Users/zhoucheng/Downloads")
	if err != nil {
		fmt.Println(err.Error())
	}
}
