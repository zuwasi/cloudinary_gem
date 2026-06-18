package cloudinary_gem

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (c *Cloudinary) SignParams(params map[string]any) (string, error) {
	if c.APISecret == "" {
		return "", errors.New("api secret is required")
	}
	return SignParams(params, c.APISecret, c.SignatureAlgorithm), nil
}

func SignParams(params map[string]any, secret string, algo Algorithm) string {
	keys := make([]string, 0, len(params))
	for k, v := range params {
		if k == "file" || k == "api_key" || k == "resource_type" || k == "cloud_name" || k == "signature" || v == nil || fmt.Sprint(v) == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, k+"="+fmt.Sprint(params[k]))
	}
	base := strings.Join(pairs, "&") + secret
	if algo == SHA256 {
		sum := sha256.Sum256([]byte(base))
		return hex.EncodeToString(sum[:])
	}
	sum := sha1.Sum([]byte(base))
	return hex.EncodeToString(sum[:])
}

func (c *Cloudinary) urlSignature(path string) string {
	sum := sha1.Sum([]byte(path + c.APISecret))
	return strings.TrimRight(base64.URLEncoding.EncodeToString(sum[:])[:8], "=")
}

func VerifyNotificationSignature(body, timestamp, signature, secret string, algo Algorithm, maxAge time.Duration) bool {
	if maxAge > 0 {
		ts, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil || time.Since(time.Unix(ts, 0)) > maxAge {
			return false
		}
	}
	msg := body + timestamp
	var got string
	if algo == SHA256 {
		sum := hmac.New(sha256.New, []byte(secret))
		sum.Write([]byte(msg))
		got = hex.EncodeToString(sum.Sum(nil))
	} else {
		sum := hmac.New(sha1.New, []byte(secret))
		sum.Write([]byte(msg))
		got = hex.EncodeToString(sum.Sum(nil))
	}
	return hmac.Equal([]byte(got), []byte(signature))
}

func VerifyAPIResponseSignature(publicID string, version int64, signature, secret string, algo Algorithm) bool {
	params := map[string]any{"public_id": publicID, "version": version}
	return hmac.Equal([]byte(SignParams(params, secret, algo)), []byte(signature))
}
