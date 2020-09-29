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
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/virink/vulwarning/common"
	"github.com/virink/vulwarning/model"
)

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
func PusherMessage(p *model.PushData) {
	var wg = &sync.WaitGroup{}
	wg.Add(4)
	go PushToFeishu(wg, p)
	go PushToFeishuV2(wg, p)
	go PushToQiwei(wg, p)
	go PushToDingding(wg, p)
	wg.Wait()
}

func makePushMessage(w *model.Warning) (p *model.PushData) {
	p = &model.PushData{
		Title: w.Title,
		Text: fmt.Sprintf(
			"%s\n\nTime : %v\nUrl  : %s  \nFrom : %s  ",
			w.Desc,
			w.Time.Format("2006-01-02 15:04:05"),
			w.Link,
			w.From,
		),
	}
	common.Logger.Debugln(p)
	return
}

// DoJob -
func DoJob(push bool) {
	for _, pn := range GetPlugins() {
		p := PluginFactry(pn)
		err := p.Crawl()
		if err != nil {
			common.Logger.Errorln(err)
			continue
		}
		for _, warn := range p.Result() {
			if model.AddWarning(warn) {
				if push {
					PusherMessage(makePushMessage(warn))
				}
				model.UpdateWarning(warn.Link)
			}
		}
	}
}
