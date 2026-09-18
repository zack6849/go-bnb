package configuration

import (
	"crypto/rand"
	"gobnb/internal/database"
	"math/big"
	"testing"
)

func TestGetStringReadsFromEnv(t *testing.T) {
	envBackedKeys := []string{
		database.DriverKey,
		database.HostnameKey,
		database.SchemaKey,
		database.SslKey,
	}
	for _, value := range envBackedKeys {
		t.Logf("Confirming key %s reads properly from environment variables", value)
		assertStringBackedEnvValue(t, value)
	}
}

func TestGetIntReadsFromEnv(t *testing.T) {
	envBackedKeys := []string{
		database.PortKey,
	}
	for _, value := range envBackedKeys {
		t.Logf("Confirming key %s reads properly from environment variables", value)
		assertIntBackedEnvValue(t, value)
	}
}

func assertIntBackedEnvValue(t *testing.T, key string) {
	num, _ := rand.Int(rand.Reader, big.NewInt(4096))
	expected := num.String()

	t.Setenv(key, expected)
	val := GetString(key, "DUMMY_VALUE")
	if val != expected {
		if val == "DUMMY_VALUE" {
			t.Errorf("Returned default value when value was explicitly set")
		}
		t.Errorf("Returned unexpected value %s when expecting %s for key %s", val, expected, key)
	}
}

func assertStringBackedEnvValue(t *testing.T, key string) {
	expected := rand.Text()
	t.Setenv(key, expected)
	val := GetString(key, "DUMMY_VALUE")
	if val != expected {
		if val == "DUMMY_VALUE" {
			t.Errorf("Returned default value when value was explicitly set")
		}
		t.Errorf("Returned unexpected value %s when expecting %s for key %s", val, expected, key)
	}
}
