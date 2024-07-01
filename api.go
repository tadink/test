package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

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

type Class struct {
	Id   int    `json:"id" gorm:"column:id"`
	Name string `json:"name" gorm:"column:name"`
}

func (t *Class) TableName() string {
	return "class"
}

type Vod struct {
	Id          int           `json:"vod_id" gorm:"column:id"`
	TypeId      int           `json:"type_id" gorm:"column:type_id"`
	Name        string        `json:"vod_name" gorm:"column:name"`
	EnName      string        `json:"vod_en" gorm:"column:en_name"`
	Sub         string        `json:"vod_sub" gorm:"column:sub"`
	Status      int           `json:"vod_status" gorm:"column:status"`
	State       string        `json:"vod_state" gorm:"column:state"`
	Letter      string        `json:"vod_letter" gorm:"column:letter"`
	Pic         string        `json:"vod_pic" gorm:"column:pic"`
	Remark      string        `json:"vod_remarks" gorm:"column:remark"`
	Score       string        `json:"vod_score" gorm:"score"`
	VodDirector string        `json:"vod_director" gorm:"-"`
	Directors   []*Director   `json:"directors" gorm:"many2many:director_vod;foreignKey:id;joinForeignKey:vod_id;References:id;joinReferences:director_id"`
	VodActor    string        `json:"vod_actor" gorm:"-"`
	Actors      []*Actor      `json:"actors" gorm:"many2many:actor_vod;foreignKey:id;joinForeignKey:vod_id;References:id;joinReferences:actor_id"`
	VodClass    string        `json:"vod_class" gorm:"-"`
	Classes     []*Class      `json:"classes" gorm:"many2many:class_vod;foreignKey:id;joinForeignKey:vod_id;References:id;joinReferences:class_id"`
	Area        string        `json:"vod_area" gorm:"column:area"`
	Language    string        `json:"vod_lang" gorm:"column:language"`
	Year        string        `json:"vod_year" gorm:"column:year"`
	VodTime     string        `json:"vod_time" gorm:"column:vod_time"`
	Content     string        `json:"vod_content" gorm:"column:content"`
	SourceName  string        `json:"source_name" gorm:"column:source_name"`
	VodPlayUrl  string        `json:"vod_play_url" gorm:"-"`
	VodPlayUrls []*VodPlayUrl `json:"vod_play_urls" gorm:"foreignKey:vod_id"`
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

type ActorVod struct {
	VodId   int
	ActorId int
}

func (a *ActorVod) TableName() string {
	return "actor_vod"
}

type DirectorVod struct {
	VodId      int
	DirectorId int
}

func (d *DirectorVod) TableName() string {
	return "director_vod"
}

type ClassVod struct {
	VodId   int
	ClassId int
}

func (c *ClassVod) TableName() string {
	return "class_vod"
}

type Result struct {
	Code int
	Msg  string
	List []*Vod
}

var db *gorm.DB

func init() {
	var err error
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)
	db, err = gorm.Open(mysql.Open("root:root@tcp(127.0.1:3306)/vod"), &gorm.Config{Logger: newLogger})

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

var (
	hour      int64
	startPage int
	isUpdate  bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "vod",
		Short: "全量抓取",
		Run: func(cmd *cobra.Command, args []string) {
			getAllVod()
		},
	}
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "更新数据",
		Run: func(cmd *cobra.Command, args []string) {
			updateVod()
		},
	}
	updateCmd.Flags().Int64VarP(&hour, "hours", "s", 0, "最近多少小时的数据")
	rootCmd.Flags().IntVarP(&startPage, "page", "p", 1, "起始页")
	rootCmd.AddCommand(updateCmd)
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func getAllVod() {
	run()
}
func updateVod() {
	isUpdate = true
	run()
}

