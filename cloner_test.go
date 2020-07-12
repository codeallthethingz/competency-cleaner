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
	err := CloneRepo()
	if err != nil {
		t.Fatal(err)
	}
	competenciesDir, err := os.Stat("/tmp/competencies/competencies")
	if err != nil {
		t.Fatal(err)
	}
	if !competenciesDir.IsDir() {
		t.Fatal("competencies directory in searchspring/competencies repo is not a directory")
	}
	contents, err := ioutil.ReadDir("/tmp/competencies/competencies")
	if err != nil {
		t.Fatal(err)
	}
	if len(contents) < 20 {
		t.Fatal("competencies directory in searchspring/competencies repo has less than 20 competencies")
	}
}

func TestProcessCompetencies(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "compcleanertest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	err = os.MkdirAll(dir+"/competencies/competencies", 0777)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(dir+"/competencies/competencies/github.md", []byte("# Github Competency\ngithub is great\n## how do you prove it?\nmake a branch\n## How do you improve it?\nread the docs"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(dir+"/competencies/competencies/mentorship.md", []byte("# "), 0644)
	if err != nil {
		t.Fatal(err)
	}
	good, bad, err := ProcessCompetencies(dir + "/competencies/competencies/")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, 1, len(good))
	require.Equal(t, "Github", good[0].Title)
	require.Equal(t, "github is great", good[0].Levels[0].Summary)
	require.Equal(t, "make a branch", good[0].Levels[0].Prove)
	require.Equal(t, "read the docs", good[0].Levels[0].Improve)
	require.Equal(t, 1, len(bad))
	require.Equal(t, 4, len(bad[0].Reasons))
	require.Equal(t, "no title", bad[0].Reasons[0])
	require.Equal(t, "level 1 missing improve", bad[0].Reasons[1])
	require.Equal(t, "level 1 missing prove", bad[0].Reasons[2])
	require.Equal(t, "level 1 must have summary", bad[0].Reasons[3])
}

func TestValidate(t *testing.T) {
	reasons := validateCompetency(&CompetencyDocument{})
	require.Equal(t, "no levels", reasons[1])
}
