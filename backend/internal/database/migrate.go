package database

import (
	"context"
	"fmt"
	"log"

	"github.com/arangodb/go-driver"
)

// Migrate ensures all collections and indices defined in the schema exist.
// This should be called on application startup.
func (a *ArangoDB) Migrate(ctx context.Context) error {
	log.Println("Starting database migration...")

	for _, def := range Schema {
		if err := a.ensureCollection(ctx, def); err != nil {
			return fmt.Errorf("failed to ensure collection %s: %w", def.Name, err)
		}
	}

	log.Println("Database migration completed successfully")
	return nil
}

// MigrateWithSchema allows passing a custom schema instead of using the default
func (a *ArangoDB) MigrateWithSchema(ctx context.Context, schema []CollectionDefinition) error {
	log.Println("Starting database migration with custom schema...")

	for _, def := range schema {
		if err := a.ensureCollection(ctx, def); err != nil {
			return fmt.Errorf("failed to ensure collection %s: %w", def.Name, err)
		}
	}

	log.Println("Database migration completed successfully")
	return nil
}

func (a *ArangoDB) ensureCollection(ctx context.Context, def CollectionDefinition) error {
	exists, err := a.Database.CollectionExists(ctx, def.Name)
	if err != nil {
		return fmt.Errorf("failed to check collection existence: %w", err)
	}

	var col driver.Collection
	if exists {
		col, err = a.Database.Collection(ctx, def.Name)
		if err != nil {
			return fmt.Errorf("failed to get collection: %w", err)
		}
		log.Printf("Collection %s already exists", def.Name)
	} else {
		options := &driver.CreateCollectionOptions{}
		if def.Type == EdgeCollection {
			options.Type = driver.CollectionTypeEdge
		} else {
			options.Type = driver.CollectionTypeDocument
		}

		col, err = a.Database.CreateCollection(ctx, def.Name, options)
		if err != nil {
			return fmt.Errorf("failed to create collection: %w", err)
		}
		log.Printf("Created collection %s (%s)", def.Name, def.Type)
	}

	// Ensure indices
	for _, idx := range def.Indices {
		if err := a.ensureIndex(ctx, col, idx); err != nil {
			return fmt.Errorf("failed to ensure index on %v: %w", idx.Fields, err)
		}
	}

	return nil
}

func (a *ArangoDB) ensureIndex(ctx context.Context, col driver.Collection, idx IndexDefinition) error {
	// Check if an index on these fields already exists with different settings
	existingIndexes, err := col.Indexes(ctx)
	if err != nil {
		return fmt.Errorf("failed to list indexes: %w", err)
	}

	for _, existing := range existingIndexes {
		// Skip system indexes (primary, edge) - they can't be modified
		indexType := existing.Type()
		if indexType == driver.PrimaryIndex || indexType == driver.EdgeIndex {
			continue
		}

		if fieldsMatch(existing.Fields(), idx.Fields) {
			log.Printf("Index on %s: %v already exists. (unique: %v, sparse: %v)", col.Name(), idx.Fields, existing.Unique(), existing.Sparse())
			// Check if settings differ
			if existing.Unique() != idx.Unique || existing.Sparse() != idx.Sparse {
				log.Printf("Index on %s: %v has different settings (unique: %v->%v, sparse: %v->%v), recreating...",
					col.Name(), idx.Fields, existing.Unique(), idx.Unique, existing.Sparse(), idx.Sparse)
				
				// Delete the old index
				if err := existing.Remove(ctx); err != nil {
					return fmt.Errorf("failed to remove old index: %w", err)
				}
				log.Printf("Removed old index on %s: %v", col.Name(), idx.Fields)
			} else {
				// Index exists with same settings, nothing to do
				log.Printf("Index on %s: %v already exists with correct settings", col.Name(), idx.Fields)
				return nil
			}
		}
	}

	// Create the index
	switch idx.Type {
	case PersistentIndex:
		_, _, err = col.EnsurePersistentIndex(ctx, idx.Fields, &driver.EnsurePersistentIndexOptions{
			Unique: idx.Unique,
			Sparse: idx.Sparse,
		})
	case HashIndex:
		_, _, err = col.EnsureHashIndex(ctx, idx.Fields, &driver.EnsureHashIndexOptions{
			Unique: idx.Unique,
			Sparse: idx.Sparse,
		})
	case SkiplistIndex:
		_, _, err = col.EnsureSkipListIndex(ctx, idx.Fields, &driver.EnsureSkipListIndexOptions{
			Unique: idx.Unique,
			Sparse: idx.Sparse,
		})
	case FulltextIndex:
		if len(idx.Fields) > 0 {
			_, _, err = col.EnsureFullTextIndex(ctx, idx.Fields, nil)
		}
	case GeoIndex:
		_, _, err = col.EnsureGeoIndex(ctx, idx.Fields, nil)
	case TTLIndex:
		if len(idx.Fields) > 0 {
			_, _, err = col.EnsureTTLIndex(ctx, idx.Fields[0], 0, nil)
		}
	default:
		return fmt.Errorf("unsupported index type: %s", idx.Type)
	}

	if err != nil {
		return err
	}

	log.Printf("Created index on %s: %v (unique: %v, sparse: %v)", col.Name(), idx.Fields, idx.Unique, idx.Sparse)
	return nil
}

// fieldsMatch checks if two field slices contain the same fields (order matters for indexes)
func fieldsMatch(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
