package schema_store

import (
	"database/sql"
	"sync"

	"github.com/moiot/gravity/pkg/config"

	"github.com/juju/errors"

	"github.com/moiot/gravity/pkg/utils"
)

type SimpleSchemaStore struct {
	sync.RWMutex
	sourceDB *sql.DB
	dbCfg    *config.DBConfig
	schemas  map[string]Schema
}

func (store *SimpleSchemaStore) IsInCache(dbName string) bool {
	if _, ok := store.schemas[dbName]; ok {
		return true
	} else {
		return false
	}
}

func (store *SimpleSchemaStore) getFromCache(dbName string, lock bool) (Schema, bool) {
	if lock {
		store.RLock()
		defer store.RUnlock()
	}

	cachedSchema, ok := store.schemas[dbName]
	if ok {
		return cachedSchema, true
	}
	return nil, false
}

func (store *SimpleSchemaStore) GetSchema(dbName string) (Schema, error) {
	if dbName == "" {
		return nil, nil
	}

	schema, ok := store.getFromCache(dbName, true)
	if ok {
		return schema, nil
	}

	store.Lock()
	defer store.Unlock()
	schema, ok = store.getFromCache(dbName, false)
	if ok {
		return schema, nil
	}
	schema, err := GetSchemaFromDB(store.sourceDB, dbName)
	if err != nil {
		return nil, errors.Trace(err)
	}
	store.schemas[dbName] = schema
	return schema, nil
}

func (store *SimpleSchemaStore) InvalidateSchemaCache(schema string) {
	store.Lock()
	defer store.Unlock()

	delete(store.schemas, schema)
}

func (store *SimpleSchemaStore) InvalidateCache() {
	// Invalidate Schema cache
	store.Lock()
	defer store.Unlock()

	// make a new map here
	store.schemas = make(map[string]Schema)
}

func (store *SimpleSchemaStore) Close() {
	if store.sourceDB != nil {
		store.sourceDB.Close()
	}
}

func NewSimpleSchemaStoreFromDBConn(db *sql.DB) (SchemaStore, error) {
	return &SimpleSchemaStore{sourceDB: db, schemas: make(map[string]Schema)}, nil
}

func NewSimpleSchemaStore(dbCfg *config.DBConfig) (*SimpleSchemaStore, error) {
	sourceDB, err := utils.CreateDBConnection(dbCfg)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &SimpleSchemaStore{dbCfg: dbCfg, schemas: make(map[string]Schema), sourceDB: sourceDB}, nil
}
