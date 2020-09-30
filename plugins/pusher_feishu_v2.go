package plugins

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/virink/vulwarning/common"
	"github.com/virink/vulwarning/model"
)

// https://getfeishu.cn/hc/zh-cn/articles/360024984973-%E5%9C%A8%E7%BE%A4%E8%81%8A%E4%B8%AD%E4%BD%BF%E7%94%A8%E6%9C%BA%E5%99%A8%E4%BA%BA#%E2%80%A2%E2%80%8B%E5%9C%A8%E7%BE%A4%E8%81%8A%E4%B8%AD%E4%BD%BF%E7%94%A8%E8%87%AA%E5%AE%9A%E4%B9%89%E6%9C%BA%E5%99%A8%E4%BA%BA

// Content -
type Content struct {
	Href string `json:"href"`
	Tag  string `json:"tag"`
	Text string `json:"text"`
}

// FeishuDataV2 -
type FeishuDataV2 struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Post struct {
			ZhCn struct {
				Content [][]Content `json:"content"`
				Title   string      `json:"title,omitempty"`
			} `json:"zh_cn"`
		} `json:"post"`
	} `json:"content"`
}

func newFeishuDataV2(p *model.PushDataV2) []byte {
	// 富文本
	s := &FeishuDataV2{MsgType: "post"}
	foundLibs := []Content{}
	found := true
	for _, name := range model.LibNames {
		if strings.Contains(p.Desc, *name) {
			if found {
				foundLibs = append(foundLibs, Content{Tag: "text", Text: "发现: "})
				found = false
			}
			foundLibs = append(foundLibs, Content{
				Tag: "a", Text: *name, Href: fmt.Sprintf("%s/%s", libURL, *name),
			})
		}
	}
	s.Content.Post.ZhCn.Title = p.Title
	s.Content.Post.ZhCn.Content = [][]Content{
		[]Content{
			Content{Tag: "text", Text: fmt.Sprintf(`来源: %s`, p.From)},
		},
		[]Content{
			Content{Tag: "text", Text: fmt.Sprintf(`时间: %s`, p.Time)},
		},
	}

	common.Logger.Debugln(p.CVE)
	if len(p.CVE) > 3 {
		s.Content.Post.ZhCn.Content = append(s.Content.Post.ZhCn.Content,
			[]Content{
				Content{Tag: "text", Text: "编号: "},
				Content{Tag: "a", Text: p.CVE, Href: fmt.Sprintf(`https://nvd.nist.gov/vuln/detail/%s`, p.CVE)},
			},
			[]Content{Content{Tag: "text", Text: fmt.Sprintf(`等级: %s`, p.CVSS)}},
			[]Content{Content{Tag: "text", Text: fmt.Sprintf(`说明: %s`, p.CVES)}},
		)
	}
	if len(foundLibs) > 0 {
		s.Content.Post.ZhCn.Content = append(s.Content.Post.ZhCn.Content, foundLibs)
	}
	s.Content.Post.ZhCn.Content = append(s.Content.Post.ZhCn.Content,
		[]Content{
			Content{Tag: "text", Text: "详情:"},
			Content{Tag: "a", Text: "点击查看", Href: p.Link},
		},
		[]Content{Content{Tag: "text", Text: p.Desc}},
	)
	data, err := json.Marshal(&s)
	if err != nil {
		common.Logger.Errorln(err)
		return nil
	}
	return data
}

// PushToFeishuV2 -
func PushToFeishuV2(wg *sync.WaitGroup, p *model.PushDataV2) {
	defer wg.Done()
	if len(common.Conf.Pusher.FeishuV2) > 0 {
		data := newFeishuDataV2(p)
		target := fmt.Sprintf(
			`https://open.feishu.cn/open-apis/bot/v2/hook/%s`,
			common.Conf.Pusher.FeishuV2,
		)
		if err := httpPostJSON(target, data); err != nil {
			common.Logger.Errorln("Push to Feishu", err)
		}
	}
}
