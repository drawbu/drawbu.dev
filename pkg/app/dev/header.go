// +build !production

package dev

import (
	_ "embed"

	"github.com/a-h/templ"
)

//go:embed autorefresh.js
var s []byte

func init() {
	Header = templ.Raw("<script>" + string(s) + "</script>")
}
