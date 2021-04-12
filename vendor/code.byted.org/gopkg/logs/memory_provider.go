package logs

import (
	"bytes"
)

type MemoryProvider struct {
	*bytes.Buffer
}

func NewMemoryProvider() *MemoryProvider {
	return &MemoryProvider{
		Buffer: &bytes.Buffer{},
	}
}

func (provider *MemoryProvider) Init() error {
	return nil
}

func (provider *MemoryProvider) SetLevel(l int) {}

func (provider *MemoryProvider) WriteMsg(msg string, level int) error {
	_, err := provider.WriteString(msg)
	return err
}

func (provider *MemoryProvider) Destroy() error {
	provider.Buffer.Reset()
	return nil
}

func (provider *MemoryProvider) Flush() error {
	return nil
}
