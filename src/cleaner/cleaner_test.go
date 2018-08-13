package main

import (
	"fmt"
	"testing"
	"io/ioutil"
	"log"
	"os"
)

var dirName string
var cntFiles int = 0

func init() {
	dirName = initFolderSetup()
}

func initFolderSetup() (string) {

	dir, err := ioutil.TempDir(".", "toclean")
	if err != nil {
		log.Fatal(err)
	}
//	defer os.RemoveAll(dir)			// comment-out to troubleshoot what is done in temp folder

	files := []struct {
		name string
	}{
		{"S01E01.mp4"},
		{"S01E01 WEB-HD.mp4"},
		{"S01E02 WEB-HD.mp4"},
		{"S01E01 WEB-HD.xVid.mp4"},
		{"S01E01  xVid.mp4"},
		{"S01E01 xVid.mp4"},
	}

	for _, f := range files {

		_, err := os.OpenFile( dir + "/" +f.name, os.O_RDONLY|os.O_CREATE, 0666)

		if err != nil {
			fmt.Println("Could not create the temp file for testing: ", f.name)
		}

		cntFiles++
	}

//	fmt.Println("Created temp: ", dir)
	return dir
}

func initFolderTearDown() {
	fmt.Println("Undo nothing")
}

func TestWeeder(t *testing.T) {
	cases := []struct {
		in, weed, want string
	}{
		{"S01E01.mp4",		"WEB-HD", "S01E01.mp4"},
		{"S01E01 WEB-HD.mp4",	"WEB-HD", "S01E01.mp4"},
		{"S01E01  xVid.mp4",	"WEB-HD", "S01E01 xVid.mp4"},
		{"S01E01  xVid.mp4",	"xVid", "S01E01.mp4"},
	}

	for _, c := range cases {
		got := Weeder(c.in, c.weed)

		if got != c.want {
			t.Errorf("Weeder(%q, %q) == %q, want %q", c.in, c.weed, got, c.want)
		}
	}
}

func TestScanner(t *testing.T) {

	Scanner(dirName, ".*", "WEB-HD")

	fmt.Println("Created temp: ", dirName)

	gotFiles := 0

    // Open the cleaned folder
    f, err := os.Open(dirName)
    if err != nil {
        log.Fatal(err)
    }

    // Build a list of all files in that folder
    files, err := f.Readdir(-1)
    f.Close()
    if err != nil {
        log.Fatal(err)
    }

    // Go through each file and
    for _, file := range files {
	name := file.Name()
	gotFiles++
	//fileMatch, _ := regexp.MatchString(filePtr, oldname)
	//if fileMatch {
		fmt.Println(name)
    	//}
    }

    if gotFiles != cntFiles {
	t.Errorf("Scanner() == %d, want %d files", gotFiles, cntFiles )
    }
}
