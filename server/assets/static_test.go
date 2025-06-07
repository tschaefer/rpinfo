/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package assets

import (
	"io/fs"
	"testing"
)

func Test_StaticContentProvidesAllAssets(t *testing.T) {
	expectedFiles := [4]string{
		"index.html",
		"normalize.min.css",
		"openapi.yml",
		"redoc.standalone.js",
	}

	actualFiles := make(map[string]bool)
	err := fs.WalkDir(StaticContent, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			t.Fatalf("error walking through static content: %v", err)
		}
		actualFiles[path] = true
		return nil
	})
	if err != nil {
		t.Fatalf("error accessing static content: %v", err)
	}

	for _, expectedFile := range expectedFiles {
		if _, exists := actualFiles[expectedFile]; !exists {
			t.Errorf("expected file %s not found in static content", expectedFile)
		}
	}
}
