package web

import "embed"

// The all: prefix embeds the placeholder, so the directory exists before tailwind generates its output.
//
//go:embed all:assets
var Files embed.FS
