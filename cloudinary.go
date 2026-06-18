package cloudinary_gem

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const Version = "0.1.0"

type Algorithm string

const (
	SHA1   Algorithm = "sha1"
	SHA256 Algorithm = "sha256"
)

type Cloudinary struct {
	CloudName          string
	APIKey             string
	APISecret          string
	Secure             bool
	SignatureAlgorithm Algorithm
	HTTPClient         *http.Client
}

type URLOptions struct {
	ResourceType string
	Type         string
	Secure       *bool
	Format       string
	Version      int64
	Width        int
	Height       int
	Crop         string
	Gravity      string
	Quality      string
	FetchFormat  string
	SignURL      bool
}

func New(cloudName, apiKey, apiSecret string) (*Cloudinary, error) {
	if cloudName == "" {
		return nil, errors.New("cloud name is required")
	}
	return &Cloudinary{CloudName: cloudName, APIKey: apiKey, APISecret: apiSecret, Secure: true, SignatureAlgorithm: SHA1, HTTPClient: http.DefaultClient}, nil
}

func NewFromEnv() (*Cloudinary, error) { return NewFromURL(os.Getenv("CLOUDINARY_URL")) }

func NewFromURL(raw string) (*Cloudinary, error) {
	if raw == "" {
		return nil, errors.New("cloudinary url is required")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(u.Scheme, "cloudinary") {
		return nil, fmt.Errorf("invalid scheme %q", u.Scheme)
	}
	key := u.User.Username()
	secret, _ := u.User.Password()
	return New(u.Host, key, secret)
}

func (c *Cloudinary) URL(publicID string, opts URLOptions) (string, error) {
	if c == nil {
		return "", errors.New("nil Cloudinary")
	}
	if publicID == "" {
		return "", errors.New("public id is required")
	}
	resourceType := defaultString(opts.ResourceType, "image")
	deliveryType := defaultString(opts.Type, "upload")
	secure := c.Secure
	if opts.Secure != nil {
		secure = *opts.Secure
	}
	scheme := "http"
	if secure {
		scheme = "https"
	}
	parts := []string{scheme + "://res.cloudinary.com", c.CloudName, resourceType, deliveryType}
	trans := transformation(opts)
	signPath := strings.Join(appendNonEmpty([]string{}, trans, version(opts.Version), addFormat(publicID, opts.Format)), "/")
	if opts.SignURL {
		if c.APISecret == "" {
			return "", errors.New("api secret is required to sign url")
		}
		parts = append(parts, "s--"+c.urlSignature(signPath)+"--")
	}
	parts = append(parts, appendNonEmpty([]string{}, trans, version(opts.Version), escapePath(addFormat(publicID, opts.Format)))...)
	return strings.Join(parts, "/"), nil
}

func transformation(o URLOptions) string {
	var p []string
	if o.Width > 0 {
		p = append(p, "w_"+strconv.Itoa(o.Width))
	}
	if o.Height > 0 {
		p = append(p, "h_"+strconv.Itoa(o.Height))
	}
	if o.Crop != "" {
		p = append(p, "c_"+o.Crop)
	}
	if o.Gravity != "" {
		p = append(p, "g_"+o.Gravity)
	}
	if o.Quality != "" {
		p = append(p, "q_"+o.Quality)
	}
	if o.FetchFormat != "" {
		p = append(p, "f_"+o.FetchFormat)
	}
	return strings.Join(p, ",")
}

func (c *Cloudinary) SignParams(params map[string]any) (string, error) {
	if c.APISecret == "" {
		return "", errors.New("api secret is required")
	}
	return signParams(params, c.APISecret, c.SignatureAlgorithm), nil
}

func signParams(params map[string]any, secret string, algo Algorithm) string {
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

type UploadResult map[string]any

func (c *Cloudinary) Upload(ctx context.Context, file io.Reader, filename string, params map[string]any) (UploadResult, error) {
	if params == nil {
		params = map[string]any{}
	}
	params["timestamp"] = time.Now().Unix()
	params["api_key"] = c.APIKey
	sig, err := c.SignParams(params)
	if err != nil {
		return nil, err
	}
	params["signature"] = sig
	body, pw := io.Pipe()
	mw := multipart.NewWriter(pw)
	go func() {
		defer pw.Close()
		defer mw.Close()
		fw, err := mw.CreateFormFile("file", filename)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		if _, err := io.Copy(fw, file); err != nil {
			pw.CloseWithError(err)
			return
		}
		for k, v := range params {
			_ = mw.WriteField(k, fmt.Sprint(v))
		}
	}()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", c.CloudName), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	return c.doJSON(req)
}

func (c *Cloudinary) Admin(ctx context.Context, method, path string, params url.Values) (map[string]any, error) {
	if params == nil {
		params = url.Values{}
	}
	u := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/%s", c.CloudName, strings.TrimLeft(path, "/"))
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.APIKey, c.APISecret)
	return c.doJSON(req)
}

func (c *Cloudinary) doJSON(req *http.Request) (map[string]any, error) {
	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var out map[string]any
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return out, fmt.Errorf("cloudinary api error: status %d", res.StatusCode)
	}
	return out, nil
}

func (c *Cloudinary) urlSignature(path string) string {
	sum := sha1.Sum([]byte(path + c.APISecret))
	return strings.TrimRight(base64.URLEncoding.EncodeToString(sum[:])[:8], "=")
}

func version(v int64) string {
	if v > 0 {
		return "v" + strconv.FormatInt(v, 10)
	}
	return ""
}
func defaultString(v, d string) string {
	if v == "" {
		return d
	}
	return v
}
func appendNonEmpty(dst []string, vals ...string) []string {
	for _, v := range vals {
		if v != "" {
			dst = append(dst, v)
		}
	}
	return dst
}
func addFormat(id, format string) string {
	if format == "" || strings.Contains(id[strings.LastIndex(id, "/")+1:], ".") {
		return id
	}
	return id + "." + strings.TrimPrefix(format, ".")
}
func escapePath(p string) string {
	parts := strings.Split(p, "/")
	for i := range parts {
		parts[i] = url.PathEscape(parts[i])
	}
	return strings.Join(parts, "/")
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
