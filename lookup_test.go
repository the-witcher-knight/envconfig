package env

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type webConfig struct {
	Host string `env:"HOST,required"`
	Port int    `env:"PORT,required"`
}
type appConfig struct {
	Web             webConfig `env:"WEB"`
	Env             string    `env:"ENV,expectedValues=development production"`
	DB              string    `env:"DB,required"`
	EnableLogging   bool      `env:"ENABLE_LOGGING"`
	APIKey          string    `env:"API_KEY,required"`
	RedisURL        string    `env:"REDIS_URL,required"`
	RedisPassword   string    `env:"REDIS_PASSWORD"`
	JWTSecret       string    `env:"JWT_SECRET,required"`
	LogLevel        string    `env:"LOG_LEVEL"`
	MaxConnections  int       `env:"MAX_CONNECTIONS"`
	ReadTimeout     int       `env:"READ_TIMEOUT"`
	WriteTimeout    int       `env:"WRITE_TIMEOUT"`
	GracefulTimeout int       `env:"GRACEFUL_TIMEOUT"`
	CacheEnabled    bool      `env:"CACHE_ENABLED"`
	CacheExpiry     int       `env:"CACHE_EXPIRY"`
	SMTPHost        string    `env:"SMTP_HOST"`
	SMTPPort        int       `env:"SMTP_PORT"`
	SMTPUser        string    `env:"SMTP_USER"`
	SMTPPassword    string    `env:"SMTP_PASSWORD"`
	AllowCORS       bool      `env:"ALLOW_CORS"`
	StaticFilesPath string    `env:"STATIC_FILES_PATH"`
	SessionTimeout  int       `env:"SESSION_TIMEOUT"`
	MaintenanceMode bool      `env:"MAINTENANCE_MODE" `
	ExternalAPIURL  string    `env:"EXTERNAL_API_URL"`
}

func Test_Lookup(t *testing.T) {
	tcs := map[string]struct {
		givenInput map[string]string
		expResult  appConfig
		expErr     error
	}{
		"success": {
			givenInput: map[string]string{
				"HOST":              "0.0.0.0",
				"PORT":              "8080",
				"DB":                "postgres://postgres:password@localhost:5432/mydb?sslmode=disable",
				"ENABLE_LOGGING":    "true",
				"ENV":               "production",
				"API_KEY":           "12345-abcde-67890-fghij",
				"REDIS_URL":         "redis://localhost:6379",
				"REDIS_PASSWORD":    "redispassword",
				"JWT_SECRET":        "mysecretjwtkey",
				"LOG_LEVEL":         "debug",
				"MAX_CONNECTIONS":   "150",
				"READ_TIMEOUT":      "30",
				"WRITE_TIMEOUT":     "30",
				"GRACEFUL_TIMEOUT":  "15",
				"CACHE_ENABLED":     "true",
				"CACHE_EXPIRY":      "600",
				"SMTP_HOST":         "smtp.mailtrap.io",
				"SMTP_PORT":         "2525",
				"SMTP_USER":         "user@mailtrap.io",
				"SMTP_PASSWORD":     "smtppassword",
				"ALLOW_CORS":        "true",
				"STATIC_FILES_PATH": "/var/www/static",
				"SESSION_TIMEOUT":   "86400",
				"MAINTENANCE_MODE":  "false",
				"EXTERNAL_API_URL":  "https://api.external-service.com/v1",
			},
			expResult: appConfig{
				Web: webConfig{
					Host: "0.0.0.0",
					Port: 8080,
				},
				DB:              "postgres://postgres:password@localhost:5432/mydb?sslmode=disable",
				EnableLogging:   true,
				Env:             "production",
				APIKey:          "12345-abcde-67890-fghij",
				RedisURL:        "redis://localhost:6379",
				RedisPassword:   "redispassword",
				JWTSecret:       "mysecretjwtkey",
				LogLevel:        "debug",
				MaxConnections:  150,
				ReadTimeout:     30,
				WriteTimeout:    30,
				GracefulTimeout: 15,
				CacheEnabled:    true,
				CacheExpiry:     600,
				SMTPHost:        "smtp.mailtrap.io",
				SMTPPort:        2525,
				SMTPUser:        "user@mailtrap.io",
				SMTPPassword:    "smtppassword",
				AllowCORS:       true,
				StaticFilesPath: "/var/www/static",
				SessionTimeout:  86400,
				MaintenanceMode: false,
				ExternalAPIURL:  "https://api.external-service.com/v1",
			},
		},
		"error - validation fail": {
			givenInput: map[string]string{
				"ENV": "SHINOBI",
			},
			expErr: errors.New("HOST is required\nPORT is required\nENV is unexpected value: SHINOBI\nDB is required\nAPI_KEY is required\nREDIS_URL is required\nJWT_SECRET is required"),
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// Given
			for k, v := range tc.givenInput {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tc.givenInput {
					os.Unsetenv(k)
				}
			}()

			// When
			var out appConfig
			err := Lookup(&out)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, out)
			}
		})
	}
}

type mockValidator struct{}

func (v mockValidator) Validate(_, _ string) error {
	return nil
}

func Test_AddValidator(t *testing.T) {
	// Given
	const mockValidatorKey = "mock"
	f := func(_ string) Validator {
		return mockValidator{}
	}

	// When
	AddValidator(mockValidatorKey, f)
	defer delete(validatorRegistry, mockValidatorKey) // Clean up

	// Then
	require.NotNil(t, validatorRegistry[mockValidatorKey])
}

func Benchmark_Lookup(b *testing.B) {
	// Setup
	values := map[string]string{
		"HOST":              "127.0.0.1",
		"PORT":              "8080",
		"DB":                "postgres://postgres:password@localhost:5432/mydb?sslmode=disable",
		"ENABLE_LOGGING":    "true",
		"ENV":               "production",
		"API_KEY":           "12345-abcde-67890-fghij",
		"REDIS_URL":         "redis://localhost:6379",
		"REDIS_PASSWORD":    "redispassword",
		"JWT_SECRET":        "mysecretjwtkey",
		"LOG_LEVEL":         "debug",
		"MAX_CONNECTIONS":   "150",
		"READ_TIMEOUT":      "30",
		"WRITE_TIMEOUT":     "30",
		"GRACEFUL_TIMEOUT":  "15",
		"CACHE_ENABLED":     "true",
		"CACHE_EXPIRY":      "600",
		"SMTP_HOST":         "smtp.mailtrap.io",
		"SMTP_PORT":         "2525",
		"SMTP_USER":         "user@mailtrap.io",
		"SMTP_PASSWORD":     "smtppassword",
		"ALLOW_CORS":        "true",
		"STATIC_FILES_PATH": "/var/www/static",
		"SESSION_TIMEOUT":   "86400",
		"MAINTENANCE_MODE":  "false",
		"EXTERNAL_API_URL":  "https://api.external-service.com/v1",
	}
	for k, v := range values {
		os.Setenv(k, v)
	}
	defer func() {
		for k := range values {
			os.Unsetenv(k)
		}
	}()

	// Run
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var out appConfig
		_ = Lookup(&out)
	}
}
