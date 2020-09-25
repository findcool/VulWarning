package plugins

import "testing"

func TestTencentTiCrawl(t *testing.T) {
	p := &PluginTencentTi{}
	p.Crawl()
}
