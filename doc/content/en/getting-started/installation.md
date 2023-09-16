---
weight: 100
title: "Installation"
description: "Setup kube-resource-relabel-webhook on your infrastructure"
icon: "article"
draft: false
toc: true
---

# Helm Installation

### 1. Add the chart repository
```shell
helm repo add kube-resource-relabel-webhook https://pdylanross.github.io/kube-resource-relabel-webhook/helm-charts
helm repo update
```

### 2. Setup values based on your [desired configuration](/kube-resource-relabel-webhook/configuration)
### 3. Install and Enjoy!
```shell
helm install kube-resource-relabel-webhook/kube-resource-relabel-webhook \
  --create-namespace --namespace relabel-webhook \
  -f values.yaml relabel-webhook
```