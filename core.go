package tan8

import (
	"github.com/bysir-zl/bygo/util/http_util"
	"github.com/bysir-zl/bygo/util"
	"github.com/bysir-zl/bjson"
	"regexp"
	"errors"
	"net/http"
	"os"
	"io"
	"path/filepath"
	"strings"
)

func GetUrl(url string) string {
	_, rsp, _ := http_util.Get(url, util.OrderKV{}, nil)
	return rsp
}

func FindImage(s string) (images []string, err error) {
	regImage, err := regexp.Compile(`var arr_img = (.*?);`)
	if err != nil {
		return
	}
	rs := regImage.FindStringSubmatch(s)
	if len(rs) == 0 {
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
		err = errors.New("no matched")
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
		err = errors.New("no matched")
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