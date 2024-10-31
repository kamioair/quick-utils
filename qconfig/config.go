package qconfig

import (
	"encoding/json"
	"fmt"
	"github.com/liaozhibinair/quick-utils/qio"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"strings"
)

var (
	_filePath string
)

func Init(filePath string, configContent []byte) {
	_filePath = filePath
	if qio.PathExists(filePath) == false {
		err := qio.WriteAllBytes(filePath, configContent, false)
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}
	viper.SetConfigFile(filePath)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func Load(module string, configStruct any) {
	value := viper.Get(module)
	if value == nil {
		save(_filePath, map[string]any{module: configStruct})
		return
	}
	js, err := json.Marshal(value)
	if err == nil {
		err = json.Unmarshal(js, configStruct)
		if err != nil {
			return
		}
	}
	currJs, err := json.Marshal(configStruct)
	str1 := strings.ToLower(string(js))
	str2 := strings.ToLower(string(currJs))
	if str1 != str2 {
		save(_filePath, map[string]any{module: configStruct})
	}
}

func Get[T any](module, key string, defValue T) T {
	// 从缓存中获取
	tKey := fmt.Sprintf("%s.%s", module, key)
	if module == "" {
		tKey = key
	}
	value := viper.Get(tKey)
	if value == nil && module != "" {
		// 尝试找默认配置
		value = viper.Get(key)
	}

	// 写入默认值
	if value == nil {
		viper.SetDefault(tKey, defValue)
		value = defValue
	}

	// 转换具体指
	js, err := json.Marshal(value)
	if err == nil {
		newObj := new(T)
		err := json.Unmarshal(js, &newObj)
		cont := string(js)
		cont = cont + ""
		if err == nil {
			return *newObj
		}
	}
	return *new(T)
}

func save(filePath string, configStruct any) error {
	yamlData, err := qio.ReadAllBytes(filePath)
	if err != nil {
		return err
	}
	// 更新YAML数据
	jsonStr, err := json.Marshal(configStruct)
	if err != nil {
		return err
	}
	updatedYAML, err := updateYAMLWithJSON(yamlData, string(jsonStr))
	if err != nil {
		return err
	}
	// 过滤多余空格
	str := string(updatedYAML)
	str = strings.Replace(str, "\n\n", "\n", -1)
	// 写回文件
	err = qio.WriteAllBytes(filePath, []byte(str), false)
	if err != nil {
		return err
	}
	return nil
}

// updateNodeValue 更新yaml.Node中的值
func updateNodeValue(node *yaml.Node, key string, value interface{}) {
	if node == nil || node.Content == nil {
		return
	}
	for i, n := range node.Content {
		if n.Kind == yaml.ScalarNode && n.Value == key {
			// 找到对应的键，更新其值
			node.Content[i+1].Value = fmt.Sprintf("%v", value)
			return
		}
	}
}

// findValueInData 在data中查找指定键的值
func findValueInData(data interface{}, keys []string) interface{} {
	if len(keys) == 0 {
		return nil
	}
	switch v := data.(type) {
	case map[string]interface{}:
		if value, ok := v[keys[0]]; ok {
			if len(keys) == 1 {
				return value
			}
			return findValueInData(value, keys[1:])
		}
	case []interface{}:
		for _, item := range v {
			if value := findValueInData(item, keys); value != nil {
				return value
			}
		}
	}
	return nil
}

// updateYAMLWithJSON 将JSON字符串转换为interface{}，并更新YAML节点
func updateYAMLWithJSON(yamlData []byte, jsonStr string) ([]byte, error) {
	// 解析JSON字符串
	var jsonData interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	if err != nil {
		return nil, err
	}

	// 解析YAML文件
	var root yaml.Node
	err = yaml.Unmarshal(yamlData, &root)
	if err != nil {
		return nil, err
	}

	// 递归更新YAML节点
	updateYAMLNode(&root, jsonData, []string{})

	// 将更新后的YAML节点转换为字节数组
	return yaml.Marshal(&root)
}

// updateYAMLNode 递归更新YAML节点
func updateYAMLNode(node *yaml.Node, data interface{}, keys []string) {
	if node == nil {
		return
	}
	switch node.Kind {
	case yaml.DocumentNode:
		updateYAMLNode(node.Content[0], data, keys)
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]
			newKeys := append(keys, keyNode.Value)
			value := findValueInData(data, newKeys)
			if value != nil {
				if valueNode.Tag != "!!map" {
					updateNodeValue(node, keyNode.Value, value)
				}
			}
			updateYAMLNode(valueNode, data, newKeys)
		}
	case yaml.SequenceNode:
		for i, n := range node.Content {
			updateYAMLNode(n, data, append(keys, fmt.Sprintf("%d", i)))
		}
	default:
		fmt.Println("config.go updateYAMLNode unhandled default case")
	}
}
