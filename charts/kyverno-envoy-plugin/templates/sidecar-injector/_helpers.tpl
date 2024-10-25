{{/* vim: set filetype=mustache: */}}

{{- define "sidecar-injector.name" -}}
{{ template "kyverno.lib.names.name" . }}-sidecar-injector
{{- end -}}

{{- define "sidecar-injector.labels" -}}
{{- template "kyverno.lib.labels.merge" (list
  (include "kyverno.lib.labels.common" .)
  (include "sidecar-injector.labels.match" .)
) -}}
{{- end -}}

{{- define "sidecar-injector.labels.match" -}}
{{- template "kyverno.lib.labels.merge" (list
  (include "kyverno.lib.labels.common.selector" .)
  (include "kyverno.lib.labels.component" "sidecar-injector")
) -}}
{{- end -}}

{{- define "sidecar-injector.role.name" -}}
{{- include "kyverno.lib.names.fullname" . -}}:sidecar-injector
{{- end -}}

{{- define "sidecar-injector.service-account.name" -}}
{{- if .Values.sidecarInjector.rbac.create -}}
    {{- default (include "sidecar-injector.name" .) .Values.sidecarInjector.rbac.serviceAccount.name -}}
{{- else -}}
    {{- required "A service account name is required when `rbac.create` is set to `false`" .Values.sidecarInjector.rbac.serviceAccount.name -}}
{{- end -}}
{{- end -}}

{{- define "sidecar-injector.serviceName" -}}
{{- printf "%s-svc" (include "kyverno.lib.names.fullname" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "sidecar-injector.image" -}}
{{- printf "%s/%s:%s" .registry .repository (default "latest" .tag) -}}
{{- end -}}
