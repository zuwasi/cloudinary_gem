package cloudinary_gem

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Cloudinary struct {
	CloudName          string
	APIKey             string
	APISecret          string
	Secure             bool
	SignatureAlgorithm Algorithm
	HTTPClient         *http.Client
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
