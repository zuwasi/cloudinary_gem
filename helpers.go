package cloudinary_gem

import (
	"fmt"
	"html"
	"sort"
	"strings"
)

type HTMLAttrs map[string]string

func ImageTag(cld *Cloudinary, publicID string, opts URLOptions, attrs HTMLAttrs) (string, error) {
	u, err := cld.URL(publicID, opts)
	if err != nil {
		return "", err
	}
	return tag("img", mergeAttrs(attrs, HTMLAttrs{"src": u}), true), nil
}
func VideoTag(cld *Cloudinary, publicID string, opts URLOptions, attrs HTMLAttrs) (string, error) {
	opts.ResourceType = "video"
	u, err := cld.URL(publicID, opts)
	if err != nil {
		return "", err
	}
	return tag("video", mergeAttrs(attrs, HTMLAttrs{"src": u}), false), nil
}
func SourceTag(cld *Cloudinary, publicID string, opts URLOptions, attrs HTMLAttrs) (string, error) {
	u, err := cld.URL(publicID, opts)
	if err != nil {
		return "", err
	}
	return tag("source", mergeAttrs(attrs, HTMLAttrs{"srcset": u}), true), nil
}
func tag(name string, attrs HTMLAttrs, selfClosing bool) string {
	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := []string{"<" + name}
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf(`%s="%s"`, k, html.EscapeString(attrs[k])))
	}
	if selfClosing {
		return strings.Join(parts, " ") + ">"
	}
	return strings.Join(parts, " ") + fmt.Sprintf("></%s>", name)
}
func mergeAttrs(a, b HTMLAttrs) HTMLAttrs {
	out := HTMLAttrs{}
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}
	return out
}
