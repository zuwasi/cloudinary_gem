# 100% Replacement Checklist

This checklist tracks the original Ruby SDK public areas and their Go replacements.

## Core SDK
- [x] Config from URL/env
- [x] Errors/API errors
- [x] URL generation core
- [x] URL signing SHA1/SHA256 short/long
- [~] Full transformation engine: chaining/named/basic option maps implemented; advanced layers/conditionals pending
- [x] API request signing
- [x] API response/notification verification
- [x] Preloaded file parsing
- [~] Auth token generation
- [~] Analytics signature parity
- [~] Responsive helpers
- [~] Download/archive URL helpers

## Upload API
- [x] signed upload
- [x] unsigned upload
- [ ] upload_large/chunking
- [~] destroy
- [~] rename
- [~] explicit
- [~] archive/zip implemented; slideshow/sprite/multi/explode pending
- [~] tags/context implemented; metadata pending
- [ ] upload parameter builder parity

## Admin API
- [x] generic Admin method
- [~] resources/tags/transformations/upload presets
- [x] metadata fields/rules wrappers
- [x] folders and upload mappings wrappers
- [x] streaming profiles wrappers
- [x] analysis/access mode/related assets/publish wrappers
- [~] account provisioning API wrappers

## Query APIs
- [x] Search builder
- [x] Search folders builder

## Framework integrations
- [~] Rails helpers converted to basic Go HTML/template helpers
- [~] Video helpers converted to basic Go HTML/template helpers
- [ ] ActiveStorage equivalent documented/implemented as Go storage interface
- [ ] CarrierWave equivalent documented/implemented as Go uploader interface
- [x] Cache interfaces and in-memory implementation

## Tests/security
- [~] Ruby utils specs translated: URL/signing/chaining/responsive slices
- [~] Uploader specs translated with httptest
- [~] Admin/search specs translated with httptest
- [ ] Helpers/cache/auth token specs translated
- [x] go test ./...
- [x] govulncheck ./...
