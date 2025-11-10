package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

func main() {
	root := "/Users/zhoucheng/Desktop/Develop/Anime-Helper" // 可以通过命令行参数传入
	output := "output.zip"

	if err := ZipWithoutGitignored(root, output); err != nil {
		fmt.Println("❌ Error:", err)
	} else {
		fmt.Println("✅ 打包完成:", output)
	}
}

func ZipWithoutGitignored(rootDir, zipFile string) error {
	// 1️⃣ 打开文件系统
	fs := osfs.New(rootDir)

	// 2️⃣ 读取所有 .gitignore 规则
	patterns, err := gitignore.ReadPatterns(fs, []string{""})
	if err != nil {
		return fmt.Errorf("读取 .gitignore 失败: %w", err)
	}
	matcher := gitignore.NewMatcher(patterns)

	// 3️⃣ 创建 zip 输出文件
	out, err := os.Create(zipFile)
	if err != nil {
		return fmt.Errorf("创建 zip 文件失败: %w", err)
	}
	defer out.Close()

	zipWriter := zip.NewWriter(out)
	defer zipWriter.Close()

	// 4️⃣ 遍历目录
	err = filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 相对路径（给 matcher 用）
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}
		if relPath == ".git" || strings.HasPrefix(relPath, ".git/") {
			return filepath.SkipDir
		}

		// 判断是否被忽略
		splitPath := strings.Split(relPath, string(filepath.Separator))
		if matcher.Match(splitPath, d.IsDir()) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil // 忽略此文件
		}

		// 只打包文件
		if !d.Type().IsRegular() {
			return nil
		}

		// 写入 zip
		return addFileToZip(zipWriter, path, relPath)
	})

	if err != nil {
		return fmt.Errorf("遍历目录失败: %w", err)
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, filePath, relPath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	w, err := zipWriter.Create(relPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, f)
	return err
}
