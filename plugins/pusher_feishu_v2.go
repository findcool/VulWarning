package plugins

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/virink/vulWarning/common"
	"github.com/virink/vulWarning/model"
)

// FeishuDataV2 -
type FeishuDataV2 struct {
	MsgType string `json:"msg_type"` // text,markdown
	Text    struct {
		Content string `json:"content"`
	}
}

func newFeishuDataV2(p *model.PushData) []byte {
	s := &FeishuDataV2{MsgType: "text"}
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

// PushToFeishuV2 -
func PushToFeishuV2(wg *sync.WaitGroup, p *model.PushData) {
	defer wg.Done()
	if len(common.Conf.Pusher.FeishuV2) > 0 {
		// common.Logger.Debugln("Push to FeishuV2")
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
