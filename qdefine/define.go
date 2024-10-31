package qdefine

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
	GetDateTime(key string) DateTime
	GetFiles(key string) []File
	GetStruct(refStruct any)
	Raw() any
}

// Date 日期
type Date uint32

// DateTime 日期+时间
type DateTime uint64
