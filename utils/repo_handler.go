package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
	ignore "github.com/sabhiram/go-gitignore"
)

func copyFile(src, dst string, perm os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func RepoToNew(repoPath string, outputPath string, keepGit bool, keepIgnore bool, packZip bool) string {
	_, err := git.PlainOpen(repoPath)
	if err != nil {
		return err.Error()
	}

	repoName := filepath.Base(repoPath)
	outputPath = filepath.Join(outputPath, repoName)

	var ign *ignore.GitIgnore
	gitignorePath := filepath.Join(repoPath, ".gitignore")
	if _, err := os.Stat(gitignorePath); err == nil {
		ign, err = ignore.CompileIgnoreFile(gitignorePath)
		if err != nil {
			return err.Error()
		}
	}

	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return err.Error()
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
		if strings.HasPrefix(relPath, ".git") && !keepGit {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 跳过 .gitignore 忽略的文件
		if ign != nil && ign.MatchesPath(relPath) && !keepIgnore {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 复制文件
		destPath := filepath.Join(outputPath, relPath)

		if info.IsDir() {
			// 创建对应的目标目录
			return os.MkdirAll(destPath, info.Mode())
		}

		// 复制文件
		if err := copyFile(path, destPath, info.Mode()); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err.Error()
	}

	return "Ok"
}
