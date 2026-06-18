package cloudinary_gem

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type AuthTokenOptions struct {
	Key       string
	ACL       []string
	URL       string
	StartTime int64
	Duration  int64
	ExpiresAt int64
	IP        string
}

func GenerateAuthToken(o AuthTokenOptions) string {
	parts := []string{}
	if o.StartTime > 0 {
		parts = append(parts, "st="+strconv.FormatInt(o.StartTime, 10))
	}
	if o.Duration > 0 {
		parts = append(parts, "exp="+strconv.FormatInt(o.StartTime+o.Duration, 10))
	} else if o.ExpiresAt > 0 {
		parts = append(parts, "exp="+strconv.FormatInt(o.ExpiresAt, 10))
	}
	if o.IP != "" {
		parts = append(parts, "ip="+o.IP)
	}
	if len(o.ACL) > 0 {
		sort.Strings(o.ACL)
		parts = append(parts, "acl="+strings.Join(o.ACL, "!"))
	}
	if o.URL != "" && len(o.ACL) == 0 {
		parts = append(parts, "url="+strings.ToLower(url.QueryEscape(o.URL)))
	}
	toSign := strings.Join(parts, "~")
	mac := hmac.New(sha256.New, []byte(o.Key))
	mac.Write([]byte(toSign))
	return "__cld_token__=" + toSign + "~hmac=" + hex.EncodeToString(mac.Sum(nil))
}
