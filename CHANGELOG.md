# Changelog

## 1.5.0-r0 (2022/02/21)

* GoReleaser 1.5.0 (#102)
* Bump github.com/alecthomas/kong from 0.2.22 to 0.4.1 (#93 #101)

## 1.2.5-r1 (2022/01/09)

* Improve CGO handling (#95)
* Handle `GOMIPS64` (#94)
* Fix `GIT_REF`

## 1.2.5-r0 (2022/01/02)

* GoReleaser 1.2.5 (#90)
* jq example (#89)
* Handle snaps and brews (#87)
* Handle alt packages (#88)
* Update goxx example (#86)
* Override compilers (#85)

## 1.2.4-r0 (2021/12/31)

* GoReleaser 1.2.4 (#84)
* Update CGO examples with [goxx](https://github.com/crazy-max/goxx) (#83)

## 1.2.2-r2 (2021/12/29)

* Lookup GoReleaser binary path (#82)
* Update demos (#81)
* Handle C and C++ compilers (#80)
* Fix binary artifact output
* Fix artifact extension
* Move syscall to golang.org/x/sys
* Move from io/ioutil to os package

## 1.2.2-r1 (2021/12/25)

* Note about `CGO_ENABLED` (#79)
* Display `go env` on debug (#78)

## 1.2.2-r0 (2021/12/24)

* GoReleaser 1.2.2 (#76)
* Allow specifying GoReleaser yaml config (#74)

## 1.1.0-r5 (2021/12/21)

* Allow list of artifact types (#70 #72)
* `--artifact-type` deprecated. Use `--artifacts` instead (#72)
* `--hooks` deprecated. Use `--pre-hooks` instead (#73)
* `--build-pre-hooks` deprecated. Use `--pre-hooks` instead (#73)
* `--build-post-hooks` deprecated. Use `--post-hooks` instead (#73)

## 1.1.0-r4 (2021/12/20)

* Fix artifact version (#68)
* CGO usage section
* Add asmflags, gcflags and tags options (#67)
* More demos (#66)
* Add go-binary option (#65)

## 1.1.0-r3 (2021/12/18)

* Do not set flags and ldflags if empty
* Fix flags type
* Demo app (#64)

## 1.1.0-r2 (2021/12/18)

* Add flags option (#63)
* Bump github.com/alecthomas/kong from 0.2.20 to 0.2.22 (#62)

## 1.1.0-r1 (2021/12/12)

* Fix typo `--build-pre-hoosk` (#60)

## 1.1.0-r0 (2021/12/12)

* GoReleaser 1.1.0 (#59)
* Add build pre and post hooks options (#58)
* Bump github.com/alecthomas/kong from 0.2.18 to 0.2.20 (#56)

## 1.0.0-r0 (2021/11/17)

* GoReleaser 1.0.0 (#53)
* xx 1.0.0 (#52)
* Bump github.com/alecthomas/kong from 0.2.17 to 0.2.18 (#51)

## 0.180.0-r0 (2021/09/26)

* GoReleaser 0.180.0 (#43)

## 0.175.0-r0 (2021/08/21)

* GoReleaser 0.175.0 (#38)
* Go 1.17 (#39)

## 0.173.2-r0 (2021/07/11)

* GoReleaser 0.173.2 (#34)

## 0.169.0-r0 (2021/06/18)

* GoReleaser 0.169.0 (#29)
* Bump github.com/alecthomas/kong from 0.2.16 to 0.2.17 (#28)

## 0.166.2-r0 (2021/05/31)

* GoReleaser 0.166.2 (#24)

## 0.165.0-r0 (2021/05/25)

* GoReleaser 0.165.0 (#20)

## 0.164.0-r4 (2021/05/08)

* Add `replacements` option
* `artifact-type` can be `archive` or `bin`

## 0.164.0-r3 (2021/05/08)

* Skip changelog

## 0.164.0-r2 (2021/05/07)

* Add `artifact-type` and `checksum` options (#11)

## 0.164.0-r1 (2021/05/06)

* Add `env` option

## 0.164.0-r0 (2021/04/25)

* GoReleaser 0.164.0

## 0.162.0-r0 (2021/04/08)

* GoReleaser 0.162.0

## 0.161.1-r0 (2021/03/25)

* GoReleaser 0.161.1

## 0.160.0-r0 (2021/03/21)

* GoReleaser 0.160.0
* Disable checksum
* Bump github.com/alecthomas/kong from 0.2.15 to 0.2.16 (#3)

## 0.159.0-r3 (2021/03/07)

* Do not fail on Git tag not found (#2)
* Improve usage examples

## 0.159.0-r2 (2021/03/07)

* Fix `GORELEASER_CURRENT_TAG`

## 0.159.0-r1 (2021/03/06)

* Add snapshot option
* Handle wrong ref

## 0.159.0-r0 (2021/03/06)

* Initial version with [GoReleaser](https://github.com/goreleaser/goreleaser) 0.159.0
