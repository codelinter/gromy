package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/mintyowl/gromy/bmk"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

var chromJSON string

const helperMsg = `
   ###############################################################
   ########                      GROMY                    ########
   ########       sorts chrom(ium) bookmarks by date      ########
   ###############################################################
   ########   Folders will get pushed to the top/bottom   ########
   ########        RESTART CHROM(IUM) TO SEE EFFECT       ########
   ###############################################################
   ########               BACKUPS CAN BE FOUND            ########
   ########  with filenames ending '__backup__<UnixTS>'   ######## 
   ###############################################################
   
   -loc=<Full path to the google-chrom(ium) bookmark json file>
   ------------------------------------------------------------
   -gc [no value required. When used, type switches to "Google Chrome" instead of the default "Chromium"]
   ------------------------------------------------------------------------------------------------------
   -order=<oldest/o> [latest is default]
   -------------------------------------
   -y [no value required. When used, will bypass the prompt completely]
   --------------------------------------------------------------------
   -restore [no value required. Restore old backups by date]
   --------------------------------------------------------------------

   EXAMPLE1: gromy 
	 #1 Targets, default, Chromium bookmarks at default location and sorts by, default, latest first
		
   EXAMPLE2: gromy -gc -order=oldest 
	 #2 targets Google-Chrome to sort by oldest first
		
   EXAMPLE3: gromy -gc -order=oldest 
	 #3 targets Google-Chrome to sort by oldest first

   EXAMPLE4: gromy.exe -order=oldest -loc='D:\Path\To\Different\Location\Bookmarks'
	 #4 Sort by oldest first at custom location
   
   EXAMPLE5: gromy --restore -gc
	 #5 Restore Google Chrome bookmarks from previous/old backups. Prompts you to choose from options
`

func handleError(err error, code int) {
	color.HiRed("%s", err.Error())
	if code > 0 {
		os.Exit(code)
	}
}

func flagUsage() {
	fmt.Fprintln(os.Stderr, helperMsg)
}

func main() {
	isGC := flag.Bool("gc", false, "Chromium or Google-Chrome")
	isRestore := flag.Bool("restore", false, "Restore from backup")
	noPrompt := flag.Bool("y", false, "Execute without prompt")
	help := flag.Bool("h", false, "Helper")
	bookmarkFileLocation := flag.String("loc", "", "Location to the bookmark json file")
	order := flag.String("order", "latest", "latest or oldest")
	flag.Usage = flagUsage
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}
	dir, err := homedir.Dir()
	if err != nil {
		handleError(err, 200)
	}
	switch runtime.GOOS {
	case "linux":
		chromJSON = setPath(dir, *bookmarkFileLocation, *isGC)
	case "windows":
		chromJSON = setPath(dir, *bookmarkFileLocation, *isGC)
	case "darwin":
		chromJSON = setPath(dir, *bookmarkFileLocation, *isGC)
	default:
		if *bookmarkFileLocation == "" {
			handleError(errors.New("Operating system not supported. Try --loc flag to provide full path to bookmark file instead"), 250)
		} else {
			chromJSON = *bookmarkFileLocation
		}
	}
	if chromJSON == "" {
		handleError(errors.New("Bookmark file location is empty"), 211)
	} else {
		stats, err := os.Stat(chromJSON)
		if err != nil {
			handleError(fmt.Errorf("'%s' path doesnot exist", chromJSON), 132)
		}
		if stats.IsDir() {
			handleError(fmt.Errorf("'%s' path is a directory", chromJSON), 132)

		}
	}
	var o int
	var ordr string
	if *order != "" {
		if strings.ToLower(*order) == "oldest" || strings.ToLower(*order) == "o" {
			o = 1
		}
	}
	if o == 0 {
		ordr = fmt.Sprint("Latest")
	} else {
		ordr = fmt.Sprint("Oldest")
	}

	if *isRestore {
		bmk.Restore(chromJSON)
		return
	}

	if !*noPrompt {
		c := color.New(color.FgHiMagenta)
		fmt.Print("File to be modified >> ")
		c.Printf("%s\n", chromJSON)
		fmt.Print("Order >> ")
		c.Printf("%s\n", ordr)
		//c.Println(ordr)
		c2 := color.New(color.FgHiRed).Add(color.Bold)
		c2.Print("Confirm modifications (y/n)?: ")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if scanner.Text() == "y" {
				app := bmk.NewApp(chromJSON)
				app.Doit(nil, o)
				return
			}
			exiting := "Exiting..."
			handleError(errors.New(exiting), 0)
		}
	} else {
		app := bmk.NewApp(chromJSON)
		app.Doit(nil, o)
	}
}