func run() {
	var pageChan = make(chan int, 5)
	var wg = new(sync.WaitGroup)
	var once = new(sync.Once)
	var finishChan = make(chan struct{})
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go worker(pageChan, finishChan, once, wg)
	}

	page := startPage
	finish := false
	for {
		select {
		case <-finishChan:
			finish = true
		case pageChan <- page:
		}
		if finish {
			break
		}
		page++
	}
	close(pageChan)
	wg.Wait()

}
func worker(pageChan chan int, finishChan chan<- struct{}, once *sync.Once, wg *sync.WaitGroup) {
	for page := range pageChan {
		vods, err := queryPage(page)
		if err != nil {
			fmt.Println(err.Error())
			once.Do(func() {
				close(finishChan)
			})
		}
		for _, vod := range vods {
			parseVod(vod)
			fmt.Println(vod.Id)
			err = AddVod(vod)
			if err != nil {
				fmt.Println(err.Error())

			}
		}

	}
	wg.Done()
}
func queryPage(page int) ([]*Vod, error) {
	requestUrl := fmt.Sprintf("https://hw8.live/api.php/provide/vod/?ac=videolist&pg=%d", page)
	u, _ := url.Parse(requestUrl)
	if isUpdate {
		if hour == 0 {
			hour = queryLastTime()
		}
		q := u.Query()
		q.Add("h", fmt.Sprintf("%d", hour))
		u.RawQuery = q.Encode()
	}

	data, err := getResponseData(u.String())
	if err != nil {
		return nil, err
	}
	var result Result
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	if len(result.List) == 0 {
		return nil, errors.New("last page")
	}
	return result.List, nil

}
func parseVod(vod *Vod) {
	vod.SourceName = "hw8"
	if vod.Year == "" {
		vod.Year = "0"
	}
	if vod.Score == "" {
		vod.Score = "0.0"
	}
	classes := strings.Split(vod.VodClass, ",")
	for _, class := range classes {
		if class == "" {
			continue
		}
		vod.Classes = append(vod.Classes, &Class{Name: class})
	}
	directors := strings.Split(vod.VodDirector, ",")
	for _, director := range directors {
		if director == "" {
			continue
		}
		vod.Directors = append(vod.Directors, &Director{Name: director})
	}
	actors := strings.Split(vod.VodActor, ",")
	for _, actor := range actors {
		if actor == "" {
			continue
		}
		vod.Actors = append(vod.Actors, &Actor{Name: actor})
	}
	playUrls := strings.Split(vod.VodPlayUrl, "#")
	for _, playUrl := range playUrls {
		parts := strings.Split(playUrl, "$")
		if len(parts) != 2 {
			continue
		}
		vod.VodPlayUrls = append(vod.VodPlayUrls, &VodPlayUrl{
			Name: parts[0],
			Url:  parts[1],
		})
	}
}

