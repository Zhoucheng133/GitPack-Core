package utils

import (
	"os"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
	ignore "github.com/sabhiram/go-gitignore"
)

func RepoToNew(repoPath string, outputPath string) error {
	_, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	var ign *ignore.GitIgnore
	gitignorePath := filepath.Join(repoPath, ".gitignore")

	if _, err := os.Stat(gitignorePath); err == nil {
		ign, err = ignore.CompileIgnoreFile(gitignorePath)
		if err != nil {
			return err
		}
	}

	err = filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(repoPath, path)
		if err != nil {
			return err
		}

		// 跳过根目录
		if relPath == "." {
			return nil
		}

		// 跳过引用文件
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		// 跳过 .git 目录
		if strings.HasPrefix(relPath, ".git") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 跳过 .gitignore 忽略的文件
		if ign != nil && ign.MatchesPath(relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// TODO 复制

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
