package utils

import (
	"bufio"
	"strings"

	"github.com/jala-R/VideoAutomatorGUI/packages/GUI/model"
)

func MarshallScript(script string) [][]string {
	var parsedScript = [][]string{}
	var paraSentences = []string{}
	reader := strings.NewReader(script)
	scriptReader := bufio.NewReader(reader)

	for line, err := scriptReader.ReadString('\n'); err == nil || len(line) != 0; line, err = scriptReader.ReadString('\n') {
		if len(line) == 0 || line[0] == '\n' {
			continue
		}
		if line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}
		if isParaEnding(line) {
			if len(paraSentences) != 0 {
				parsedScript = append(parsedScript, paraSentences)
				paraSentences = []string{}
			}
		} else {
			paraSentences = append(paraSentences, line)
		}
	}

	if len(paraSentences) != 0 {
		parsedScript = append(parsedScript, paraSentences)
	}

	return parsedScript
}

func UnmarshallScript(script [][]string) string {
	var out = ""

	for _, para := range script {
		for _, sentence := range para {
			out += (sentence + "\n")
		}
		out += (model.EOL) + "\n"
	}

	return out
}

func isParaEnding(line string) bool {
	return line == model.EOL
}
