package cloudinary_gem

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type UploadResult map[string]any

func (c *Cloudinary) Upload(ctx context.Context, file io.Reader, filename string, params map[string]any) (UploadResult, error) {
	if params == nil {
		params = map[string]any{}
	}
	if _, ok := params["timestamp"]; !ok {
		params["timestamp"] = time.Now().Unix()
	}
	params["api_key"] = c.APIKey
	sig, err := c.SignParams(params)
	if err != nil {
		return nil, err
	}
	params["signature"] = sig
	body, pw := io.Pipe()
	mw := multipart.NewWriter(pw)
	go func() {
		defer pw.Close()
		defer mw.Close()
		fw, err := mw.CreateFormFile("file", filename)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		if _, err := io.Copy(fw, file); err != nil {
			pw.CloseWithError(err)
			return
		}
		for k, v := range params {
			_ = mw.WriteField(k, fmt.Sprint(v))
		}
	}()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.uploadURL(defaultString(fmt.Sprint(params["resource_type"]), "image")), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	return c.doJSON(req)
}

func (c *Cloudinary) UnsignedUpload(ctx context.Context, file io.Reader, filename, uploadPreset string, params map[string]any) (UploadResult, error) {
	if params == nil {
		params = map[string]any{}
	}
	params["upload_preset"] = uploadPreset
	body, pw := io.Pipe()
	mw := multipart.NewWriter(pw)
	go func() {
		defer pw.Close()
		defer mw.Close()
		fw, err := mw.CreateFormFile("file", filename)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		if _, err := io.Copy(fw, file); err != nil {
			pw.CloseWithError(err)
			return
		}
		for k, v := range params {
			_ = mw.WriteField(k, fmt.Sprint(v))
		}
	}()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.uploadURL("image"), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	return c.doJSON(req)
}

func (c *Cloudinary) uploadURL(resourceType string) string {
	return fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/%s/upload", c.CloudName, defaultString(resourceType, "image"))
}
