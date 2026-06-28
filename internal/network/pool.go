package network

import (
	"sync"
)

const BufferSize = 32 * 1024

var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, BufferSize)
		return &b
	},
}

func GetBuffer() *[]byte {
	return bufPool.Get().(*[]byte)
}

func PutBuffer(b *[]byte) {
	bufPool.Put(b)
}
