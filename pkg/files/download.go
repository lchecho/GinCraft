package files

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	netUrl "net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type DownloadResult struct {
	Origin   string
	FilePath string
	Err      error
}

func DownloadFiles(dir string, urls ...string) ([]DownloadResult, error) {
	if len(urls) == 0 {
		return nil, nil
	}

	// 创建保存目录
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %v", err)
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	result := make([]DownloadResult, len(urls))
	// 创建信号量控制并发数
	semaphore := make(chan struct{}, 10)
	var wg sync.WaitGroup

	for i, url := range urls {
		wg.Add(1)
		go func(index int, url string) {
			defer wg.Done()
			// 获取信号量
			semaphore <- struct{}{}
			defer func() {
				<-semaphore
			}()

			r := DownloadResult{
				Origin: url,
			}
			r.FilePath, r.Err = downloadFile(url, dir, client)
			result[index] = r
		}(i, url)
	}
	wg.Wait()

	return result, nil
}

func downloadFile(url string, dir string, client *http.Client) (string, error) {
	// 发送 GET 请求
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("下载文件, 请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载文件, 失败状态码: %d", resp.StatusCode)
	}

	// 解析完整的 URL 字符串
	u, err := netUrl.Parse(url)
	if err != nil {
		return "", fmt.Errorf("下载文件, 解析URL失败: %v", err)
	}
	// 构建完整的保存路径
	var ext string
	fileName := filepath.Base(u.Path)
	ss := strings.Split(fileName, ".")
	if len(ss) > 1 {
		ext = ss[len(ss)-1]
	}
	filePath := fmt.Sprintf("%s/%s.%s", dir, uuid.New().String(), ext)

	// 创建本地文件
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("下载文件, 创建文件失败: %v", err)
	}
	defer file.Close()

	// 复制响应体到文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("下载文件, 保存文件失败: %v", err)
	}

	return filePath, nil
}
