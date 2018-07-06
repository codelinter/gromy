package main

import "path/filepath"

func setPath(dir, loc string, isGC bool) string {
	var finalPathToBookmark string
	gccfg := "Library/Application Support/Google/Chrome/Default/Bookmarks"
	ccfg := "Library/Application Support/Chromium/Default/Bookmarks"
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
