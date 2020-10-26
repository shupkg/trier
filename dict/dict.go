package dict

import (
	"time"
)

type Dict interface {
	Set(key string, valueKind Kind, value interface{})
	Get(key string, valueKind Kind) (value interface{})
	Del(key string)

	SetString(key string, value string)          //设置键字符串值
	SetInt(key string, value int64)              //设置键整数值
	SetFloat(key string, value float64)          //设置键浮点数值
	SetBool(key string, value bool)              //设置键布尔值
	SetTime(key string, value time.Time)         //设置键时间值
	SetDuration(key string, value time.Duration) //设置键时长值

	GetString(key string) (value string)          //获取键的字符串值
	GetInt(key string) (value int64)              //获取键的整数值
	GetFloat(key string) (value float64)          //获取键的浮点数值
	GetBool(key string) (value bool)              //获取键的布尔值
	GetTime(key string) (value time.Time)         //获取键的时间值
	GetDuration(key string) (value time.Duration) //获取键的时长值

	Copy(stat Dict)               //用新表覆盖现有表
	AddInt(key string, stat Dict) //叠加新表计数
	Keys() []string               //键集合
	Types() map[string]Kind       //键类型

	Clone() Dict //克隆

	MarshalBinary() (data []byte, err error) //编码方法
	UnmarshalBinary(data []byte) error       //解码方法

	Each(eachFunc func(key string, kind Kind, value interface{}) bool)
}
