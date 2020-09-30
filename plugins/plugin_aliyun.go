package plugins

import (
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/virink/vulwarning/common"
	"github.com/virink/vulwarning/model"
)

// PluginAliyun -
type PluginAliyun struct {
	c   *colly.Collector
	res []*model.Warning
}

// Result -
func (p *PluginAliyun) Result() []*model.Warning {
	return p.res
}

// Crawl -
func (p *PluginAliyun) Crawl() error {
	p.c = newCustomCollector([]string{"help.aliyun.com"})

	p.c.OnRequest(func(r *colly.Request) {
		common.Logger.Debugln("Crawling [Aliyun]", r.URL)
	})

	p.c.OnHTML("div#se-knowledge", func(e *colly.HTMLElement) {
		bi := 0
		e.ForEach("p", func(i int, ee *colly.HTMLElement) {
			if ee.Text == "漏洞描述" {
				bi = i + 1
			}
			if bi > 0 && bi == i {
				for _, w := range p.res {
					if w.Link == e.Request.URL.String() {
						w.Desc = ee.Text
						// CVE, CVSS, DESC := GetCVE(w.Desc)
						// w.CVE = CVE
						// w.CVSS = CVSS
						// w.CVES = DESC
						break
					}
				}
			}
		})
	})

	p.c.OnHTML("li.y-clear", func(e *colly.HTMLElement) {
		title := e.ChildText("a[href]")
		if strings.Contains(title, "漏洞预警") {
			link := e.Request.AbsoluteURL(e.ChildAttr("a[href]", "href"))
			_time := e.ChildText("span")
			_time = _time[:len(_time)-8]

			p.res = append(p.res, &model.Warning{
				Title:    title,
				Link:     link,
				From:     "aliyun",
				Time:     getTime("2006-01-0215:04:05", _time),
				CreateAt: time.Now(),
			})
			p.c.Visit(link)
			common.Logger.Debugln("Crwaled [Aliyun]", title, _time)
		}
	})
	p.c.Visit("https://help.aliyun.com/noticelist/9213612.html")
	p.c.Wait()
	return nil
}
