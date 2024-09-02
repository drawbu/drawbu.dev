package static

import "embed"

// content holds our static articles
//go:embed robots.txt generated.css iosevka-comfy-fixed-bold.woff2 iosevka-comfy-fixed-regular.woff2
var Files embed.FS
