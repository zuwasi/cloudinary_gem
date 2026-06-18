# 100% Replacement Checklist

This checklist tracks the original Ruby SDK public areas and their Go replacements.

## Core SDK
- [x] Config from URL/env
- [x] Errors/API errors
- [x] URL generation core
- [x] URL signing SHA1/SHA256 short/long
- [~] Full transformation engine
- [x] API request signing
- [x] API response/notification verification
- [x] Preloaded file parsing
- [~] Auth token generation
- [~] Analytics signature parity
- [ ] Responsive helpers
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
- [~] metadata fields/rules wrappers
- [~] folders wrappers; upload mappings pending
- [~] streaming profiles wrappers
- [~] analysis wrapper; access mode/related assets/publish pending
- [ ] account provisioning API

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
- [~] Ruby utils specs translated
- [ ] Uploader specs translated with httptest
- [ ] Admin/search specs translated with httptest
- [ ] Helpers/cache/auth token specs translated
- [x] go test ./...
- [x] govulncheck ./...
