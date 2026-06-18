package cloudinary_gem

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

func BuildArray(values []string) string { return strings.Join(values, ",") }
func EncodeContext(values map[string]string) string {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+strings.ReplaceAll(values[k], "|", "\\|"))
	}
	return strings.Join(parts, "|")
}
func ResourceTypeForFormat(format string) string {
	switch strings.ToLower(strings.TrimPrefix(format, ".")) {
	case "jpg", "jpeg", "png", "gif", "webp", "bmp", "tiff", "ico", "svg", "avif", "heic":
		return "image"
	case "mp4", "webm", "mov", "avi", "m3u8", "mp3", "wav", "ogg":
		return "video"
	default:
		return "raw"
	}
}
func IsRemote(raw string) bool {
	return strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") || strings.HasPrefix(raw, "s3://") || strings.HasPrefix(raw, "gs://")
}
func RandomPublicID() string { return fmt.Sprintf("go_%d", time.Now().UnixNano()) }
func intString(v int) string { return fmt.Sprintf("%d", v) }
func nowUnix() int64         { return time.Now().Unix() }
