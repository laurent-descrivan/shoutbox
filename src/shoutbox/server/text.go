package server

import (
	"io/ioutil"
	"math/rand"
	"regexp"
	"strings"
	"sync"
)

var lines []string
var lineIdx int
var linesMutex sync.Mutex

func init() {
	ReloadText()
	lineIdx = 0
}

func ReloadText() {
	bytes, err := ioutil.ReadFile(LINES_FILE_PATH)
	if err != nil {
		bytes = []byte{}
	}

	txt := regexp.MustCompile(`(?m)(^\s*$|^\s+|\s+$)`).ReplaceAllLiteralString(string(bytes), "")
	txt = regexp.MustCompile(`\n+`).ReplaceAllLiteralString(txt, "\n")
	txt = strings.TrimSpace(txt)

	splittedLines := strings.Split(txt, "\n")

	linesMutex.Lock()
	defer linesMutex.Unlock()
	lines = splittedLines
	lineIdx = len(lines) - 1
}

func GetRandomLine() string {
	if len(lines) == 0 {
		return "<no text>"
	}
	linesMutex.Lock()
	defer linesMutex.Unlock()
	lineIdx := rand.Int31n(int32(len(lines)))
	return lines[lineIdx]
}

func GetNextLine() string {
	if len(lines) == 0 {
		return "<no text>"
	}
	linesMutex.Lock()
	defer linesMutex.Unlock()
	result := lines[lineIdx]
	lineIdx = (lineIdx + 1) % len(lines)
	return result
}
