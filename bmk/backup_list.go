package bmk

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var tFormat = "Jan 02 2006 3:04:05 PM"

func getTimeList(root string) []string {
	ext := "__backup__"
	var finalList []string
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Dir(path) != root {
			return nil
		}
		if !strings.Contains(filepath.Base(path), ext) {
			return nil
		}
		a := strings.Split(filepath.Base(path), ext)
		tm := a[1]
		iVal, _ := strconv.Atoi(tm)
		finalList = append(finalList, time.Unix(int64(iVal), 0).Format(tFormat))
		return nil
	})
	return finalList
}

func formatList(root string) []string {
	return getTimeList(root)
}
