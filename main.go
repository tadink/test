package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TypeModel struct {
	Id       int    `json:"id" gorm:"column:id"`
	ParentId int    `json:"parent_id" gorm:"column:parent_id"`
	Name     string `json:"name" gorm:"column:name"`
	Url      string `json:"url" gorm:"-"`
}

func (t *TypeModel) TableName() string {
	return "vod_type"
}

type Director struct {
	Id   int    `json:"id" gorm:"column:id"`
	Name string `json:"name" gorm:"column:name"`
}

func (d *Director) TableName() string {
	return "director"
}

type Actor struct {
	Id   int    `json:"id" gorm:"column:id"`
	Name string `json:"name" gorm:"column:name"`
}

func (a *Actor) TableName() string {
	return "actor"
}

type Tag struct {
	Id   int    `json:"id" gorm:"column:id"`
	Name string `json:"name" gorm:"column:name"`
}

func (t *Tag) TableName() string {
	return "tag"
}

type Vod struct {
	Id          int          `json:"id"`
	Name        string       `json:"name"`
	Image       string       `json:"image"`
	Remark      string       `json:"remark"`
	Rate        float64      `json:"rate"`
	Alias       string       `json:"alias"`
	Directors   []Director   `json:"directors" gorm:"many2many:director_vod;foreignKey:id;joinForeignKey:vod_id;References:id;joinReferences:director_id"`
	Actors      []Actor      `json:"actors" gorm:"many2many:actor_vod;foreignKey:id;joinForeignKey:vod_id;References:id;joinReferences:actor_id"`
	Tags        []Tag        `json:"tags" gorm:"many2many:tag_vod;foreignKey:id;joinForeignKey:vod_id;References:id;joinReferences:tag_id"`
	Area        string       `json:"area"`
	Language    string       `json:"language"`
	PublishTime int          `json:"publish_time"`
	UpdateTime  string       `json:"update_time"`
	Description string       `json:"description"`
	SourceUrl   string       `json:"source_url"`
	VodPlayUrls []VodPlayUrl `json:"vod_play_urls" gorm:"foreignKey:vod_id"`
	TypeId      int          `json:"type_id"`
}

func (t *Vod) TableName() string {
	return "vod"
}

type VodPlayUrl struct {
	Id    int
	Name  string `json:"name" gorm:"column:name"`
	Url   string `json:"url" gorm:"column:url"`
	VodId int    `json:"vod_id" gorm:"column:vod_id"`
}

func (t *VodPlayUrl) TableName() string {
	return "vod_play_url"
}

