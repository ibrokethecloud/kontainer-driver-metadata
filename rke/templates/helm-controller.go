package templates

const HelmControllerTemplate = `
{{- if eq .Scope "namespace" }}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: helm
  namespace: helm-controller
---
apiVersion: v1
kind: Namespace
metadata:
  name: helm-controller
  labels:
    name: helm-controller
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: helm-role-binding
  namespace: helm-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: admin
subjects:
- namespace: helm-controller
  kind: ServiceAccount
  name: helm
{{- else }}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: helm
  namespace: kube-system 
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: helm-role-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- namespace: kube-system
  kind: ServiceAccount
  name: helm
{{- end }}
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: helmcharts.helm.cattle.io
  namespace: {{ if eq .Scope "cluster" }}kube-system{{- else -}}helm-controller{{- end }}
spec:
  group: helm.cattle.io
  version: v1
  names:
    kind: HelmChart
    plural: helmcharts
    singular: helmchart
  scope: {{ if eq .Scope "namespace" -}}Namespaced{{- else -}}Cluster{{- end }}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: helm-controller
  namespace:  {{ if eq .Scope "cluster" }}kube-system{{- else -}}helm-controller{{- end }}
  labels:
    app: helm-controller
spec:
  replicas: 1
  selector:
    matchLabels:
      app: helm-controller
  template:
    metadata:
      labels:
        app: helm-controller
    spec:
      serviceAccountName: helm
      containers:
        - name: helm-controller
          image: {{ .Image }}
          command: ["helm-controller"]
{{- if eq .Scope "namespace" }}
          args: ["--namespace", "helm-controller"]
{{ end -}}
`