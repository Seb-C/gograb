package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
)

func main() {
	fileName := os.Getenv("GOFILE")
	if fileName == "" {
		log.Fatalf("Missing env variable GOFILE from go:generate")
	}

	lineNumberAsString := os.Getenv("GOLINE")
	if lineNumberAsString == "" {
		log.Fatalf("Missing env variable GOLINE from go:generate")
	}

	sourceFile, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Cannot read the source file: %s", err)
	}

	source, err := getSource(sourceFile)
	if err != nil {
		log.Fatal(err)
	}

	newSource, err := findAndReplaceTargets(sourceFile, func(params []byte, content []byte) ([]byte, error) {
		fmt.Println("BEGIN")
		fmt.Println(string(params))
		fmt.Println("CONTENT")
		fmt.Println(string(content))
		fmt.Println("END")

		return []byte("replaced"), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	_ = source
	_ = newSource
}

func getSource(sourceFile []byte) ([]byte, error) {
	matchBegin := "//[ \t]*gograb:source[ \t]*"
	matchEnd := "\n\t*//[ \t]*gograb:end"
	sourceRegexp := regexp.MustCompile("(?s)" + matchBegin + "\r?\n(.*?)" + matchEnd)

	matches := sourceRegexp.FindAllSubmatch(sourceFile, 2)
	if len(matches) == 0 {
		return nil, errors.New("Could not find source")
	}
	if len(matches) > 1 {
		return nil, errors.New("Found multiple sources")
	}

	blockMatch := matches[0]
	if len(blockMatch) != 2 {
		panic(fmt.Sprintf("Got %d capture groups, expected 1", len(blockMatch)-1))
	}

	return blockMatch[1], nil
}

// TODO test this function
// TODO implement code to replace a target given a regexp and code block
func findAndReplaceTargets(
	sourceFileContent []byte,
	iterator func(params []byte, content []byte) ([]byte, error),
) ([]byte, error) {
	matchBegin := "//[ \t]*gograb:target[ \t]*(.*?)[ \t]*"
	matchEnd := "\n\t*//[ \t]*gograb:end"
	targetRegexp := regexp.MustCompile("(?s)" + matchBegin + "\r?\n(.*?)" + matchEnd)

	matches := targetRegexp.FindAllSubmatchIndex(sourceFileContent, -1)
	if len(matches) == 0 {
		return nil, errors.New("Could not find any target")
	}

	// Looping in reverse so we can safely replace the contents one by one
	targetFileContent := sourceFileContent
	for matchIndex := len(matches)-1; matchIndex >= 0; matchIndex-- {
		match := matches[matchIndex]
		if groups, expect := (len(match)/2)-1, 2; groups != expect {
			panic(fmt.Sprintf("Got %d capture groups for match %d, expected %d", matchIndex, groups, expect))
		}

		result, err := iterator(
			sourceFileContent[match[2]:match[3]],
			sourceFileContent[match[4]:match[5]],
		)
		if len(matches) == 0 {
			return nil, fmt.Errorf("Could not replace target %d: %w", matchIndex, err)
		}

		targetFileContent = append(
			append(targetFileContent[:match[4]], result...),
			targetFileContent[match[5]:]...,
		)
	}

	fmt.Println(string(targetFileContent))

	return nil, nil
}
