package cloudinary_gem

import (
	"encoding/base64"
	"fmt"
	"strings"
)

const analyticsChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

type Analytics struct{ Product, SDKCode, SDKVersion, TechVersion string }

func SDKAnalyticsSignature(a Analytics) string {
	if a.Product == "" {
		a.Product = "B"
	}
	if a.SDKCode == "" {
		a.SDKCode = "G"
	}
	if a.SDKVersion == "" {
		a.SDKVersion = Version
	}
	if a.TechVersion == "" {
		a.TechVersion = "go"
	}
	payload := fmt.Sprintf("B%s%s%s%s", a.Product, a.SDKCode, a.SDKVersion, a.TechVersion)
	return base64.RawURLEncoding.EncodeToString([]byte(payload))
}
func SDKAnalyticsQueryParam(a Analytics) string { return "_a=" + SDKAnalyticsSignature(a) }
func decodeAnalyticsForTest(sig string) string {
	b, _ := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(sig, "_a="))
	return string(b)
}
