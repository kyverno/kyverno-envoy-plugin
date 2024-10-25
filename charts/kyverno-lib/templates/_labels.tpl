{{/* vim: set filetype=mustache: */}}

{{- define "kyverno.lib.labels.merge" -}}
{{- $labels := dict -}}
{{- range . -}}
  {{- $labels = merge $labels (fromYaml .) -}}
{{- end -}}
{{- with $labels -}}
  {{- toYaml $labels -}}
{{- end -}}
{{- end -}}

{{- define "kyverno.lib.labels.helm" -}}
helm.sh/chart: {{ template "kyverno.lib.chart.name" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{- define "kyverno.lib.labels.version" -}}
app.kubernetes.io/version: {{ template "kyverno.lib.chart.version" . }}
{{- end -}}

{{- define "kyverno.lib.labels.component" -}}
app.kubernetes.io/component: {{ . }}
{{- end -}}

{{- define "kyverno.lib.labels.name" -}}
app.kubernetes.io/name: {{ . }}
{{- end -}}

{{- define "kyverno.lib.labels.common" -}}
{{- template "kyverno.lib.labels.merge" (list
  (include "kyverno.lib.labels.helm" .)
  (include "kyverno.lib.labels.version" .)
) -}}
{{- end -}}

{{- define "kyverno.lib.labels.common.selector" -}}
app.kubernetes.io/part-of: {{ template "kyverno.lib.names.fullname" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}
