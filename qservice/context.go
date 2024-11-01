package qservice

import (
	"encoding/json"
	"fmt"
	"github.com/gobeam/stringy"
	"github.com/liaozhibinair/quick-utils/qdefine"
	easyCon "github.com/qiu-tec/easy-con.golang"
	"reflect"
	"strconv"
	"strings"
)

type context struct {
	pack   easyCon.PackReq
	values *values
}

type values struct {
	InputMaps   []map[string]interface{}
	InputRaw    interface{}
	OutputValue interface{}
}

func newControl(pack easyCon.PackReq) (*context, error) {
	ctx := &context{
		pack: pack,
		values: &values{
			InputMaps: make([]map[string]interface{}, 0),
		},
	}
	if pack.Content != nil {
		var content []byte
		switch pack.Content.(type) {
		case string:
			str := pack.Content.(string)
			content = []byte(str)
		default:
			js, err := json.Marshal(pack.Content)
			if err != nil {
				return nil, err
			}
			content = js
		}
		err := ctx.values.load(content)
		if err != nil {
			return nil, err
		}
	}
	return ctx, nil
}

func (c *context) GetString(key string) string {
	value := c.values.getValue(key)
	// 返回
	if value == nil {
		return ""
	}
	str := ""
	switch value.(type) {
	case string:
		str = fmt.Sprintf("%s", value)
	default:
		temp, err := json.Marshal(value)
		if err != nil {
			str = fmt.Sprintf("%v", value)
		} else {
			str = string(temp)
		}
	}
	return str
}

func (c *context) GetInt(key string) int {
	num, err := strconv.Atoi(c.GetString(key))
	if err != nil {
		panic(err)
	}
	return num
}

func (c *context) GetUInt(key string) uint64 {
	num, err := strconv.ParseUint(c.GetString(key), 10, 64)
	if err != nil {
		panic(err)
	}
	return num
}

func (c *context) GetByte(key string) byte {
	num, err := strconv.ParseInt(c.GetString(key), 10, 8)
	if err != nil {
		panic(err)
	}
	return byte(num)
}

func (c *context) GetBool(key string) bool {
	value := strings.ToLower(c.GetString(key))
	if value == "true" || value == "1" {
		return true
	}
	return false
}

func (c *context) GetDate(key string) qdefine.Date {
	model := struct {
		Time qdefine.Date
	}{}
	js := fmt.Sprintf("{\"Time\":\"%s\"}", c.GetString(key))
	err := json.Unmarshal([]byte(js), &model)
	if err != nil {
		panic(err)
	}
	return model.Time
}

func (c *context) GetDateTime(key string) qdefine.DateTime {
	model := struct {
		Time qdefine.DateTime
	}{}
	js := fmt.Sprintf("{\"Time\":\"%s\"}", c.GetString(key))
	err := json.Unmarshal([]byte(js), &model)
	if err != nil {
		panic(err)
	}
	return model.Time
}

func (c *context) GetFiles(key string) []qdefine.File {
	value := c.values.getValue(key)
	// 返回
	if files, ok := value.([]qdefine.File); ok {
		return files
	}
	return nil
}

func (c *context) GetStruct(refStruct any) {
	val := c.values.getValue("")
	// 先转为json
	js, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}
	// 再反转
	err = json.Unmarshal(js, refStruct)
	if err != nil {
		panic(err)
	}
}

func (c *context) Raw() any {
	return c.values.InputRaw
}

func (d *values) load(content []byte) error {
	var obj interface{}
	err := json.Unmarshal(content, &obj)
	if err != nil {
		return err
	}
	maps := make([]map[string]interface{}, 0)
	kind := reflect.TypeOf(obj).Kind()
	if kind == reflect.Slice {
		for _, o := range obj.([]interface{}) {
			maps = append(maps, o.(map[string]interface{}))
		}
	} else if kind == reflect.Map || kind == reflect.Struct {
		maps = append(maps, obj.(map[string]interface{}))
	} else {
		maps = append(maps, map[string]interface{}{"": obj})
	}
	d.InputRaw = obj
	d.InputMaps = maps
	return nil
}

func (d *values) getValue(key string) interface{} {
	if len(d.InputMaps) == 0 {
		return nil
	}
	var value interface{}
	if v, ok := d.InputMaps[0][key]; ok {
		// 如果存在
		value = v
	} else {
		str := stringy.New(key).CamelCase().ToLower()
		// 如果不存在，尝试查找
		for k, v := range d.InputMaps[0] {
			if str == stringy.New(k).CamelCase().ToLower() {
				value = v
				break
			}
		}
	}
	return value
}
