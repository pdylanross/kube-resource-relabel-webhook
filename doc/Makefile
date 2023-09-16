THEME_VERSION := v0.41.1
THEME := hugo-geekdoc

.PHONY: setup
setup: doc-assets template-pages

.PHONY: doc-assets
doc-assets: themes/hugo-geekdoc

themes/hugo-geekdoc:
	mkdir -p themes/hugo-geekdoc/ ; \
	curl -sSL "https://github.com/thegeeklab/$(THEME)/releases/download/${THEME_VERSION}/$(THEME).tar.gz" | tar -xz -C themes/hugo-geekdoc/ --strip-components=1

.PHONY: clean
clean:
	rm -rf themes/hugo-geekdoc && \
	rm -rf public

	rm -f ./content/contribution.md \
		./content/code-of-conduct.md \
		./content/changelog.md \

.PHONY: template-pages
template-pages:
	gomplate --include *.md \
		--input-dir content_templates --output-dir content \
		--datasource contributing=../CONTRIBUTING.md?type=text/plain \
		--datasource code-of-conduct=../CODE_OF_CONDUCT.md?type=text/plain \
		--datasource changelog=../CHANGELOG.md?type=text/plain

.PHONY: serve
serve: setup
	hugo serve