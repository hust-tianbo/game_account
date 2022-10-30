package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/hust-tianbo/go_lib/log"
	"gopkg.in/yaml.v3"
)

var GConfig *SystemConfig

type SystemConfig struct {
	APPID     string               `yaml:"app_id"`
	APPSecret string               `yaml:"app_secret"`
	Log       map[string]yaml.Node `yaml:"log"` // 日志配置
}

type YamlNodeDecoder struct {
	Node *yaml.Node
}

// Decode 解析yaml node配置
func (d *YamlNodeDecoder) Decode(conf interface{}) error {

	if d.Node == nil {
		return errors.New("yaml node empty")
	}
	return d.Node.Decode(conf)
}

func InitConfig() bool {
	return true
}

func init() {
	yamlFile, err := ioutil.ReadFile("sys.yaml")
	if err != nil {
		fmt.Errorf("[ReadConfig]read yaml failed:%+v", err)
		panic("no sys.yaml")
		return
	}

	fmt.Printf("[ReadConfig]read sys config success:%+v\n", string(yamlFile))

	unmarshalErr := yaml.Unmarshal(yamlFile, GConfig)
	if unmarshalErr != nil {
		fmt.Errorf("[ReadConfig]unmarshal failed:%+v", unmarshalErr)
		panic("unmarshal failed")
		return
	}

	for name, node := range GConfig.Log {
		fmt.Printf("[ReadConfig]Setup log\n:%+v", name)
		log.DefaultLogFactory.Setup(name, &YamlNodeDecoder{Node: &node})
	}

	return
}
