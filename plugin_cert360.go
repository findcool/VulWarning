package main

import (
	"time"

	"github.com/gocolly/colly"
)

// PluginCert360 -
type PluginCert360 struct {
	c   *colly.Collector
	res []*Warings
}

// Result -
func (p *PluginCert360) Result() []*Warings {
	return p.res
}

// Crawl -
func (p *PluginCert360) Crawl() error {
	f := newFeedCrawl()
	items := f.parseFeed("https://cert.360.cn/feed")
	for _, item := range items {
		// https://cert.360.cn/report/detail?id=d42e9ec786a8fa79dd23ffc188d187fa
		p.res = append(p.res, &Warings{
			Title:    item.Title,
			Link:     item.Link,
			Index:    item.Link,
			From:     "cert360",
			Desc:     item.Desc,
			Time:     time.Unix(item.PubDate, 0),
			CreateAt: time.Now(),
		})
		logger.Debugln("Crwaled [Cert360]", item.Title, item.PubDate)
	}
	return nil
}
