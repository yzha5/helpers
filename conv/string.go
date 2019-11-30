package conv

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

//将字符串数组转换为一个字符串 可带分隔符
//假如：参数 strA:["abc","def"] spl:| bra:["(",")"] ==> "(abc)|(def)"
//str 字符串数组
//spl 分隔符
//bra 使用括号将每个数组元素括起来
func ArrayToString(strA []string, spl string, bra [2]string) (str string) {
	for i, s := range strA {
		if i+1 == len(strA) {
			str = fmt.Sprintf("%s%s%s%s", str, bra[0], s, bra[1])
		} else {
			str = fmt.Sprintf("%s%s%s%s%s", str, bra[0], s, bra[1], spl)
		}
	}
	return str
}

//json字符串 转 map[string]interface{}
func JsonToMap(json_ string) (map_ map[string]interface{}, err error) {
	err = json.Unmarshal([]byte(json_), &map_)
	if err != nil {
		return nil, errors.New("无法转换，数据格式对不上！")
	}
	return
}

//struct 转 map
func StructToMap(obj interface{}) (map_ map[string]interface{}, err error) {
	//通过反射获取传入obj的struct类型和值
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	//pointer to value
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	map_ = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		if reflect.TypeOf(v.Field(i).Interface()).Kind() == reflect.Struct {
			//判断基础类型是否为 struct
			m, e := StructToMap(v.Field(i).Interface())

			if e != nil {
				return nil, errors.New("无法转换数据")
			}
			map_[t.Field(i).Name] = m
		} else if reflect.TypeOf(v.Field(i).Interface()).Kind() == reflect.Ptr {
			//判断基础类型是否为 pointer
			if v.Field(i).Elem().Kind() == reflect.Struct {
				//指针后的类型是否为 struct
				m, e := StructToMap(v.Field(i).Elem().Interface())

				if e != nil {
					return nil, errors.New("无法转换数据")
				}
				map_[t.Field(i).Name] = m
			} else {
				map_[t.Field(i).Name] = v.Field(i).Elem()
			}
		} else {
			map_[t.Field(i).Name] = v.Field(i).Interface()
		}
	}
	return
}
