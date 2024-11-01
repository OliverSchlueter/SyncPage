package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Repository struct {
	Owner     string `json:"owner"`
	Repo      string `json:"repo"`
	AuthToken string `json:"auth_token"`
}

type WorkflowRun struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Conclusion   string `json:"conclusion"`
	HTMLURL      string `json:"html_url"`
	DisplayTitle string `json:"display_title"`
}

type Artifact struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type WorkflowRunsResponse struct {
	WorkflowRuns []WorkflowRun `json:"workflow_runs"`
}

type ArtifactsResponse struct {
	Artifacts []Artifact `json:"artifacts"`
}

func (r *Repository) GetLatestWorkflowRun(workflowName string) (*WorkflowRun, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs", r.Owner, r.Repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if r.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+r.AuthToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var runsResponse WorkflowRunsResponse
	if err := json.NewDecoder(resp.Body).Decode(&runsResponse); err != nil {
		return nil, err
	}

	for _, run := range runsResponse.WorkflowRuns {
		if run.Name == workflowName {
			return &run, nil
		}
	}

	return nil, fmt.Errorf("no workflow run found with name: %s", workflowName)
}

func (r *Repository) GetArtifacts(run *WorkflowRun) ([]Artifact, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs/%d/artifacts", r.Owner, r.Repo, run.ID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if r.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+r.AuthToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artifactsResponse ArtifactsResponse
	if err := json.NewDecoder(resp.Body).Decode(&artifactsResponse); err != nil {
		return nil, err
	}

	return artifactsResponse.Artifacts, nil
}

func (a *Artifact) Download(repo *Repository) ([]byte, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/artifacts/%d/zip", repo.Owner, repo.Repo, a.ID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if repo.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+repo.AuthToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}
