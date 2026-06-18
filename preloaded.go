package cloudinary_gem

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type PreloadedFile struct {
	ResourceType string
	Type         string
	Version      int64
	PublicID     string
	Format       string
	Signature    string
}

var preloadedRE = regexp.MustCompile(`^([^/]+)/([^/]+)/v([0-9]+)/(.+)#([^#]+)$`)

func ParsePreloadedFile(s string) (*PreloadedFile, error) {
	m := preloadedRE.FindStringSubmatch(s)
	if m == nil {
		return nil, errors.New("invalid preloaded file")
	}
	ver, _ := strconv.ParseInt(m[3], 10, 64)
	id := m[4]
	format := ""
	if i := strings.LastIndex(id, "."); i >= 0 {
		format = id[i+1:]
		id = id[:i]
	}
	return &PreloadedFile{ResourceType: m[1], Type: m[2], Version: ver, PublicID: id, Format: format, Signature: m[5]}, nil
}

func (p PreloadedFile) Verify(secret string, algo Algorithm) bool {
	return VerifyAPIResponseSignature(p.PublicID, p.Version, p.Signature, secret, algo)
}
func (p PreloadedFile) String() string {
	id := p.PublicID
	if p.Format != "" {
		id += "." + p.Format
	}
	return fmt.Sprintf("%s/%s/v%d/%s#%s", p.ResourceType, p.Type, p.Version, id, p.Signature)
}
