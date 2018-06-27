/*
 * Revision History:
 *     Initial: 2018/06/20        Li Zebang
 */

package github

import (
	"context"
	"sync"

	"github.com/google/go-github/github"
)

// RepositoryInformation -
type RepositoryInformation struct {
	Owner       *string    `json:"owner"`
	Avatar      *string    `json:"avatar"`
	Repo        *string    `json:"repo"`
	Description *string    `json:"description"`
	Topics      []string   `json:"topics"`
	Languages   []language `json:"languages"`
	Readme      *string    `json:"readme"`
}

type language struct {
	Language   string `json:"language"`
	Proportion int    `json:"proportion"`
}

type repositoriesGetReturn struct {
	repository *github.Repository
	response   *github.Response
	err        error
}

func (cp ChannelPool) repositoriesGet(ctx context.Context, tag, owner, repo string, wg *sync.WaitGroup, ret *repositoriesGetReturn) {
	defer wg.Done()
	client, err := cp.Get(tag)
	if err != nil {
		ret.err = err
		return
	}
	defer cp.Put(client)

	repository, response, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		ret.err = err
		return
	}

	if response != nil {
		err = client.HandleResponse(response)
		if err != nil {
			ret.err = err
			return
		}
	}

	ret.repository = repository
	ret.response = response
}

type repositoriesGetReadmeReturn struct {
	repositoryContent *github.RepositoryContent
	response          *github.Response
	err               error
}

func (cp ChannelPool) repositoriesGetReadme(ctx context.Context, tag, owner, repo string, opt *github.RepositoryContentGetOptions, wg *sync.WaitGroup, ret *repositoriesGetReadmeReturn) {
	defer wg.Done()
	client, err := cp.Get(tag)
	if err != nil {
		ret.err = err
		return
	}
	defer cp.Put(client)

	repositoryContent, response, err := client.Repositories.GetReadme(ctx, owner, repo, opt)
	if err != nil {
		ret.err = err
		return
	}

	if response != nil {
		err = client.HandleResponse(response)
		if err != nil {
			ret.err = err
			return
		}
	}

	ret.repositoryContent = repositoryContent
	ret.response = response
}

type repositoriesListAllTopicsReturn struct {
	topics   []string
	response *github.Response
	err      error
}

func (cp ChannelPool) repositoriesListAllTopics(ctx context.Context, tag, owner, repo string, wg *sync.WaitGroup, ret *repositoriesListAllTopicsReturn) {
	defer wg.Done()
	client, err := cp.Get(tag)
	if err != nil {
		ret.err = err
		return
	}
	defer cp.Put(client)

	topics, response, err := client.Repositories.ListAllTopics(ctx, owner, repo)
	if err != nil {
		ret.err = err
		return
	}

	if response != nil {
		err = client.HandleResponse(response)
		if err != nil {
			ret.err = err
			return
		}
	}

	ret.topics = topics
	ret.response = response
}

type repositoriesListLanguagesReturn struct {
	languages map[string]int
	response  *github.Response
	err       error
}

func (cp ChannelPool) repositoriesListLanguages(ctx context.Context, tag, owner, repo string, wg *sync.WaitGroup, ret *repositoriesListLanguagesReturn) {
	defer wg.Done()
	client, err := cp.Get(tag)
	if err != nil {
		ret.err = err
		return
	}
	defer cp.Put(client)

	languages, response, err := client.Repositories.ListLanguages(ctx, owner, repo)
	if err != nil {
		ret.err = err
		return
	}

	if response != nil {
		err = client.HandleResponse(response)
		if err != nil {
			ret.err = err
			return
		}
	}

	ret.languages = languages
	ret.response = response
}

// GetRepositoryInformation -
func (cp ChannelPool) GetRepositoryInformation(tag, owner, repo string) (*RepositoryInformation, error) {
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	wg.Add(4)

	var (
		getReturn           = &repositoriesGetReturn{}
		getReadmeReturn     = &repositoriesGetReadmeReturn{}
		listAllTopicsReturn = &repositoriesListAllTopicsReturn{}
		listLanguagesReturn = &repositoriesListLanguagesReturn{}
	)
	go cp.repositoriesGet(ctx, tag, owner, repo, wg, getReturn)
	go cp.repositoriesGetReadme(ctx, tag, owner, repo, nil, wg, getReadmeReturn)
	go cp.repositoriesListAllTopics(ctx, tag, owner, repo, wg, listAllTopicsReturn)
	go cp.repositoriesListLanguages(ctx, tag, owner, repo, wg, listLanguagesReturn)

	wg.Wait()

	if getReturn.err != nil {
		return nil, getReturn.err
	}
	if getReadmeReturn.err != nil {
		return nil, getReadmeReturn.err
	}
	if listAllTopicsReturn.err != nil {
		return nil, listAllTopicsReturn.err
	}
	if listLanguagesReturn.err != nil {
		return nil, listLanguagesReturn.err
	}

	var (
		ls    = make([]language, len(listLanguagesReturn.languages))
		sum   int
		index int
	)
	for key, val := range listLanguagesReturn.languages {
		ls[index] = language{
			Language:   key,
			Proportion: val,
		}
		sum += val
		index++
	}
	for index := range ls {
		ls[index].Proportion /= sum
	}

	return &RepositoryInformation{
		Owner:       &owner,
		Avatar:      getReturn.repository.Owner.AvatarURL,
		Repo:        &repo,
		Description: getReturn.repository.Description,
		Topics:      listAllTopicsReturn.topics,
		Readme:      getReadmeReturn.repositoryContent.DownloadURL,
	}, nil
}
