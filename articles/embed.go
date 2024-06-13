package articles

import "embed"

// content holds our static articles
//go:embed *.md
var Articles embed.FS