func getResponseData(requestUrl string) ([]byte, error) {
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, errors.Join(errors.New("request error"), err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.Join(errors.New("get response error"), err)
	}
	return io.ReadAll(response.Body)
}
func AddVod(v *Vod) error {
	if len(v.Classes) > 0 {
		for _, class := range v.Classes {
			item := queryClass(class.Name)
			if item != nil {
				class.Id = item.Id
				continue
			}
			err := db.Model(&Class{}).Where("name=?", class.Name).FirstOrCreate(class).Error
			if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
				return err
			}
		}
	}
	if len(v.Directors) > 0 {
		for _, director := range v.Directors {
			item := queryDirector(director.Name)
			if item != nil {
				director.Id = item.Id
				continue
			}
			err := db.Model(&Director{}).Where("name=?", director.Name).FirstOrCreate(&director).Error
			if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
				return err
			}
		}

	}
	if len(v.Actors) > 0 {
		for _, actor := range v.Actors {
			item := queryActor(actor.Name)
			if item != nil {
				actor.Id = item.Id
				continue
			}
			err := db.Model(&Actor{}).Where("name=?", actor.Name).FirstOrCreate(&actor).Error
			if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
				return err
			}
		}

	}
	go func(requestUrl string) {
		err := downloadPic(requestUrl)
		if err != nil {
			fmt.Println(err.Error())
		}
	}(v.Pic)
	picUrl, _ := url.Parse(v.Pic)

	err := db.Model(&Vod{}).
		Omit("Classes", "Actors", "Directors", "VodPlayUrls").
		Assign(&Vod{
			Remark:   v.Remark,
			Score:    v.Score,
			State:    v.State,
			Status:   v.Status,
			EnName:   v.EnName,
			Sub:      v.Sub,
			Language: v.Language,
			Pic:      picUrl.Path,
		}).
		FirstOrCreate(v).Error
	if err != nil {
		return err
	}
	if len(v.VodPlayUrls) > 0 {
		err = deletePlayUrls(v.Id)
		if err != nil {
			return err
		}
		for _, playUrl := range v.VodPlayUrls {
			playUrl.VodId = v.Id
		}
		err = addPlayUrls(v.VodPlayUrls)
		if err != nil {
			return err
		}
	}
	for _, actor := range v.Actors {
		err := db.Model(&ActorVod{}).Where("actor_id=? and vod_id=?", actor.Id, v.Id).Create(&ActorVod{
			v.Id,
			actor.Id,
		}).Error
		if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
			return err
		}
	}
	for _, director := range v.Directors {
		err := db.Model(&DirectorVod{}).Where("director_id=? and vod_id=?", director.Id, v.Id).FirstOrCreate(&DirectorVod{DirectorId: director.Id, VodId: v.Id}).Error
		if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
			return err
		}
	}
	for _, class := range v.Classes {
		err := db.Model(&ClassVod{}).Where("class_id=? and vod_id=?", class.Id, v.Id).FirstOrCreate(&ClassVod{
			ClassId: class.Id,
			VodId:   v.Id,
		}).Error
		if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
			return err
		}
	}
	return nil
}
func queryClass(name string) *Class {
	var result = new(Class)
	err := db.Model(&Class{}).Where("name=?", name).First(result).Error
	if err != nil {
		return nil
	}
	return result
}
func queryDirector(name string) *Director {
	var result = new(Director)
	err := db.Model(&Director{}).Where("name=?", name).First(result).Error
	if err != nil {
		return nil
	}
	return result
}
func queryActor(name string) *Actor {
	var result = new(Actor)
	err := db.Model(&Actor{}).Where("name=?", name).First(&result).Error
	if err != nil {
		return nil
	}
	return result
}
func queryLastTime() int64 {
	var result Vod
	err := db.Model(&Vod{}).Select("vod_time").Limit(1).Order("vod_time desc").First(&result).Error
	if err != nil {
		return 0
	}
	t, err := time.Parse("2006-01-02 15:04:05", result.VodTime)
	if err != nil {
		return 0
	}
	return int64(math.Floor(time.Now().Sub(t).Hours() + 2))

}
func downloadPic(requestUrl string) error {
	picUrl, err := url.Parse(requestUrl)
	if err != nil {
		return err
	}
	dir := filepath.Dir(picUrl.Path)
	if strings.Index(dir, "/") == 0 {
		dir = dir[1:]
	}
	if !PathExists(dir) {
		err = os.MkdirAll(dir, os.FileMode(0766))
		if err != nil {
			return err
		}
	}
	fileName := filepath.Base(picUrl.Path)
	f, err := os.OpenFile(dir+"/"+fileName, os.O_CREATE|os.O_RDWR, os.FileMode(0766))
	if err != nil {
		return err
	}
	request, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, response.Body)
	if err != nil {
		return err
	}
	return nil
}
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
func deletePlayUrls(vodId int) error {
	return db.Model(&VodPlayUrl{}).Where("vod_id=?", vodId).Delete(&VodPlayUrl{}).Error
}
func addPlayUrls(items []*VodPlayUrl) error {
	return db.Model(&VodPlayUrl{}).Create(&items).Error
}
