package main

import (
	"gitpack-core/utils"

	"C"
)

//export RepoToNew
func RepoToNew(repoPath *C.char, outputPath *C.char, keepGit C.int, keepIgnore C.int) C.int {
	err := utils.RepoToNew(C.GoString(repoPath), C.GoString(outputPath), int(keepGit) != 0, int(keepIgnore) != 0)

	if err != nil {
		return 1
	}
	return 0
}

func main() {}
