package cloudinary_gem

import (
	"fmt"
	"sort"
	"strings"
)

type Transformation map[string]any

type TransformationChain []Transformation

func GenerateTransformationString(chain TransformationChain) string {
	parts := make([]string, 0, len(chain))
	for _, t := range chain {
		if s := transformationMapString(t); s != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, "/")
}

func transformationMapString(t Transformation) string {
	if len(t) == 0 {
		return ""
	}
	if named, ok := t["transformation"]; ok {
		return namedTransformation(named)
	}
	order := []struct{ key, code string }{{"if", "if"}, {"crop", "c"}, {"height", "h"}, {"width", "w"}, {"aspect_ratio", "ar"}, {"background", "b"}, {"border", "bo"}, {"color", "co"}, {"default_image", "d"}, {"delay", "dl"}, {"density", "dn"}, {"duration", "du"}, {"effect", "e"}, {"fetch_format", "f"}, {"flags", "fl"}, {"gravity", "g"}, {"keyframe_interval", "ki"}, {"overlay", "l"}, {"opacity", "o"}, {"prefix", "p"}, {"quality", "q"}, {"radius", "r"}, {"streaming_profile", "sp"}, {"underlay", "u"}, {"angle", "a"}, {"x", "x"}, {"y", "y"}, {"zoom", "z"}}
	parts := []string{}
	for _, item := range order {
		if v, ok := t[item.key]; ok && fmt.Sprint(v) != "" {
			parts = append(parts, item.code+"_"+transformValue(v))
		}
	}
	vars := []string{}
	for k, v := range t {
		if strings.HasPrefix(k, "$") {
			vars = append(vars, k+"_"+fmt.Sprint(v))
		}
	}
	sort.Strings(vars)
	parts = append(vars, parts...)
	return strings.Join(parts, ",")
}

func namedTransformation(v any) string {
	switch x := v.(type) {
	case string:
		return "t_" + x
	case []string:
		return "t_" + strings.Join(x, ".")
	case []any:
		parts := []string{}
		for _, p := range x {
			parts = append(parts, fmt.Sprint(p))
		}
		return "t_" + strings.Join(parts, ".")
	default:
		return "t_" + fmt.Sprint(v)
	}
}

func transformValue(v any) string {
	switch x := v.(type) {
	case []string:
		return strings.Join(x, ".")
	case []int:
		parts := []string{}
		for _, n := range x {
			parts = append(parts, fmt.Sprint(n))
		}
		return strings.Join(parts, ".")
	case []any:
		parts := []string{}
		for _, p := range x {
			parts = append(parts, fmt.Sprint(p))
		}
		return strings.Join(parts, ".")
	case map[string]any:
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := []string{}
		for _, k := range keys {
			parts = append(parts, k+":"+fmt.Sprint(x[k]))
		}
		return strings.Join(parts, ":")
	default:
		return fmt.Sprint(v)
	}
}

func (c *Cloudinary) URLWithTransformations(publicID string, chain TransformationChain, opts URLOptions) (string, error) {
	base := GenerateTransformationString(chain)
	extra := TransformationString(opts)
	if base != "" && extra != "" {
		base += "/" + extra
	} else if extra != "" {
		base = extra
	}
	return c.urlWithTransformationString(publicID, opts, base)
}

func (c *Cloudinary) urlWithTransformationString(publicID string, opts URLOptions, trans string) (string, error) {
	saved := opts
	opts.Width, opts.Height, opts.Crop, opts.Gravity, opts.Quality, opts.FetchFormat, opts.Radius, opts.Angle, opts.Effect, opts.Overlay, opts.Underlay = 0, 0, "", "", "", "", "", "", "", "", ""
	u, err := c.URL(publicID, opts)
	if err != nil || trans == "" {
		return u, err
	}
	prefix, suffix, ok := strings.Cut(u, "/"+addFormat(publicID, saved.Format))
	if !ok {
		return u, nil
	}
	return prefix + "/" + trans + "/" + addFormat(publicID, saved.Format) + suffix, nil
}
