package cloudinary_gem

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type rewriteTransport struct {
	base string
	rt   http.RoundTripper
}

func (t rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	req.URL.Host = strings.TrimPrefix(t.base, "http://")
	return t.rt.RoundTrip(req)
}

func testClient(handler http.HandlerFunc) (*Cloudinary, func()) {
	s := httptest.NewServer(handler)
	c, _ := New("demo", "key", "secret")
	c.HTTPClient = &http.Client{Transport: rewriteTransport{base: s.URL, rt: http.DefaultTransport}}
	return c, s.Close
}

func TestAdminWrappersUseExpectedPaths(t *testing.T) {
	seen := []string{}
	c, closeFn := testClient(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Method+" "+r.URL.Path)
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	})
	defer closeFn()
	_, _ = c.UploadMappings(t.Context(), nil)
	_, _ = c.MetadataFields(t.Context(), nil)
	_, _ = c.StreamingProfiles(t.Context(), nil)
	_, _ = c.DeleteUploadPreset(t.Context(), "preset")
	want := []string{"GET /v1_1/demo/upload_mappings", "GET /v1_1/demo/metadata_fields", "GET /v1_1/demo/streaming_profiles", "DELETE /v1_1/demo/upload_presets/preset"}
	for i := range want {
		if seen[i] != want[i] {
			t.Fatalf("want %v got %v", want, seen)
		}
	}
}

func TestUploaderWrappersPostSignedForms(t *testing.T) {
	var body string
	c, closeFn := testClient(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(bodyBytes)
		body = string(bodyBytes)
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	})
	defer closeFn()
	_, err := c.Rename(t.Context(), "from", "to", nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(body, "from_public_id=from") || !strings.Contains(body, "to_public_id=to") || !strings.Contains(body, "signature=") {
		t.Fatalf("bad signed form %s", body)
	}
}

func TestAccountAPIPaths(t *testing.T) {
	var seen string
	c, closeFn := testClient(func(w http.ResponseWriter, r *http.Request) {
		seen = r.Method + " " + r.URL.Path
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	})
	defer closeFn()
	_, _ = c.Account("acct").SubAccounts(t.Context(), nil)
	if seen != "GET /v1_1/demo/provisioning/accounts/acct/sub_accounts" {
		t.Fatalf("bad path %s", seen)
	}
}
