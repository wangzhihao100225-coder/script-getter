package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	GitHubUser = "wangzhihao100225-coder" 
	GitHubRepo = "script-getter" 
)

const UA = "Surge/3035 CFNetwork/1410.0.3 Darwin/22.6.0"

type SyncTask struct {
	ConfURL  string
	JSURL    string
	ConfFile string
	JSFile   string
}

func fetchWithUA(url string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("URL 不能为空")
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", UA)
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
	tasks := []SyncTask{
		{
			ConfURL:  "https://ddgksf2013.top/rewrite/BiliBiliAds.conf",
			JSURL:    "https://ddgksf2013.top/scripts/bilibili.ads.js",
			ConfFile: "BiliBiliAds.conf",
			JSFile:   "bilibili.ads.js",
		},
		{
			// 百度网盘：配置和脚本是同一个 URL
			ConfURL:  "https://ddgksf2013.top/scripts/bdpan.ads.js",
			JSURL:    "https://ddgksf2013.top/scripts/bdpan.ads.js",
			ConfFile: "BaiduPanAds.conf", // 存为 .conf 方便 Stash 识别
			JSFile:   "bdpan.ads.js",
		},
		{
			// 知乎：配置和脚本是同一个 URL
			ConfURL:  "https://ddgksf2013.top/scripts/zhihu.ads.js",
			JSURL:    "https://ddgksf2013.top/scripts/zhihu.ads.js",
			ConfFile: "ZhihuAds.conf", // 存为 .conf 方便 Stash 识别
			JSFile:   "zhihu.ads.js",
		},
		{
			ConfURL:  "https://ddgksf2013.top/rewrite/XiaoHongShuAds.conf",
			JSURL:    "https://ddgksf2013.top/scripts/redbook.ads.js",
			ConfFile: "XiaoHongShuAds.conf",
			JSFile:   "redbook.ads.js",
		},
	}

	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(t SyncTask) {
			defer wg.Done()
			
			// 1. 下载脚本源文件
			jsContent, err := fetchWithUA(t.JSURL)
			if err != nil {
				fmt.Printf("❌ %s (脚本) 下载失败: %v\n", t.JSFile, err)
				return
			}
			os.WriteFile(t.JSFile, []byte(jsContent), 0644)

			// 2. 下载配置源文件（虽然对于百度/知乎来说内容一样）
			confContent, err := fetchWithUA(t.ConfURL)
			if err != nil {
				fmt.Printf("❌ %s (配置) 下载失败: %v\n", t.ConfFile, err)
				return
			}

			// 3. 替换逻辑：把配置里指向 ddgksf2013 的链接换成你 GitHub 的 Raw 链接
			myJSURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/%s", GitHubUser, GitHubRepo, t.JSFile)
			newConfContent := strings.ReplaceAll(confContent, t.JSURL, myJSURL)

			// 4. 保存为 .conf 文件供 Stash 订阅
			os.WriteFile(t.ConfFile, []byte(newConfContent), 0644)
			fmt.Printf("✅ %s 同步并替换完成\n", t.ConfFile)
		}(task)
	}
	wg.Wait()
	fmt.Println("🎉 完美！所有脚本和“披着JS外衣的配置”均已清洗完成！")
}
