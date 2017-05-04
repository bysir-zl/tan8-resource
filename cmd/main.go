package main

import (
	"flag"
	"github.com/bysir-zl/bygo/log"
	"github.com/bysir-zl/tan8-resource"
	"regexp"
	"strings"
)

func main() {
	var types, downloadPath, musicUrl, musicId string
	flag.StringVar(&downloadPath, "path", "./download/", "download path")
	flag.StringVar(&musicUrl, "url", "", "tan8 music page url")
	flag.StringVar(&musicId, "id", "", "tan8 music id")
	flag.StringVar(&types, "type", "one", "one or page")
	flag.Parse()

	switch types {
	case "one":
		getOne(downloadPath, musicUrl, musicId)
	case "page":
		getIndex(musicUrl, downloadPath)
	}

	return
}

func getOne(downloadPath, musicUrl, musicId string) {
	if musicUrl == "" && musicId == "" {
		log.Error("tan8", "url or id must exist")
		return
	}

	if musicUrl != "" {
		musicId = regexp.MustCompile(`http://www.tan8.com/jitapu-(\d+?).html`).FindStringSubmatch(musicUrl)[1]
	}

	if len(downloadPath) != 0 && !strings.HasSuffix(downloadPath, "/") {
		downloadPath = downloadPath + "/"
	}
	err := tan8.WorkOne(musicId, downloadPath)
	if err != nil {
		log.Error("tan8", err)
		return
	}
	log.Info("tan8", "download success")
}

func getIndex(pageUrl, downloadPath string) {
	if pageUrl == "" {
		log.Error("tan8", "url must exist")
		return
	}

	err := tan8.WorkPage(pageUrl, downloadPath)
	if err != nil {
		log.Error("tan8", err)
		return
	}

	log.Info("tan8", "download success")
}
