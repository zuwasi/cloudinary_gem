package cloudinary_gem

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type URLOptions struct {
	ResourceType       string
	Type               string
	Secure             *bool
	PrivateCDN         bool
	CNAME              string
	SecureDistribution string
	URLSuffix          string
	UseRootPath        bool
	Format             string
	Version            int64
	Width              int
	Height             int
	Crop               string
	Gravity            string
	Quality            string
	FetchFormat        string
	Radius             string
	Angle              string
	Effect             string
	Overlay            string
	Underlay           string
	SignURL            bool
	LongURLSignature   bool
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
	prefix, err := c.urlPrefix(scheme, resourceType, deliveryType, opts)
	if err != nil {
		return "", err
	}
	parts := []string{prefix}
	trans := TransformationString(opts)
	publicPath := addFormat(publicID, opts.Format)
	if opts.URLSuffix != "" {
		if strings.ContainsAny(opts.URLSuffix, "/.") {
			return "", errors.New("url suffix must not contain / or .")
		}
		publicPath = publicID + "/" + opts.URLSuffix
		if opts.Format != "" {
			publicPath += "." + strings.TrimPrefix(opts.Format, ".")
		}
	}
	signPath := strings.Join(appendNonEmpty([]string{}, trans, version(opts.Version), addFormat(publicID, opts.Format)), "/")
	if opts.SignURL {
		if c.APISecret == "" {
			return "", errors.New("api secret is required to sign url")
		}
		parts = append(parts, "s--"+c.urlSignature(signPath, opts.LongURLSignature)+"--")
	}
	parts = append(parts, appendNonEmpty([]string{}, trans, version(opts.Version), escapePath(publicPath))...)
	return strings.Join(parts, "/"), nil
}

func (c *Cloudinary) urlPrefix(scheme, resourceType, deliveryType string, opts URLOptions) (string, error) {
	if opts.URLSuffix != "" && deliveryType != "upload" && deliveryType != "private" && deliveryType != "authenticated" {
		return "", errors.New("url suffix is only supported for upload/private/authenticated delivery")
	}
	host := "res.cloudinary.com"
	includeCloudName := true
	if opts.CNAME != "" && scheme == "http" {
		host = opts.CNAME
	} else if opts.SecureDistribution != "" && scheme == "https" {
		host = opts.SecureDistribution
		includeCloudName = !opts.PrivateCDN
	} else if opts.PrivateCDN {
		host = c.CloudName + "-res.cloudinary.com"
		includeCloudName = false
	}
	segments := []string{scheme + "://" + host}
	if includeCloudName {
		segments = append(segments, c.CloudName)
	}
	if opts.UseRootPath {
		if resourceType != "image" || deliveryType != "upload" {
			return "", errors.New("root path is only supported for image/upload")
		}
		return strings.Join(segments, "/"), nil
	}
	if opts.URLSuffix != "" {
		segments = append(segments, urlSuffixResource(resourceType, deliveryType))
		return strings.Join(segments, "/"), nil
	}
	segments = append(segments, resourceType, deliveryType)
	return strings.Join(segments, "/"), nil
}

func urlSuffixResource(resourceType, deliveryType string) string {
	if resourceType == "raw" {
		return "files"
	}
	if resourceType == "video" {
		return "videos"
	}
	if deliveryType == "private" {
		return "private_images"
	}
	if deliveryType == "authenticated" {
		return "authenticated_images"
	}
	return "images"
}

func TransformationString(o URLOptions) string {
	var p []string
	if o.Crop != "" {
		p = append(p, "c_"+o.Crop)
	}
	if o.Height > 0 {
		p = append(p, "h_"+strconv.Itoa(o.Height))
	}
	if o.Width > 0 {
		p = append(p, "w_"+strconv.Itoa(o.Width))
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
	if o.Radius != "" {
		p = append(p, "r_"+o.Radius)
	}
	if o.Angle != "" {
		p = append(p, "a_"+o.Angle)
	}
	if o.Effect != "" {
		p = append(p, "e_"+o.Effect)
	}
	if o.Overlay != "" {
		p = append(p, "l_"+o.Overlay)
	}
	if o.Underlay != "" {
		p = append(p, "u_"+o.Underlay)
	}
	return strings.Join(p, ",")
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

func Bool(v bool) *bool { return &v }

func (o URLOptions) String() string { return fmt.Sprintf("%+v", struct{ URLOptions }{o}) }
