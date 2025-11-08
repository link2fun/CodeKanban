package git

import (
	"reflect"
	"testing"
)

func TestParseWorktreeList(t *testing.T) {
	input := `worktree /repo
HEAD 40c7e09a8e6d05218f61d07eb05b525f49f302e5
branch refs/heads/main

worktree /repo/feature
HEAD 5da41358595c294c5b4af4a3e163192f7ca2ce50
branch refs/heads/feature/demo
`

	got := parseWorktreeList(input)
	want := []WorktreeInfo{
		{
			Path:       "/repo",
			Branch:     "main",
			HeadCommit: "40c7e09",
		},
		{
			Path:       "/repo/feature",
			Branch:     "feature/demo",
			HeadCommit: "5da4135",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected worktree list: %#v", got)
	}
}
