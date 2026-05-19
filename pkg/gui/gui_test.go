package gui

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/stretchr/testify/assert"
)

func TestBreakingChangesSinceExcludesFutureReleases(t *testing.T) {
	lastVersion := &types.VersionNumber{Major: 0, Minor: 61, Patch: 1}
	currentVersion := &types.VersionNumber{Major: 0, Minor: 61, Patch: 1}
	notesByVersion := map[string]string{
		"0.61.1": "current",
		"0.62.0": "future",
	}

	assert.Empty(t, breakingChangesSince(lastVersion, currentVersion, notesByVersion))
}

func TestBreakingChangesSinceIncludesReleasesUpToCurrentVersion(t *testing.T) {
	lastVersion := &types.VersionNumber{Major: 0, Minor: 60, Patch: 0}
	currentVersion := &types.VersionNumber{Major: 0, Minor: 62, Patch: 0}
	notesByVersion := map[string]string{
		"0.61.0": "first",
		"0.62.0": "second",
		"0.63.0": "future",
	}

	assert.Equal(t, []string{"first", "second"}, breakingChangesSince(lastVersion, currentVersion, notesByVersion))
}
