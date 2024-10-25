{{/* vim: set filetype=mustache: */}}

{{- define "kyverno.lib.chart.name" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "kyverno.lib.chart.version" -}}
  {{- .Chart.Version | replace "+" "_" -}}
{{- end -}}
