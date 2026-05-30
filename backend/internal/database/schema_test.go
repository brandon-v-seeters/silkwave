package database

import "testing"

func TestSchemaHasPerArtistReleaseSlugIndex(t *testing.T) {
	t.Parallel()

	for _, collection := range Schema {
		if collection.Name != "Releases" {
			continue
		}

		for _, index := range collection.Indices {
			if index.Unique && fieldsEqual(index.Fields, []string{"artistKey", "slug"}) {
				return
			}
		}

		t.Fatal("Releases collection is missing unique (artistKey, slug) index")
	}

	t.Fatal("Releases collection is missing from schema")
}

func fieldsEqual(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}

	return true
}
