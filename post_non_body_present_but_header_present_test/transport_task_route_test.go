package post_non_body_present_but_header_present_test

import (
	"log/slog"
	"net/http"
	"strings"
	"testing"
)

func TestTransportTaskRouteCreation(t *testing.T) {
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

	// Pruebas para validar las cabeceras requeridas en la ruta /route_creation
	headerTests := []struct {
		name      string
		headers   http.Header
		wantError bool
		errorMsgs []string
	}{
		{
			name: "All required headers present - /route_creation",
			headers: setupHeaders(map[string]string{
				"X-cmRef":    "SIMPLIROUTE",
				"X-country":  "CL",
				"X-commerce": "CORP",
			}),
			wantError: false,
		},
		{
			name: "Missing required header X-cmRef - /route_creation",
			headers: setupHeaders(map[string]string{
				"X-country":  "CL",
				"X-commerce": "CORP",
			}),
			wantError: true,
			errorMsgs: []string{
				"Header 'X-cmRef' is required but not present",
			},
		},
		{
			name: "Missing required header X-country - /route_creation",
			headers: setupHeaders(map[string]string{
				"X-cmRef":    "SIMPLIROUTE",
				"X-commerce": "CORP",
			}),
			wantError: true,
			errorMsgs: []string{
				"Header 'X-country' is required but not present",
			},
		},
		{
			name: "Missing required header X-commerce - /route_creation",
			headers: setupHeaders(map[string]string{
				"X-cmRef":   "SIMPLIROUTE",
				"X-country": "CL",
			}),
			wantError: true,
			errorMsgs: []string{
				"Header 'X-commerce' is required but not present",
			},
		},
		{
			name:      "Missing all required headers - /route_creation",
			headers:   setupHeaders(map[string]string{}),
			wantError: true,
			errorMsgs: []string{
				"Header 'X-cmRef' is required but not present",
				"Header 'X-country' is required but not present",
				"Header 'X-commerce' is required but not present",
			},
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
				// Dividir el mensaje de error en líneas y compararlas con las esperadas
				trimmedErrs := strings.Split(strings.TrimSpace(err.Error()), "\n")

				if len(trimmedErrs) != len(tt.errorMsgs) {
					t.Errorf("Test %s: expected %d errors, got %d", tt.name, len(tt.errorMsgs), len(trimmedErrs))
				}

				for i, expected := range tt.errorMsgs {
					if i < len(trimmedErrs) {
						if strings.TrimSpace(trimmedErrs[i]) != strings.TrimSpace(expected) {
							t.Errorf("Test %s: expected error message '%s', got '%s'", tt.name, expected, trimmedErrs[i])
						}
					}
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
