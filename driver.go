package compclean

import "fmt"

func Drive() error {
	dir, err := CloneRepo()
	if err != nil {
		return err
	}
	_, bad, err := ProcessCompetencies(dir)
	if err != nil {
		return err
	}
	if len(bad) > 0 {
		msg := ""
		for _, b := range bad {
			msg += fmt.Sprintf("https://github.com/searchspring/competencies/blob/master/%s    %s\n", b.Document.ID, b.Reasons)
		}
		return fmt.Errorf("Found non-zero number of bad competencies: %d bad\n%s", len(bad), msg)
	}

	return nil
}
