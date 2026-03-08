{{if .Command -}}
```sh
{{trimSuffix "\n" .Command}}
```
{{- else if .EmbedScript -}}
```{{.ScriptLanguage}}
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

{{if .UseFencedCodeBlockForOutput -}}
```
{{trimSuffix "\n" .CombinedOutput}}
```
{{- else -}}
{{trimSuffix "\n" .CombinedOutput}}
{{- end}}

</details>
{{- else -}}
{{if .UseFencedCodeBlockForOutput -}}
Output:

```
{{trimSuffix "\n" .CombinedOutput}}
```
{{- else -}}
{{trimSuffix "\n" .CombinedOutput}}
{{- end}}
{{- end}}
{{- end -}}
