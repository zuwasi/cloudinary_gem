package cloudinary_gem

import "testing"

func TestRubySpecURLDistributionParity(t *testing.T) {
	c, _ := New("demo", "key", "secret")
	cases := []struct {
		name string
		id   string
		opts URLOptions
		want string
	}{
		{"default secure", "test", URLOptions{}, "https://res.cloudinary.com/demo/image/upload/test"},
		{"cname http", "test", URLOptions{Secure: Bool(false), CNAME: "example.com"}, "http://example.com/demo/image/upload/test"},
		{"private cdn", "test", URLOptions{PrivateCDN: true}, "https://demo-res.cloudinary.com/image/upload/test"},
		{"secure distribution", "test", URLOptions{SecureDistribution: "assets.example.com"}, "https://assets.example.com/demo/image/upload/test"},
		{"private cdn secure distribution", "test", URLOptions{PrivateCDN: true, SecureDistribution: "assets.example.com"}, "https://assets.example.com/image/upload/test"},
		{"format", "test", URLOptions{Format: "jpg"}, "https://res.cloudinary.com/demo/image/upload/test.jpg"},
		{"url suffix", "test", URLOptions{URLSuffix: "hello"}, "https://res.cloudinary.com/demo/images/test/hello"},
		{"url suffix private", "test", URLOptions{URLSuffix: "hello", PrivateCDN: true, Format: "jpg"}, "https://demo-res.cloudinary.com/images/test/hello.jpg"},
		{"url suffix raw", "test", URLOptions{URLSuffix: "hello", PrivateCDN: true, ResourceType: "raw"}, "https://demo-res.cloudinary.com/files/test/hello"},
		{"url suffix video", "test", URLOptions{URLSuffix: "hello", PrivateCDN: true, ResourceType: "video"}, "https://demo-res.cloudinary.com/videos/test/hello"},
		{"root path", "test", URLOptions{UseRootPath: true}, "https://res.cloudinary.com/demo/test"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := c.URL(tc.id, tc.opts)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("want %s got %s", tc.want, got)
			}
		})
	}
}

func TestRubySpecSignedURLParity(t *testing.T) {
	c, _ := New("test", "123456789012345", "AbcdEfghIjklmnopq1234567890")
	u, err := c.URL("some_public_id.jpg", URLOptions{Type: "authenticated", SignURL: true, Overlay: "text:Helvetica_50:test+text"})
	if err != nil {
		t.Fatal(err)
	}
	want := "https://res.cloudinary.com/test/image/authenticated/s--j5Z1ILxd--/l_text:Helvetica_50:test+text/some_public_id.jpg"
	if u != want {
		t.Fatalf("want %s got %s", want, u)
	}

	c2, _ := New("test123", "a", "b")
	u, err = c2.URL("sample.jpg", URLOptions{SignURL: true})
	if err != nil {
		t.Fatal(err)
	}
	if u != "https://res.cloudinary.com/test123/image/upload/s--v2fTPYTu--/sample.jpg" {
		t.Fatalf("unexpected sha1 url %s", u)
	}

	c2.SignatureAlgorithm = SHA256
	u, err = c2.URL("sample.jpg", URLOptions{SignURL: true})
	if err != nil {
		t.Fatal(err)
	}
	if u != "https://res.cloudinary.com/test123/image/upload/s--2hbrSMPO--/sample.jpg" {
		t.Fatalf("unexpected sha256 url %s", u)
	}
	u, err = c2.URL("sample.jpg", URLOptions{SignURL: true, LongURLSignature: true})
	if err != nil {
		t.Fatal(err)
	}
	if u != "https://res.cloudinary.com/test123/image/upload/s--2hbrSMPOjj5BJ4xV7SgFbRDevFaQNUFf--/sample.jpg" {
		t.Fatalf("unexpected long sha256 url %s", u)
	}
}

func TestRubySpecURLValidationParity(t *testing.T) {
	c, _ := New("demo", "key", "secret")
	if _, err := c.URL("test", URLOptions{URLSuffix: "hello/world"}); err == nil {
		t.Fatal("expected invalid suffix error")
	}
	if _, err := c.URL("test", URLOptions{URLSuffix: "hello", Type: "facebook"}); err == nil {
		t.Fatal("expected invalid type error")
	}
	if _, err := c.URL("test", URLOptions{UseRootPath: true, ResourceType: "raw"}); err == nil {
		t.Fatal("expected invalid root path error")
	}
}
