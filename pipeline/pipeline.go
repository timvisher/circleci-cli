package pipeline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/CircleCI-Public/circleci-cli/git"
	"github.com/CircleCI-Public/circleci-cli/paths"
	"github.com/CircleCI-Public/circleci-cli/version"
	"github.com/pkg/errors"
)

// CircleCI provides various `<< pipeline.x >>` values to be used in your config, but sometimes we need to fabricate those values when validating config.
type Values map[string]string

func FabricatedValues() Values {

	revision := git.Revision()
	gitUrl := "https://github.com/CircleCI-Public/circleci-cli"
	projectType := "github"

	// If we encounter an error infering project, skip this and use defaults.
	if remote, err := git.InferProjectFromGitRemotes(); err == nil {
		switch remote.VcsType {
		case git.GitHub:
			gitUrl = fmt.Sprintf("https://github.com/%s/%s", remote.Organization, remote.Project)
			projectType = "github"
		case git.Bitbucket:
			gitUrl = fmt.Sprintf("https://bitbucket.org/%s/%s", remote.Organization, remote.Project)
			projectType = "bitbucket"
		}
	}

	return map[string]string{
		"id":                "00000000-0000-0000-0000-000000000001",
		"number":            "1",
		"project.git_url":   gitUrl,
		"project.type":      projectType,
		"git.tag":           git.Tag(),
		"git.branch":        git.Branch(),
		"git.revision":      revision,
		"git.base_revision": revision,
	}
}

// TODO: type Parameters map[string]string

// KeyVal is a data structure specifically for passing pipeline data to GraphQL which doesn't support free-form maps.
type KeyVal struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

// PrepareForGraphQL takes a golang homogenous map, and transforms it into a list of keyval pairs, since GraphQL does not support homogenous maps.
func PrepareForGraphQL(kvMap Values) []KeyVal {
	// we need to create the slice of KeyVals in a deterministic order for testing purposes
	keys := make([]string, 0, len(kvMap))
	for k := range kvMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	kvs := make([]KeyVal, 0, len(kvMap))
	for _, k := range keys {
		kvs = append(kvs, KeyVal{Key: k, Val: kvMap[k]})
	}
	return kvs
}

func Trigger(token string) error {

	remote, err := git.InferProjectFromGitRemotes()

	if err != nil {
		return errors.Wrap(err, "this command must be run from inside a git repositopry")
	}

	url := fmt.Sprintf("https://circleci.com/api/v2/project/%s/%s/%s/pipeline",
		strings.ToLower(string(remote.VcsType)),
		remote.Organization,
		remote.Project)

	parameters := map[string]string{
		"branch": "422-windows-installer",
	}

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(parameters); err != nil {
		return errors.Wrap(err, "Failed to encode JSON body for POST")
	}

	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Circle-Token", token)
	req.Header.Set("User-Agent", version.UserAgent())

	client := &http.Client{}
	resp, err := client.Do(req)

	var response struct {
		Number int
		State  string
		Id     string
	}

	if err != nil {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return errors.Wrap(err, "json decode error")
	}

	fmt.Printf("Pipeline %d created\n", response.Number)

	fmt.Println(paths.ProjectUrl(remote))

	return nil
}
