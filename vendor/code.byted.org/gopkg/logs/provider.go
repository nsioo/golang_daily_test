package logs

type LogProvider interface {
	Init() error
	SetLevel(l int)
	WriteMsg(msg string, level int) error
	Destroy() error
	Flush() error
}

type KVLogProvider interface {
	LogProvider

	// headers: PSM, Level, Time, Date, LogID ...
	// kvs: user-defined Key and Values
	WriteMsgKVs(level int, msg string, headers map[string]string, kvs map[string]string) error
}

type LogProviderPlus interface {
	LogProvider
	Level() int
	Write(*LogMsg) error
}
