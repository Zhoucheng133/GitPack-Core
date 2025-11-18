package utils

import (
	"io"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	ignore "github.com/sabhiram/go-gitignore"
)

func RepoCheck(repoPath string) int {
	_, err := git.PlainOpen(repoPath)
	if err != nil {
		return 0
	}
	gitignorePath := filepath.Join(repoPath, ".gitignore")
	if _, err := os.Stat(gitignorePath); err == nil {
		_, err = ignore.CompileIgnoreFile(gitignorePath)
		if err != nil {
			return 0
		}
	}
	return 1
}

func copyFile(src, dst string, perm os.FileMode) error {
	// 打开源文件
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// 创建目标目录
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	// 创建目标文件
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer out.Close()

	// 拷贝内容
	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	// 强制写入
	return out.Sync()
}

func RepoToNew(repoPath string, outputPath string, keepGit bool, packZip bool) string {

	repoName := filepath.Base(repoPath)

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err.Error()
	}

	// 获取 HEAD
	headRef, err := repo.Head()
	if err != nil {
		return err.Error()
	}

	commit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return err.Error()
	}

	tree, err := commit.Tree()
	if err != nil {
		return err.Error()
	}

	// 遍历 Git 记录的所有文件（等价于 git ls-files）
	err = tree.Files().ForEach(func(f *object.File) error {
		src := filepath.Join(repoPath, f.Name)
		dst := filepath.Join(outputPath, f.Name)

		info, err := os.Stat(src)
		if err != nil {
			return err
		}

		return copyFile(src, dst, info.Mode())
	})

	if err != nil {
		return err.Error()
	}

	return "Ok"
}
