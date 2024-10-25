{{/* vim: set filetype=mustache: */}}

{{- define "kyverno.sidecar-injector.name" -}}
{{ template "kyverno.names.name" . }}-sidecar-injector
{{- end -}}

{{- define "kyverno.sidecar-injector.labels" -}}
{{- template "kyverno.labels.merge" (list
  (include "kyverno.labels.common" .)
  (include "kyverno.sidecar-injector.labels.match" .)
) -}}
{{- end -}}

{{- define "kyverno.sidecar-injector.labels.match" -}}
{{- template "kyverno.labels.merge" (list
  (include "kyverno.labels.match.common" .)
  (include "kyverno.labels.component" "sidecar-injector")
) -}}
{{- end -}}

{{- define "kyverno.sidecar-injector.role.name" -}}
{{- include "kyverno.names.fullname" . -}}:sidecar-injector
{{- end -}}

{{- define "kyverno.sidecar-injector.service-account.name" -}}
{{- if .Values.sidecarInjector.rbac.create -}}
    {{- default (include "kyverno.sidecar-injector.name" .) .Values.sidecarInjector.rbac.serviceAccount.name -}}
{{- else -}}
    {{- required "A service account name is required when `rbac.create` is set to `false`" .Values.sidecarInjector.rbac.serviceAccount.name -}}
{{- end -}}
{{- end -}}

{{- define "kyverno.sidecar-injector.serviceName" -}}
{{- printf "%s-svc" (include "kyverno.names.fullname" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "kyverno.sidecar-injector.caCertificatesConfigMapName" -}}
{{- printf "%s-ca-certificates" (include "kyverno.sidecar-injector.name" .) -}}
{{- end -}}

{{- define "kyverno.sidecar-injector.image" -}}
{{- printf "%s/%s:%s" .registry .repository (default "latest" .tag) -}}
{{- end -}}
