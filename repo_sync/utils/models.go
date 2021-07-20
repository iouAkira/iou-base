package utils

type ReposConfig struct {
	Repos []struct {
		RepoName    string `json:"repo_name"`
		RepoURL     string `json:"repo_url"`
		RepoBranch  string `json:"repo_branch"`
		RepoPrivate bool   `json:"repo_private,omitempty"`
		GitAccount  string `json:"git_account,omitempty"`
		GitToken    string `json:"git_token,omitempty"`
	} `json:"repos"`
}
