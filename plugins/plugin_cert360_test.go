package plugins

import (
	"fmt"
	"testing"
)

func TestCert360Crawl(t *testing.T) {
	p := &PluginCert360{}
	p.Crawl()
	res := p.Result()
	if len(res) > 0 {
		for _, r := range res {
			fmt.Println(r)
		}
	}
}
