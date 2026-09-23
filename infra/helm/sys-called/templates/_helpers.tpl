{{- define "sys-called.name" -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end }}

{{- define "sys-called.labels" -}}
app.kubernetes.io/part-of: sys-called
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version }}
{{- end }}

{{- define "sys-called.selector" -}}
app.kubernetes.io/instance: {{ .root.Release.Name }}
app.kubernetes.io/name: {{ .component }}
{{- end }}

{{/* A digest pins the exact image CI tested and signed, and wins over the tag. */}}
{{- define "sys-called.image" -}}
{{- if .Values.image.digest -}}
{{ printf "%s@%s" .Values.image.repository .Values.image.digest }}
{{- else -}}
{{ printf "%s:%s" .Values.image.repository (default .Chart.AppVersion .Values.image.tag | toString) }}
{{- end -}}
{{- end }}

{{- define "sys-called.secretName" -}}
{{ default (printf "%s-secrets" (include "sys-called.name" .)) .Values.secrets.existingSecret }}
{{- end }}

{{- define "sys-called.postgresName" -}}
{{ printf "%s-postgres" (include "sys-called.name" .) }}
{{- end }}

{{/* More than one pod may serve at once. */}}
{{- define "sys-called.multiReplica" -}}
{{- if or .Values.autoscaling.enabled (gt (int .Values.app.replicas) 1) -}}true{{- end -}}
{{- end }}

{{- define "sys-called.ticketsCache" -}}
{{- if eq (toString .Values.app.ticketsCache) "auto" -}}
{{ ternary "false" "true" (eq (include "sys-called.multiReplica" .) "true") }}
{{- else -}}
{{ .Values.app.ticketsCache | toString }}
{{- end -}}
{{- end }}

{{/*
Keeps a generated secret value stable across upgrades: an explicit value wins,
then the one already in the cluster, and only then a newly generated one.
*/}}
{{- define "sys-called.secretValue" -}}
{{- $existing := lookup "v1" "Secret" .root.Release.Namespace (include "sys-called.secretName" .root) -}}
{{- if .value -}}
{{ .value | b64enc }}
{{- else if and $existing (index $existing.data .key) -}}
{{ index $existing.data .key }}
{{- else -}}
{{ .generate | b64enc }}
{{- end -}}
{{- end }}

{{- define "sys-called.podSecurityContext" -}}
runAsNonRoot: true
runAsUser: 65532
runAsGroup: 65532
fsGroup: 65532
seccompProfile:
  type: RuntimeDefault
{{- end }}

{{- define "sys-called.containerSecurityContext" -}}
allowPrivilegeEscalation: false
readOnlyRootFilesystem: true
capabilities:
  drop: ["ALL"]
{{- end }}

{{/* Environment shared by the migrate init container and the app. */}}
{{- define "sys-called.env" -}}
- name: APP_ENV
  value: {{ .Values.app.env | quote }}
- name: LOG_LEVEL
  value: {{ .Values.app.logLevel | quote }}
- name: DB_MAX_CONNS
  value: {{ .Values.app.dbMaxConns | quote }}
{{- if .Values.postgresql.enabled }}
- name: DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "sys-called.secretName" . }}
      key: postgres-password
- name: DATABASE_URL
  value: {{ printf "postgres://%s:$(DB_PASSWORD)@%s:5432/%s?sslmode=disable" .Values.postgresql.user (include "sys-called.postgresName" .) .Values.postgresql.database | quote }}
{{- else }}
- name: DATABASE_URL
  valueFrom:
    secretKeyRef:
      name: {{ required "externalDatabase.existingSecret is required when postgresql.enabled=false" .Values.externalDatabase.existingSecret }}
      key: {{ .Values.externalDatabase.urlKey }}
{{- end }}
{{- end }}
