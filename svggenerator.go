package main

import "text/template"

const blank = `<?xml version="1.0"?><svg xmlns="http://www.w3.org/2000/svg" width="200" height="200"><path fill="none" stroke="#999" stroke-width="2" d="M1,1V199H199V1z"/></svg>`

var tmpl = `{{ define "Sparkline" }}<?xml version="1.0"?>
<svg xmlns="http://www.w3.org/2000/svg" width="{{ .Opts.ImgWidth }}" height="{{ .Opts.ImgHeight }}">
	<path d="{{ .Line }}" stroke="{{ .Opts.LineColor }}" stroke-width="{{ .Opts.LineWidth }}" fill="none"/>
	{{ if .Opts.IsFilled }}<path d="{{ .Line }} {{ .Closure }}" fill="{{ .Opts.FillColor }}"/>{{ end }}
</svg>
{{ end }}`

var svgTpl = template.Must(template.New("").Parse(tmpl))
