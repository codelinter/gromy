package bmk

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"

	s "gopkg.in/AlecAivazis/survey.v1"
)

var tFormat = "Jan 02 2006 3:04:05 PM"

type fBackup struct {
	Editr  string
	Backup string //`survey:"backup"`
	Tester string
}

type backedUpFile struct {
	name       string
	timeFormat int
}

var finalList = make([]backedUpFile, 0)

func getTimeList(root string) []backedUpFile {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		// errors not important in this case
		abs, _ := filepath.Abs(filepath.Dir(path))
		absRoot, _ := filepath.Abs(root)
		if abs != absRoot {
			return nil
		}
		if !strings.Contains(filepath.Base(path), backupExt) {
			return nil
		}
		a := strings.Split(filepath.Base(path), backupExt)
		if len(a) != 2 {
			return nil
		}
		iVal, e := strconv.Atoi(a[1])
		if e != nil {
			handleError(e, 1)
		}
		bF := backedUpFile{
			name:       path,
			timeFormat: iVal,
		}
		finalList = append(finalList, bF)
		return nil
	})
	if err != nil {
		return nil
	}
	return finalList
}
func timeInFormat(tm int) string {
	return time.Unix(int64(tm), 0).Format(tFormat)
}

func formatList(root string) []backedUpFile {
	return getTimeList(root)
}

func getType(jsonPath string) string {
	var typ int
	if !strings.Contains(strings.ToLower(jsonPath), "chromium") {
		typ = 1
	}
	var chrom = "Chromium"
	if typ == 1 {
		chrom = "Google Chrome"
	}
	return chrom
}

func getFormattedList(jsonPath string, ans *fBackup, finalList []backedUpFile) ([]string, error) {
	chrom := getType(jsonPath)
	var iList []int
	var list []string
	for _, bF := range finalList {
		iList = append(iList, bF.timeFormat)
	}
	sort.Sort((sort.Reverse(sort.IntSlice(iList))))
	for _, l := range iList {
		list = append(list, timeInFormat(l))
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("empty list")
	}
	msg := fmt.Sprintf("Choose backup for %s", chrom)
	var qs = []*s.Question{
		{
			Name: "backup",
			Prompt: &s.Select{
				Message: msg,
				Options: list,
				Default: list[0],
			},
		},
	}

	err := s.Ask(qs, ans)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func getIndex(list []string, ans *fBackup) int {
	var idx = -1
	for i, l := range list {
		if l == ans.Backup {
			idx = i
		}
	}
	return idx
}

// Restore tries to restore bookmark from the backup
func Restore(jsonPath string) {
	fList := formatList(filepath.Dir(jsonPath))
	ans := &fBackup{}
	list, err := getFormattedList(jsonPath, ans, fList)
	if err != nil {
		// TODO
		return
	}
	idx := getIndex(list, ans)
	if idx == -1 {
		handleError(fmt.Errorf("Unexpected behavior"), 222)
	}
	color.Magenta("\nRestoring %s to %s\n", finalList[idx].name, jsonPath)
	b, err := ioutil.ReadFile(finalList[idx].name)
	if err != nil {
		handleError(err, 210)
	}
	writeDataToFile(jsonPath, b)
	printSuccess(getType(jsonPath))
}
