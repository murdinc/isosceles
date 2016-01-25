package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Canary(t *testing.T) {
	assert.True(t, true, "Canary test passing")
}
