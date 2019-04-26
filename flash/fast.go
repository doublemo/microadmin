package flash

const (
	FlashMessageTypeError   string = "danger"
	FlashMessageTypeSuccess        = "wuccess"
	FlashMessageTypeWarning        = "warning "
)

type Fast map[string][]string

func (f Fast) Error(msg ...string) {
	f[FlashMessageTypeError] = msg
}

func (f Fast) Success(msg ...string) {
	f[FlashMessageTypeSuccess] = msg
}

func (f Fast) Warning(msg ...string) {
	f[FlashMessageTypeWarning] = msg
}

func (f Fast) IsError() bool {
	if m, ok := f[FlashMessageTypeError]; ok && len(m) > 0 {
		return true
	}
	return false
}

func NewFast() Fast{
	return make(Fast)
}