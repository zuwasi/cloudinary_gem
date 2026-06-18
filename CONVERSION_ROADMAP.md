# Ruby-to-Go Conversion Roadmap

Source repository: https://github.com/cloudinary/cloudinary_gem
Reference checkout: `C:\projects\ruby_to_golang\cloudinary_gem_original`
Target module: `github.com/zuwasi/cloudinary_gem`

## Goal

Convert the Cloudinary Ruby SDK into a Go SDK with equivalent core behavior where Go has an equivalent runtime concept. Ruby/Rails-specific integrations are converted into idiomatic Go equivalents or documented as not directly applicable.

## Source inventory

Original SDK library files: 42 Ruby/Rake files
Original SDK spec files: 58 Ruby spec files

## Conversion map

| Ruby area | Source files | Go target | Status |
|---|---|---|---|
| Base config | `base_config.rb`, `config.rb`, `account_config.rb` | `config.go` | implemented |
| Exceptions | `exceptions.rb` | Go `error` values/types | partial |
| URL/transformation utilities | `utils.rb`, `helper.rb`, `video_helper.rb` | `url.go`, Ruby parity tests | partial: distribution, suffix, root path, signing, and basic transformations |
| API signing/signature verification | `utils.rb`, `preloaded_file.rb`, `auth_token.rb` | `sign.go`, `auth_token.go`, `preloaded.go` | partial: SHA1/SHA256 URL signatures and basic request signatures |
| Upload API | `uploader.rb`, `base_api.rb` | `upload.go`, `client.go` | partial |
| Admin API | `api.rb`, `account_api.rb` | `admin.go`, `account.go` | partial |
| Search APIs | `search.rb`, `search_folders.rb` | `search.go` | not started |
| Downloader/static helpers | `downloader.rb`, `static.rb`, rake tasks | `download.go`, examples/CLI later | not started |
| Cache/breakpoints | `cache*.rb` | `cache.go` | not started |
| Analytics signature | `analytics.rb` | `analytics.go` | not started |
| Rails/ActiveStorage/CarrierWave | `active_storage/*`, `carrier_wave*`, `railtie.rb`, `engine.rb` | Go examples/interfaces where applicable | not directly portable |
| Tests | `spec/*` | Go `_test.go` files | partial |

## Definition of done for full conversion

1. Every Ruby library file is mapped to a Go implementation, Go equivalent, or explicit non-portable note.
2. Public SDK operations have Go tests derived from the Ruby specs.
3. `go test ./...` passes.
4. `govulncheck ./...` reports no vulnerabilities.
5. README documents Go API coverage and Ruby incompatibilities.

## Current limitation

This repository currently contains a partial Go SDK. It is not yet a complete feature-equivalent conversion of all 42 Ruby library files and 58 Ruby spec files.
