package core

import (
	"hash/fnv"
	"math/rand"
)

const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	ColorGray    = "\033[90m"
)

var colors = []string{
	ColorRed, ColorGreen, ColorYellow,
	ColorBlue, ColorMagenta, ColorCyan,
	ColorWhite, ColorGray,
}

var patterns = []string{
	`
  ██  
 ████ 
██████
██████
 ████ 
  ██  `,
	`
██████
█    █
█    █
█    █
█    █
██████`,
	`
█    █
██  ██
██████
██████
██  ██
█    █`,
	`
 █████
█     █
███████
█     █
█     █
██████ `,
	`
  ██  
 █ █ █
█  █  █
█  █  █
 █ █ █
  ██  `,
	`
██████
     █
██████
█     
██████
██████`,
}

func GeneratePFP(feedRef string) string {
	h := fnv.New32a()
	h.Write([]byte(feedRef))
	seed := int64(h.Sum32())

	rg := rand.New(rand.NewSource(seed))

	patternIdx := rg.Intn(len(patterns))
	pattern := patterns[patternIdx]

	fgIdx := rg.Intn(len(colors))
	fg := colors[fgIdx]

	bgIdx := (fgIdx + rg.Intn(len(colors)-1) + 1) % len(colors)
	bg := colors[bgIdx]

	return renderColoredPFP(pattern, fg, bg)
}

func renderColoredPFP(pattern, fg, bg string) string {
	lines := splitLines(pattern)[1:]

	var result string
	for _, line := range lines {
		coloredLine := ""
		for _, ch := range line {
			if ch == '█' {
				coloredLine += fg + string(ch)
			} else {
				coloredLine += bg + " "
			}
		}
		coloredLine += ColorReset
		result += coloredLine + "\n"
	}

	return result
}

func splitLines(s string) []string {
	var lines []string
	var line string
	for _, ch := range s {
		if ch == '\n' {
			lines = append(lines, line)
			line = ""
		} else {
			line += string(ch)
		}
	}
	if len(line) > 0 {
		lines = append(lines, line)
	}
	return lines
}
