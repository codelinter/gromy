package bmk

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	s "gopkg.in/AlecAivazis/survey.v1"
)

var tFormat = "Jan 02 2006 3:04:05 PM"

func getTimeList(root string) []string {
	var finalList = make([]string, 0)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
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
		iVal, err := strconv.Atoi(a[1])
		if err != nil {
			handleError(err, 1)
		}
		finalList = append(finalList, time.Unix(int64(iVal), 0).Format(tFormat))
		return nil
	})
	return finalList
}

func formatList(root string) []string {
	return getTimeList(root)
}

// Restore tries to restore bookmark from the backup
func Restore(jsonPath string) {
	fList := formatList(filepath.Dir(jsonPath))
	var typ int
	if !strings.Contains(strings.ToLower(jsonPath), "chromium") {
		typ = 1
	}
	setupPrompt(fList, typ)
}

type fBackup struct {
	Editr  string
	Backup string //`survey:"backup"`
	Tester string
}

func setupPrompt(list []string, typ int) {
	var chrom = "Chromium"
	if typ == 1 {
		chrom = "Google Chrome"
	}
	msg := fmt.Sprintf("Choose backup for %s", chrom)
	var qs = []*s.Question{{
		Name: "backup",
		Prompt: &s.Select{
			Message: msg,
			Options: list,
			Default: list[0],
		},
	},
	}
	ans := &fBackup{}
	err := s.Ask(qs, ans)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("FINAL %v\n", ans.Backup)
}
