package static

import "embed"

// content holds our static articles
//go:embed robots.txt generated.css
var Files embed.FS
