package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/armarquez/gather-code/version"
)

var (
	inputPath   string
	outputFile  string
	extensions  string
	versionFlag bool
	debug       bool
)

func init() {
	flag.StringVar(&inputPath, "input-path", "", "Path to traverse for code files")
	flag.StringVar(&inputPath, "i", "", "Path to traverse for code files (shorthand)")
	flag.StringVar(&outputFile, "output-file", "", "Output file path (default is standard output)")
	flag.StringVar(&outputFile, "o", "", "Output file path (shorthand)")
	flag.StringVar(&extensions, "extensions", "go", "Comma-separated list of file extensions to search (default is 'go')")
	flag.StringVar(&extensions, "e", "go", "Comma-separated list of file extensions to search (shorthand)")
	flag.BoolVar(&versionFlag, "version", false, "Print version information and exit")
	flag.BoolVar(&versionFlag, "v", false, "Print version information and exit")
	flag.BoolVar(&debug, "debug", false, "Enable debug mode to print file paths and decisions")
	flag.Parse()

	if versionFlag {
		fmt.Println("Build Date:", version.BuildDate)
		fmt.Println("Git Commit:", version.GitCommit)
		fmt.Println("Version:", version.Version)
		fmt.Println("Go Version:", version.GoVersion)
		fmt.Println("OS / Arch:", version.OsArch)
		os.Exit(0)
	}

	if inputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: input-path is required")
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	// Define the file extensions to look for
	var extList []string
	if extensions != "" {
		extList = strings.Split(extensions, ",")
	}
	if debug {
		fmt.Printf("Extensions: %v\n", extList)
	}

	// Open output file or use standard output
	var out *os.File
	if outputFile != "" {
		var err error
		out, err = os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer out.Close()
	} else {
		out = os.Stdout
	}

	// Walk through the input path
	err := filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if hasExtension(path, extList) {
			if debug {
				fmt.Printf("Checking file: %s\n", path)
				fmt.Println("-- Decision: Include")
			}
			err := copyFileContents(out, path)
			if err != nil {
				return fmt.Errorf("error copying file contents: %v", err)
			}
			fmt.Fprintln(out, "-------------")
		} else {
			if debug {
				fmt.Printf("Checking file: %s\n", path)
				fmt.Println("-- Decision: Skip")
			}
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error traversing input path: %v\n", err)
		os.Exit(1)
	}
}

// Check if the file has one of the desired extensions
func hasExtension(path string, extensions []string) bool {
	if len(extensions) == 0 {
		return true // If no extensions specified, match all
	}
	ext := strings.ToLower(filepath.Ext(path))
	if debug {
		fmt.Printf("-- Extension: %s\n", ext)
	}
	for _, e := range extensions {
		if ext == "."+strings.ToLower(e) {
			return true
		}
	}
	return false
}

// Copy the contents of the file to the output file
func copyFileContents(out *os.File, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "File: "+filePath)
	fmt.Fprintln(out, "```")
	_, err = out.Write(data)
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "```")
	return nil
}
