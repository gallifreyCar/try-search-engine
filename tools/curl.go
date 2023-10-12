package tools

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gallifreyCar/try-search-engine/model"
	"github.com/imroc/req/v3"
	"strings"
	"time"
)

const (
	Success int = iota
	NetworkError
	HtmlError
	OthersError
)

// 重用client，4秒超时，不跟随重定向
var client = req.C().SetTimeout(time.Second * 4).SetRedirectPolicy(req.NoRedirectPolicy())

func Curl(status model.Status) (doc *goquery.Document, res int) {

	resp, err := client.R().SetHeader("User-Agent", "Sogou web model/4.0(+http://www.sogou.com/docs/help/webmasters.htm#07)").Get(status.Url)
	if err != nil {
		fmt.Println("访问失败", status.Url, err)
		document, _ := goquery.NewDocumentFromReader(strings.NewReader(""))
		return document, NetworkError
	}
	html := resp.String()
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println("HTML解析失败", status.Url, err)
		document, _ := goquery.NewDocumentFromReader(strings.NewReader(""))
		return document, HtmlError
	}
	return doc, Success

}
