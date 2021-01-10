{{/* vim: set filetype=mustache: */}}

{{/*
    Define the Application Name
*/}}
{{- define "app.name" -}}
    {{- .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
    Create a Fully Qualified Application Name.
    By default, it uses the <Helm Release Name.Chart Name>
    However, if the Release Name contains the Chart Name, only the Release Name is Used.
*/}}
{{- define "app.fullname" -}}
    {{- if contains .Chart.Name .Release.Name -}}
        {{- .Release.Name -}}
    {{- else -}}
        {{- printf "%s-%s" .Chart.Name .Release.Name  | trunc 63 | trimSuffix "-" -}}
    {{- end -}}
{{- end -}}


{{/*
    Create chart name and version as used by the chart label.
*/}}
{{- define "app.chart" -}}
    {{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}