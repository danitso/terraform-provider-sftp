package main

import (
	"testing"
)

// TestResourceFileInstantiation tests whether the ResourceFile instance can be instantiated.
func TestResourceFileInstantiation(t *testing.T) {
	s := resourceFile()

	if s == nil {
		t.Fatalf("Cannot instantiate ResourceFile")
	}
}

// TestResourceFileSchema tests the ResourceFile schema.
func TestResourceFileSchema(t *testing.T) {
	s := resourceFile()

	requiredKeys := []string{
		dataSourceDiskLabelKey,
		dataSourceDiskServerIDKey,
		dataSourceDiskSizeKey,
	}

	for _, v := range requiredKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in ResourceFile.Schema: Missing argument \"%s\"", v)
		}

		if s.Schema[v].Required != true {
			t.Fatalf("Error in ResourceFile.Schema: Argument \"%s\" is not required", v)
		}
	}

	attributeKeys := []string{
		dataSourceDiskPrimaryKey,
	}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in dataSourceDisk.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in dataSourceDisk.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
