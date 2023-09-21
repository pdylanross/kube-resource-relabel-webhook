# Changelog

## [0.3.0](https://github.com/pdylanross/kube-resource-relabel-webhook/compare/v0.2.5...v0.3.0) (2023-09-21)


### Features

* added cert-manager support ([c995df4](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/c995df416df91fb6112c6eaf98d303f8b3ed1b59))
* added is-type condition ([09ee7c2](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/09ee7c2d0c9ec0f34686876971ac96eb8bbff870))
* **doc:** fleshed out docs around helm & config ([df5a67f](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/df5a67f05a2afab567f40a987d1a71db2d00d08e))
* **test:** added testing for ingress resources ([2110fdb](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/2110fdb8a990381cbb3282cae4cd08151ea4343f))


### Bug Fixes

* **integration-test:** add endpoint readiness gate ([7042069](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/70420690f173a71f31d9b227b1dca27b1185a8e2))
* linter errors ([ed24dd1](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/ed24dd1c520657d53b5affdac607cc643086ac5e))
* multiple patches overriding each other ([dcd476e](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/dcd476ef31b0050d2460f4bc3bf27a352f16b3f9))

## [0.2.5](https://github.com/pdylanross/kube-resource-relabel-webhook/compare/v0.2.4...v0.2.5) (2023-09-15)


### Bug Fixes

* **integration-tests:** add small delay after helm install ([e94008d](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/e94008db34e1ded5a87955dfb16bdbbe31a20375))
* permissions for release pipeline ([a4b3124](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/a4b31249f10f226f909b338b350d60790ceccdf6))

## [0.2.4](https://github.com/pdylanross/kube-resource-relabel-webhook/compare/v0.2.3...v0.2.4) (2023-09-15)


### Bug Fixes

* **ci:** use docker build and push ([c85ed70](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/c85ed70093d049ad10b4bf15bd984a4d91209829))
* cleanup docs ([93d29be](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/93d29befad35727b4151e6a76396285df822fec4))
* optimized multiplatform container build ([0c80d60](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/0c80d600935b07205f0e1307d25a951387ad1b2d))
* re-added build target for integration tests ([588a10b](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/588a10bd2e9694062c251e525623e2b15fb432ac))
* set chart version with release pls ([c85ed70](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/c85ed70093d049ad10b4bf15bd984a4d91209829))

## [0.2.3](https://github.com/pdylanross/kube-resource-relabel-webhook/compare/v0.2.2...v0.2.3) (2023-09-15)


### Bug Fixes

* dockerfile cache go mod correctly ([a68307b](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/a68307b44306e941c138c1cd0208b1a6add0b4a0))


## [0.2.0](https://github.com/pdylanross/kube-resource-relabel-webhook/compare/v0.1.3...v0.2.0) (2023-09-14)


### Features

* added helm release ([410aba7](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/410aba7155bacb6bc10dbac93ce2b1f60f75e716))


## [0.1.1](https://github.com/pdylanross/kube-resource-relabel-webhook/compare/v0.1.0...v0.1.1) (2023-09-09)


### Bug Fixes

* added on tag workflow to push images ([2a61838](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/2a618384f5593bb38c28c75ea62f86903c6ca3eb))

## [0.1.0](https://github.com/pdylanross/kube-resource-relabel-webhook/compare/v0.0.1...v0.1.0) (2023-09-09)


### Features

* added ensure-label action ([ec4f150](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/ec4f15061e40cba35afbd2a033c29a3c3d5801b8))
* added integration tests ([818ea1d](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/818ea1d665e4dbe11e2b9e4bf70fab906fe60828))
* created helm chart ([17ed25d](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/17ed25d176364c1d9187a225d04d45bf780f9809))
* initial server and logging implementation ([49464fd](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/49464fdcda43e9ae9e37a03cc68b25fc75705b30))
* relabel config and basic structure laid out ([af8ef06](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/af8ef065daddd62b1be12b71bca63e505a8afc63))
* webhook handler implementation ([2519fbd](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/2519fbd0778265c0d432629cfd79086d19626956))


### Bug Fixes

* cleaned up linter warnings ([fa68491](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/fa6849108d38b45d2907ed27b75e80096eb418ea))
* rename sanitizeKeyForJSONPatch ([2b82ede](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/2b82ede1a6ddb7aee29618964a4e124380230551))
* use full 1.21.0 for go mod ([80ba26a](https://github.com/pdylanross/kube-resource-relabel-webhook/commit/80ba26a5e08f8a5b1b010b3fd253594d6317e6a8))
