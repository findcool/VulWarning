package plugins

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/virink/vulwarning/common"
	"github.com/virink/vulwarning/model"
)

// TODO: libURL
const libURL = "http://xxxxxx"

var tr *http.Transport

func init() {
	tr = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns: 100,
	}
}

// TranslateResp -
type TranslateResp struct {
	ResponseData struct {
		TranslatedText string `json:"translatedText"`
	} `json:"responseData"`
}

func translate(text string, lang ...string) string {
	_lang := "zh"
	if len(lang) > 0 && len(lang[0]) > 0 {
		_lang = lang[0]
	}
	_params := url.Values{}
	_params.Set("q", text)
	_params.Set("langpair", fmt.Sprintf("en|%s", _lang))
	data, err := httpGet(fmt.Sprintf(
		"https://api.mymemory.translated.net/get?%s",
		_params.Encode(),
	))
	if err != nil {
		common.Logger.Errorln(err)
		return text
	}
	var resp TranslateResp
	if err = json.Unmarshal(data, &resp); err != nil {
		common.Logger.Errorln(err)
		return text
	}
	if strings.Contains(resp.ResponseData.TranslatedText, "QUERY LENGTH LIMIT EXCEDEED") {
		return text
	}
	return resp.ResponseData.TranslatedText
}

func newCustomCollector(domains []string) *colly.Collector {
	var c *colly.Collector
	c = colly.NewCollector(
		colly.UserAgent("Vul Warnings Bot"),
		colly.MaxDepth(2),
		colly.Async(true),
		// colly.Debugger(&debug.LogDebugger{}),
	)
	c.AllowedDomains = domains

	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	return c
}

func getTime(_timeFormat, _time string) time.Time {
	t, err := time.Parse(_timeFormat, _time)
	if err != nil {
		common.Logger.Println(err.Error())
		t = time.Now()
	}
	return t
}

// MD5 -
func MD5(text string) string {
	ctx := md5.New()
	_, _ = ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

func httpGet(targetURL string) (body []byte, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)
	if req, err = http.NewRequest("GET", targetURL, nil); err != nil {
		return nil, err
	}
	// req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Transport: tr, Timeout: 5 * time.Second}
	if resp, err = client.Do(req); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	}
	return body, nil
}

func httpPostJSON(target string, data []byte) (err error) {
	var (
		req  *http.Request
		resp *http.Response
	)
	if req, err = http.NewRequest("POST", target, bytes.NewBuffer(data)); err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Transport: tr, Timeout: 5 * time.Second}
	if resp, err = client.Do(req); err != nil {
		return err
	}
	defer resp.Body.Close()

	if common.DebugMode {
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			common.Logger.Debugln(err)
		} else {
			common.Logger.Debugln(string(body))
		}
	}

	return nil
}

// PusherMessage -
func PusherMessage(p *model.PushDataV2) {
	var wg = &sync.WaitGroup{}
	wg.Add(3)
	go PushToFeishuV2(wg, p)
	go PushToQiwei(wg, p)
	go PushToDingding(wg, p)
	wg.Wait()
}

func makePushMessage(w *model.Warning) (p *model.PushDataV2) {
	p = &model.PushDataV2{
		From:  w.From,
		Link:  w.Link,
		Title: w.Title,
		Desc:  w.Desc,
		CVE:   w.CVE,
		CVES:  w.CVES,
		CVSS:  w.CVSS,
		Time:  w.Time.Format("2006-01-02 15:04:05"),
	}
	common.Logger.Debugln(p)
	return
}

// CrawlCVE -
func CrawlCVE(cve string) (CVSS string, DESC string) {
	// target := fmt.Sprintf("https://nvd.nist.gov/vuln/detail/%s", cve)
	c := newCustomCollector([]string{"nvd.nist.gov"})
	c.OnRequest(func(r *colly.Request) {
		common.Logger.Debugln("Crawling [CVEDetail]", r.URL)
	})
	c.OnHTML("p[data-testid=vuln-description]", func(e *colly.HTMLElement) {
		DESC = e.Text
	})
	c.OnHTML("a[data-testid=vuln-cvss3-panel-score]", func(e *colly.HTMLElement) {
		CVSS = e.Text
	})
	c.Visit(fmt.Sprintf("https://nvd.nist.gov/vuln/detail/%s", cve))
	c.Wait()
	return
}

// GetCVE -
func GetCVE(text string) (CVE string) {
	match := regexp.MustCompile(`(?mi)cve-?\s?(\d{4})-?(\d+)`).FindAllStringSubmatch(text, -1)
	if len(match) > 0 && len(match[0]) > 2 {
		CVE = fmt.Sprintf("CVE-%s-%s", match[0][1], match[0][2])
	}
	return
}

// GetCVEDetail -
func GetCVEDetail(CVE string) (CVSS, Detail string) {
	CVSS, Detail = CrawlCVE(CVE)
	if len(Detail) > 0 {
		Detail = translate(Detail)
	}
	return
}

// DoJob -
func DoJob(push bool) {
	model.RefreshLib()
	for _, pn := range GetPlugins() {
		p := PluginFactry(pn)
		err := p.Crawl()
		if err != nil {
			common.Logger.Errorln(err)
			continue
		}
		for _, warn := range p.Result() {
			// Is Exists
			if model.WarningIsExistsByLink(warn.Link) {
				continue
			}
			// Get CVE Detail
			CVE := GetCVE(fmt.Sprintf("%s %s", warn.Title, warn.Desc))
			if len(CVE) > 0 {
				if w, ok := model.FindWarningByCVE(CVE); ok {
					warn.CVE = CVE
					warn.CVES = w.CVES
					warn.CVSS = w.CVSS
				} else {
					CVSS, CVES := GetCVEDetail(CVE)
					warn.CVE = CVE
					warn.CVES = CVES
					warn.CVSS = CVSS
				}
			}
			warn.Title = strings.ReplaceAll(warn.Title, "{CVE}", CVE)
			if model.AddWarning(warn) {
				if push {
					PusherMessage(makePushMessage(warn))
				}
				model.UpdateWarning(warn.Link)
			}
		}
	}
}
