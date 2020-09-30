package plugins

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/virink/vulwarning/common"
	"github.com/virink/vulwarning/model"
)

/**
Search API
https://docs.github.com/en/rest/reference/search
Rate limit: 10 requests per minute.
*/

// GithubSearchResult -
type GithubSearchResult struct {
	Items []struct {
		Description string `json:"description"`
		UpdatedAt   string `json:"updated_at"`
		CreatedAt   string `json:"created_at"`
		FullName    string `json:"full_name"`
		Name        string `json:"name"`
		SvnURL      string `json:"svn_url"`
	} `json:"items"`
	TotalCount int64 `json:"total_count"`
}

// PluginGithubCVE -
type PluginGithubCVE struct {
	res []*model.Warning
}

// Result -
func (p *PluginGithubCVE) Result() []*model.Warning {
	return p.res
}

// Crawl -
func (p *PluginGithubCVE) Crawl() error {
	target := fmt.Sprintf(
		"https://api.github.com/search/repositories?q=CVE-%d&sort=updated&per_page=10",
		time.Now().Year(),
	)
	data, err := httpGet(target)
	if err != nil {
		return err
	}
	var gsr GithubSearchResult
	err = json.Unmarshal(data, &gsr)
	if err != nil {
		return err
	}
	for _, item := range gsr.Items {
		p.res = append(p.res, &model.Warning{
			Title:    "Found [{CVE}] on GitHub",
			Link:     item.SvnURL,
			From:     "githubcve",
			Desc:     item.Description,
			Time:     time.Unix(ParsePubDate(item.UpdatedAt), 0),
			CreateAt: time.Now(),
		})
		common.Logger.Debugln("Crwaled [GithubCVE]", item.FullName, item.UpdatedAt)
	}
	return nil
}
