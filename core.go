package tan8

import (
	"errors"
	"fmt"
	"github.com/bysir-zl/bjson"
	"github.com/bysir-zl/bygo/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func GetUrlContent(url string) string {
	rsp, _ := http.Get(url)
	bs, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return ""
	}
	return string(bs)
	//_, rsp, _ := http_util.Get(url, util.OrderKV{}, nil)
	//return rsp
}

func GetMusicPageContent(musicId string) string {
	url := fmt.Sprintf("http://www.tan8.com/jitapu-%s.html", musicId)
	return GetUrlContent(url)
}

// 分页 获取此页面的所有musicId
func WorkPage(pageUrl, basePath string) (err error) {
	// 第一页
	firstPage := GetUrlContent(pageUrl)
	// 找出最大的页码
	rs := regexp.MustCompile(`&amp;per_page=(\d+?)"`).FindAllStringSubmatch(firstPage, -1)
	maxPerPage := 0
	for _, s := range rs {
		if len(s) > 1 {
			perPage, _ := strconv.Atoi(s[1])
			if perPage > maxPerPage {
				maxPerPage = perPage
			}
		}
	}

	// 存200个待处理的musicId
	musicIdChan := make(chan string, 200)
	musicIds, err := FindIndexMusicIds(firstPage)
	if err != nil {
		return
	}
	for _, v := range musicIds {
		musicIdChan <- v
	}

	go func() {
		for perPage := 20; perPage <= maxPerPage; perPage += 20 {
			pageContent := GetUrlContent(pageUrl + "?per_page=" + strconv.Itoa(perPage))
			ids, err := FindIndexMusicIds(pageContent)
			if err != nil {
				log.Error("tan8-FindIndexMusicIds", err)
				continue
			}

			for _, id := range ids {
				musicIdChan <- id
			}
		}

		close(musicIdChan)
	}()

	// 并行30个下载
	c := make(chan struct{}, 30)
	for musicId := range musicIdChan {
		c <- struct{}{}
		go func(musicId string) {
			log.Info("tan8-WorkOne", musicId, "downloading ")
			err := WorkOne(musicId, basePath)
			if err != nil {
				log.Error("tan8-WorkOne", musicId, err)
			} else {
				log.Info("tan8-WorkOne", musicId, "download success ")
			}

			<-c
		}(musicId)
	}

	return
}

func FindIndexMusicIds(s string) (musicIds []string, err error) {
	regImage, err := regexp.Compile(`/jitapu-(\d+?)\.html`)
	if err != nil {
		return
	}
	rs := regImage.FindAllStringSubmatch(s, -1)
	if len(rs) == 0 {
		err = errors.New("no matched MusicId, content:" + s)
		return
	}

	musicIds = []string{}
	for _, r := range rs {
		if len(r) > 1 {
			musicIds = append(musicIds, r[1:]...)
		}
	}
	if len(musicIds) == 0 {
		err = errors.New("no matched MusicId, content:" + s)
		return
	}

	return
}

func FindImage(s string) (images []string, err error) {
	regImage, err := regexp.Compile(`var arr_img = (.*?);`)
	if err != nil {
		return
	}
	rs := regImage.FindStringSubmatch(s)
	if len(rs) == 0 {
		err = errors.New("no matched Image, content:" + s)
		return
	}
	//log.Info("test", rs[1])

	p := rs[1]
	bj, err := bjson.New([]byte(p))
	if err != nil {
		return
	}
	img := bj.Index(0).Pos("img")
	l := img.Len()
	images = make([]string, l)
	for i := 0; i < l; i++ {
		images[i] = img.Index(i).String()
	}
	return
}

func FindMp3(s string) (mp3 string, err error) {
	regImage, err := regexp.Compile(`(http://oss\.tan8\.com/.*?\.mp3\?v=\d+?)'`)
	if err != nil {
		return
	}
	rs := regImage.FindStringSubmatch(s)
	if len(rs) == 0 {
		err = errors.New("no matched mp3, content:" + s)
		return
	}
	//log.Info("test", rs[1])

	p := rs[1]

	mp3 = p
	return
}

func FindName(s string) (name string, err error) {
	regName, err := regexp.Compile(`var yuepu_info_title = "(.*?)";`)
	if err != nil {
		return
	}
	rs := regName.FindStringSubmatch(s)
	if len(rs) == 0 {
		err = errors.New("no matched Name content:" + s)
		return
	}
	//log.Info("test", rs[1])

	musicName := rs[1]

	//
	regSinger, err := regexp.Compile(`var yuepu_info_singer = "(.*?)";`)
	if err != nil {
		return
	}
	rs = regSinger.FindStringSubmatch(s)
	if len(rs) == 0 {
		err = errors.New("no matched")
		return
	}
	//log.Info("test", rs[1])

	singerName := rs[1]

	name = musicName + " - " + singerName
	return
}

func DownFile(url, path, name string) (err error) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
	fileName := filepath.Base(url)
	fileName = strings.Split(fileName, "?")[0]
	if name != "" {
		fileName = name
	}

	os.MkdirAll(path, os.ModePerm)

	file, err := os.Create(path + fileName)
	if err != nil {
		return
	}
	_, err = io.Copy(file, res.Body)
	return
}

func WorkOne(musicId, basePath string) (err error) {
	s := GetMusicPageContent(musicId)

	is, err := FindImage(s)
	if err != nil {
		return
	}
	r, err := FindMp3(s)
	if err != nil {
		return
	}
	name, err := FindName(s)
	if err != nil {
		return
	}
	name = name + " . " + musicId

	path := basePath + name + "/"
	err = DownFile(r, path, "")
	if err != nil {
		return
	}

	for _, img := range is {
		name := strings.Split(img, "image_")[1]
		name = strings.Split(name, "?")[0]
		err := DownFile(img, path, name)
		if err != nil {
			log.Error("tan8-DownFile", err)
		}
	}

	return
}
