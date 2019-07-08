package ksmglog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecord_Hash(t *testing.T) {
	record := Record{
		ID:          111,
		Description: "description",
	}

	assert.NoError(t, record.Hash())
	assert.Equal(t, "8d0cdef44129b9ad51cf04d2c9142be7", record.HashString)
}
