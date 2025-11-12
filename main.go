package main

import (
	"gitpack-core/utils"

	"C"
)
import "fmt"

//export RepoToNew
func RepoToNew(repoPath *C.char, outputPath *C.char, keepGit C.int) *C.char {
	return C.CString(utils.RepoToNew(C.GoString(repoPath), C.GoString(outputPath), int(keepGit) != 0, false))
}

//export RepoToZip
func RepoToZip(repoPath *C.char, outputPath *C.char, keepGit C.int) *C.char {
	return C.CString(utils.RepoToNew(C.GoString(repoPath), C.GoString(outputPath), int(keepGit) != 0, true))
}

func main() {
	fmt.Println(utils.RepoToNew("/Users/zhoucheng/Desktop/Develop/DAV-Server", "/Users/zhoucheng/Downloads", false, false))
}
