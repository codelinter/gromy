package bmk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/fatih/color"
)

const (
	backupExt = "__backup__"
	success   = `[✓]`
	point     = `[•]`
	failure   = `[✘]`
)

// App encapsulates this app
type App struct {
	filename string
}

// Top needs a doc
type Top struct {
	Checksum string  `json:"checksum"`
	Roots    roots   `json:"roots"`
	Version  float64 `json:"version"`
}

type roots struct {
	BookmarkBar bm `json:"bookmark_bar"`
	Other       struct {
		Children     []interface{} `json:"children"`
		DateAdded    string        `json:"date_added"`
		DateModified string        `json:"date_modified"`
		ID           string        `json:"id"`
		Name         string        `json:"name"`
		Type         string        `json:"type"`
	} `json:"other"`
	Synced struct {
		Children     []interface{} `json:"children"`
		DateAdded    string        `json:"date_added"`
		DateModified string        `json:"date_modified"`
		ID           string        `json:"id"`
		Name         string        `json:"name"`
		Type         string        `json:"type"`
	} `json:"synced"`
}

type bm struct {
	Children     []bm   `json:"children"`
	DateAdded    string `json:"date_added"`
	DateModified string `json:"date_modified"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	Type         string `json:"type"`
}

type bookmarks []bm

type byLatestFirst struct {
	bookmarks
}

type byOldestFirst struct {
	bookmarks
}

func writeDataToFile(filename string, data []byte) {
	e := ioutil.WriteFile(filename, data, 0744)
	if e != nil {
		handleError(e, 210)
	}

}

// WriteToFile writes 'data' to file with 'filename
func (a *ff) WriteToFile(filename string, contents []byte) {
	writeDataToFile(filename, contents)
}

// FileContents gets the contents of original bookmark file
func (a *ff) FileContents(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

type ff struct{}

// Filer is an interface that is implemented by the client
type Filer interface {
	FileContents(string) ([]byte, error)
	WriteToFile(string, []byte)
}

func handleError(err error, code int) {
	e := fmt.Sprintf("%s %s", failure, err.Error())
	color.HiRed("\n\t%s", e)
	os.Exit(code)
}

func (by byLatestFirst) Less(i, j int) bool {
	if by.bookmarks[j].Type == "folder" {
		return false
	}
	return by.bookmarks[j].DateAdded < by.bookmarks[i].DateAdded
}

func (by byOldestFirst) Less(i, j int) bool {
	if by.bookmarks[j].Type == "folder" {
		return false
	}
	return by.bookmarks[j].DateAdded > by.bookmarks[i].DateAdded
}
func (bm bookmarks) Len() int      { return len(bm) }
func (bm bookmarks) Swap(i, j int) { bm[i], bm[j] = bm[j], bm[i] }

// DealChildren recusively updates the children
func DealChildren(ch *[]bm, order int) {
	for _, c := range *ch {
		if len(c.Children) > 0 {
			if order == 0 {
				sort.Sort(byLatestFirst{bookmarks: c.Children})
			} else {
				sort.Sort(byOldestFirst{bookmarks: c.Children})
			}
			DealChildren(&c.Children, order)
		}
	}
}

// NewApp returns an app that implements App interface
func NewApp(filename string) *App {
	return &App{
		filename: filename,
	}
}

// Doit executes everything
func (a *App) Doit(filer Filer, order ...int) {
	if filer == nil {
		filer = &ff{}
	}
	fc, e := filer.FileContents(a.filename)
	if e != nil {
		handleError(e, 240)
	}
	top := a.ResolveJSON(fc)
	if top == nil {
		handleError(errors.New("ResloveJSON failed"), 84)
	}
	if top.Version != 1 {
		handleError(fmt.Errorf("Unsupported version %v", top.Version), 213)
	}
	var o int
	if len(order) > 0 {
		o = order[0]
	}
	DealChildren(&top.Roots.BookmarkBar.Children, o)
	res, e := json.MarshalIndent(top, "", "  ")
	if e != nil {
		handleError(e, 84)
	}
	now := time.Now().Unix()
	backup := a.filename + backupExt + strconv.Itoa(int(now))

	filer.WriteToFile(backup, fc)

	filer.WriteToFile(a.filename, res)
	printSuccess(getType(a.filename))
	fmt.Printf("\t%s Previous bookmark is backed up in the file '%s'\n", success, backup)
}
func printSuccess(typ string) {
	color.Green("\n\t%s%s\n", success, " SUCCESS")
	fmt.Printf("\t%s Restart %s to see effect\n", point, typ)
}

// ResolveJSON populates 'top' from the json
func (a *App) ResolveJSON(fc []byte) *Top {
	v := Top{}
	e := json.Unmarshal(fc, &v)
	if e != nil {
		//handleError(e, 84)
		return nil
	}
	return &v
}
