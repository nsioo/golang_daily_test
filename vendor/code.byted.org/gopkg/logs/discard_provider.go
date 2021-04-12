package logs

type DiscardProvider struct{}

func NewDiscardProvider() *DiscardProvider { return &DiscardProvider{} }

func (dp *DiscardProvider) Init() error { return nil }

func (dp *DiscardProvider) SetLevel(l int) {}

func (dp *DiscardProvider) WriteMsg(msg string, level int) error { return nil }

func (dp *DiscardProvider) Destroy() error { return nil }

func (dp *DiscardProvider) Flush() error { return nil }
