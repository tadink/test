package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var blockes = []string{
	"baidu.com",
	"1688.com",
	"qcc.com",
	"gov.cn",
	"zhipin.com",
	"tianyancha.com",
	"sogou.com",
	"douyin.com",
	"sina.cn",
	"zhaopinku.cn",
	"google.com",
	"coovee.com",
	"job5156.com",
	"1637.com",
	"amap.com",
	"163.com",
	"cfi.cn",
	"cfi.net.cn",
	"amz123.com",
	"facebook.com",
	"sohu.com",
}
var cities = []string{
	"北京",
	"上海",
	"重庆",
	"广州",
	"深圳",
	"广东",
	"杭州",
	"浙江",
	"武汉",
	"湖北",
	"成都",
	"四川",
	"南京",
	"江苏",
	"沈阳",
	"辽宁",
	"长沙",
	"湖南",
	"石家庄",
	"河北",
	"郑州",
	"河南",
	"济南",
	"山东",
	"哈尔滨",
	"黑龙江",
	"长春",
	"吉林",
	"西安",
	"陕西",
	"福州",
	"福建",
	"合肥",
	"安徽",
	"南昌",
	"江西",
	"昆明",
	"云南",
	"呼和浩特",
	"内蒙古",
	"南宁",
	"广西",
	"太原",
	"山西",
	"乌鲁木齐",
	"新疆",
	"贵阳",
	"贵州",
	"兰州",
	"甘肃",
	"西宁",
	"青海",
	"海口",
	"海南",
	"银川",
	"宁夏",
	"拉萨",
	"西藏",
}
var keywords = []string{
	"机械有限公司",
	"五金有限公司",
	"电气有限公司",
	"机床有限公司",
	"电子有限公司",
	"灯饰有限公司",
	"环保有限公司",
	"印刷有限公司",
	"工程机械有限公司",
	"汽配有限公司",
	"泵阀有限公司",
	"纸业有限公司",
	"仪器仪表有限公司",
	"安防有限公司",
	"汽车有限公司",
	"汽修有限公司",
	"通信有限公司",
	"过滤有限公司",
	"消防有限公司",
	"加工有限公司",
	"汽车用品有限公司",
	"丝印特印有限公司",
	"包装有限公司",
	"LED有限公司",
	"水工业有限公司",
	"二手设备有限公司",
	"广电有限公司",
	"耐火材料有限公司",
	"焊接切割有限公司",
	"暖通空调有限公司",
	"添加剂有限公司",
	"纸管有限公司",
	"太阳能有限公司",
	"物流设备有限公司",
	"热泵有限公司",
	"工控有限公司",
	"信息有限公司",
	"紧固件有限公司",
	"IT有限公司",
	"家电有限公司",
	"礼品有限公司",
	"家居有限公司",
	"运动休闲有限公司",
	"家具有限公司",
	"办公有限公司",
	"酒店有限公司",
	"美容美发有限公司",
	"教育装备有限公司",
	"服装有限公司",
	"服饰有限公司",
	"古玩有限公司",
	"玩具有限公司",
	"制鞋有限公司",
	"音响灯光有限公司",
	"小家电有限公司",
	"皮具有限公司",
	"零食有限公司",
	"智能有限公司",
	"二手有限公司",
	"珠宝有限公司",
	"影音有限公司",
	"宠物有限公司",
	"母婴有限公司",
	"消毒产品有限公司",
	"商务服务有限公司",
	"生活服务有限公司",
	"广告有限公司",
	"教育有限公司",
	"物流有限公司",
	"交通运输有限公司",
	"网站建设有限公司",
	"展会有限公司",
	"维修有限公司",
	"项目有限公司",
	"创业有限公司",
	"船舶有限公司",
	"二手回收有限公司",
	"翻译有限公司",
	"建材有限公司",
	"能源有限公司",
	"冶金有限公司",
	"纺织有限公司",
	"化工有限公司",
	"表面处理有限公司",
	"房地产有限公司",
	"超硬材料有限公司",
	"塑料有限公司",
	"钢铁有限公司",
	"橡胶有限公司",
	"皮革有限公司",
	"丝网有限公司",
	"涂料有限公司",
	"石材有限公司",
	"石油有限公司",
	"卫浴有限公司",
	"陶瓷有限公司",
	"玻璃有限公司",
	"养殖有限公司",
	"水果批发有限公司",
	"食品有限公司",
	"食品机械有限公司",
	"农机有限公司",
	"园林有限公司",
	"农化有限公司",
	"饲料有限公司",
	"茶叶有限公司",
	"种子有限公司",
	"蔬菜有限公司",
	"光电有限公司",
}
var visitedUrl = make(map[string]bool)
var resultFile *os.File

func main() {

	var err error
	resultFile, err = os.OpenFile("result.json", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	data, _ := os.ReadFile("cookies.txt")
	var cookies []*http.Cookie
	if len(data) > 0 {
		cs := strings.Split(string(data), ";")
		for _, c := range cs {
			cook := strings.SplitN(c, "=", 2)
			cookies = append(cookies, &http.Cookie{
				Name:   cook[0],
				Value:  cook[1],
				Domain: ".google.com",
			})

		}
	}

	http.DefaultClient.Jar, _ = cookiejar.New(nil)
	u, _ := url.Parse("https://www.google.com")
	http.DefaultClient.Jar.SetCookies(u, cookies)
	for _, city := range cities {
		for _, keyword := range keywords {
			for i := 0; i < 30; i++ {
				err := googleSearch(city+keyword, i*10)
				if err != nil {
					fmt.Println(err.Error())
					time.Sleep(time.Second)
					break
				}
			}
			os.WriteFile("已完成关键词.txt", []byte(keyword), os.ModePerm)
			fmt.Println(city + keyword)
		}
	}

}
func googleSearch(keyword string, start int) error {
	requestUrl := fmt.Sprintf("https://www.google.com/search?q=%s&start=%d", keyword, start)
	request, _ := http.NewRequest("GET", requestUrl, nil)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	if len(doc.Find("a[data-ved]").Nodes) < 1 {
		ret, _ := doc.Html()
		os.WriteFile(keyword+".html", []byte(ret), os.ModePerm)
		return fmt.Errorf("未搜索到结果:%s,%d", keyword, start)
	}
	doc.Find("a[data-ved]").Each(func(i int, s *goquery.Selection) {
		href := s.AttrOr("href", "")
		if strings.Index(href, "/url") != 0 {
			return
		}
		u, err := url.Parse(href)
		if err != nil {
			return
		}
		q := u.Query().Get("q")
		u, err = url.Parse(q)
		if err != nil {
			return
		}
		if u.Scheme == "" || u.Host == "" {
			return
		}
		if isBlock(u.Host) {
			return
		}
		err = fetchCompanySite(u.Scheme + "://" + u.Host)
		if err != nil {
			fmt.Println(err.Error())
		}
	})
	return nil
}
func fetchCompanySite(requestUrl string) error {
	if visitedUrl[requestUrl] {
		return nil
	}
	defer func() { visitedUrl[requestUrl] = true }()
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	replacer := strings.NewReplacer(" ", "", "\n", "", " ", "", "\r", "")
	title := replacer.Replace(doc.Find("title").Text())
	_, err = resultFile.WriteString(fmt.Sprintf("%s,%s\n", requestUrl, title))
	if err != nil {
		return err
	}
	return nil

}
func isBlock(host string) bool {
	for _, block := range blockes {
		if strings.Contains(host, block) {
			return true
		}
	}
	return false
}
