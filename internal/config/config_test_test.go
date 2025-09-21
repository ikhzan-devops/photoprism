package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_InitializeTestData(t *testing.T) {
	c := NewConfig(CliTestContext())

	err := c.InitializeTestData()
	assert.NoError(t, err)
}

func TestConfig_AssertTestData(t *testing.T) {
	c := NewConfig(CliTestContext())
	// Ensure fixtures are initialized, then verify required directories.
	if err := c.InitializeTestData(); err != nil {
		t.Fatalf("InitializeTestData failed: %v", err)
	}
	c.AssertTestData(t)
}
