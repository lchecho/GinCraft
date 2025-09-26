package files

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func SaveUploadFile(fileHeader *multipart.FileHeader, dir string) (string, error) {
	// 1. 获取原始文件名
	originalName := fileHeader.Filename

	// 2. 生成安全的文件名（防止路径穿越攻击）
	safeName := filepath.Base(originalName)
	safeName = strings.ReplaceAll(safeName, " ", "_") // 替换空格

	// 3. 构建完整的保存路径
	savePath := filepath.Join(dir, safeName)

	// 4. 确保上传目录存在
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", fmt.Errorf("创建上传目录失败: %v", err)
	}

	// 5. 检查文件是否已存在，如果存在则重命名
	counter := 1
	originalSavePath := savePath
	for {
		if _, err = os.Stat(savePath); os.IsNotExist(err) {
			break
		}
		// 文件已存在，生成新名称
		ext := filepath.Ext(originalSavePath)
		nameWithoutExt := strings.TrimSuffix(originalSavePath, ext)
		savePath = fmt.Sprintf("%s_(%d)%s", nameWithoutExt, counter, ext)
		counter++
	}

	// 6. 保存文件
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", fmt.Errorf("保存上传文件失败: %v", err)
	}

	return savePath, nil
}
