package server

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strings"
	"sync"
)

var lines []string
var linesMutex sync.Mutex

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
}

func GetRandomLine() string {
	if len(lines) == 0 {
		return "<no text>"
	}
	linesMutex.Lock()
	defer linesMutex.Unlock()
	i := rand.Int31n(int32(len(lines)))
	fmt.Println("i", i)
	return lines[i]
}
