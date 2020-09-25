package plugins

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/virink/vulWarning/common"
	"github.com/virink/vulWarning/model"
)

// PluginQianxinTi -
type PluginQianxinTi struct {
	c   *colly.Collector
	res []*model.Warning
}

// Result -
func (p *PluginQianxinTi) Result() []*model.Warning {
	return p.res
}

// Crawl -
func (p *PluginQianxinTi) Crawl() error {
	p.c = newCustomCollector([]string{"ti.qianxin.com"})

	p.c.OnRequest(func(r *colly.Request) {
		common.Logger.Debugln("Crawling [QianxinTi]", r.URL)
	})

	// TODO: Get Content From __NUXT__
	// p.c.OnHTML("script", func(e *colly.HTMLElement) {
	// 	common.Logger.Debugln(e.Text)
	// })

	p.c.OnHTML("div.art-container", func(e *colly.HTMLElement) {
		title := e.ChildText("div.text-box > div.title-home > a")
		_time := e.ChildText("div.text-box > div.author > p")
		if len(_time) > 10 {
			_time = _time[:10]
		}
		desc := e.ChildText("div.text-box > div.brief")
		p.res = append(p.res, &model.Warning{
			Title:    title,
			Link:     fmt.Sprintf(`%s?404=%s`, e.Request.URL.String(), MD5(e.Request.URL.String())),
			Desc:     desc,
			From:     "qianxin_ti",
			Time:     getTime("2006-01-02", _time),
			CreateAt: time.Now(),
		})
		common.Logger.Debugln("Crwaled [QianxinTI]", title, _time)
	})

	p.c.Visit("https://ti.qianxin.com/advisory/category/%E6%BC%8F%E6%B4%9E%E9%80%9A%E5%91%8A/")
	p.c.Wait()

	return nil
}
