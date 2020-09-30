package plugins

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/virink/vulwarning/common"
	"github.com/virink/vulwarning/model"
)

func calcSign(signToken string) string {
	timestamp := time.Now().UnixNano() / 1000 / 1000
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, signToken)
	h := hmac.New(sha256.New, []byte(signToken))
	h.Write([]byte(stringToSign))
	sumValue := h.Sum(nil)
	return fmt.Sprintf(
		"&timestamp=%d&sign=%s",
		timestamp,
		url.QueryEscape(base64.StdEncoding.EncodeToString(sumValue)),
	)
}

// DingdingData -
type DingdingData struct {
	MsgType string `json:"msgtype"` // text,markdown
	Text    struct {
		Content string `json:"content"`
	} `json:"text,omitempty"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown,omitempty"`
}

func newDingdingData(p *model.PushDataV2) []byte {
	// https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq
	s := &DingdingData{MsgType: "markdown"}
	s.Markdown.Title = p.Title

	// From:  p.From,
	// CVE:   p.CVE, fmt.Sprintf(`https://nvd.nist.gov/vuln/detail/%s`, p.CVE)
	// CVES:  p.CVES,
	// CVSS:  p.CVSS,
	// Time:  p.Time,
	// Link:  p.Link,
	// Desc:  p.Desc,
	desc := p.Desc
	for _, name := range model.LibNames {
		if strings.Contains(desc, *name) {
			desc = strings.ReplaceAll(desc, *name, fmt.Sprintf("[%s](%s/%s)", *name, libURL, *name))
		}
	}

	s.Markdown.Text = fmt.Sprintf("**来源:** %s\n**编号:** https://nvd.nist.gov/vuln/detail/%s\n**等级:** %s\n**说明:** %s\n**时间:** %s\n[查看详情](%s)\n\n%s", p.From, p.CVE, p.CVSS, p.CVES, p.Time, p.Link, desc)
	data, err := json.Marshal(&s)
	if err != nil {
		common.Logger.Errorln(err)
		return nil
	}
	return data
}

// PushToDingding -
func PushToDingding(wg *sync.WaitGroup, p *model.PushDataV2) {
	defer wg.Done()
	if len(common.Conf.Pusher.Dingding) > 0 && len(common.Conf.Pusher.DingdingSign) > 0 {
		// common.Logger.Debugln("Push to Dingding")
		data := newDingdingData(p)
		target := fmt.Sprintf(
			`https://oapi.dingtalk.com/robot/send?access_token=%s%s`,
			common.Conf.Pusher.Dingding,
			calcSign(common.Conf.Pusher.DingdingSign),
		)
		if err := httpPostJSON(target, data); err != nil {
			common.Logger.Errorln("Push to Dingding", err)
		}
	}
}
