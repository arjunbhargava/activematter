package main

import (
	"fmt"
	"github.com/arjunbhargava/activematter/xmlutils"
	"os"
	"path/filepath"
)

func parseDirectory(s string) {
	fmt.Println("Directory name is: " + s)
	fileList, err := filepath.Glob(s + "*.gz")
	if err != nil {
		fmt.Println(err)
	} else {
		for _, record := range fileList {
			xmlutils.ParseOAIXML(record)
			return
		}
	}
}

func parseFile(record string) {
	xmlutils.ParseOAIXML(record)
}

/* Entry point /wrapper for parsing records. Takes either a directory
or a single file */

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: parse_records XML_dir_name or parse_recrds XML_file_name")
		return
	}

	fi, err := os.Stat(os.Args[1])

	if err != nil {
		fmt.Printf("Error: File not found")
		return
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		fmt.Println("Parsing directory: " + os.Args[1])
		parseDirectory(os.Args[1])
	case mode.IsRegular():
		fmt.Println("Parsing file: " + os.Args[1])
		parseFile(os.Args[1])
	}
}
