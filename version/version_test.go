/*
Copyright (c) 2025 Tobias Schäfer. All rights reserved.
Licensed under the MIT license, see LICENSE in the project root for details.
*/
package version

import (
	"testing"
)

func Test_EmptyVersionMeansReleaseReturnsDev(t *testing.T) {
	Version = ""
	output := Release()
	expected := "dev"
	if output != expected {
		t.Errorf("Version is not from Release - want: %s, got: %s\n", expected, output)
	}
}

func Test_VersionReturnedFromRelease(t *testing.T) {
	Version = "testing-manual"
	output := Release()
	expected := Version
	if output != expected {
		t.Errorf("Version is not from Release - want: %s, got: %s\n", expected, output)
	}
}

func Test_EmptyGitCommitMeansCommitReturnsEmpty(t *testing.T) {
	GitCommit = ""
	output := Commit()
	expected := ""
	if output != expected {
		t.Errorf("GitCommit is not from Commit - want: %s, got: %s\n", expected, output)
	}
}

func Test_GitCommitReturnedFromCommit(t *testing.T) {
	GitCommit = "testing-manual"
	output := Commit()
	expected := GitCommit
	if output != expected {
		t.Errorf("GitCommit is not from Commit - want: %s, got: %s\n", expected, output)
	}
}
