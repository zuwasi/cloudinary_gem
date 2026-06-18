package cloudinary_gem

import (
	"testing"
	"time"
)

func TestNewFromURL(t *testing.T) {
	c, err := NewFromURL("cloudinary://key:secret@demo")
	if err != nil {
		t.Fatal(err)
	}
	if c.CloudName != "demo" || c.APIKey != "key" || c.APISecret != "secret" {
		t.Fatalf("unexpected config: %#v", c)
	}
}

func TestURL(t *testing.T) {
	c, _ := New("demo", "key", "secret")
	u, err := c.URL("folder/sample", URLOptions{Width: 100, Height: 150, Crop: "fill", Format: "jpg"})
	if err != nil {
		t.Fatal(err)
	}
	want := "https://res.cloudinary.com/demo/image/upload/c_fill,h_150,w_100/folder/sample.jpg"
	if u != want {
		t.Fatalf("want %s got %s", want, u)
	}
}

func TestSignParamsSHA1(t *testing.T) {
	c, _ := New("demo", "key", "abcd")
	sig, err := c.SignParams(map[string]any{"public_id": "sample", "timestamp": 1315060510})
	if err != nil {
		t.Fatal(err)
	}
	if sig != "c3470533147774275dd37996cc4d0e68fd03cd4f" {
		t.Fatalf("unexpected signature %s", sig)
	}
}

func TestSignParamsSHA256(t *testing.T) {
	c, _ := New("demo", "key", "abcd")
	c.SignatureAlgorithm = SHA256
	sig, err := c.SignParams(map[string]any{"public_id": "sample", "timestamp": 1315060510})
	if err != nil {
		t.Fatal(err)
	}
	if sig != "0d4fe14b2b4a3f68a97ccc5097c43908b623d24293c296826a9390c14d891509" {
		t.Fatalf("unexpected signature %s", sig)
	}
}

func TestVerifyNotificationSignatureRejectsExpired(t *testing.T) {
	if VerifyNotificationSignature("{}", "1", "bad", "secret", SHA1, time.Second) {
		t.Fatal("expected expired signature to fail")
	}
}
