# kube-resource-relabel-webhook
[![Go Report Card](https://goreportcard.com/badge/github.com/pdylanross/kube-resource-relabel-webhook)](https://goreportcard.com/report/github.com/pdylanross/kube-resource-relabel-webhook)
![GitHub](https://img.shields.io/github/license/pdylanross/kube-resource-relabel-webhook)
![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/pdylanross/kube-resource-relabel-webhook/ci.yml)
![GitHub release (with filter)](https://img.shields.io/github/v/release/pdylanross/kube-resource-relabel-webhook)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](code_of_conduct.md)


Have you ever been a platform engineer implementing a new component only to realize that you then needed to go add a ton of specific labels to various workloads owned by over one hundred different maintainers?

I have. I'm lazy. I'd rather automate the pain away.

kube-resource-relabel-webhook does exactly as it sounds like - it's a kubernetes mutation webhook that adds new labels and annotations based on a user provided rules-based configuration.

### Read more about it on the [documentation site](https://pdylanross.github.io/kube-resource-relabel-webhook/). 

## Contribution

Please, thx. 

[Read the guide](CONTRIBUTING.md)

## License

[Unlicense](LICENSE). Or - "I dont care, just don't sue me'