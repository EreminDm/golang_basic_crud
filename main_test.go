package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	conn, db, dbtype := envf()
	assert.Equal(t, "localhost:27017", conn, "Expected is not equals to actual, got %s", conn)
	assert.Equal(t, "person", db, "Expected is not equals to actual, got %s", db)
	assert.Equal(t, "mongo", dbtype, "Expected is not equals to actual, got %s", dbtype)
}
