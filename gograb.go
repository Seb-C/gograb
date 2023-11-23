package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const doc = `Invalid arguments. Expected Syntax:
	//go:generate blockName expression replacement
Where:
	blockName is the name of the source data defined with //gograb:begin blockName
	expression is a regular expression to apply to the source data
	replacement is the replacement for the regex (see regexp.ReplaceAllString)
`

// TODO document and make an option for the substitution
const regexVariableSymbol = "@"

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
	if err != nil {
		log.Fatal(err)
	}

	newBlockContent := transformBlockAsRequested(blockContent, regex, replacement)

	replacedContent, err := replaceTargetBlockContent(sourceFile, newBlockContent, lineNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(replacedContent))
}

func getBlockContent(sourceFile []byte, blockName string) ([]byte, error) {
	matchBegin := "//[ \t]*gograb:begin " + blockName + "[ \t]*"
	matchEnd := "\n\t*//[ \t]*gograb:end"
	blockFinder := regexp.MustCompile("(?s)" + matchBegin + "\r?\n(.*?)" + matchEnd)

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

func replaceTargetBlockContent(
	sourceFile []byte,
	newContent []byte,
	startLineNumber int,
) ([]byte, error) {
	matchEnd := "(\n\t*//[ \t]*gograb:end.*?)$"
	startLineNumberString := strconv.Itoa(startLineNumber)
	targetFinder := regexp.MustCompile(
		"(?s)^(([^\n]*\n?){" + startLineNumberString + "}?)(.*?)" + matchEnd,
	)

	matches := targetFinder.FindAllSubmatch(sourceFile, 2)
	if len(matches) == 0 {
		return nil, fmt.Errorf(
			"Could not find a target between line %d and //gograb:end in source file",
			startLineNumber,
		)
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf(
			"Found multiple targets between line %d and //gograb:end in source file",
			startLineNumber,
		)
	}

	targetMatch := matches[0]
	if len(targetMatch) != 5 {
		panic(fmt.Errorf("Got %d capture groups, expected 4", len(targetMatch) - 1))
	}

	return append(append(targetMatch[1], newContent...), targetMatch[4]...), nil
}

func transformBlockAsRequested(
	blockContent []byte,
	regex *regexp.Regexp,
	replacement string,
) []byte {
	formattedReplacement := strings.ReplaceAll(replacement, regexVariableSymbol, "$")
	return regex.ReplaceAll(blockContent, []byte(formattedReplacement))
}
