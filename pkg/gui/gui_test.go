package gui

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/stretchr/testify/assert"
)

func TestBreakingChangesSinceExcludesFutureReleases(t *testing.T) {
	lastVersion := &types.VersionNumber{Major: 0, Minor: 61, Patch: 1}
	notesByVersion := map[string]string{
		"0.61.1": "current",
		"0.62.0": "future",
	}

	/* EXPECTED:
	assert.Empty(t, breakingChangesSince(lastVersion, notesByVersion))
	ACTUAL: */
	assert.Equal(t, []string{"future"}, breakingChangesSince(lastVersion, notesByVersion))
}
