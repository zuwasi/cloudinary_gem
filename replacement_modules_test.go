package cloudinary_gem

import (
	"strings"
	"testing"
)

func TestSearchBuilder(t *testing.T) {
	c, _ := New("demo", "key", "secret")
	m := c.NewSearch().Expression("resource_type:image").MaxResults(10).SortBy("created_at", "asc").Aggregate("format").WithField("context").Fields("public_id", "format").TTL(300).ToMap()
	if m["expression"] != "resource_type:image" || m["max_results"] != 10 || m["ttl"] != 300 {
		t.Fatalf("unexpected search map %#v", m)
	}
	if len(m["aggregate"].([]string)) != 1 || len(m["with_field"].([]string)) != 1 {
		t.Fatalf("unexpected arrays %#v", m)
	}
}

func TestMemoryCache(t *testing.T) {
	c := NewMemoryCache()
	key := BreakpointCacheKey("id", "image", "upload", "c_fill", "jpg")
	c.Set(key, "100,200")
	if got, ok := c.Get(key); !ok || got != "100,200" {
		t.Fatalf("cache miss %q %v", got, ok)
	}
	c.Flush()
	if _, ok := c.Get(key); ok {
		t.Fatal("expected flush")
	}
}

func TestHelpers(t *testing.T) {
	c, _ := New("demo", "key", "secret")
	img, err := ImageTag(c, "sample", URLOptions{Format: "jpg"}, HTMLAttrs{"alt": "Sample"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(img, `src="https://res.cloudinary.com/demo/image/upload/sample.jpg"`) || !strings.Contains(img, `alt="Sample"`) {
		t.Fatalf("bad image tag %s", img)
	}
}

func TestUtils(t *testing.T) {
	if BuildArray([]string{"a", "b"}) != "a,b" {
		t.Fatal("bad array")
	}
	if EncodeContext(map[string]string{"b": "2", "a": "1|x"}) != `a=1\|x|b=2` {
		t.Fatal("bad context")
	}
	if ResourceTypeForFormat("jpg") != "image" || ResourceTypeForFormat("mp4") != "video" || ResourceTypeForFormat("zip") != "raw" {
		t.Fatal("bad resource type")
	}
	if !IsRemote("https://example.com/a.jpg") || IsRemote("local.jpg") {
		t.Fatal("bad remote detection")
	}
}

func TestAnalytics(t *testing.T) {
	q := SDKAnalyticsQueryParam(Analytics{SDKVersion: "1.2.3", TechVersion: "go1.25"})
	if !strings.HasPrefix(q, "_a=") {
		t.Fatalf("bad analytics %s", q)
	}
	if !strings.Contains(decodeAnalyticsForTest(q), "1.2.3") {
		t.Fatalf("bad decoded analytics %s", q)
	}
}

func TestDownloadURLs(t *testing.T) {
	c, _ := New("demo", "key", "secret")
	u, err := c.PrivateDownloadURL("folder/file", "jpg", URLOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(u, "/image/authenticated/s--") || !strings.HasSuffix(u, "/folder/file.jpg") {
		t.Fatalf("bad private url %s", u)
	}
}
