package qdefine

import "reflect"

// File 文件
type File struct {
	Name string // 文件名
	Size int64  // 文件大小
	Data []byte // 内容
}

// Context 上下文
type Context interface {
	GetString(key string) string
	GetInt(key string) int
	GetUInt(key string) uint64
	GetByte(key string) byte
	GetBool(key string) bool
	GetDate(key string) Date
	GetTime(key string) DateTime
	GetStruct(key string, objType reflect.Type) any
	GetList(listType reflect.Type) any
	GetFiles(key string) []File
	GetReturnValue() interface{}
	SetNewReturnValue(newValue interface{})
}
