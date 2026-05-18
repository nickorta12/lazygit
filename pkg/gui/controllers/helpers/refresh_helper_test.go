package helpers

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/hosting_service"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestGetGithubBaseRemote(t *testing.T) {
	cases := []struct {
		name             string
		githubRemotes    []githubRemoteInfo
		configuredRemote string
		expected         string
	}{
		{
			name:             "configured remote wins",
			githubRemotes:    makeGithubRemoteInfoList("origin", "upstream", "fork"),
			configuredRemote: "fork",
			expected:         "fork",
		},
		{
			name:             "configured remote not in github remotes returns nil",
			githubRemotes:    makeGithubRemoteInfoList("origin"),
			configuredRemote: "missing",
			expected:         "",
		},
		{
			name:             "single github remote is auto-picked",
			githubRemotes:    makeGithubRemoteInfoList("myremote"),
			configuredRemote: "",
			expected:         "myremote",
		},
		{
			name:             "upstream is preferred when multiple github remotes exist",
			githubRemotes:    makeGithubRemoteInfoList("origin", "upstream", "fork"),
			configuredRemote: "",
			expected:         "upstream",
		},
		{
			name:             "no upstream and multiple remotes returns nil",
			githubRemotes:    makeGithubRemoteInfoList("origin", "fork"),
			configuredRemote: "",
			expected:         "",
		},
		{
			name:             "empty list returns nil",
			githubRemotes:    nil,
			configuredRemote: "",
			expected:         "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := getGithubBaseRemote(c.githubRemotes, c.configuredRemote)
			if c.expected == "" {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, c.expected, result.remote.Name)
			}
		})
	}
}

func TestGetAuthenticatedPullRequestRemotes(t *testing.T) {
	githubRemotes := []githubRemoteInfo{
		makeGithubRemoteInfo("origin", "github", "github.com"),
		makeGithubRemoteInfo("fork", "github", "github.com"),
		makeGithubRemoteInfo("enterprise", "github", "ghe.example.com"),
		makeGithubRemoteInfo("missing-auth", "github", "no-token.example.com"),
		makeGithubRemoteInfo("gitea", "gitea", "gitea.example.com"),
	}

	callsByHost := map[string]int{}
	result := getAuthenticatedPullRequestRemotes(githubRemotes, func(provider string, host string) string {
		callsByHost[host]++
		switch host {
		case "github.com":
			return "github-token"
		case "ghe.example.com":
			return "ghe-token"
		default:
			return ""
		}
	})

	assert.Equal(t, []githubRemoteInfo{
		makeAuthenticatedGithubRemoteInfo("origin", "github", "github.com", "github-token"),
		makeAuthenticatedGithubRemoteInfo("fork", "github", "github.com", "github-token"),
		makeAuthenticatedGithubRemoteInfo("enterprise", "github", "ghe.example.com", "ghe-token"),
		makeAuthenticatedGithubRemoteInfo("gitea", "gitea", "gitea.example.com", ""),
	}, result)
	// Two remotes share github.com; the lookup runs only once.
	assert.Equal(t, map[string]int{
		"github.com":           1,
		"ghe.example.com":      1,
		"no-token.example.com": 1,
		"gitea.example.com":    1,
	}, callsByHost)
}

func makeGithubRemoteInfoList(names ...string) []githubRemoteInfo {
	return lo.Map(names, func(name string, _ int) githubRemoteInfo {
		return makeGithubRemoteInfo(name, "github", name)
	})
}

func makeGithubRemoteInfo(name string, provider string, webDomain string) githubRemoteInfo {
	return githubRemoteInfo{
		remote: &models.Remote{Name: name},
		serviceInfo: hosting_service.ServiceInfo{
			Provider:  provider,
			RepoName:  name,
			WebDomain: webDomain,
		},
	}
}

func makeAuthenticatedGithubRemoteInfo(name string, provider string, webDomain string, authToken string) githubRemoteInfo {
	info := makeGithubRemoteInfo(name, provider, webDomain)
	info.authToken = authToken
	return info
}
