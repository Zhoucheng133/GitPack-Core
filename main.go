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

//export  RepoCheck
func RepoCheck(repoPath *C.char) C.int {
	return C.int(utils.RepoCheck(C.GoString(repoPath)))
}

func main() {
	fmt.Println(utils.RepoCheck("/Users/zhoucheng/Desktop/Develop"))
}
