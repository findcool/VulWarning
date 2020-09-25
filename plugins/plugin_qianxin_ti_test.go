package plugins

import (
	"fmt"
	"testing"
)

func TestPluginQianxinTiCrawl(t *testing.T) {
	p := &PluginQianxinTi{}
	p.Crawl()

	for _, x := range p.Result() {
		fmt.Println(x)
	}

}
