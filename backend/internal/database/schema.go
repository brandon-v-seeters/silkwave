package database

// CollectionType represents the type of ArangoDB collection
type CollectionType string

const (
	DocumentCollection CollectionType = "Document"
	EdgeCollection     CollectionType = "Edge"
)

// IndexType represents the type of ArangoDB index
type IndexType string

const (
	PersistentIndex IndexType = "persistent"
	HashIndex       IndexType = "hash"
	SkiplistIndex   IndexType = "skiplist"
	FulltextIndex   IndexType = "fulltext"
	GeoIndex        IndexType = "geo"
	TTLIndex        IndexType = "ttl"
)

// IndexDefinition defines an index on a collection
type IndexDefinition struct {
	Fields []string  `json:"fields"`
	Unique bool      `json:"unique"`
	Sparse bool      `json:"sparse"` // Sparse indexes skip documents where indexed fields are null/missing
	Type   IndexType `json:"type"`
}

// CollectionDefinition defines a collection and its indices
type CollectionDefinition struct {
	Name    string            `json:"name"`
	Type    CollectionType    `json:"type"`
	Indices []IndexDefinition `json:"indices"`
}

// Schema defines the database schema
var Schema = []CollectionDefinition{
	{
		Name: "Users",
		Type: DocumentCollection,
		Indices: []IndexDefinition{
			{Fields: []string{"username"}, Unique: true, Sparse: true, Type: PersistentIndex},
			{Fields: []string{"email"}, Unique: true, Type: PersistentIndex},
		},
	},
	{
		Name: "Artists",
		Type: DocumentCollection,
		Indices: []IndexDefinition{
			{Fields: []string{"slug"}, Unique: true, Type: PersistentIndex},
			{Fields: []string{"name"}, Unique: true, Type: PersistentIndex},
		},
	},
	{
		Name: "UsersArtists",
		Type: EdgeCollection,
		Indices: []IndexDefinition{
			{Fields: []string{"_from", "_to"}, Unique: true, Type: PersistentIndex},
		},
	},
	{
		Name: "Follows",
		Type: EdgeCollection,
		Indices: []IndexDefinition{
			{Fields: []string{"_from", "_to"}, Unique: true, Type: PersistentIndex},
			{Fields: []string{"_to"}, Unique: false, Type: PersistentIndex},
		},
	},
	{
		Name: "Subscriptions",
		Type: DocumentCollection,
		Indices: []IndexDefinition{
			{Fields: []string{"artistKey"}, Unique: false, Type: PersistentIndex},
			{Fields: []string{"artistKey", "currency", "price"}, Unique: true, Type: PersistentIndex},
		},
	},
	{
		Name: "Subscribers",
		Type: DocumentCollection,
		Indices: []IndexDefinition{
			{Fields: []string{"artistKey"}, Unique: false, Type: PersistentIndex},
			{Fields: []string{"subscriberKey"}, Unique: false, Type: PersistentIndex},
			{Fields: []string{"subscriptionKey"}, Unique: false, Type: PersistentIndex},
			{Fields: []string{"subscriberKey", "subscriptionKey"}, Unique: true, Type: PersistentIndex},
		},
	},
	{
		Name: "Releases",
		Type: DocumentCollection,
		Indices: []IndexDefinition{
			{Fields: []string{"id"}, Unique: true, Type: PersistentIndex},
			{Fields: []string{"slug"}, Unique: false, Type: PersistentIndex},
			{Fields: []string{"artistKey", "slug"}, Unique: true, Type: PersistentIndex},
			{Fields: []string{"artistKey", "status"}, Unique: false, Type: PersistentIndex},
			{Fields: []string{"artistKey", "publishAt"}, Unique: false, Type: PersistentIndex},
			{Fields: []string{"status", "publishAt"}, Unique: false, Type: PersistentIndex},
		},
	},
	{
		Name: "Tracks",
		Type: DocumentCollection,
		Indices: []IndexDefinition{
			{Fields: []string{"releaseKey", "id"}, Unique: false, Type: PersistentIndex},
		},
	},
}
