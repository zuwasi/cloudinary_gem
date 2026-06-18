package cloudinary_gem

import (
	"context"
	"net/http"
	"net/url"
)

type AccountAPI struct {
	Cloudinary *Cloudinary
	AccountID  string
}

func (c *Cloudinary) Account(accountID string) *AccountAPI {
	return &AccountAPI{Cloudinary: c, AccountID: accountID}
}
func (a *AccountAPI) call(ctx context.Context, method, path string, params url.Values) (map[string]any, error) {
	return a.Cloudinary.Admin(ctx, method, "provisioning/accounts/"+url.PathEscape(a.AccountID)+"/"+path, params)
}
func (a *AccountAPI) CreateSubAccount(ctx context.Context, name string, cloudName string, custom map[string]any, enabled *bool) (map[string]any, error) {
	p := map[string]any{"name": name}
	if cloudName != "" {
		p["cloud_name"] = cloudName
	}
	if custom != nil {
		p["custom_attributes"] = custom
	}
	if enabled != nil {
		p["enabled"] = *enabled
	}
	return a.Cloudinary.postJSON(ctx, "provisioning/accounts/"+url.PathEscape(a.AccountID)+"/sub_accounts", p)
}
func (a *AccountAPI) SubAccounts(ctx context.Context, params url.Values) (map[string]any, error) {
	return a.call(ctx, http.MethodGet, "sub_accounts", params)
}
func (a *AccountAPI) SubAccount(ctx context.Context, id string) (map[string]any, error) {
	return a.call(ctx, http.MethodGet, "sub_accounts/"+url.PathEscape(id), nil)
}
func (a *AccountAPI) DeleteSubAccount(ctx context.Context, id string) (map[string]any, error) {
	return a.call(ctx, http.MethodDelete, "sub_accounts/"+url.PathEscape(id), nil)
}
func (a *AccountAPI) Users(ctx context.Context, params url.Values) (map[string]any, error) {
	return a.call(ctx, http.MethodGet, "users", params)
}
func (a *AccountAPI) User(ctx context.Context, id string) (map[string]any, error) {
	return a.call(ctx, http.MethodGet, "users/"+url.PathEscape(id), nil)
}
func (a *AccountAPI) DeleteUser(ctx context.Context, id string) (map[string]any, error) {
	return a.call(ctx, http.MethodDelete, "users/"+url.PathEscape(id), nil)
}
func (a *AccountAPI) UserGroups(ctx context.Context, params url.Values) (map[string]any, error) {
	return a.call(ctx, http.MethodGet, "user_groups", params)
}
func (a *AccountAPI) AccessKeys(ctx context.Context, subAccountID string, params url.Values) (map[string]any, error) {
	return a.call(ctx, http.MethodGet, "sub_accounts/"+url.PathEscape(subAccountID)+"/access_keys", params)
}
