package cloudinary_gem

import (
	"fmt"
	"net/url"
	"strings"
)

func (c *Cloudinary) PrivateDownloadURL(publicID, format string, opts URLOptions) (string, error) {
	opts.Type = "authenticated"
	opts.Format = format
	opts.SignURL = true
	return c.URL(publicID, opts)
}
func (c *Cloudinary) DownloadArchiveURL(params map[string]any) (string, error) {
	return c.apiDownloadURL("download_archive", params)
}
func (c *Cloudinary) DownloadZipURL(params map[string]any) (string, error) {
	return c.apiDownloadURL("download_archive", merge(params, map[string]any{"target_format": "zip"}))
}
func (c *Cloudinary) apiDownloadURL(action string, params map[string]any) (string, error) {
	if params == nil {
		params = map[string]any{}
	}
	params["api_key"] = c.APIKey
	params["timestamp"] = nowUnix()
	sig, err := c.SignParams(params)
	if err != nil {
		return "", err
	}
	q := url.Values{}
	for k, v := range params {
		q.Set(k, toString(v))
	}
	q.Set("signature", sig)
	return "https://api.cloudinary.com/v1_1/" + c.CloudName + "/image/" + action + "?" + q.Encode(), nil
}
func toString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	default:
		return intfString(x)
	}
}
func intfString(v any) string { return strings.Trim(strings.ReplaceAll(fmt.Sprint(v), " ", ","), "[]") }
