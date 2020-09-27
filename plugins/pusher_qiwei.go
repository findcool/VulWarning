package plugins

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/virink/vulwarning/common"
	"github.com/virink/vulwarning/model"
)

// QiweiData -
type QiweiData struct {
	MsgType string `json:"msgtype"` // text,markdown
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

func newQiweiData(p *model.PushData) []byte {
	// https://work.weixin.qq.com/api/doc/90000/90136/91770
	s := &QiweiData{MsgType: "text"}
	s.Text.Content = p.Text
	if len(p.Title) > 0 {
		s.Text.Content = fmt.Sprintf("%s\n\n%s", p.Title, s.Text.Content)
	}
	data, err := json.Marshal(&s)
	if err != nil {
		common.Logger.Errorln(err)
		return nil
	}
	return data
}

// PushToQiwei -
func PushToQiwei(wg *sync.WaitGroup, p *model.PushData) {
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
