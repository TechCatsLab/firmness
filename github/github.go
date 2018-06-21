/*
 * Revision History:
 *     Initial: 2018/06/20        Li Zebang
 */

// todo: rate limmit, oauth
package github

import (
	"context"
	"errors"
	"sort"

	"github.com/google/go-github/github"
)

var client = github.NewClient(nil)

// RepositoryInformation -
type RepositoryInformation struct {
	Owner       *string  `json:"owner"`
	AvatarURL   *string  `json:"avatar_url"`
	Repo        *string  `json:"repo"`
	Description *string  `json:"description"`
	Languages   []string `json:"languages"`
	Readme      *string  `json:"readme"`
}

// GetRepositoryInformation -
func GetRepositoryInformation(owner, repo string, n int) (*RepositoryInformation, error) {
	ctx := context.Background()

	repository, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	readme, _, err := client.Repositories.GetReadme(ctx, owner, repo, nil)
	if err != nil {
		return nil, err
	}

	_, ls, err := TopNLanguage(owner, repo, n)
	if err != nil {
		return nil, err
	}

	return &RepositoryInformation{
		Owner:       &owner,
		AvatarURL:   repository.Owner.AvatarURL,
		Repo:        &repo,
		Description: repository.Description,
		Languages:   ls,
		Readme:      readme.DownloadURL,
	}, nil
}

type language struct {
	language string
	bytes    int
}

type languages []language

func (l languages) Len() int           { return len(l) }
func (l languages) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l languages) Less(i, j int) bool { return l[i].bytes > l[j].bytes }

func (l languages) toStringSlice() []string {
	ls := make([]string, len(l))

	for key := range l {
		ls[key] = l[key].language
	}

	return ls
}

// TopNLanguage -
func TopNLanguage(owner, repo string, n int) (int, []string, error) {
	ctx := context.Background()

	if n < 1 {
		return 0, nil, errors.New("n can not be less than 1")
	}

	l, _, err := client.Repositories.ListLanguages(ctx, owner, repo)
	if err != nil {
		return 0, nil, err
	}

	index := 0
	ls := languages(make([]language, len(l)))
	for key, val := range l {
		ls[index] = language{
			language: key,
			bytes:    val,
		}
		index++
	}
	sort.Sort(ls)

	if len(ls) > n {
		return n, ls.toStringSlice()[:n], nil
	}

	return len(ls), ls.toStringSlice(), nil
}
