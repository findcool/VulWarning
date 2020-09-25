package plugins

import "testing"

func TestAliyunCrawl(t *testing.T) {
	p := &PluginAliyun{}
	p.Crawl()
}
