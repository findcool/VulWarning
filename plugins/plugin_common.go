package plugins

import (
	"github.com/virink/vulWarning/model"
)

// Plugin -
type Plugin interface {
	Crawl() error
	Result() []*model.Warning
}

// GetPlugins -
func GetPlugins() []string {
	return []string{
		"aliyun",
		"cert360",
		"tencentti",
		"qianxinti",
		"githubcve",
	}
}

// PluginFactry -
func PluginFactry(name string) Plugin {
	switch name {
	case "aliyun":
		return &PluginAliyun{}
	case "cert360":
		return &PluginCert360{}
	case "tencentti":
		return &PluginTencentTi{}
	case "qianxinti":
		return &PluginQianxinTi{}
	case "githubcve":
		return &PluginGithubCVE{}
	default:
		return nil
	}
}

func getFeature(content string) {
	// TODO: getFeature
	/*
		timestamp
		cve
		cvss
		vul_type
		version
		product
		summary
		ref
	*/
}

func secThief() {
	// TODO: https://sec.thief.one/atom.xml
}

func freebuf() {
	// TODO: https://search.freebuf.com/search/?search=%E6%BC%8F%E6%B4%9E%20%E9%A2%84%E8%AD%A6#article
}

func openwall() {
	// TODO: https://www.openwall.com/lists/oss-security/
}
