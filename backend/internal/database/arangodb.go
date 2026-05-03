package database

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/brandon-v-seeters/go-silk-wave/internal/config"
)

type ArangoDB struct {
	Client   driver.Client
	Database driver.Database
}

func NewArangoDB(cfg *config.Config) (*ArangoDB, error) {
	ctx := context.Background()

	// Create HTTP connection
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{cfg.ArangoEndpoint},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	// Create client with authentication
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(cfg.ArangoUsername, cfg.ArangoPassword),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	// Check if database exists, create if not
	exists, err := client.DatabaseExists(ctx, cfg.ArangoDatabase)
	if err != nil {
		return nil, fmt.Errorf("failed to check database existence: %w", err)
	}

	var db driver.Database
	if !exists {
		db, err = client.CreateDatabase(ctx, cfg.ArangoDatabase, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
	} else {
		db, err = client.Database(ctx, cfg.ArangoDatabase)
		if err != nil {
			return nil, fmt.Errorf("failed to get database: %w", err)
		}
	}

	return &ArangoDB{
		Client:   client,
		Database: db,
	}, nil
}

func (a *ArangoDB) Close() error {
	// ArangoDB driver doesn't require explicit close
	return nil
}

func (a *ArangoDB) GetCollection(ctx context.Context, name string) (driver.Collection, error) {
	exists, err := a.Database.CollectionExists(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check collection existence: %w", err)
	}

	if !exists {
		col, err := a.Database.CreateCollection(ctx, name, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create collection: %w", err)
		}
		return col, nil
	}

	col, err := a.Database.Collection(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}
	return col, nil
}

// Query executes an AQL query and returns results as a slice of the generic type T.
// Pass nil for bindVars if no variables are needed.
func Query[T any](ctx context.Context, db *ArangoDB, query string, bindVars map[string]interface{}) ([]T, error) {
	cursor, err := db.Database.Query(ctx, query, bindVars)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer cursor.Close()

	var results []T
	for cursor.HasMore() {
		var item T
		_, err := cursor.ReadDocument(ctx, &item)
		if err != nil {
			return nil, fmt.Errorf("failed to read document: %w", err)
		}
		results = append(results, item)
	}

	return results, nil
}

// QueryOne executes an AQL query and returns a single result of the generic type T.
// Returns an error if no results are found.
func QueryOne[T any](ctx context.Context, db *ArangoDB, query string, bindVars map[string]interface{}) (*T, error) {
	cursor, err := db.Database.Query(ctx, query, bindVars)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer cursor.Close()

	if !cursor.HasMore() {
		return nil, fmt.Errorf("no results found")
	}

	var item T
	_, err = cursor.ReadDocument(ctx, &item)
	if err != nil {
		return nil, fmt.Errorf("failed to read document: %w", err)
	}

	return &item, nil
}

// Transaction wraps an ArangoDB transaction for executing multiple queries atomically.
type Transaction struct {
	tx       driver.TransactionID
	db       *ArangoDB
	database driver.Database
}

// BeginTransaction starts a new transaction with the specified collections.
// readCollections: collections that will be read from
// writeCollections: collections that will be written to
func (a *ArangoDB) BeginTransaction(ctx context.Context, readCollections, writeCollections []string) (*Transaction, error) {
	txID, err := a.Database.BeginTransaction(ctx, driver.TransactionCollections{
		Read:  readCollections,
		Write: writeCollections,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &Transaction{
		tx:       txID,
		db:       a,
		database: a.Database,
	}, nil
}

// Commit commits the transaction.
func (t *Transaction) Commit(ctx context.Context) error {
	if err := t.database.CommitTransaction(ctx, t.tx, nil); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// Abort aborts the transaction.
func (t *Transaction) Abort(ctx context.Context) error {
	if err := t.database.AbortTransaction(ctx, t.tx, nil); err != nil {
		return fmt.Errorf("failed to abort transaction: %w", err)
	}
	return nil
}

// TransactionContext returns a context with the transaction ID attached.
func (t *Transaction) TransactionContext(ctx context.Context) context.Context {
	return driver.WithTransactionID(ctx, t.tx)
}

// Query executes an AQL query within the transaction and returns results as a slice of T.
func TxQuery[T any](ctx context.Context, tx *Transaction, query string, bindVars map[string]interface{}) ([]T, error) {
	txCtx := tx.TransactionContext(ctx)
	cursor, err := tx.database.Query(txCtx, query, bindVars)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer cursor.Close()

	var results []T
	for cursor.HasMore() {
		var item T
		_, err := cursor.ReadDocument(txCtx, &item)
		if err != nil {
			return nil, fmt.Errorf("failed to read document: %w", err)
		}
		results = append(results, item)
	}

	return results, nil
}

// TxQueryOne executes an AQL query within the transaction and returns a single result.
func TxQueryOne[T any](ctx context.Context, tx *Transaction, query string, bindVars map[string]interface{}) (*T, error) {
	txCtx := tx.TransactionContext(ctx)
	cursor, err := tx.database.Query(txCtx, query, bindVars)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer cursor.Close()

	if !cursor.HasMore() {
		return nil, fmt.Errorf("no results found")
	}

	var item T
	_, err = cursor.ReadDocument(txCtx, &item)
	if err != nil {
		return nil, fmt.Errorf("failed to read document: %w", err)
	}

	return &item, nil
}

// TxExec executes a write query within the transaction (INSERT, UPDATE, REMOVE).
func TxExec(ctx context.Context, tx *Transaction, query string, bindVars map[string]interface{}) error {
	txCtx := tx.TransactionContext(ctx)
	cursor, err := tx.database.Query(txCtx, query, bindVars)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	defer cursor.Close()
	return nil
}

// WithTransaction executes a function within a transaction, automatically committing on success
// or aborting on error.
func WithTransaction(ctx context.Context, db *ArangoDB, readCollections, writeCollections []string, fn func(tx *Transaction) error) error {
	tx, err := db.BeginTransaction(ctx, readCollections, writeCollections)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		// Attempt to abort, but return the original error
		_ = tx.Abort(ctx)
		return err
	}

	return tx.Commit(ctx)
}
