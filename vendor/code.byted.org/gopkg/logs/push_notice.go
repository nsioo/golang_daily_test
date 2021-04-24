package logs

import (
	"context"
	"sync"
)

const (
	noticeCtxKey = "K_NOTICE"
)

func GetNotice(ctx context.Context) *NoticeKVs {
	i := ctx.Value(noticeCtxKey)
	if ntc, ok := i.(*NoticeKVs); ok {
		return ntc
	}
	return nil
}

type NoticeKVs struct {
	kvs []interface{}
	sync.Mutex
}

func NewNoticeKVs() *NoticeKVs {
	return &NoticeKVs{
		kvs: make([]interface{}, 0, 16),
	}
}

func (l *NoticeKVs) PushNotice(k, v interface{}) {
	l.Lock()
	l.kvs = append(l.kvs, k, v)
	l.Unlock()
}

func (l *NoticeKVs) KVs() []interface{} {
	l.Lock()
	kvs := l.kvs
	l.Unlock()
	return kvs
}
