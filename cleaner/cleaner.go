package main

// Import of basic modules to play around printing, file commands, regexp and cli options
import (
    "fmt"
    "log"
    "os"
    "regexp"
    "flag"
)

func main() {

    // Managing global values and expose as CLI parameters 
    dirname	:= "."	 // folder where to list and clean files
    filePtr	:= "*"	 // pattern to only modify files matching that pattern
    weed	:= "WebRip" // example string to weed out of file name

    var pWeed	= flag.String("weed",	"HDTV",	"String you want to weed out of all files in folder")
    var pFile	= flag.String("file",	"*",	"Pattern to only rename files matching this. Defauly: all files")
    var pDir	= flag.String("dirname",".",	"Exact directory path where to rename files. Default: current folder")

    // Parsing the command line arguments as defined above
    flag.Parse()

    if pWeed != nil {
	weed	= *pWeed;
    }
    if pFile != nil {
	filePtr	= *pFile;
    }
    if pDir != nil {
	dirname	= *pDir;
    }

    fmt.Println(weed)
    fmt.Println(filePtr)
    fmt.Println(dirname)
	
    f, err := os.Open(dirname)
    if err != nil {
        log.Fatal(err)
    }
    files, err := f.Readdir(-1)
    f.Close()
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
	oldname := file.Name()

	fileMatch, _ := regexp.MatchString(filePtr, oldname)
	
	if fileMatch {
		fmt.Println(file.Name())

		re := regexp.MustCompile(weed)
		newname := re.ReplaceAllLiteralString(oldname, "")

		if newname != oldname {
			fmt.Println(">>", newname)
			err :=	os.Rename(oldname, newname)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("Nothing to clean", oldname)
		}
    	}
    }

}

