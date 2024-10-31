package qconvert

// ToAny
//
//	@Description: 将上下文转为指定类型
//	@param raw 原始对象
//	@return T 指定类型
func ToAny[T any](raw interface{}) T {
	return *new(T)
}