var typeList = []TypeModel{
	{20, 0, "电影", ""},
	{22, 20, "冒险片", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{24, 20, "剧情片", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{26, 20, "动作片", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{28, 20, "动画电影", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{32, 20, "喜剧片", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{34, 20, "奇幻片", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{36, 20, "恐怖片", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{38, 20, "悬疑片", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{40, 20, "惊悚片", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{42, 20, "歌舞片", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{44, 20, "灾难片", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{46, 20, "爱情片", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{48, 20, "科幻片", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{50, 20, "犯罪片", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{52, 20, "经典片", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{54, 20, "网络电影", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{56, 20, "战争片", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{58, 20, "伦理片", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{60, 0, "电视剧", ""},
	{62, 60, "欧美剧", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{64, 60, "日剧", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{66, 60, "韩剧", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{68, 60, "台剧", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{70, 60, "泰剧", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{72, 60, "国产剧", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{74, 60, "港剧", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{76, 60, "新马剧", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{78, 60, "其他剧", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{80, 0, "动漫", ""},
	{96, 80, "欧美动漫", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{98, 80, "日韩动漫", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{100, 80, "国产动漫", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{102, 80, "新马泰动漫", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{104, 80, "港台动漫", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{106, 80, "其他动漫", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{82, 0, "综艺", ""},
	{108, 82, "国产综艺", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{110, 82, "日韩综艺", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{112, 82, "欧美综艺", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{114, 82, "新马泰综艺", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{116, 82, "港台综艺", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{118, 82, "其他综艺", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{84, 0, "体育", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{86, 0, "纪录片", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
	{120, 0, "短剧", "https://hw8.live/index.php/vod/type/id/%d/page/%d.html"},
}
var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(mysql.Open("root:root@tcp(127.0.1:3306)/vod"))
	if err != nil {
		log.Fatal(err.Error())
	}
	sqlDb, err := db.DB()
	if err != nil {
		log.Fatal(err.Error())
	}
	sqlDb.SetConnMaxIdleTime(time.Hour)
	sqlDb.SetConnMaxLifetime(24 * time.Hour)
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(20)
}

func main() {

	var itemChan = make(chan TypeModel, 5)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			for item := range itemChan {
				queryType(&item)
			}
			wg.Done()
		}()
	}
	for _, item := range typeList {
		if item.Url == "" {
			continue
		}
		itemChan <- item

	}
	wg.Wait()
}

func queryType(model *TypeModel) {
	var page = 1
	for {
		detailUrls, err := queryTypePage(model.Url, model.Id, page)
		if err != nil {
			fmt.Println(err.Error(), model, page)
			break
		}
		for _, u := range detailUrls {
			v, err := queryDetail(u)
			if err != nil {
				fmt.Println(err.Error())
			}
			v.TypeId = model.Id
			err = AddVod(v)
			if err != nil {
				fmt.Println("add vod error", err.Error())
			}

		}

	}

}
func queryTypePage(url string, id int, page int) ([]string, error) {
	requestUrl := fmt.Sprintf(url, id, page)
	doc, err := getDocument(requestUrl)
	if err != nil {
		return nil, err
	}
	var detailUrls = make([]string, 0)
	doc.Find(".xing_vb4 a").Each(func(i int, selection *goquery.Selection) {
		href := selection.AttrOr("href", "")
		if href == "" {
			return
		}
		detailUrls = append(detailUrls, href)
	})
	return detailUrls, nil

}
func queryDetail(requestUrl string) (*Vod, error) {
	requestUrl = fmt.Sprintf("https://hw8.live%s", requestUrl)
	doc, err := getDocument(requestUrl)
	if err != nil {
		return nil, err
	}
	idRegexp, _ := regexp.Compile(`(\d+)\.html`)
	matches := idRegexp.FindStringSubmatch(requestUrl)
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, err
	}

	vodBox := doc.Find(".vodBox")
	image := vodBox.Find(".vodImg img").AttrOr("src", "")
	name := vodBox.Find(".vodInfo .vodh h2").Text()
	remark := vodBox.Find(".vodInfo .vodh span").Text()
	rate := vodBox.Find(".vodInfo .vodh label").Text()
	alias := vodBox.Find(".vodinfobox ul li:nth-child(1) span").Text()
	director := vodBox.Find(".vodinfobox ul li:nth-child(2) span").Text()
	actor := vodBox.Find(".vodinfobox ul li:nth-child(3) span").Text()
	tag := vodBox.Find(".vodinfobox ul li:nth-child(4) span").Text()
	area := vodBox.Find(".vodinfobox ul li:nth-child(5) span").Text()
	language := vodBox.Find(".vodinfobox ul li:nth-child(6) span").Text()
	publishTime := vodBox.Find(".vodinfobox ul li:nth-child(7) span").Text()
	updateTime := vodBox.Find(".vodinfobox ul li:nth-child(8) span").Text()
	description := doc.Find(".ibox .vodplayinfo").First().Text()
	var vodPlayUrls = make([]VodPlayUrl, 0)
	doc.Find(".vodplayinfo ul li span").Each(func(i int, selection *goquery.Selection) {
		if selection.Text() == "" {
			return
		}
		playUrl := strings.Split(selection.Text(), "$")
		if len(playUrl) != 2 {
			return
		}
		vodPlayUrls = append(vodPlayUrls, VodPlayUrl{
			Name: playUrl[0],
			Url:  playUrl[1],
		})

	})
	rateNum, _ := strconv.ParseFloat(rate, 64)
	publishTimeNum, _ := strconv.Atoi(publishTime)
	directors := strings.Split(director, ",")
	ds := make([]Director, 0)
	for _, d := range directors {
		ds = append(ds, Director{Name: d})
	}
	actors := strings.Split(actor, ",")
	as := make([]Actor, 0)
	for _, a := range actors {
		as = append(as, Actor{Name: a})
	}
	tags := strings.Split(tag, ",")
	ts := make([]Tag, 0)
	for _, t := range tags {
		ts = append(ts, Tag{Name: t})
	}
	v := &Vod{
		Id:          id,
		Name:        name,
		Image:       image,
		Remark:      remark,
		Rate:        rateNum,
		Alias:       alias,
		Directors:   ds,
		Actors:      as,
		Tags:        ts,
		Area:        area,
		Language:    language,
		PublishTime: publishTimeNum,
		UpdateTime:  updateTime,
		Description: description,
		SourceUrl:   requestUrl,
		VodPlayUrls: vodPlayUrls,
	}
	return v, nil
}

func getDocument(requestUrl string) (*goquery.Document, error) {
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, errors.Join(errors.New("request error"), err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.Join(errors.New("get response error"), err)
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, errors.Join(errors.New("parse doc error"), err)
	}
	return doc, nil
}
func AddVod(v *Vod) error {
	//if len(v.Tags) > 0 {
	//	for _, tag := range v.Tags {
	//		err := db.Model(&Tag{}).FirstOrCreate(tag).Error
	//		if err != nil {
	//			return err
	//		}
	//	}
	//
	//}
	//if len(v.Directors) > 0 {
	//	db.Model(&Director{}).FirstOrCreate(v.Directors)
	//}
	//db.Model()
	return db.Model(&Vod{}).FirstOrCreate(v).Error
}
