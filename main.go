package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

// ======= 必须修改以下两行 =======
const (
	GitHubUser = "wangzhihao100225-coder" 
	GitHubRepo = "script-getter"     
)

const UA = "Surge/3035 CFNetwork/1410.0.3 Darwin/22.6.0"

func fetchWithUA(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", UA)
	req.Header.Set("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

func main() {
	var wg sync.WaitGroup
	var jsContent, confContent string
	var jsErr, confErr error

	jsURL := "https://ddgksf2013.top/scripts/bilibili.ads.js"
	confURL := "https://ddgksf2013.top/rewrite/BiliBiliAds.conf"

	wg.Add(2)
	go func() {
		defer wg.Done()
		jsContent, jsErr = fetchWithUA(jsURL)
	}()
	go func() {
		defer wg.Done()
		confContent, confErr = fetchWithUA(confURL)
	}()
	wg.Wait()

	if jsErr != nil || confErr != nil {
		fmt.Printf("下载失败: JS(%v), Conf(%v)\n", jsErr, confErr)
		return
	}

	// 1. 保存原始脚本
	os.WriteFile("bilibili.ads.js", []byte(jsContent), 0644)

	// 2. 处理配置文件：把墨鱼的 JS 链接换成你自己的 GitHub Raw 链接
	myJSURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/bilibili.ads.js", GitHubUser, GitHubRepo)
	newConfContent := strings.ReplaceAll(confContent, jsURL, myJSURL)

	// 3. 保存修改后的配置文件
	os.WriteFile("BiliBiliAds.conf", []byte(newConfContent), 0644)

	fmt.Println("✅ 脚本与配置已同步并完成链接替换")
}