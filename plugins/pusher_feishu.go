package plugins

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/virink/vulwarning/common"
	"github.com/virink/vulwarning/model"
)

// FeishuData -
type FeishuData struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text"`
}

func newFeishuData(p *model.PushData) []byte {
	s := &FeishuData{
		Title: p.Title,
		Text:  p.Text,
	}
	data, err := json.Marshal(&s)
	if err != nil {
		common.Logger.Errorln(err)
		return nil
	}
	return data
}

// PushToFeishu -
func PushToFeishu(wg *sync.WaitGroup, p *model.PushData) {
	defer wg.Done()
	if len(common.Conf.Pusher.Feishu) > 0 {
		// common.Logger.Debugln("Push to Feishu")
		data := newFeishuData(p)
		target := fmt.Sprintf(
			`https://open.feishu.cn/open-apis/bot/hook/%s`,
			common.Conf.Pusher.Feishu,
		)
		if err := httpPostJSON(target, data); err != nil {
			common.Logger.Errorln("Push to Feishu", err)
		}
	}
}
