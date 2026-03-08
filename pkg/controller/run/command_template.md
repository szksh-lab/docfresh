{{if .Command -}}
```sh
{{trimSuffix "\n" .Command}}
```
{{- else if .EmbedScript -}}
```{{.CommandLanguage}}
{{trimSuffix "\n" .Content}}
```
{{- else -}}
```sh
{{join " " .Shell}} {{trimSuffix "\n" .Script}}
```
{{- end}}

{{if not .Quiet -}}
{{if .DetailsTagSummary -}}
<details>
<summary>{{.DetailsTagSummary}}</summary>

{{if .CodeBlock -}}
```{{.OutputLanguage}}
{{trimSuffix "\n" .CombinedOutput}}
```
{{- else -}}
{{trimSuffix "\n" .CombinedOutput}}
{{- end}}

</details>
{{- else -}}
{{if .CodeBlock -}}
Output:

```{{.OutputLanguage}}
{{trimSuffix "\n" .CombinedOutput}}
```
{{- else -}}
{{trimSuffix "\n" .CombinedOutput}}
{{- end}}
{{- end}}
{{- end -}}
