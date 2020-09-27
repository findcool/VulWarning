package common

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Config -
type Config struct {
	MySQL struct {
		Host    string `yaml:"host"`
		User    string `yaml:"user"`
		Pass    string `yaml:"pass"`
		Name    string `yaml:"name"`
		Charset string `yaml:"charset"`
		Prefix  string `yaml:"prefix"`
	} `yaml:"mysql"`
	Pusher struct {
		Qiwei        string `yaml:"qiwei"`
		Dingding     string `yaml:"dingding"`
		DingdingSign string `yaml:"dingding_sign"`
		Feishu       string `yaml:"feishu"`
		FeishuV2     string `yaml:"feishu_v2"`
		Custom       string `yaml:"custom"`
	} `yaml:"pusher"`
	Server struct {
		Debug bool   `yaml:"debug"`
		Spec  string `yaml:"spec"`
	} `yaml:"server"`
}

// TemplateConfig -
func TemplateConfig() []byte {
	conf := &Config{}

	conf.MySQL.Charset = "utf8mb4"
	conf.MySQL.Host = "127.0.0.1"
	conf.MySQL.User = "root"
	conf.MySQL.Pass = "123456"
	conf.MySQL.Name = "vulwarning"

	conf.Pusher.Qiwei = "693axxx6-7aoc-4bc4-97a0-xxxxxxxxxxxx"
	conf.Pusher.Dingding = "fb9b5f1a04ac4305a7da1axxxxxxxxxx"
	conf.Pusher.Feishu = "fb9b5f1a04ac4305a7da1axxxxxxxxxx"

	conf.Server.Debug = true
	conf.Server.Spec = "* */10 * * * *"

	// yamlData
	data, err := yaml.Marshal(conf)
	if err != nil {
		return nil
	}
	return data
}

// LoadConfig -
func LoadConfig(fn string) (Config, error) {
	var data []byte
	var err error
	//  && os.IsNotExist(err)
	if _, err = os.Stat(fn); err != nil {
		return Conf, err
	}
	if data, err = ioutil.ReadFile(fn); err != nil {
		return Conf, err
	}
	if err = yaml.Unmarshal(data, &Conf); err != nil {
		return Conf, err
	}
	return Conf, nil
}
