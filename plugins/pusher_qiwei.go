package plugins

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/virink/vulwarning/common"
	"github.com/virink/vulwarning/model"
)

// QiweiData -
type QiweiData struct {
	MsgType  string `json:"msgtype"` // text,markdown
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown"`
}

func newQiweiData(p *model.PushDataV2) []byte {
	// https://work.weixin.qq.com/api/doc/90000/90136/91770
	s := &QiweiData{MsgType: "markdown"}
	desc := p.Desc
	for _, name := range model.LibNames {
		if strings.Contains(desc, *name) {
			desc = strings.ReplaceAll(desc, *name, fmt.Sprintf("[%s](%s/%s)", *name, libURL, *name))
		}
	}
	s.Markdown.Content = fmt.Sprintf("## %s\n\n**来源:** %s\n**编号:** https://nvd.nist.gov/vuln/detail/%s\n**等级:** %s\n**说明:** %s\n**时间:** %s\n[查看详情](%s)\n\n%s", p.Title, p.From, p.CVE, p.CVSS, p.CVES, p.Time, p.Link, desc)
	data, err := json.Marshal(&s)
	if err != nil {
		common.Logger.Errorln(err)
		return nil
	}
	return data
}

// PushToQiwei -
func PushToQiwei(wg *sync.WaitGroup, p *model.PushDataV2) {
	defer wg.Done()
	if len(common.Conf.Pusher.Qiwei) > 0 {
		// common.Logger.Debugln("Push to Qiwei")
		data := newQiweiData(p)
		target := fmt.Sprintf(
			`https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s`,
			common.Conf.Pusher.Qiwei,
		)
		if err := httpPostJSON(target, data); err != nil {
			common.Logger.Errorln("Push to Qiwei", err)
		}
	}
}
