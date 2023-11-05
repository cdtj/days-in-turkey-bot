package assets

import "embed"

/*
Since we cannot embed ../. dir, I'm using this hack to embed my assets
*/

//go:embed i18n/*.toml
var I18n embed.FS

//go:embed country/*.toml
var Countries embed.FS
