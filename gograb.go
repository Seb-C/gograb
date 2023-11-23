package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

const doc = `Invalid arguments. Expected Syntax:
	//go:generate blockName expression replacement
Where:
	blockName is the name of the source data defined with //gograb:begin blockName
	expression is a regular expression to apply to the source data
	replacement is the replacement for the regex (see regexp.ReplaceAllString)
`

// TODO document and make an option for the substitution
const RegexVariableSymbol = "@"

func main() {
	if len(os.Args) != 4 {
		log.Fatalf(doc)
	}

	blockName := os.Args[1]
	expression := os.Args[2]
	replacement := os.Args[3]
	if blockName == "" || expression == "" || replacement == "" {
		log.Fatalf(doc)
	}

	regex, err := regexp.Compile(expression)
	if err != nil {
		log.Fatalf("Failed to compile regular expression: %s", err)
	}

	fileName := os.Getenv("GOFILE")
	if fileName == "" {
		log.Fatalf("Missing env variable GOFILE from go:generate")
	}

	lineNumberAsString := os.Getenv("GOLINE")
	if lineNumberAsString == "" {
		log.Fatalf("Missing env variable GOLINE from go:generate")
	}

	lineNumber, err := strconv.Atoi(lineNumberAsString)
	if err != nil {
		log.Fatalf(
			"Invalid line number %q from env variable GOLINE from go:generate: %s",
			lineNumberAsString,
			err,
		)
	}

	sourceFile, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Cannot read the source file: %s", err)
	}

	blockContent, err := getBlockContent(sourceFile, blockName)
	if blockContent == nil {
		log.Fatalf("%s", err)
	}

	fmt.Println(string(blockContent))
	_ = lineNumber
	_ = regex
}

func getBlockContent(sourceFile []byte, blockName string) ([]byte, error) {
	matchBegin := "//[ \t]*gograb:begin " + blockName + "[ \t]*"
	matchEnd := "\t*//[ \t]*gograb:end"
	blockFinder := regexp.MustCompile("(?sm)" + matchBegin + "\r?\n(.*?)" + matchEnd)

	matches := blockFinder.FindAllSubmatch(sourceFile, 2)
	if len(matches) == 0 {
		return nil, fmt.Errorf("Could not find block %q in source file", blockName)
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("Found multiple %q blocks in source file", blockName)
	}

	blockMatch := matches[0]
	if len(blockMatch) != 2 {
		panic(fmt.Errorf("Got %d capture groups, expected 1", len(blockMatch) - 1))
	}

	return blockMatch[1], nil
}
