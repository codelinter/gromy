package main

import "path/filepath"

func setPath(dir, loc string, isGC bool) string {
	var finalPathToBookmark string
	gccfg := `AppData\Local\Google\Chrome\User Data\Default\Bookmarks`
	ccfg := `AppData\Local\chromium\User Data\Default\Bookmarks`
	if isGC {
		finalPathToBookmark = filepath.Join(dir, gccfg)
	} else {
		finalPathToBookmark = filepath.Join(dir, ccfg)
	}
	// *bookmarkFileLocation overrides everyone else
	if loc != "" {
		finalPathToBookmark = loc
	}
	return finalPathToBookmark
}
