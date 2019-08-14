/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestDataSourceRemoteFileInstantiation tests whether the ResourceFile instance can be instantiated.
func TestDataSourceRemoteFileInstantiation(t *testing.T) {
	s := dataSourceRemoteFile()

	if s == nil {
		t.Fatalf("Cannot instantiate ResourceFile")
	}
}

// TestDataSourceRemoteFileSchema tests the ResourceFile schema.
func TestDataSourceRemoteFileSchema(t *testing.T) {
	s := dataSourceRemoteFile()

	requiredKeys := []string{
		mkDataSourceRemoteFileHost,
		mkDataSourceRemoteFilePath,
	}

	for _, v := range requiredKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in dataSourceRemoteFile.Schema: Missing argument \"%s\"", v)
		}

		if s.Schema[v].Required != true {
			t.Fatalf("Error in dataSourceRemoteFile.Schema: Argument \"%s\" is not required", v)
		}
	}

	attributeKeys := []string{
		mkDataSourceRemoteFileContents,
		mkDataSourceRemoteFileLastModified,
		mkDataSourceRemoteFileSize,
	}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in dataSourceRemoteFile.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in dataSourceRemoteFile.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
