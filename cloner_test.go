package compclean

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepoContract(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	dir, err := CloneRepo()
	if err != nil {
		t.Fatal(err)
	}
	competenciesDir, err := os.Stat(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !competenciesDir.IsDir() {
		t.Fatal("competencies directory in searchspring/competencies repo is not a directory")
	}
	contents, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(contents) < 20 {
		t.Fatal("competencies directory in searchspring/competencies repo has less than 20 competencies")
	}
}

func TestProcessCompetencies(t *testing.T) {
	dir := createTempCompetencies(false)
	defer os.RemoveAll(dir)

	good, bad, err := ProcessCompetencies(dir + "/competencies/competencies/")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, 1, len(good))
	require.Equal(t, "competencies/github.md", good[0].Path)
	require.Equal(t, "Github", good[0].Title)
	require.Equal(t, "github is great", good[0].Levels[0].Summary)
	require.Equal(t, "make a branch", good[0].Levels[0].Prove)
	require.Equal(t, "read the docs", good[0].Levels[0].Improve)
	require.Equal(t, 1, len(bad))
	require.Equal(t, 4, len(bad[0].Reasons))
	require.Equal(t, "competencies/mentorship.md", bad[0].Document.Path)
	require.Equal(t, "no title", bad[0].Reasons[0])
	require.Equal(t, "level 1 missing improve", bad[0].Reasons[1])
	require.Equal(t, "level 1 missing prove", bad[0].Reasons[2])
	require.Equal(t, "level 1 must have summary", bad[0].Reasons[3])
}

func TestValidate(t *testing.T) {
	reasons := validateCompetency(&CompetencyDocument{})
	require.Equal(t, "no levels", reasons[1])
}

func TestWithChannel(t *testing.T) {
	dir := createTempCompetencies(true)
	defer os.RemoveAll(dir)

	good := make(chan *CompetencyDocument)
	bad := make(chan *BadCompetencyDocument)
	errChan := make(chan error)
	quit := make(chan int)
	go func() {
		ProcessCompetenciesWithChannel(dir+"/competencies/competencies/", good, bad, errChan, quit)
	}()
	goodCount, badCount, errCount := 0, 0, 0
	finished := false
	for !finished {
		select {
		case <-good:
			goodCount++
		case <-bad:
			badCount++
		case <-errChan:
			errCount++
		case <-quit:
			finished = true
		}
	}
	require.Equal(t, 1, goodCount)
	require.Equal(t, 1, badCount)
	require.Equal(t, 1, errCount)
}

func createTempCompetencies(createEmpty bool) string {
	dir, err := ioutil.TempDir("/tmp", "compcleanertest")
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(dir+"/competencies/competencies", 0777)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(dir+"/competencies/competencies/github.md", []byte("# Github Competency\ngithub is great\n## how do you prove it?\nmake a branch\n## How do you improve it?\nread the docs"), 0644)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(dir+"/competencies/competencies/mentorship.md", []byte("# "), 0644)
	if err != nil {
		panic(err)
	}
	if createEmpty {
		err = ioutil.WriteFile(dir+"/competencies/competencies/zempty.md", []byte(""), 0644)
		if err != nil {
			panic(err)
		}
	}
	return dir
}
