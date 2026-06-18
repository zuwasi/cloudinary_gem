package cloudinary_gem

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Cloudinary) ResourceByAssetID(ctx context.Context, assetID string, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "resources/by_asset_id/"+url.PathEscape(assetID), params)
}
func (c *Cloudinary) ResourcesByTag(ctx context.Context, tag string, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "resources/image/tags/"+url.PathEscape(tag), params)
}
func (c *Cloudinary) ResourcesByContext(ctx context.Context, key, value string, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "resources/image/context/"+url.PathEscape(key)+"/"+url.PathEscape(value), params)
}
func (c *Cloudinary) DeleteResourcesByPrefix(ctx context.Context, prefix string, params url.Values) (map[string]any, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("prefix", prefix)
	return c.Admin(ctx, http.MethodDelete, "resources/image/upload", params)
}
func (c *Cloudinary) DeleteResourcesByTag(ctx context.Context, tag string, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodDelete, "resources/image/tags/"+url.PathEscape(tag), params)
}
func (c *Cloudinary) DeleteAllResources(ctx context.Context, params url.Values) (map[string]any, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("all", "true")
	return c.Admin(ctx, http.MethodDelete, "resources/image/upload", params)
}
func (c *Cloudinary) Transformation(ctx context.Context, transformation string, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "transformations/"+url.PathEscape(transformation), params)
}
func (c *Cloudinary) DeleteTransformation(ctx context.Context, transformation string) (map[string]any, error) {
	return c.Admin(ctx, http.MethodDelete, "transformations/"+url.PathEscape(transformation), nil)
}
func (c *Cloudinary) UploadPreset(ctx context.Context, name string, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "upload_presets/"+url.PathEscape(name), params)
}
func (c *Cloudinary) DeleteUploadPreset(ctx context.Context, name string) (map[string]any, error) {
	return c.Admin(ctx, http.MethodDelete, "upload_presets/"+url.PathEscape(name), nil)
}
func (c *Cloudinary) RootFolders(ctx context.Context, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "folders", params)
}
func (c *Cloudinary) Subfolders(ctx context.Context, path string, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "folders/"+url.PathEscape(path), params)
}
func (c *Cloudinary) DeleteFolder(ctx context.Context, path string) (map[string]any, error) {
	return c.Admin(ctx, http.MethodDelete, "folders/"+url.PathEscape(path), nil)
}
func (c *Cloudinary) CreateFolder(ctx context.Context, name string) (map[string]any, error) {
	return c.postJSON(ctx, "folders/"+url.PathEscape(name), nil)
}
func (c *Cloudinary) MetadataFields(ctx context.Context, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "metadata_fields", params)
}
func (c *Cloudinary) MetadataRules(ctx context.Context, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "metadata_rules", params)
}
func (c *Cloudinary) StreamingProfiles(ctx context.Context, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "streaming_profiles", params)
}
func (c *Cloudinary) StreamingProfile(ctx context.Context, name string, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, "streaming_profiles/"+url.PathEscape(name), params)
}
func (c *Cloudinary) DeleteStreamingProfile(ctx context.Context, name string) (map[string]any, error) {
	return c.Admin(ctx, http.MethodDelete, "streaming_profiles/"+url.PathEscape(name), nil)
}
func (c *Cloudinary) Analyze(ctx context.Context, inputType, analysisType string, params url.Values) (map[string]any, error) {
	return c.Admin(ctx, http.MethodGet, fmt.Sprintf("analysis/%s/%s", url.PathEscape(inputType), url.PathEscape(analysisType)), params)
}
