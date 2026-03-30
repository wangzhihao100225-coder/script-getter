package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

// ======= 务必修改这两行 =====
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
			ConfURL:  "https://ddgksf2013.top/rewrite/BaiduPanAds.conf",
			JSURL:    "https://ddgksf2013.top/scripts/bdpan.ads.js",
			ConfFile: "BaiduPanAds.conf",
			JSFile:   "bdpan.ads.js",
		},
		{
			ConfURL:  "https://ddgksf2013.top/rewrite/ZhihuAds.conf",
			JSURL:    "https://ddgksf2013.top/scripts/zhihu.ads.js",
			ConfFile: "ZhihuAds.conf",
			JSFile:   "zhihu.ads.js",
		},
		{
			// 🆕 新增：小红书
			ConfURL:  "https://ddgksf2013.top/rewrite/XiaoHongShuAds.conf",
			JSURL:    "https://ddgksf2013.top/scripts/xiaohongshu.ads.js",
			ConfFile: "XiaoHongShuAds.conf",
			JSFile:   "xiaohongshu.ads.js",
		},
	}

	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(t SyncTask) {
			defer wg.Done()
			
			// 下载 JS
			jsContent, err := fetchWithUA(t.JSURL)
			if err != nil {
				fmt.Printf("❌ %s 下载失败: %v\n", t.JSFile, err)
				return
			}
			os.WriteFile(t.JSFile, []byte(jsContent), 0644)

			// 下载 Conf
			confContent, err := fetchWithUA(t.ConfURL)
			if err != nil {
				fmt.Printf("❌ %s 下载失败: %v\n", t.ConfFile, err)
				return
			}

			// 链接替换逻辑
			myJSURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/%s", GitHubUser, GitHubRepo, t.JSFile)
			newConfContent := strings.ReplaceAll(confContent, t.JSURL, myJSURL)

			os.WriteFile(t.ConfFile, []byte(newConfContent), 0644)
			fmt.Printf("✅ %s 已就绪\n", t.ConfFile)
		}(task)
	}
	wg.Wait()
	fmt.Println("🎉 全家桶同步任务全部完成！")
}
