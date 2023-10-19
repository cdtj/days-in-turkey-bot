package assets

import "embed"

/*
Since we cannot embed ../. dir, I'm using this hack to embed my assets
*/

//go:embed l10n/*.toml
var L10n embed.FS
