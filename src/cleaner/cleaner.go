package main

// Import of basic modules to play around printing, file commands, regexp and cli options
import (
    "fmt"
    "log"
    "os"
    "regexp"
    "strings"
    "flag"
)

// Function which is getting a name of a file, removing any extensions, normalising as well
// Input  : "Black.Panther.HDTV.BRRip.2017.mp4"
// Output : "Black Panther.mp4"

func Weeder(fileName string, weed string) (string) {

	// initiate our return variable
	newname := fileName

	// Breaking down filename in clear parts
	re := regexp.MustCompile(`^(.*)\.(\S{3,4})$`)
	findings := re.FindStringSubmatch(fileName)

	if findings != nil {
		baseName:= findings[1]
		ext	:= findings[2]

		// Replace any "." by a space
		re = regexp.MustCompile(`\.`)
		baseName = re.ReplaceAllString(baseName, " ")

		// Main weeding out work... based on parameter
		// Building up the regexp "(string)" to weed out of the filename
		s := []string{ "(" , weed , ")" }
		repl := strings.Join( s, "" )
		re = regexp.MustCompile(repl)
		baseName = re.ReplaceAllString(baseName, "")

		// Removing duplicate spaces
		re = regexp.MustCompile(`\s+`)
		baseName = re.ReplaceAllString(baseName, " ")

		// Removing any trailing space
		re = regexp.MustCompile(`\s+$`)
		baseName = re.ReplaceAllString(baseName, "")

		// Reconstruct the filename
		newname = strings.Join( []string{ baseName, ".", ext}, "" )
	}

	return newname
}

// We need a small utility to remove dodgy endings in file names
// Usually coming from internet sources each with their own extensions, tags..

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

    // Display the options used in processing below
    fmt.Println(weed)
    fmt.Println(filePtr)
    fmt.Println(dirname)

    // Open the target folder
    f, err := os.Open(dirname)
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
    // 1- confirm it is matching our pattern
    // 2- process the weeding out bad strings to compute cleaner name
    // 3- rename the file to the new name (if changed)
    for _, file := range files {
	oldname := file.Name()

	fileMatch, _ := regexp.MatchString(filePtr, oldname)
	
	if fileMatch {
		fmt.Println(file.Name())

		//re := regexp.MustCompile(weed)
		//newname := re.ReplaceAllLiteralString(oldname, "")

		newname := Weeder(oldname, weed)
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

