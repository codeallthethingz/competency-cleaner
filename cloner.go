package compclean

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
)

func CloneRepo() (string, error) {
	os.RemoveAll("/tmp/competencies")
	_, err := git.PlainClone("/tmp/competencies", false, &git.CloneOptions{
		URL: "https://github.com/searchspring/competencies",
	})

	return "/tmp/competencies/competencies/", err
}

type BadCompetencyDocument struct {
	Document *CompetencyDocument
	Reasons  []string
}

func ProcessCompetencies(dir string) ([]*CompetencyDocument, []*BadCompetencyDocument, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}
	good := []*CompetencyDocument{}
	bad := []*BadCompetencyDocument{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			contents, err := ioutil.ReadFile(dir + file.Name())
			if err != nil {
				return nil, nil, err
			}
			doc, err := Convert(string(contents))
			if err != nil {
				return nil, nil, err
			}
			doc.Path = "competencies/" + file.Name()
			reasons := validateCompetency(doc)
			if len(reasons) == 0 {
				good = append(good, doc)
			} else {
				bad = append(bad, &BadCompetencyDocument{Document: doc, Reasons: reasons})
			}

		}
	}
	return good, bad, nil
}

func validateCompetency(doc *CompetencyDocument) []string {
	reasons := []string{}
	if doc.Title == "" {
		reasons = append(reasons, "no title")
	}
	if len(doc.Levels) == 0 {
		reasons = append(reasons, "no levels")
	}
	for i, level := range doc.Levels {
		if level.Improve == "" {
			reasons = append(reasons, fmt.Sprintf("level %d missing improve", i+1))
		}
		if level.Prove == "" {
			reasons = append(reasons, fmt.Sprintf("level %d missing prove", i+1))
		}
		if i == 0 && level.Summary == "" {
			reasons = append(reasons, fmt.Sprintf("level 1 must have summary"))
		}

	}
	return reasons
}
