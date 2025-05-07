package tar

import (
	"time"

	"github.com/spf13/cobra"

	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	// flag
	path string

	allCmd = &cobra.Command{
		Use: "all",
		// Args: ,
		Run: func(cmd *cobra.Command, args []string) {
			CompressPath(path)
		},
	}
)

func init() {
	allCmd.Flags().StringVarP(&path, "path", "p", ".", "Directory to be archive")
}

func CompressPath(srcDir string) error {
	destFile := fmt.Sprintf("%s.tar.gz", time.Now().AddDate(0, 0, -1).Format(time.DateOnly))

	file, err := os.Create(destFile)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer file.Close()

	// 创建 gzip writer
	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	// 创建 tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// 遍历源目录中的所有文件和子目录
	filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("遍历路径 '%s' 失败: %w", path, err)
		}

		// 跳过源目录本身
		if path == srcDir {
			return nil
		}

		// 获取文件 header
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return fmt.Errorf("创建文件 header '%s' 失败: %w", path, err)
		}

		// 修改 header 中的 Name，使其相对于源目录
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return fmt.Errorf("获取相对路径 '%s' 失败: %w", path, err)
		}
		header.Name = relPath

		// 写入 header
		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("写入 header '%s' 失败: %w", path, err)
		}

		// 如果是普通文件，则写入文件内容
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("打开文件 '%s' 失败: %w", path, err)
			}
			defer file.Close()

			if _, err := io.Copy(tarWriter, file); err != nil {
				return fmt.Errorf("写入文件 '%s' 内容失败: %w", path, err)
			}
		}

		return nil
	})

	return nil
}
