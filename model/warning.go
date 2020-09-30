package model

import (
	"fmt"
	"time"

	"github.com/virink/vulwarning/common"
)

// PushData -
type PushData struct {
	Title, Text string
}

// PushDataV2 -
type PushDataV2 struct {
	From  string
	Link  string
	Title string
	Desc  string
	CVE   string
	CVES  string
	CVSS  string
	Time  string

	Text string
}

// Warning -
type Warning struct {
	ID       uint   `gorm:"primary_key;AUTO_INCREMENT;not null"`
	From     string `gorm:"type:varchar(255)"`              // 情报平台
	Link     string `gorm:"type:varchar(250);unique_index"` // 情报链接
	Title    string `gorm:"type:varchar(255)"`
	Desc     string `gorm:"type:text"` // 情报描述/简介
	CVE      string `gorm:"type:varchar(20)"`
	CVSS     string `gorm:"type:varchar(100)"`
	CVES     string `gorm:"type:text"`
	Time     time.Time
	CreateAt time.Time
	Send     bool
	// Index    string `gorm:"type:varchar(255)"`
}

// FindWarningByLink -
func FindWarningByLink(link string) (out Warning, ok bool) {
	ok = !db.First(&out, Warning{Link: link}).RecordNotFound()
	return
}

// FindWarningByCVE -
func FindWarningByCVE(CVE string) (out Warning, ok bool) {
	ok = !db.First(&out, Warning{CVE: CVE}).RecordNotFound()
	return
}

// WarningIsExistsByLink -
func WarningIsExistsByLink(link string) bool {
	_, ok := FindWarningByLink(link)
	return ok
}

// AddWarning -
func AddWarning(w *Warning) bool {
	stmt := db.Create(w)
	if stmt.Error != nil {
		common.Logger.Errorln(stmt.Error)
	}
	return stmt.RowsAffected > 0
	// if !WarningIsExistsByLink(w.Link) {}
	// return false
}

// UpdateWarning -
func UpdateWarning(link string) error {
	return db.Model(Warning{}).
		Where(&Warning{Link: link}).
		Updates(Warning{Send: true}).Error
}

// AddWarnings -
func AddWarnings(ws []*Warning) (ps []*PushDataV2, err error) {
	ps = make([]*PushDataV2, 0)
	for _, w := range ws {
		var out Warning
		if db.First(&out, Warning{Link: w.Link}).RecordNotFound() {
			// w.Send = true
			if err = db.Create(w).Error; err != nil {
				common.Logger.Errorln(err)
			}
			// PushData
			ps = append(ps, &PushDataV2{
				From:  w.From,
				Link:  w.Link,
				Title: w.Title,
				Desc:  w.Desc,
				CVE:   w.CVE,
				CVES:  w.CVES,
				CVSS:  w.CVSS,
				Time:  w.Time.Format("2006-01-02 15:04:05"),
				// compatible
				Text: fmt.Sprintf(
					"%s\nTime : %v\nUrl  : %s  \nFrom : %s  ",
					w.Desc,
					w.Time.Format("2006-01-02 15:04:05"),
					w.Link,
					w.From,
				),
			})
		}
	}
	return ps, nil
}
