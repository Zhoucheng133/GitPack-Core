package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"

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
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return err.Error()
	}

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

	// 创建输出目录：outputPath/<repoName>
	repoName := filepath.Base(repoPath)
	targetRoot := filepath.Join(outputPath, repoName)
	if err := os.MkdirAll(targetRoot, 0755); err != nil {
		return err.Error()
	}

	// 遍历 Git 记录的所有文件
	err = tree.Files().ForEach(func(f *object.File) error {
		src := filepath.Join(repoPath, f.Name)
		dst := filepath.Join(targetRoot, f.Name)

		if !keepGit {
			rel, err := filepath.Rel(".git", f.Name)
			if err == nil && rel == "." || !strings.HasPrefix(rel, "..") {
				// 文件在 .git 目录下，跳过
				return nil
			}
		}

		info, err := os.Stat(src)
		if err != nil {
			return err
		}

		return copyFile(src, dst, info.Mode())
	})

	if err != nil {
		return err.Error()
	}

	// 如果 keepGit 为 true，需要单独复制 .git 文件夹
	if keepGit {
		gitSrc := filepath.Join(repoPath, ".git")
		gitDst := filepath.Join(targetRoot, ".git")
		err := copyDir(gitSrc, gitDst)
		if err != nil {
			return err.Error()
		}
	}

	return "Ok"
}

// 递归复制目录
func copyDir(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath, info.Mode())
	})
}
