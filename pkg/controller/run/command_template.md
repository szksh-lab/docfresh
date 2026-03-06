{{if .Command -}}
```sh
{{trimSuffix "\n" .Command}}
```
{{- else -}}
```
{{join " " .Shell}} {{trimSuffix "\n" .Script}}
```
{{- end}}

```
{{trimSuffix "\n" .CombinedOutput}}
```
