package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/hust-tianbo/go_lib/log"
	"gopkg.in/yaml.v3"
)

type SystemConfig struct {
	Log map[string]yaml.Node `yaml:"log"` // 日志配置
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

func init() {
	yamlFile, err := ioutil.ReadFile("sys.yaml")
	if err != nil {
		fmt.Errorf("[ReadConfig]read yaml failed:%+v", err)
		panic("no sys.yaml")
		return
	}

	fmt.Printf("[ReadConfig]read sys config success:%+v\n", string(yamlFile))
	config := &SystemConfig{}

	unmarshalErr := yaml.Unmarshal(yamlFile, config)
	if unmarshalErr != nil {
		fmt.Errorf("[ReadConfig]unmarshal failed:%+v", unmarshalErr)
		panic("unmarshal failed")
		return
	}

	for name, node := range config.Log {
		fmt.Printf("[ReadConfig]Setup log\n:%+v", name)
		log.DefaultLogFactory.Setup(name, &YamlNodeDecoder{Node: &node})
	}

	return
}
