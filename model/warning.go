package model

import (
	"fmt"
	"time"

	"github.com/virink/vulWarning/common"
)

// PushData -
type PushData struct {
	Title, Text string
}

// Warning -
type Warning struct {
	ID       uint   `gorm:"primary_key;AUTO_INCREMENT;not null"`
	From     string `gorm:"type:varchar(255)"`              // 情报平台
	Link     string `gorm:"type:varchar(250);unique_index"` // 情报链接
	Index    string `gorm:"type:varchar(255)"`
	Title    string `gorm:"type:varchar(255)"`
	Desc     string `gorm:"type:text"` // 情报描述/简介
	Time     time.Time
	CreateAt time.Time
	Send     bool
}

// FindWarningByLink -
func FindWarningByLink(link string) (out Warning, ok bool) {
	ok = !db.First(&out, Warning{Link: link}).RecordNotFound()
	return
}

// WarningIsExistsByLink -
func WarningIsExistsByLink(link string) bool {
	_, ok := FindWarningByLink(link)
	return ok
}

// AddWarning -
func AddWarning(w *Warning) (err error) {
	if !WarningIsExistsByLink(w.Link) {
		if err = db.Create(w).Error; err != nil {
			common.Logger.Errorln(err)
			return err
		}
	}
	return nil
}

// UpdateWarning -
func UpdateWarning(link string) error {
	return db.Model(Warning{}).
		Where(&Warning{Link: link}).
		Updates(Warning{Send: true}).Error
}

// AddWarnings -
func AddWarnings(ws []*Warning) (ps []*PushData, err error) {
	ps = make([]*PushData, 0)
	// tx := db.Begin()
	for _, w := range ws {
		var out Warning
		if db.First(&out, Warning{Link: w.Link}).RecordNotFound() {
			// w.Send = true
			if err = db.Create(w).Error; err != nil {
				common.Logger.Errorln(err)
			}
			text := fmt.Sprintf(
				"%s\nTime : %v\nUrl  : %s  \nFrom : %s  ",
				w.Desc,
				w.Time.Format("2006-01-02 15:04:05"),
				w.Link,
				w.From,
			)
			common.Logger.Debugln(text)
			ps = append(ps, &PushData{Title: w.Title, Text: text})
		}
	}
	// tx.Commit()
	return ps, nil
}
