package cloudinary_gem

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

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

func (c *Cloudinary) Ping(ctx context.Context) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "ping", nil)
}
func (c *Cloudinary) Usage(ctx context.Context) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "usage", nil)
}
func (c *Cloudinary) ResourceTypes(ctx context.Context) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "resources", nil)
}
func (c *Cloudinary) Resources(ctx context.Context, resourceType, deliveryType string, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, fmt.Sprintf("resources/%s/%s", defaultString(resourceType, "image"), defaultString(deliveryType, "upload")), params)
}
func (c *Cloudinary) Resource(ctx context.Context, publicID string, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "resources/image/upload/"+url.PathEscape(publicID), params)
}
func (c *Cloudinary) DeleteResources(ctx context.Context, publicIDs []string, params url.Values) (map[string]any, error) {
	if params == nil {
		params = url.Values{}
	}
	for _, id := range publicIDs {
		params.Add("public_ids[]", id)
	}
	return c.Admin(ctx, http.MethodDelete, "resources/image/upload", params)
}
func (c *Cloudinary) Tags(ctx context.Context, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "tags/image", params)
}
func (c *Cloudinary) Transformations(ctx context.Context, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "transformations", params)
}
func (c *Cloudinary) UploadPresets(ctx context.Context, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "upload_presets", params)
}
