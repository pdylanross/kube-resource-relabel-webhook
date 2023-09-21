---
weight: 210
title: "Rules"
---

Relabel Rules consist of Conditions and Actions. If all conditions are met, then the actions are performed. 

{{<expand "Rules Schema">}}

| Name       | Type        | Description                                                        |
|------------|-------------|--------------------------------------------------------------------|
| name       | string      | A unique name for the rule. Used mainly for debugging              |
| conditions | condition[] | The list of conditions that must be satisfied for this rule        |
| actions    | action[]    | The list of actions that will be performed when conditions are met |

{{<expand "Examples">}}
Add karpenter eviction prevention to airflow pods.
```yaml
name: "airflow-k8s-pod-operator-do-not-evict"
conditions:
  - type: is-type
    value:
      kind: pod
      version: v1
  - type: has-label
    value:
      keys:
        - dag_id
actions:
  - type: ensure-annotation
    value:
      karpenter.sh/do-not-evict: "true"
```
Ensure that label based logging systems have the correct label. 
```yaml
name: "cleanup-fluentd-labels"
conditions:
  - type: is-type
    value:
      kind: pod
      version: v1
  - type: has-annotation
    value:
      match:
        fluentd.active: "true"
actions:
  - type: ensure-label
    value:
      fluentd.active: "true"
```
{{</expand>}}
{{</expand>}}

{{<expand "Condition Schema">}}

| Name  | Type   | Description                                              |
|-------|--------|----------------------------------------------------------|
| type  | string | Type of condition                                        |
| value | any    | Condition configuration. This changes per condition type |

{{<expand "Examples">}}
Ensure the resource is a pod.
```yaml
type: is-type
value:
  kind: pod
  version: v1
```
Check that a resource has a label key
```yaml
type: has-label
value:
  keys:
    - dag_id
```
Check if a resource has an exact annotation match
```yaml
type: has-annotation
value:
  match:
    fluentd.active
```
{{</expand>}}
{{</expand>}}

{{<expand "Action Schema">}}

| Name  | Type   | Description                                        |
|-------|--------|----------------------------------------------------|
| type  | string | Type of action                                     |
| value | any    | Action configuration. This changes per action type |


{{<expand "Examples">}}
Ensure the resource has a label
```yaml
type: ensure-label
value:
  fluentd.active: "true"
```
Ensure the resource has an annotation
```yaml
type: ensure-annotation
value:
  nginx.ingress.kubernetes.io/default-backend: some-svc
```
{{</expand>}}
{{</expand>}}

# Condition Types

{{<expand "Check Resource Type">}}
#### type: `is-type`
#### Values Schema
| Name    | Type   | Descripton                   |
|---------|--------|------------------------------|
| group   | string | The kubernetes api group     |
| version | string | The kubernetes api version   |
| kind    | string | The kubernetes resource kind |

#### About
Check if the resource is of a specific kubernetes type. Missing / empty fields always evaluate to true.

{{<expand "Examples">}}
Ensure the resource is a pod.
```yaml
type: is-type
value:
  kind: pod
  version: v1
```
Any ingress resource.
```yaml
type: is-type
value:
  kind: ingress
```
All Istio networking resources
```yaml
type: is-type
value:
  group: networking.istio.io
```
{{</expand>}}
{{</expand>}}

{{<expand "Has Annotation">}}
#### type: `has-annotation`
#### Values Schema
| Name   | Type              | Descripton                         |
|--------|-------------------|------------------------------------|
| keys   | string[]          | Match on annotation keys           |
| values | string[]          | Match on annotation values         |
| kind   | map[string]string | Exact key:value annotation matches |

#### About
Check if the resource has specific annotations. Empty or null values are ignored. 

{{<expand "Examples">}}
Exact key-value match
```yaml
type: has-annotation
value:
  match:
    fluentd.active: "true"
```
Match on annotation key
```yaml
type: has-annotation
value:
  keys:
    - fluentd.active
```
Match on annotation value
```yaml
type: has-annotation
value:
  values:
    - "true"
```
{{</expand>}}
{{</expand>}}

{{<expand "Has Label">}}
#### type: `has-label`
#### Values Schema
| Name   | Type              | Descripton                    |
|--------|-------------------|-------------------------------|
| keys   | string[]          | Match on label keys           |
| values | string[]          | Match on label values         |
| kind   | map[string]string | Exact key:value label matches |

#### About
Check if the resource has specific labels. Empty or null values are ignored.

{{<expand "Examples">}}
Exact key-value match
```yaml
type: has-label
value:
  match:
    fluentd.active: "true"
```
Match on label key
```yaml
type: has-label
value:
  keys:
    - fluentd.active
```
Match on label value
```yaml
type: has-label
value:
  values:
    - "true"
```
{{</expand>}}
{{</expand>}}

## Actions

{{<expand "Ensure Annotation">}}
#### type: `ensure-annotation`
#### Values: `map[string]string`

#### About
Apply the given annotations to the object

{{<expand "Examples">}}
Exact key-value match
```yaml
type: ensure-annotation
value:
  karpenter.sh/do-not-evict: "true"
```
{{</expand>}}
{{</expand>}}

{{<expand "Ensure Label">}}
#### type: `ensure-label`
#### Values: `map[string]string`

#### About
Apply the given labels to the object

{{<expand "Examples">}}
Exact key-value match
```yaml
type: ensure-label
value:
  fluentd.active: "true"
```
{{</expand>}}
{{</expand>}}