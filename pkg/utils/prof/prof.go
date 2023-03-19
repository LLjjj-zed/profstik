package prof

import (
	"strings"
	"sync"
)

// BufferPool  定义了一个名为 BufferPool 的同步池（sync.Pool）。
// 此池被用于复用字符串缓冲区（strings.Builder）实例，以减少内存分配和垃圾回收的负担。
// 当需要创建新的字符串缓冲区时，此池会自动调用 New 函数创建一个新的实例，
// 否则会从池中获取一个已存在的实例
var BufferPool = sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
}

// Concat 使用strings.Builder.WriteString（）拼接字符串提高性能
func Concat(s ...string) string {
	buf := BufferPool.Get().(*strings.Builder)
	defer BufferPool.Put(buf)
	for i := 0; i < len(s); i++ {
		buf.WriteString(s[i])
	}
	//这样写会产生内存逃逸，所以使用defer方式
	//str := buf.String()
	//return str
	defer buf.Reset()
	return buf.String()
}
