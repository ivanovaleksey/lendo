package test

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	_, b, _, _ = runtime.Caller(0)
	Root       = filepath.Join(filepath.Dir(b),  "../..")
)

func LoadAPIEnv(t *testing.T) {
	const filename = ".env.api"
	require.NoError(t, godotenv.Load(filepath.Join(Root, filename)))
}

func LoadRegistryEnv(t *testing.T) {
	const filename = ".env.registry"
	require.NoError(t, godotenv.Load(filepath.Join(Root, filename)))
}
