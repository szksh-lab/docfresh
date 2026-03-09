{{if not .HideCommand -}}
{{if .Command -}}
{{codeFence .Command}}{{.CommandLanguage | default "sh"}}
{{trimSuffix "\n" .Command}}
{{codeFence .Command}}
{{- else if .EmbedScript -}}
{{codeFence .Content}}{{.CommandLanguage | default "sh"}}
{{trimSuffix "\n" .Content}}
{{codeFence .Content}}
{{- else -}}
{{codeFence .Script}}{{.CommandLanguage | default "sh"}}
{{join " " .Shell}} {{trimSuffix "\n" .Script}}
{{codeFence .Script}}
{{- end}}

{{end -}}
{{if not .HideOutput -}}
{{if .DetailsTagSummary -}}
<details>
<summary>{{.DetailsTagSummary}}</summary>

{{if .CodeBlock -}}
{{codeFence .CombinedOutput}}{{.OutputLanguage}}
{{trimSuffix "\n" .CombinedOutput}}
{{codeFence .CombinedOutput}}
{{- else -}}
{{trimSuffix "\n" .CombinedOutput}}
{{- end}}

</details>
{{- else -}}
{{if .CodeBlock -}}
{{if not .HideCommand -}}
Output:

{{end -}}
{{codeFence .CombinedOutput}}{{.OutputLanguage}}
{{trimSuffix "\n" .CombinedOutput}}
{{codeFence .CombinedOutput}}
{{- else -}}
{{trimSuffix "\n" .CombinedOutput}}
{{- end}}
{{- end}}
{{- end -}}
