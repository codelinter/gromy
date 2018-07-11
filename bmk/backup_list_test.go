package bmk

import (
	"reflect"
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

func Test_formatList(t *testing.T) {
	type args struct {
		root string
	}
	tests := []struct {
		name string
		args args
		want []backedUpFile
	}{
		{
			name: "test timelist",
			args: args{
				root: "./testdata/list"},
			want: []backedUpFile{backedUpFile{name: "testdata/list/Bookmarks__backup__1530810771", timeFormat: 1530810771}, backedUpFile{name: "testdata/list/Bookmarks__backup__1530810899", timeFormat: 1530810899}, backedUpFile{name: "testdata/list/Bookmarks__backup__1530866086", timeFormat: 1530866086}, backedUpFile{name: "testdata/list/Bookmarks__backup__1530869846", timeFormat: 1530869846}, backedUpFile{name: "testdata/list/Bookmarks__backup__1530869947", timeFormat: 1530869947}, backedUpFile{name: "testdata/list/Bookmarks__backup__1530810771", timeFormat: 1530810771}, backedUpFile{name: "testdata/list/Bookmarks__backup__1530810899", timeFormat: 1530810899}, backedUpFile{name: "testdata/list/Bookmarks__backup__1530866086", timeFormat: 1530866086}, backedUpFile{name: "testdata/list/Bookmarks__backup__1530869846", timeFormat: 1530869846}, backedUpFile{name: "testdata/list/Bookmarks__backup__1530869947", timeFormat: 1530869947}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatList(tt.args.root); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("formatList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getType(t *testing.T) {
	type args struct {
		jsonPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test google chrome",
			args: args{
				jsonPath: "/path/to/google-chrome/Default/Bookmarks",
			},
			want: "Google Chrome",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getType(tt.args.jsonPath); got != tt.want {
				t.Errorf("getType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFormattedList(t *testing.T) {
	type args struct {
		jsonPath  string
		ans       *fBackup
		finalList []backedUpFile
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr string
	}{
		{
			name: "formatted list",
			args: args{
				jsonPath:  "",
				ans:       &fBackup{},
				finalList: make([]backedUpFile, 0),
			},
			want:    make([]string, 0),
			wantErr: "empty list",
		},
		{
			name: "formatted list",
			args: args{
				jsonPath:  "",
				ans:       &fBackup{},
				finalList: []backedUpFile{backedUpFile{name: "testdata/list/Bookmarks__backup__1530810771", timeFormat: 1530810771}, backedUpFile{name: "testdata/list/Bookmarks__backup__1530810899", timeFormat: 1530810899}, backedUpFile{name: "testdata/list/Bookmarks__backup__1530866086", timeFormat: 1530866086}, backedUpFile{name: "testdata/list/Bookmarks__backup__1530869846", timeFormat: 1530869846}, backedUpFile{name: "testdata/list/Bookmarks__backup__1530869947", timeFormat: 1530869947}, backedUpFile{name: "testdata/list/Bookmarks__backup__1530810771", timeFormat: 1530810771}, backedUpFile{name: "testdata/list/Bookmarks__backup__1530810899", timeFormat: 1530810899}, backedUpFile{name: "testdata/list/Bookmarks__backup__1530866086", timeFormat: 1530866086}, backedUpFile{name: "testdata/list/Bookmarks__backup__1530869846", timeFormat: 1530869846}, backedUpFile{name: "testdata/list/Bookmarks__backup__1530869947", timeFormat: 1530869947}},
			},
			want:    make([]string, 0),
			wantErr: "EOF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFormattedList(tt.args.jsonPath, tt.args.ans, tt.args.finalList)
			if err.Error() != tt.wantErr {
				t.Errorf("getFormattedList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("getFormattedList() = %v, want %v", got, tt.want)
			}
		})
	}
}
