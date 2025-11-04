{{/* vim: set filetype=mustache: */}}

{{- define "validation-webhook.name" -}}
{{ template "kyverno.lib.names.name" . }}-validation
{{- end -}}

{{- define "validation-webhook.service-account.name" -}}
{{- if .Values.rbac.create -}}
  {{- default (include "validation-webhook.name" .) .Values.rbac.serviceAccount.name -}}
{{- else -}}
  {{- required "A service account name is required when `rbac.create` is set to `false`" .Values.rbac.serviceAccount.name -}}
{{- end -}}
{{- end -}}

{{- define "validation-webhook.image" -}}
{{- printf "%s/%s:%s" .registry .repository (default "latest" .tag) -}}
{{- end -}}
