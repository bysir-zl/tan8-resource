package tan8

import (
	"testing"
	"github.com/bysir-zl/bygo/log"
	"time"
)

func TestGetUrl(t *testing.T) {
	s := GetMusicPageContent("53142")

	is, err := FindImage(s)
	t.Log(err, is)
	r, err := FindMp3(s)
	t.Log(err, r)
	name, err := FindName(s)
	t.Log(err, name)

	path := "./download/" + name + "/"
	err = DownFile(r, path, "")
	if err != nil {
		t.Log(err)
	}
	for _, img := range is {
		err = DownFile(img, path, "")
		if err != nil {
			t.Log(err)
		}
	}

}

func TestWork(t *testing.T) {
	WorkOne("10512", "./download/")
}

func TestWorkPage(t *testing.T) {
	err := WorkPage("http://www.tan8.com/guitar-291-0-collects-1-0.html", "./download/")
	if err != nil {
		t.Error(err)
	}
}

func TestChan(t *testing.T) {
	var c = make(chan int, 10)

	go func() {
		c <- 1
		c <- 2
		c <- 3
		c <- 4
		c <- 5
		time.Sleep(time.Second * 3)
		close(c)
	}()

	for v:=range c{
		time.Sleep(time.Second * 3)

		log.Info("test",v)
	}
}
