package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sangrita-tech/golang-pkg/pkg/config"
	"github.com/stretchr/testify/require"
)

type testCfg struct {
	Port int    `yaml:"port" env:"PORT" env-required:"true"`
	Name string `yaml:"name" env:"NAME"`
}

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	p := filepath.Join(dir, name)
	err := os.WriteFile(p, []byte(content), 0o600)
	require.NoError(t, err)

	return p
}

func Test_Load_EnvOnly_ReturnsCfg(t *testing.T) {
	t.Setenv("PORT", "1234")
	t.Setenv("NAME", "from-env")

	cfg, err := config.Load[testCfg]("")

	require.NoError(t, err)
	require.Equal(t, 1234, cfg.Port)
	require.Equal(t, "from-env", cfg.Name)
}

func Test_Load_EnvOnlyMissingRequired_ReturnsError(t *testing.T) {
	os.Unsetenv("PORT")
	t.Setenv("NAME", "x")

	_, err := config.Load[testCfg]("")

	require.Error(t, err)
}

func Test_Load_MissingFileFallsBackToEnv_ReturnsCfg(t *testing.T) {
	t.Setenv("PORT", "1111")
	t.Setenv("NAME", "fallback-env")
	missing := filepath.Join(t.TempDir(), "nope.yaml")

	cfg, err := config.Load[testCfg](missing)

	require.NoError(t, err)
	require.Equal(t, 1111, cfg.Port)
	require.Equal(t, "fallback-env", cfg.Name)
}

func Test_Load_MissingFileFallsBackToEnvMissingRequired_ReturnsError(t *testing.T) {
	os.Unsetenv("PORT")
	t.Setenv("NAME", "fallback-env")
	missing := filepath.Join(t.TempDir(), "nope.yaml")

	_, err := config.Load[testCfg](missing)

	require.Error(t, err)
}

func Test_Load_FileProvidedReadsFile_ReturnsCfg(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("NAME")
	dir := t.TempDir()
	path := writeTempFile(t, dir, "cfg.yaml", "port: 2222\nname: from-file\n")

	cfg, err := config.Load[testCfg](path)

	require.NoError(t, err)
	require.Equal(t, 2222, cfg.Port)
	require.Equal(t, "from-file", cfg.Name)
}

func Test_Load_FileProvidedEnvOverrides_ReturnsCfgWithOverrides(t *testing.T) {
	t.Setenv("PORT", "3333")
	t.Setenv("NAME", "from-env")
	dir := t.TempDir()
	path := writeTempFile(t, dir, "cfg.yaml", "port: 2222\nname: from-file\n")

	cfg, err := config.Load[testCfg](path)

	require.NoError(t, err)
	require.Equal(t, 3333, cfg.Port)
	require.Equal(t, "from-env", cfg.Name)
}

func Test_Load_ParseError_ReturnsErrorAndZero(t *testing.T) {
	t.Setenv("PORT", "9999")
	t.Setenv("NAME", "env")
	dir := t.TempDir()
	path := writeTempFile(t, dir, "cfg.yaml", "port: [\n")

	cfg, err := config.Load[testCfg](path)

	require.Error(t, err)
	require.Equal(t, testCfg{}, cfg)
}
