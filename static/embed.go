package static

import "embed"

// content holds our static articles
//go:embed robots.txt generated.css *.woff2
var Files embed.FS
