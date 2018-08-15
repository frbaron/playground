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

func printInfo(msg string) {
	if bVerbose > 0 {
		fmt.Println(msg)
	}
}

func Scanner(dirName string, filePtr string, weed string) {

    // Display the options used in processing below
    printInfo(weed)
    printInfo(filePtr)
    printInfo(dirName)

    // Open the target folder
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
    // 1- confirm it is matching our pattern
    // 2- process the weeding out bad strings to compute cleaner name
    // 3- rename the file to the new name (if changed)
    for _, file := range files {
	oldname := file.Name()

	fileMatch, _ := regexp.MatchString(filePtr, oldname)

	if fileMatch {
		printInfo(file.Name())

		//re := regexp.MustCompile(weed)
		//newname := re.ReplaceAllLiteralString(oldname, "")

		newname := Weeder(oldname, weed)

		if newname != oldname {
			// Check if file already exists before renaming/overwriting it
			fileInfo, _ := os.Stat( dirName + "/" + newname )

			if fileInfo == nil {
				fmt.Println(">>" + newname)
				err :=	os.Rename(dirName+"/"+oldname, dirName+"/"+newname)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				printInfo("File already exists:" + newname)
			}
		} else {
			printInfo("Nothing to clean" + oldname)
		}
    	}
    }
}


// Global flag to display progress messages (silent by default)
var bVerbose int = 0

// We need a small utility to remove dodgy endings in file names
// Usually coming from internet sources each with their own extensions, tags..

func main() {

  // Managing global values and expose as CLI parameters
  dirName	:= "."	 // folder where to list and clean files
  filePtr	:= "*"	 // pattern to only modify files matching that pattern
  weed	:= "WebRip" // example string to weed out of file name

  var pVerb	= flag.Int("v",		0,	"Boolean option to activate display of progress messages")
  var pWeed	= flag.String("weed",	"HDTV",	"String you want to weed out of all files in folder")
  var pFile	= flag.String("file",	"*",	"Pattern to only rename files matching this. Defauly: all files")
  var pDir	= flag.String("dirName",".",	"Exact directory path where to rename files. Default: current folder")

  // Parsing the command line arguments as defined above
  flag.Parse()

  if pVerb != nil {
    bVerbose = *pVerb;
  }
  if pDir != nil {
    dirName	= *pDir;
  }
  if pFile != nil {
    filePtr	= *pFile;
  }
  if pWeed != nil {
    weed	= *pWeed;
  }

  Scanner(dirName, filePtr, weed)
}
