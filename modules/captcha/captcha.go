package captcha

import (
	"time"

	"github.com/dchest/captcha"
)

func init() {
	captcha.SetCustomStore(captcha.NewMemoryStore(1000, 3*time.Minute))
}

func New() string {
	d := struct {
		CaptchaId string
	}{
		captcha.New(),
	}
	return d.CaptchaId
}

func VerifyString(captchaID string, captchaSolution string) bool {
	return captcha.VerifyString(captchaID, captchaSolution)
}

func Reload(captchaID string) {
	captcha.Reload(captchaID)
}
