/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package assets

import (
	"embed"
	"io/fs"
)

//go:embed static/*
var staticFiles embed.FS

var StaticContent fs.FS

func init() {
	var err error

	if StaticContent, err = fs.Sub(staticFiles, "static"); err != nil {
		panic(err)
	}
}
