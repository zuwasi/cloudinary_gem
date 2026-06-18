package cloudinary_gem

import "testing"

func TestTransformationStringExtended(t *testing.T) {
	got := TransformationString(URLOptions{Width: 100, Height: 150, Crop: "fill", Gravity: "face", Radius: "20", Effect: "sepia"})
	want := "w_100,h_150,c_fill,g_face,r_20,e_sepia"
	if got != want {
		t.Fatalf("want %s got %s", want, got)
	}
}

func TestPreloadedFileRoundTrip(t *testing.T) {
	sig := SignParams(map[string]any{"public_id": "folder/file", "version": 1234}, "secret", SHA1)
	p, err := ParsePreloadedFile("image/upload/v1234/folder/file.jpg#" + sig)
	if err != nil {
		t.Fatal(err)
	}
	if p.PublicID != "folder/file" || p.Format != "jpg" || p.Version != 1234 {
		t.Fatalf("unexpected preloaded file: %#v", p)
	}
	if !p.Verify("secret", SHA1) {
		t.Fatal("expected signature to verify")
	}
	if p.String() != "image/upload/v1234/folder/file.jpg#"+sig {
		t.Fatalf("round trip mismatch: %s", p.String())
	}
}

func TestGenerateAuthToken(t *testing.T) {
	got := GenerateAuthToken(AuthTokenOptions{Key: "secret", ACL: []string{"/image/*"}, StartTime: 100, Duration: 60})
	if got == "" || got[:13] != "__cld_token__" {
		t.Fatalf("unexpected token %q", got)
	}
}
