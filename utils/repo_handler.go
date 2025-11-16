package utils

import (
	"archive/zip"
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

func RepoToNew(repoPath string, outputPath string, keepGit bool, packZip bool) string {
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

	var zipWriter *zip.Writer
	var zipFile *os.File

	if packZip {
		zipFilePath := filepath.Join(filepath.Dir(outputPath), repoName+".zip")
		zipFile, err = os.Create(zipFilePath)
		if err != nil {
			return err.Error()
		}
		defer zipFile.Close()
		zipWriter = zip.NewWriter(zipFile)
		defer zipWriter.Close()
	} else {
		if err := os.MkdirAll(outputPath, 0755); err != nil {
			return err.Error()
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

		// 跳过符号链接
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
		if ign != nil && ign.MatchesPath(relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if packZip {
			if info.IsDir() {
				// zip 中的目录条目
				zipDir := repoName + "/" + relPath + "/"
				zipDir = strings.ReplaceAll(zipDir, "\\", "/")
				_, err := zipWriter.Create(zipDir)
				return err
			}

			// 普通文件
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			zipPath := repoName + "/" + relPath
			zipPath = strings.ReplaceAll(zipPath, "\\", "/")
			header.Name = zipPath
			header.Method = zip.Deflate

			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}

			_, err = io.Copy(writer, file)
			return err
		}

		// --- 普通复制模式 ---
		destPath := filepath.Join(outputPath, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return copyFile(path, destPath, info.Mode())
	})

	if err != nil {
		return err.Error()
	}

	return "Ok"
}
