package cloudinary_gem

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Cloudinary) CreateTransformation(ctx context.Context, name string, definition any) (map[string]any, error) {
	return c.postJSON(ctx, "transformations/"+url.PathEscape(name), map[string]any{"transformation": definition})
}
func (c *Cloudinary) UpdateTransformation(ctx context.Context, name string, updates map[string]any) (map[string]any, error) {
	return c.postJSON(ctx, "transformations/"+url.PathEscape(name), updates)
}
func (c *Cloudinary) CreateUploadPreset(ctx context.Context, params map[string]any) (map[string]any, error) {
	return c.postJSON(ctx, "upload_presets", params)
}
func (c *Cloudinary) UpdateUploadPreset(ctx context.Context, name string, params map[string]any) (map[string]any, error) {
	return c.postJSON(ctx, "upload_presets/"+url.PathEscape(name), params)
}
func (c *Cloudinary) UploadMappings(ctx context.Context, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "upload_mappings", params)
}
func (c *Cloudinary) UploadMapping(ctx context.Context, name string, params url.Values) (map[string]any, error) {
	path := "upload_mappings"
	if name != "" {
		path += "/" + url.PathEscape(name)
	}
	return c.Admin(ctx, http.MethodGet, path, params)
}
func (c *Cloudinary) CreateUploadMapping(ctx context.Context, name string, params map[string]any) (map[string]any, error) {
	return c.postJSON(ctx, "upload_mappings/"+url.PathEscape(name), params)
}
func (c *Cloudinary) UpdateUploadMapping(ctx context.Context, name string, params map[string]any) (map[string]any, error) {
	return c.postJSON(ctx, "upload_mappings/"+url.PathEscape(name), params)
}
func (c *Cloudinary) DeleteUploadMapping(ctx context.Context, name string) (map[string]any, error) {
	return c.Admin(ctx, http.MethodDelete, "upload_mappings/"+url.PathEscape(name), nil)
}
func (c *Cloudinary) CreateStreamingProfile(ctx context.Context, name string, params map[string]any) (map[string]any, error) {
	return c.postJSON(ctx, "streaming_profiles/"+url.PathEscape(name), params)
}
func (c *Cloudinary) UpdateStreamingProfile(ctx context.Context, name string, params map[string]any) (map[string]any, error) {
	return c.postJSON(ctx, "streaming_profiles/"+url.PathEscape(name), params)
}
func (c *Cloudinary) AddMetadataField(ctx context.Context, field map[string]any) (map[string]any, error) {
	return c.postJSON(ctx, "metadata_fields", field)
}
func (c *Cloudinary) MetadataField(ctx context.Context, externalID string, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "metadata_fields/"+url.PathEscape(externalID), params)
}
func (c *Cloudinary) UpdateMetadataField(ctx context.Context, externalID string, field map[string]any) (map[string]any, error) {
	return c.postJSON(ctx, "metadata_fields/"+url.PathEscape(externalID), field)
}
func (c *Cloudinary) DeleteMetadataField(ctx context.Context, externalID string) (map[string]any, error) {
	return c.Admin(ctx, http.MethodDelete, "metadata_fields/"+url.PathEscape(externalID), nil)
}
func (c *Cloudinary) AddMetadataRule(ctx context.Context, rule map[string]any) (map[string]any, error) {
	return c.postJSON(ctx, "metadata_rules", rule)
}
func (c *Cloudinary) UpdateMetadataRule(ctx context.Context, externalID string, rule map[string]any) (map[string]any, error) {
	return c.postJSON(ctx, "metadata_rules/"+url.PathEscape(externalID), rule)
}
func (c *Cloudinary) DeleteMetadataRule(ctx context.Context, externalID string) (map[string]any, error) {
	return c.Admin(ctx, http.MethodDelete, "metadata_rules/"+url.PathEscape(externalID), nil)
}
func (c *Cloudinary) AddRelatedAssets(ctx context.Context, publicID string, assets []string) (map[string]any, error) {
	return c.postJSON(ctx, "resources/related_assets/image/upload/"+url.PathEscape(publicID), map[string]any{"assets_to_relate": assets})
}
func (c *Cloudinary) DeleteRelatedAssets(ctx context.Context, publicID string, assets []string) (map[string]any, error) {
	return c.postJSON(ctx, "resources/related_assets/image/upload/"+url.PathEscape(publicID)+"/delete", map[string]any{"assets_to_unrelate": assets})
}
func (c *Cloudinary) UpdateResourcesAccessModeByPrefix(ctx context.Context, accessMode, prefix string, params map[string]any) (map[string]any, error) {
	return c.postJSON(ctx, "resources/access_mode", merge(params, map[string]any{"access_mode": accessMode, "prefix": prefix}))
}
func (c *Cloudinary) UpdateResourcesAccessModeByTag(ctx context.Context, accessMode, tag string, params map[string]any) (map[string]any, error) {
	return c.postJSON(ctx, "resources/access_mode", merge(params, map[string]any{"access_mode": accessMode, "tag": tag}))
}
func (c *Cloudinary) UpdateResourcesAccessModeByIDs(ctx context.Context, accessMode string, publicIDs []string, params map[string]any) (map[string]any, error) {
	return c.postJSON(ctx, "resources/access_mode", merge(params, map[string]any{"access_mode": accessMode, "public_ids": publicIDs}))
}
func (c *Cloudinary) PublishByPrefix(ctx context.Context, prefix string, params map[string]any) (map[string]any, error) {
	return c.postJSON(ctx, "resources/publish_resources", merge(params, map[string]any{"prefix": prefix}))
}
func (c *Cloudinary) PublishByTag(ctx context.Context, tag string, params map[string]any) (map[string]any, error) {
	return c.postJSON(ctx, "resources/publish_resources", merge(params, map[string]any{"tag": tag}))
}
func (c *Cloudinary) PublishByIDs(ctx context.Context, publicIDs []string, params map[string]any) (map[string]any, error) {
	return c.postJSON(ctx, "resources/publish_resources", merge(params, map[string]any{"public_ids": publicIDs}))
}
func (c *Cloudinary) GetBreakpoints(ctx context.Context, publicID string, opts URLOptions) (map[string]any, error) {
	u, err := c.URL(publicID, opts)
	if err != nil {
		return nil, err
	}
	return c.Admin(ctx, http.MethodGet, "breakpoints/"+url.QueryEscape(u), nil)
}
func (c *Cloudinary) VisualSearch(ctx context.Context, params map[string]any) (map[string]any, error) {
	return c.postJSON(ctx, "resources/visual_search", params)
}
func apiPath(format string, args ...any) string { return fmt.Sprintf(format, args...) }
