package config_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/asmazovec/team-agile/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestMustRead_NilOption_ShouldNotPanic(t *testing.T) {
	assert.NotPanics(t,
		func() { config.MustRead(nil) },
	)
}

type OriginMock struct {
	Called bool
}

func (om OriginMock) ShouldBeCalled() config.Origin {
	return func(_ *config.AppConfig) error {
		om.Called = true
		return nil
	}
}

func (om OriginMock) WithError(err error) config.Origin {
	return func(_ *config.AppConfig) error {
		return err
	}
}

func (om OriginMock) WithAddress(addr string) config.Origin {
	return func(cfg *config.AppConfig) error {
		cfg.HTTPPrimaryServer.Address = addr
		return nil
	}
}

func TestMustRead_WithAddress_ShouldSetup(t *testing.T) {
	address := "address"

	cfg := config.MustRead(OriginMock{}.WithAddress(address))

	assert.Equal(t, address, cfg.HTTPPrimaryServer.Address)
}

func TestMustRead_WithManyOrigins_ShouldUseValueFromLast(t *testing.T) {
	address := "address"
	m := OriginMock{}

	cfg := config.MustRead(m.WithAddress("another "+address), m.WithAddress(address))

	assert.Equal(t, address, cfg.HTTPPrimaryServer.Address)
}

func TestMustRead_OriginWithErr_ShouldPanic(t *testing.T) {
	m := OriginMock{}
	assert.Panics(t,
		func() {
			config.MustRead(m.WithError(errors.New("error")))
		})
}

func TestFromConfig_EmptyPath_ShouldNotError(t *testing.T) {
	f := config.FromEnv("")

	err := f(nil)

	assert.NoError(t, err)
}

func TestFromConfig_EnvVars_ShouldBeSet(t *testing.T) {
	val := "address"
	t.Setenv("HTTP_ADDRESS", val)

	cfg := config.MustRead(config.FromEnv(""))

	assert.Equal(t, val, cfg.HTTPPrimaryServer.Address)
}

func createEnvConfig(file, key, val string) {
	f, _ := os.Create(file)
	_, _ = fmt.Fprintf(f, "%s=%s\n", key, val)
	_ = f.Close()
}

func TestFromConfig_ConfigVars_ShouldBeSet(t *testing.T) {
	file := filepath.Join(t.TempDir(), ".env")
	val := "address"
	createEnvConfig(file, "HTTP_ADDRESS", val)

	cfg := config.MustRead(config.FromEnv(file))

	assert.Equal(t, val, cfg.HTTPPrimaryServer.Address)
}

func TestFromConfig_EnvConfig_ShouldSetEnvVariables(t *testing.T) {
	file := filepath.Join(t.TempDir(), ".env")
	val := "address"
	createEnvConfig(file, "HTTP_ADDRESS", val)

	config.MustRead(config.FromEnv(file))

	assert.Equal(t, val, os.Getenv("HTTP_ADDRESS"))
}
