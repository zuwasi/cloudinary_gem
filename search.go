package cloudinary_gem

import (
	"context"
	"net/http"
	"net/url"
)

type Search struct {
	cld      *Cloudinary
	endpoint string
	payload  map[string]any
}

func (c *Cloudinary) NewSearch() *Search {
	return &Search{cld: c, endpoint: "resources/search", payload: map[string]any{}}
}
func (c *Cloudinary) NewSearchFolders() *Search {
	return &Search{cld: c, endpoint: "folders/search", payload: map[string]any{}}
}

func (s *Search) Expression(v string) *Search { s.payload["expression"] = v; return s }
func (s *Search) MaxResults(v int) *Search    { s.payload["max_results"] = v; return s }
func (s *Search) NextCursor(v string) *Search { s.payload["next_cursor"] = v; return s }
func (s *Search) SortBy(field, dir string) *Search {
	s.payload["sort_by"] = []map[string]string{{field: defaultString(dir, "desc")}}
	return s
}
func (s *Search) Aggregate(v string) *Search { s.appendString("aggregate", v); return s }
func (s *Search) WithField(v string) *Search { s.appendString("with_field", v); return s }
func (s *Search) Fields(v ...string) *Search { s.payload["fields"] = v; return s }
func (s *Search) TTL(v int) *Search          { s.payload["ttl"] = v; return s }
func (s *Search) ToMap() map[string]any {
	out := map[string]any{}
	for k, v := range s.payload {
		out[k] = v
	}
	return out
}
func (s *Search) Execute(ctx context.Context) (map[string]any, error) {
	return s.cld.postJSON(ctx, s.endpoint, s.payload)
}
func (s *Search) URL(ttl int, nextCursor string) (string, error) {
	q := url.Values{}
	if ttl > 0 {
		q.Set("ttl", intString(ttl))
	}
	if nextCursor != "" {
		q.Set("next_cursor", nextCursor)
	}
	res, err := s.cld.Admin(context.Background(), http.MethodPost, s.endpoint+"/url", q)
	if err != nil {
		return "", err
	}
	if u, ok := res["url"].(string); ok {
		return u, nil
	}
	return "", nil
}
func (s *Search) appendString(k, v string) {
	cur, _ := s.payload[k].([]string)
	s.payload[k] = append(cur, v)
}
