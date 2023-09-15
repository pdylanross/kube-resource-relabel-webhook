THEME_VERSION := v0.41.1
THEME := hugo-geekdoc

.PHONY: doc-assets
doc-assets: themes/hugo-geekdoc

themes/hugo-geekdoc:
	mkdir -p themes/hugo-geekdoc/ ; \
	curl -sSL "https://github.com/thegeeklab/$(THEME)/releases/download/${THEME_VERSION}/$(THEME).tar.gz" | tar -xz -C themes/hugo-geekdoc/ --strip-components=1

.PHONY: clean
clean:
	rm -rf themes/hugo-geekdoc && \
	rm -rf public