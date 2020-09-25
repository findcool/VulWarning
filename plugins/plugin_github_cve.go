package plugins

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/gocolly/colly"
	"github.com/virink/vulWarning/common"
	"github.com/virink/vulWarning/model"
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
		URL         string `json:"url"`
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

// CrawlCVE -
func (p *PluginGithubCVE) CrawlCVE(cve string) (CVSS string, DESC string) {
	// target := fmt.Sprintf("https://nvd.nist.gov/vuln/detail/%s", cve)
	c := newCustomCollector([]string{"nvd.nist.gov"})
	c.OnRequest(func(r *colly.Request) {
		common.Logger.Debugln("Crawling [CVEDetail]", r.URL)
	})
	c.OnHTML("p[data-testid=vuln-description]", func(e *colly.HTMLElement) {
		common.Logger.Debugln("Crawling [CVEDetail]", e.Text)
		DESC = e.Text
	})
	c.OnHTML("a[data-testid=vuln-cvss3-panel-score]", func(e *colly.HTMLElement) {
		common.Logger.Debugln("Crawling [CVEDetail]", e.Text)
		CVSS = e.Text
	})
	c.Visit(fmt.Sprintf("https://nvd.nist.gov/vuln/detail/%s", cve))
	c.Wait()
	return
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
		var CVE, CVSS, DESC string
		// Is Exists
		if model.WarningIsExistsByLink(item.URL) {
			continue
		}
		// Get CVE Detail
		text := fmt.Sprintf("%s %s", item.FullName, item.Description)
		match := regexp.MustCompile(`(?mi)cve-?\s?(\d{4})-?(\d+)`).FindAllStringSubmatch(text, -1)
		if len(match) > 0 && len(match[0]) > 2 {
			CVE = fmt.Sprintf("CVE-%s-%s", match[0][1], match[0][2])
			CVSS, DESC = p.CrawlCVE(CVE)
			if len(DESC) > 0 {
				DESC = translate(DESC)
			}
		}
		text = fmt.Sprintf(
			"名称: %s\n描述: %s\n编号: %s\n等级: %s\n说明: %s",
			item.FullName,
			item.Description,
			CVE,
			CVSS,
			DESC,
		)
		p.res = append(p.res, &model.Warning{
			Title:    fmt.Sprintf("Found [%s] on GitHub", CVE),
			Link:     item.URL,
			Index:    item.URL,
			From:     "githubcve",
			Desc:     text,
			Time:     time.Unix(ParsePubDate(item.UpdatedAt), 0),
			CreateAt: time.Now(),
		})
		common.Logger.Debugln("Crwaled [GithubCVE]", item.FullName, item.UpdatedAt)
	}
	return nil
}
