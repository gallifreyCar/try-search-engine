package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gallifreyCar/try-search-engine/database"
	"github.com/gallifreyCar/try-search-engine/model"
	"github.com/gallifreyCar/try-search-engine/tools"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func main() {
	//初始化配置
	err := initEnv()
	if err != nil {
		fmt.Println("初始化配置失败", err)
		return
	}

	//开始爬取数据
	nextStep(time.Now())

}

// 循环爬取数据
func nextStep(startTime time.Time) {
	//访问数据库，获取要爬取的url

	// 从数据库中获取未爬取的100条URL
	var status []model.Status
	database.DbOne.Where("craw_done = ?", 0).Limit(10).Find(&status)
	fmt.Println("--------------------------------------------------")
	fmt.Println("开始爬取", len(status), "条URL")

	// 开始爬取
	chs := make([]chan int, len(status))
	for i, v := range status {
		chs[i] = make(chan int)
		go craw(chs[i], v)
	}

	// 等待爬取完成,统计结果
	// 1:success 成功
	// 2:networkError 网络错误
	// 3:htmlError HTML解析错误
	// 4:OthersError 其他错误
	var result = make(map[int]int)
	for _, ch := range chs {
		res := <-ch
		if _, ok := result[res]; ok {
			result[res]++
		} else {
			result[res] = 1
		}
	}
	fmt.Println("成功", result[0], "条")
	fmt.Println("网络错误", result[1], "条")
	fmt.Println("HTML解析错误", result[2], "条")
	fmt.Println("其他错误", result[3], "条")
	fmt.Println("跑完一轮", time.Now().Unix()-startTime.Unix(), "秒")

	fmt.Println("--------------------------------------------------")

	//跑下一轮
	//time.Sleep(30 * time.Second)
	//nextStep(time.Now())
}

// 爬取数据
func craw(ch chan int, status model.Status) {
	// 调用Curl函数去对URL发起请求，获取响应内容
	doc, res := tools.Curl(status)
	// 失败则直接返回
	if res != 0 {
		fmt.Println("爬取失败", status.Url, res)
		ch <- res
		return
	}
	// 成功先更新Status表
	status.CrawDone = 1
	status.CrawTime = time.Now()
	database.DbOne.Save(&status)

	// 再更新Page表
	var page model.Page

	database.DbOne.Where(model.Page{ID: status.ID}).FirstOrCreate(&page)
	page.Url = status.Url
	page.Host = status.Host
	page.Title = strings.TrimSpace(doc.Find("title").Text())
	page.Text = handleText(strings.TrimSpace(doc.Text()))
	page.CrawDone = status.CrawDone
	page.CrawTime = status.CrawTime
	database.DbOne.Save(&page)
	//if r.RowsAffected == 0 {
	//	fmt.Println("更新Page表失败", status.ID)
	//	ch <- tools.OthersError
	//	return
	//}

	//解析doc中的超链接
	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		// 获取href属性
		getUrl, exists := selection.Attr("href")
		if !exists {
			return
		}
		// 过滤掉非法的url
		parse, err := url.Parse(getUrl)
		if err != nil || parse.Scheme == "" || parse.Host == "" {

			return
		}
		// 过滤掉非http和https的url
		if parse.Scheme != "http" && parse.Scheme != "https" {
			return
		}
		// 保存到数据库Status表
		status = model.Status{}
		status.Url = getUrl
		status.Host = parse.Host
		database.DbOne.Save(&status)

		// 保存到数据库Page表
		page = model.Page{}
		page.Url = getUrl
		page.Host = parse.Host
		page.Title = strings.TrimSpace(selection.Text())
		page.CrawDone = status.CrawDone
		page.Path = parse.Path
		page.Query = parse.RawQuery
		page.Scheme = parse.Scheme
		database.DbOne.Save(&page)

		database.DbOne.Where(model.Page{Url: getUrl}).FirstOrCreate(&page)

	})

	// 返回成功
	ch <- 0
}

// 初始化变量
func initEnv() error {
	err := database.InitDB()
	if err != nil {
		return err
	}
	return nil
}

func handleText(input string) string {
	if input == "" {
		return ""
	}
	//正则替换
	//\s 表示任何空白字符，包括空格、制表符、换行符等。
	//\p{Zs} 表示 Unicode 中的分隔符空格。
	//{1,} 表示匹配1个或更多次。
	reg := regexp.MustCompile(`[\s\p{Zs}]{1,}`)
	return reg.ReplaceAllString(input, "-")
}
