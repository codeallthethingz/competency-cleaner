package compclean

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func Drive() error {
	dir, err := CloneRepo()
	if err != nil {
		return err
	}
	good, bad, err := ProcessCompetencies(dir)
	if err != nil {
		return err
	}
	if len(bad) > 0 {
		msg := ""
		for _, b := range bad {
			msg += fmt.Sprintf("https://github.com/searchspring/competencies/blob/master/%s    %s\n", b.Document.Path, b.Reasons)
		}
		return fmt.Errorf("Found non-zero number of bad competencies: %d bad\n%s", len(bad), msg)
	}
	content, err := json.Marshal(good)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("competencies.json", content, 0644)
	if err != nil {
		return err
	}
	return nil
}
