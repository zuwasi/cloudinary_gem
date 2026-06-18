# cloudinary_gem for Go

A Go port of the core Cloudinary Ruby SDK behaviors: configuration from `CLOUDINARY_URL`, delivery URL generation, request signing, uploads, and Admin API calls.

This repository keeps the original Cloudinary Ruby SDK MIT license.

## Install

```bash
go get github.com/zuwasi/cloudinary_gem
```

## Usage

```go
cld, err := cloudinary_gem.NewFromURL("cloudinary://api_key:api_secret@demo")
if err != nil { panic(err) }

url, _ := cld.URL("sample.jpg", cloudinary_gem.URLOptions{
    Width: 100,
    Height: 150,
    Crop: "fill",
    Format: "jpg",
})
```
