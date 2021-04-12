package utils

import (
	"context"
	"fmt"
	"strconv"
)

var (
	unknown = "-"
)

func Value2Str(v interface{}) string {
	switch tv := v.(type) {
	case nil:
		return ""
	case bool:
		if tv == true {
			return "true"
		} else {
			return "false"
		}
	case string:
		return tv
	case []byte:
		return string(tv)
	case error:
		return tv.Error()
	case int:
		return strconv.Itoa(tv)
	case int16:
		return strconv.FormatInt(int64(tv), 10)
	case int32:
		return strconv.FormatInt(int64(tv), 10)
	case int64:
		return strconv.FormatInt(int64(tv), 10)
	case uint:
		return strconv.FormatUint(uint64(tv), 10)
	case uint16:
		return strconv.FormatUint(uint64(tv), 10)
	case uint32:
		return strconv.FormatUint(uint64(tv), 10)
	case uint64:
		return strconv.FormatUint(uint64(tv), 10)
	case float32:
		return strconv.FormatFloat(float64(tv), 'f', 3, 32)
	case float64:
		return strconv.FormatFloat(float64(tv), 'f', 3, 32)
	default:
		return fmt.Sprint(v)
	}
}

// In generally, logID was injected by RPC framework.
func LogIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return unknown
	}
	// K_LOGID define on code.byted.org/kite/kitutil. avoid dependency directly.
	val := ctx.Value("K_LOGID")
	if val != nil {
		logID := val.(string)
		return logID
	}
	return unknown
}

func SpanIDFromContext(ctx context.Context) uint64 {
	if ctx == nil {
		return 0
	}
	// K_SPANID is injected by bytedtrace (https://code.byted.org/bytedtrace/bytedtrace-client-go)
	val := ctx.Value("K_SPANID")
	if val != nil {
		spanID, valid := val.(uint64)
		if valid {
			return spanID
		}
	}
	return 0
}
