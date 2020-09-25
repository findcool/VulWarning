package plugins

import (
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/virink/vulWarning/common"
)

// Items - Items Database module
type Items struct {
	ID      uint   `gorm:"AUTO_INCREMENT;primary_key;unique_index" json:"id"`
	Title   string `gorm:"type:text" db:"title" json:"title"` // identify
	Link    string `gorm:"type:text" db:"link" json:"link"`   // link
	PubDate int64  `db:"pubdate" json:"pubdate"`              // pubdate
	Desc    string `gorm:"type:text" json:"description"`      // description
}

var dateLayout = []string{
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
	"2006-01-02T15:04:05Z",
	"2016-01-02T15:04:05.000-07:00",
	"2006-01-02T15:04:05Z07:00",
	"02 Jan 06 15:04 MST",
	"02 Jan 06 15:04 -0700",
	"Mon, 02 Jan 2006 15:04:05",
	"Mon, 02 Jan 2006 15:04:05 MST",
	"Mon Jan 02 15:04:05 -0700 2006",
	"Mon, 14 Feb 2006 15:04:05",
	"Mon, 02 Jan 2006 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 +0700",
	"Mon,02 Jan 2006 15:04:05 -0700",
}

// FeedCrawl -
type FeedCrawl struct {
	fp *gofeed.Parser
}

// ParsePubDate -
func ParsePubDate(s string) int64 {
	for _, layout := range dateLayout {
		if stamp, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return stamp.Unix()
		}
	}
	return -1
}

func (f *FeedCrawl) parseFeed(target string) []*Items {
	feed, err := f.fp.ParseURL(target)
	if err != nil {
		common.Logger.Errorln(target, err)
		return nil
	}
	common.Logger.Debugln(feed.Title)
	items := make([]*Items, len(feed.Items))
	for i, item := range feed.Items {
		link := item.Link
		if link == "" {
			link = item.GUID
		}
		var pubDate int64
		for _, date := range []string{item.Published, item.Updated} {
			pubDate = ParsePubDate(date)
			if pubDate != -1 {
				break
			}
		}
		if pubDate == -1 {
			common.Logger.Errorln("PubDate template error", item.Published, target)
			pubDate = time.Now().Unix()
		}
		items[i] = &Items{
			Title:   item.Title,
			Link:    link,
			PubDate: pubDate,
			Desc:    item.Description,
		}
		common.Logger.Debugln(item.Title, link)
	}
	return items
}

func newFeedCrawl() *FeedCrawl {
	f := &FeedCrawl{}
	f.fp = gofeed.NewParser()
	f.fp.Client = &http.Client{
		Transport: tr,
	}
	return f
}
