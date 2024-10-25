{{/* vim: set filetype=mustache: */}}

{{- define "kyverno.labels.merge" -}}
{{- $labels := dict -}}
{{- range . -}}
  {{- $labels = merge $labels (fromYaml .) -}}
{{- end -}}
{{- with $labels -}}
  {{- toYaml $labels -}}
{{- end -}}
{{- end -}}

{{- define "kyverno.labels.helm" -}}
{{- if not .Values.templating.enabled -}}
helm.sh/chart: {{ template "kyverno.chart.name" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}
{{- end -}}

{{- define "kyverno.labels.version" -}}
app.kubernetes.io/version: {{ template "kyverno.chart.version" . }}
{{- end -}}

{{- define "kyverno.labels.common" -}}
{{- template "kyverno.labels.merge" (list
  (include "kyverno.labels.helm" .)
  (include "kyverno.labels.version" .)
  (toYaml .Values.customLabels)
) -}}
{{- end -}}

{{- define "kyverno.labels.component" -}}
app.kubernetes.io/component: {{ . }}
{{- end -}}

{{- define "kyverno.labels.name" -}}
app.kubernetes.io/name: {{ . }}
{{- end -}}

{{- define "kyverno.labels.match.common" -}}
app.kubernetes.io/part-of: {{ template "kyverno.names.fullname" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}
