package cloudinary_gem

import "testing"

func TestTransformationChain(t *testing.T) {
	got := GenerateTransformationString(TransformationChain{{"x": 100, "y": 100, "width": 200, "crop": "fill"}, {"radius": 10}, {"transformation": []string{"blip", "blop"}}})
	want := "c_fill,w_200,x_100,y_100/r_10/t_blip.blop"
	if got != want {
		t.Fatalf("want %s got %s", want, got)
	}
}

func TestURLWithTransformations(t *testing.T) {
	c, _ := New("demo", "key", "secret")
	u, err := c.URLWithTransformations("sample", TransformationChain{{"x": 100, "y": 100, "width": 200, "crop": "fill"}, {"radius": 10}}, URLOptions{Crop: "crop", Width: 100})
	if err != nil {
		t.Fatal(err)
	}
	want := "https://res.cloudinary.com/demo/image/upload/c_fill,w_200,x_100,y_100/r_10/c_crop,w_100/sample"
	if u != want {
		t.Fatalf("want %s got %s", want, u)
	}
}

func TestSrcSet(t *testing.T) {
	c, _ := New("demo", "key", "secret")
	s, err := SrcSet(c, "sample", []SrcSetBreakpoint{{Width: 200}, {Width: 100}}, URLOptions{Crop: "scale", Format: "jpg"})
	if err != nil {
		t.Fatal(err)
	}
	want := "https://res.cloudinary.com/demo/image/upload/c_scale,w_100/sample.jpg 100w,https://res.cloudinary.com/demo/image/upload/c_scale,w_200/sample.jpg 200w"
	if s != want {
		t.Fatalf("want %s got %s", want, s)
	}
}
