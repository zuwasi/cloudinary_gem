package cloudinary_gem

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c *Cloudinary) postJSON(ctx context.Context, path string, payload map[string]any) (map[string]any, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/%s", c.CloudName, strings.TrimLeft(path, "/")), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.APIKey, c.APISecret)
	return c.doJSON(req)
}

func (c *Cloudinary) Destroy(ctx context.Context, publicID string, params map[string]any) (map[string]any, error) {
	return c.uploadAPICall(ctx, "destroy", merge(params, map[string]any{"public_id": publicID}))
}
func (c *Cloudinary) Rename(ctx context.Context, fromPublicID, toPublicID string, params map[string]any) (map[string]any, error) {
	return c.uploadAPICall(ctx, "rename", merge(params, map[string]any{"from_public_id": fromPublicID, "to_public_id": toPublicID}))
}
func (c *Cloudinary) Explicit(ctx context.Context, publicID string, params map[string]any) (map[string]any, error) {
	return c.uploadAPICall(ctx, "explicit", merge(params, map[string]any{"public_id": publicID}))
}
func (c *Cloudinary) AddTag(ctx context.Context, tag string, publicIDs []string, params map[string]any) (map[string]any, error) {
	return c.tagAPICall(ctx, "add", tag, publicIDs, params)
}
func (c *Cloudinary) RemoveTag(ctx context.Context, tag string, publicIDs []string, params map[string]any) (map[string]any, error) {
	return c.tagAPICall(ctx, "remove", tag, publicIDs, params)
}
func (c *Cloudinary) ReplaceTag(ctx context.Context, tag string, publicIDs []string, params map[string]any) (map[string]any, error) {
	return c.tagAPICall(ctx, "replace", tag, publicIDs, params)
}
func (c *Cloudinary) RemoveAllTags(ctx context.Context, publicIDs []string, params map[string]any) (map[string]any, error) {
	return c.tagAPICall(ctx, "remove_all", "", publicIDs, params)
}
func (c *Cloudinary) AddContext(ctx context.Context, contextMap map[string]string, publicIDs []string, params map[string]any) (map[string]any, error) {
	return c.uploadAPICall(ctx, "context", merge(params, map[string]any{"command": "add", "context": EncodeContext(contextMap), "public_ids": strings.Join(publicIDs, ",")}))
}
func (c *Cloudinary) RemoveAllContext(ctx context.Context, publicIDs []string, params map[string]any) (map[string]any, error) {
	return c.uploadAPICall(ctx, "context", merge(params, map[string]any{"command": "remove_all", "public_ids": strings.Join(publicIDs, ",")}))
}
func (c *Cloudinary) CreateArchive(ctx context.Context, params map[string]any) (map[string]any, error) {
	return c.uploadAPICall(ctx, "generate_archive", params)
}
func (c *Cloudinary) CreateZip(ctx context.Context, params map[string]any) (map[string]any, error) {
	return c.uploadAPICall(ctx, "generate_archive", merge(params, map[string]any{"target_format": "zip"}))
}

func (c *Cloudinary) uploadAPICall(ctx context.Context, action string, params map[string]any) (map[string]any, error) {
	if params == nil {
		params = map[string]any{}
	}
	params["api_key"] = c.APIKey
	if _, ok := params["timestamp"]; !ok {
		params["timestamp"] = nowUnix()
	}
	sig, err := c.SignParams(params)
	if err != nil {
		return nil, err
	}
	params["signature"] = sig
	form := url.Values{}
	for k, v := range params {
		form.Set(k, fmt.Sprint(v))
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/%s", c.CloudName, action), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.doJSON(req)
}

func (c *Cloudinary) tagAPICall(ctx context.Context, command, tag string, publicIDs []string, params map[string]any) (map[string]any, error) {
	base := map[string]any{"command": command, "public_ids": strings.Join(publicIDs, ",")}
	if tag != "" {
		base["tag"] = tag
	}
	return c.uploadAPICall(ctx, "tags", merge(params, base))
}

func merge(a, b map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}
	return out
}
