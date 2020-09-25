package plugins

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/virink/vulWarning/common"
	"github.com/virink/vulWarning/model"
)

// DingdingData -
type DingdingData struct {
	MsgType string `json:"msgtype"` // text,markdown
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

func newDingdingData(p *model.PushData) []byte {
	// https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq
	s := &DingdingData{MsgType: "text"}
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

// PushToDingding -
func PushToDingding(wg *sync.WaitGroup, p *model.PushData) {
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
