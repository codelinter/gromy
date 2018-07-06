package bmk

import (
	"fmt"
	"io/ioutil"
	"testing"
)

type testApp struct {
	counter                                    int
	originalFileName, fileName                 string
	backupContent, newContent, originalContent []byte
	haveLocation, wantedLocation               string
}

func (a *testApp) WriteToFile(filename string, contents []byte) {
	if a.counter == 0 {
		a.counter++
		a.backupContent = contents
		return
	}
	a.fileName = filename
	a.newContent = contents
	//ioutil.WriteFile("./testdata/2.json", contents, 0744)
}
func (a *testApp) FileContents(filename string) ([]byte, error) {
	gc, e := ioutil.ReadFile(a.haveLocation)
	if e != nil {
		return nil, e
	}
	a.originalContent = gc
	return gc, nil
}

func TestBmk(t *testing.T) {
	//wantedLocation := "./testdata/2.json"
	t.Run("TestBookmark", func(t *testing.T) {
		tApp := &testApp{
			haveLocation:   "./testdata/1.json",
			wantedLocation: "./testdata/2.json",
		}
		app := NewApp("f1")
		app.Doit(tApp)
		fc, e := ioutil.ReadFile(tApp.wantedLocation)
		if e != nil {
			fmt.Println("ERR", e)
		}
		if string(tApp.newContent) != string(fc) {
			t.Errorf("Wanted %s \n But Got %s ", string(fc), string(tApp.newContent))
		}
		if string(tApp.backupContent) != string(tApp.originalContent) {
			t.Errorf("Wanted %s \n But Got %s ", string(tApp.originalContent), string(tApp.backupContent))
		}
	})
}
