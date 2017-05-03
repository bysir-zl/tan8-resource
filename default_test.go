package tan8

import (
	"testing"
)

func TestGetUrl(t *testing.T) {
	s := GetUrl("http://www.tan8.com/jitapu-54165.html")

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