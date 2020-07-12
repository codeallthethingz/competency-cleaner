package compclean

import (
	"fmt"
	"regexp"
	"strings"
)

type CompetencyDocument struct {
	Title       string   `json:"title"`
	TitleSearch string   `json:"titleSearch"`
	Levels      []*Level `json:"levels"`
}

type Level struct {
	Prove   string `json:"prove"`
	Improve string `json:"improve"`
	Summary string `json:"summary"`
}

func Convert(text string) (*CompetencyDocument, error) {
	doc := &CompetencyDocument{
		Levels: []*Level{},
	}
	lines := strings.Split(text, "\n")

	title, err := findFirstNonBlankLine(lines)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(strings.TrimSpace(strings.ToLower(title)), "competency") {
		index := strings.LastIndex(strings.ToLower(title), "competency")
		title = title[0:index]
	}
	if strings.HasPrefix(title, "#") {
		title = strings.TrimSpace(title[1:])
	}
	doc.Title = title
	doc.TitleSearch = strings.ToLower(title)

	currentMatcher := &Level{}
	matching := false
	matchingProve := false
	matchingImprove := false
	mainHeadingRegex := regexp.MustCompile(`^#[^#].*$`)
	proveRegex := regexp.MustCompile(`(?i)##[^#]*how do you prove it.*`)
	improveRegex := regexp.MustCompile(`(?i)##[^#]*how do you improve it.*`)
	for _, line := range lines {
		fmt.Print("\n" + line + " --> ")
		if mainHeadingRegex.Match([]byte(line)) && matching {
			doc.Levels = append(doc.Levels, currentMatcher)
			currentMatcher = &Level{}
			fmt.Print("matched header")
			matching = true
			matchingProve = false
			matchingImprove = false
			continue
		}
		if mainHeadingRegex.Match([]byte(line)) {
			fmt.Print("matched header")
			matching = true
			matchingProve = false
			matchingImprove = false
			continue
		}
		if matching && proveRegex.Match([]byte(line)) {
			fmt.Print("matched prove")
			matchingProve = true
			matchingImprove = false
			continue
		}
		if matching && improveRegex.Match([]byte(line)) {
			fmt.Print("matched improve")
			matchingProve = false
			matchingImprove = true
			continue
		}
		if matching && matchingProve {
			fmt.Print("prove line")
			currentMatcher.Prove += line + "\n"
			continue
		}
		if matching && matchingImprove {
			fmt.Print("improve line")
			currentMatcher.Improve += line + "\n"
			continue
		}
		if matching {
			fmt.Print("summary line")
			currentMatcher.Summary += line + "\n"
		}
	}
	doc.Levels = append(doc.Levels, currentMatcher)
	for _, level := range doc.Levels {
		level.Improve = strings.TrimSpace(level.Improve)
		level.Prove = strings.TrimSpace(level.Prove)
		level.Summary = strings.TrimSpace(level.Summary)
	}
	return doc, nil
}

func findFirstNonBlankLine(lines []string) (string, error) {
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			return line, nil
		}
	}
	return "", fmt.Errorf("blank file")
}
