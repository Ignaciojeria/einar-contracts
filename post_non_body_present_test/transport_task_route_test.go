package post_non_body_present_test

import (
	"log/slog"
	"net/http"
	"strings"
	"testing"
)

func TestTransportTaskRoute(t *testing.T) {
	// Crear una instancia de la estructura que implementa el método ValidateBytes
	schema, err := NewTransportTaskRoute()
	if err != nil {
		t.Fatalf("Failed to create TransportTaskRoute schema: %v", err)
	}

	// Validar el JSON resultante contra el esquema
	err = schema.ValidateBodyBytes(fakeExampleBody)
	if err != nil {
		slog.Error("Validation failed", slog.String("error", err.Error()))
		t.Errorf("Validation failed: %v", err)
	}

	// Si tienes más pruebas, puedes agregarlas aquí
	tests := []struct {
		name      string
		input     []byte
		wantError bool
	}{
		{
			name:      "Valid JSON",
			input:     fakeExampleBody,
			wantError: false,
		},
		// Agrega más casos de prueba según sea necesario
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.ValidateBodyBytes(tt.input)
			if (err != nil) != tt.wantError {
				slog.Error("Test validation error", slog.String("test", tt.name), slog.String("error", err.Error()))
				t.Errorf("Test %s: expected error = %v, got %v", tt.name, tt.wantError, err != nil)
			}
		})
	}

	// Pruebas para validar las cabeceras no requeridas
	headerTests := []struct {
		name      string
		headers   http.Header
		wantError bool
		errorMsg  string
	}{
		{
			name: "All optional headers present",
			headers: setupHeaders(map[string]string{
				"X-cmRef":    "SIMPLIROUTE",
				"X-country":  "CL",
				"X-commerce": "CORP",
			}),
			wantError: false,
		},
		{
			name: "Missing optional header X-cmRef",
			headers: setupHeaders(map[string]string{
				"X-country":  "CL",
				"X-commerce": "CORP",
			}),
			wantError: false,
		},
		{
			name:    "Missing all optional headers",
			headers: setupHeaders(map[string]string{
				// No headers set
			}),
			wantError: false,
		},
	}

	for _, tt := range headerTests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.ValidateRequestHeaders(tt.headers)
			if (err != nil) != tt.wantError {
				slog.Error("Header validation error", slog.String("test", tt.name), slog.String("error", err.Error()))
				t.Errorf("Test %s: expected error = %v, got %v", tt.name, tt.wantError, err != nil)
			}
			if err != nil && tt.wantError {
				// Asegurarse de que los mensajes de error sean comparados correctamente
				trimmedErr := strings.TrimSpace(err.Error())
				trimmedExpected := strings.TrimSpace(tt.errorMsg)
				if trimmedErr != trimmedExpected {
					t.Errorf("Test %s: expected error message '%s', got '%s'", tt.name, trimmedExpected, trimmedErr)
				}
			}
		})
	}
}

// setupHeaders es una función auxiliar para configurar las cabeceras de manera correcta
func setupHeaders(headersMap map[string]string) http.Header {
	headers := http.Header{}
	for key, value := range headersMap {
		headers.Set(key, value)
	}
	return headers
}

var fakeExampleBody = []byte(`
{
}
`)
