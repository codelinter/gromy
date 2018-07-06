package main

import "path/filepath"

func setPath(dir, loc string, isGC bool) string {
	var finalPathToBookmark string
	gccfg := ".config/google-chrome/Default/Bookmarks"
	ccfg := ".config/chromium/Default/Bookmarks"
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
