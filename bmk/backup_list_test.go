package bmk

import (
	"testing"
)

func TestList(t *testing.T) {
	t.Run("Test for backup files with timestamp", func(t *testing.T) {
		exp := []string{"Jul 05 2018 10:42:51 PM", "Jul 05 2018 10:44:59 PM", "Jul 06 2018 2:04:46 PM", "Jul 06 2018 3:07:26 PM", "Jul 06 2018 3:09:07 PM"}
		res := formatList("./testdata/list")
		if len(res) != len(exp) {
			t.Error(res)
		}
	})
}
