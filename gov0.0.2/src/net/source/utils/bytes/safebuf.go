package bytes

import (
	"bytes"
	"sync"
)

/*
*  线程安全的加了 读写锁
*
*/
type SafeBuff struct {
	buf bytes.Buffer
	wrlock sync.RWMutex
}

func (b *SafeBuff) Write(p []byte) (n int, err error) {
	b.wrlock.Lock()
	defer b.wrlock.Unlock()
	return b.buf.Write(p)
}

func (b *SafeBuff) Len() int{
  return  b.buf.Len()
}

func (b *SafeBuff) Read(p []byte) (n int, err error) {
	b.wrlock.RLock()
	defer b.wrlock.RUnlock()
	return b.buf.Read(p)
}

func (b *SafeBuff) Reset() {
	b.wrlock.Lock()
	defer b.wrlock.Unlock()
	b.buf.Reset()
}


func NewBuffer() *SafeBuff {
	s := SafeBuff{}
	//s.buf.Write(buf)
	return &s
}
