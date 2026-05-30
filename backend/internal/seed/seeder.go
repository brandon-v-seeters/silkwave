package seed

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/brandon-v-seeters/go-silk-wave/internal/auth"
	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
)

type Options struct {
	Namespace string
	Password  string
	ResetSeed bool
}

type Seeder struct {
	db        *database.ArangoDB
	passwords *auth.PasswordService
}

type Summary struct {
	Namespace string
	Password  string
	ResetSeed bool
	Deleted   map[string]int
	Upserted  map[string]int
}

func NewSeeder(db *database.ArangoDB, passwords *auth.PasswordService) *Seeder {
	return &Seeder{db: db, passwords: passwords}
}

func (s *Seeder) Run(ctx context.Context, opts Options) (Summary, error) {
	if opts.Namespace == "" {
		opts.Namespace = DefaultNamespace
	}
	if opts.Password == "" {
		opts.Password = DefaultPassword
	}

	if err := s.db.Migrate(ctx); err != nil {
		return Summary{}, fmt.Errorf("run migrations: %w", err)
	}

	summary := Summary{
		Namespace: opts.Namespace,
		Password:  opts.Password,
		ResetSeed: opts.ResetSeed,
		Deleted:   map[string]int{},
		Upserted:  map[string]int{},
	}

	if opts.ResetSeed {
		deleted, err := s.resetSeed(ctx, opts.Namespace)
		if err != nil {
			return Summary{}, err
		}
		summary.Deleted = deleted
	}

	passwordHash, err := s.passwords.HashPassword(opts.Password)
	if err != nil {
		return Summary{}, fmt.Errorf("hash seed password: %w", err)
	}

	catalog := buildCatalog(opts.Namespace, passwordHash, time.Now().UTC())
	if err := s.upsertCatalog(ctx, catalog, summary.Upserted); err != nil {
		return Summary{}, err
	}

	return summary, nil
}

func (s *Seeder) upsertCatalog(ctx context.Context, catalog catalog, counts map[string]int) error {
	groups := []struct {
		collection string
		docs       []document
	}{
		{"Users", catalog.Users},
		{"Artists", catalog.Artists},
		{"UsersArtists", catalog.UserArtists},
		{"Releases", catalog.Releases},
		{"Tracks", catalog.Tracks},
		{"Subscriptions", catalog.Subscriptions},
		{"Subscribers", catalog.Subscribers},
	}

	for _, group := range groups {
		for _, doc := range group.docs {
			if err := s.upsert(ctx, group.collection, doc); err != nil {
				return fmt.Errorf("upsert %s/%v: %w", group.collection, doc["_key"], err)
			}
			counts[group.collection]++
		}
	}

	return nil
}

func (s *Seeder) upsert(ctx context.Context, collection string, doc document) error {
	key, ok := doc["_key"].(string)
	if !ok || key == "" {
		return fmt.Errorf("seed document for %s is missing _key", collection)
	}
	namespace := namespaceFor(doc)
	selector, err := selectorFor(collection, doc)
	if err != nil {
		return err
	}

	query := /*aql*/ `
		LET existing = FIRST(
			FOR doc IN @@collection
				FILTER MATCHES(doc, @selector)
				LIMIT 1
				RETURN doc
		)
		FILTER existing == null OR existing.seed.namespace == @namespace
		UPSERT @selector
			INSERT @doc
			UPDATE MERGE(OLD, UNSET(@doc, "_key", "_id", "_rev"))
		IN @@collection
		RETURN NEW._key
	`

	if _, err := database.QueryOne[string](ctx, s.db, query, map[string]interface{}{
		"@collection": collection,
		"doc":         doc,
		"namespace":   namespace,
		"selector":    selector,
	}); err != nil {
		return err
	}

	return nil
}

func namespaceFor(doc document) string {
	meta, ok := doc["seed"].(seedMeta)
	if !ok {
		return DefaultNamespace
	}
	return meta.Namespace
}

func selectorFor(collection string, doc document) (document, error) {
	switch collection {
	case "Users":
		return selector(doc, "email")
	case "Artists":
		return selector(doc, "slug")
	case "UsersArtists":
		return selector(doc, "_from", "_to")
	case "Releases":
		return selector(doc, "artistKey", "slug")
	case "Tracks":
		return selector(doc, "releaseKey", "id")
	case "Subscriptions":
		return selector(doc, "artistKey", "currency", "price")
	case "Subscribers":
		return selector(doc, "subscriberKey", "subscriptionKey")
	default:
		return selector(doc, "_key")
	}
}

func selector(doc document, fields ...string) (document, error) {
	out := document{}
	for _, field := range fields {
		value, ok := doc[field]
		if !ok {
			return nil, fmt.Errorf("seed document is missing selector field %q", field)
		}
		out[field] = value
	}
	return out, nil
}

func (s *Seeder) resetSeed(ctx context.Context, namespace string) (map[string]int, error) {
	counts := map[string]int{}
	collections := []string{"Follows", "Subscribers", "Subscriptions", "Tracks", "Releases", "UsersArtists", "Artists", "Users"}

	for _, collection := range collections {
		q := /*aql*/ `
			FOR doc IN @@collection
				FILTER doc.seed.namespace == @namespace
				REMOVE doc IN @@collection
				RETURN OLD._key
		`
		keys, err := database.Query[string](ctx, s.db, q, map[string]interface{}{
			"@collection": collection,
			"namespace":   namespace,
		})
		if err != nil {
			return nil, fmt.Errorf("reset seed collection %s: %w", collection, err)
		}
		counts[collection] = len(keys)
	}

	return counts, nil
}

func (s Summary) String() string {
	var builder strings.Builder

	builder.WriteString("Seeded Silkwave dev universe.\n\n")
	builder.WriteString(fmt.Sprintf("Namespace: %s\n", s.Namespace))
	if s.ResetSeed {
		builder.WriteString("Reset: seed-owned records were deleted before upsert.\n")
	}
	builder.WriteString("\nLogins:\n")
	builder.WriteString(fmt.Sprintf("artist@silkwave.test / %s\n", s.Password))
	builder.WriteString(fmt.Sprintf("fan@silkwave.test / %s\n", s.Password))
	builder.WriteString(fmt.Sprintf("subscriber@silkwave.test / %s\n", s.Password))

	builder.WriteString("\nPrimary URLs:\n")
	builder.WriteString("Artist: /artist/framer\n")
	builder.WriteString("Release: /artist/framer/releases/moon-in-my-sky\n")
	builder.WriteString("Drafts: /upload/drafts\n")

	if len(s.Deleted) > 0 {
		builder.WriteString("\nDeleted seed records:\n")
		appendCounts(&builder, s.Deleted)
	}

	builder.WriteString("\nUpserted seed records:\n")
	appendCounts(&builder, s.Upserted)

	builder.WriteString("\nNotes:\n")
	builder.WriteString("- Media paths are fake v1 storage paths. No files were uploaded.\n")
	builder.WriteString("- Seeded Follows are skipped until notification demos need them.\n")
	builder.WriteString("- Purchases are skipped because there is no commerce contract yet.\n")

	return builder.String()
}

func appendCounts(builder *strings.Builder, counts map[string]int) {
	order := []string{"Users", "Artists", "UsersArtists", "Follows", "Releases", "Tracks", "Subscriptions", "Subscribers"}
	for _, collection := range order {
		if count, ok := counts[collection]; ok {
			builder.WriteString(fmt.Sprintf("- %s: %d\n", collection, count))
		}
	}
}
