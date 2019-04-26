package token

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"sort"
	"strings"
)

// TKS 签名
type TKS []string

func (t TKS) Len() int { return len(t) }

func (t TKS) Less(i, j int) bool { return t[i] > t[j] }

func (t TKS) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t *TKS) Push(s ...string) {
	*t = append(*t, s...)
	sort.Sort(*t)
}

func (t TKS) Marshal(key string) string {
	s := strings.Join(t, "&")
	mac := hmac.New(sha1.New, []byte(key+"&"))
	mac.Write([]byte(s))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func NewTKS() TKS {
	return make(TKS, 0)
}
