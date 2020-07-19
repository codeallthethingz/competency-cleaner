package compclean

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessor(t *testing.T) {
	doc, err := Convert("\n#Github Competency  \nfirst\nsummary\n\n## How do you prove it?\nprove it\ntext\n## How do you improve it?\nimprove text\n# Github - Level 2\nlevel 2 summary\n##How do you prove it?\nprove 2 text\n## How do you improve it?\nimprove 2 text\n")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, "Github", doc.Title)
	require.Equal(t, "github", doc.TitleSearch)
	require.Equal(t, 2, len(doc.Levels))
	require.Equal(t, "prove it\ntext", doc.Levels[0].Prove)
	require.Equal(t, "improve text", doc.Levels[0].Improve)
	require.Equal(t, "first\nsummary", doc.Levels[0].Summary)
	require.Equal(t, "prove 2 text", doc.Levels[1].Prove)
	require.Equal(t, "improve 2 text", doc.Levels[1].Improve)
	require.Equal(t, "level 2 summary", doc.Levels[1].Summary)
}

func TestPrefixRemoved(t *testing.T) {
	doc, err := Convert("\n#Competency Github Competency  \nfirst\nsummary\n\n## How do you prove it?\nprove it\ntext\n## How do you improve it?\nimprove text\n# Github - Level 2\nlevel 2 summary\n##How do you prove it?\nprove 2 text\n## How do you improve it?\nimprove 2 text\n")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, "Github", doc.Title)
	doc, err = Convert("\n#Competency - Github Competency  \nfirst\nsummary\n\n## How do you prove it?\nprove it\ntext\n## How do you improve it?\nimprove text\n# Github - Level 2\nlevel 2 summary\n##How do you prove it?\nprove 2 text\n## How do you improve it?\nimprove 2 text\n")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, "Github", doc.Title)
}

func TestNoSummary(t *testing.T) {
	doc, err := Convert("\n#Github Competency\n\n## How do you prove it?\nprove it\ntext\n## How do you improve it?\nimprove text\n# Github - Level 2\n##How do you prove it?\nprove 2 text\n## How do you improve it?\nimprove 2 text\n")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, "Github", doc.Title)
	require.Equal(t, "github", doc.TitleSearch)
	require.Equal(t, 2, len(doc.Levels))
	require.Equal(t, "prove it\ntext", doc.Levels[0].Prove)
	require.Equal(t, "improve text", doc.Levels[0].Improve)
	require.Equal(t, "", doc.Levels[0].Summary)
	require.Equal(t, "prove 2 text", doc.Levels[1].Prove)
	require.Equal(t, "improve 2 text", doc.Levels[1].Improve)
	require.Equal(t, "", doc.Levels[1].Summary)
}

func TestBlankFile(t *testing.T) {
	_, err := Convert("   \n\t\n")
	require.NotNil(t, err, "expected error with blank file")
}
