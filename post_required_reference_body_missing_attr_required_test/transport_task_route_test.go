package post_required_reference_body_missing_attr_required_test

import (
	"log/slog"
	"net/http"
	"strings"
	"testing"
)

func TestTransportTaskRoute(t *testing.T) {

	schema, err := NewTransportTaskRoute()
	if err != nil {
		t.Fatalf("Failed to create TransportTaskRoute schema: %v", err)
	}

	// Validar el JSON resultante contra el esquema
	err = schema.ValidateBodyBytes(fakeExampleBody)
	if err == nil {
		t.Error("Validation should have failed, but it passed")
		slog.Error("Validation failed to catch missing required fields")
	} else {
		expectedErrorSubstring := `"businessIdentifiers" value is required`
		if !strings.Contains(err.Error(), expectedErrorSubstring) {
			t.Errorf("Expected error to contain '%s', got '%s'", expectedErrorSubstring, err.Error())
		} else {
			slog.Info("Validation correctly caught missing required fields")
		}
	}

	// Pruebas para validar las cabeceras requeridas
	headerTests := []struct {
		name      string
		headers   http.Header
		wantError bool
		errorMsg  string
	}{
		{
			name: "All required headers present",
			headers: setupHeaders(map[string]string{
				"X-cmRef":    "SIMPLIROUTE",
				"X-country":  "CL",
				"X-commerce": "CORP",
			}),
			wantError: false,
		},
		{
			name: "Missing required header X-cmRef",
			headers: setupHeaders(map[string]string{
				"X-country":  "CL",
				"X-commerce": "CORP",
			}),
			wantError: true,
			errorMsg:  "Header 'X-cmRef' is required but not present",
		},
		{
			name: "Missing required header X-country",
			headers: setupHeaders(map[string]string{
				"X-cmRef":    "SIMPLIROUTE",
				"X-commerce": "CORP",
			}),
			wantError: true,
			errorMsg:  "Header 'X-country' is required but not present",
		},
		{
			name: "Missing required header X-commerce",
			headers: setupHeaders(map[string]string{
				"X-cmRef":   "SIMPLIROUTE",
				"X-country": "CL",
			}),
			wantError: true,
			errorMsg:  "Header 'X-commerce' is required but not present",
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
	"referenceId": "simplirouteRouteID",
	"vehicle": {
	  "plate": "CXC-48900",
	  "references": [
		{
		  "type": "geosortId",
		  "value": "3191"
		}
	  ]
	},
	"carrier": {
	  "name": "carrier_name",
	  "document": {
		"type": "RUT",
		"value": "18666636-4"
	  },
	  "references": [
		{
		  "type": "geosortId",
		  "value": "3191"
		}
	  ]
	},
	"visits": [
	  {
		"visit": {
		  "visitId": "bad44a88-454b-431a-9251-7ef63104b11e",
		  "deliverySequence": 1
		},
		"orders": [
		  {
			"referenceId": "20454021",
			"type": "DELIVERY_ORDER",
			"references": [
			  {
				"type": "anyConsumerOrderReference",
				"value": "anyConsumerOrderReferenceValue"
			  }
			],
			"origin": {
			  "nodeInfo": {
				"nodeId": "aad44a88-454b-431a-9251-7ef63104b112",
				"operatorId": "string",
				"operatorType": "string",
				"references": [
				  {
					"type": "facility",
					"value": "7100"
				  }
				]
			  }
			},
			"destination": {
			  "addressInfo": {
				"contact": {
				  "fullName": "Walther Kallina",
				  "contacts": [
					{
					  "type": "phone",
					  "value": "977575871"
					}
				  ],
				  "documents": [
					{
					  "type": "RUT",
					  "value": "9358315-9"
					}
				  ]
				},
				"politicalAreaId": "string",
				"district": "string",
				"addressLine1": "Santa Elena Calle del Sol",
				"addressLine2": "44",
				"addressLine3": "string",
				"latitude": -33.2180344,
				"longitude": -70.73417529999999,
				"zipCode": "string",
				"timeZone": "string"
			  }
			},
			"items": [
			  {
				"itemId": "242705",
				"description": "JUNTACREST ULTRAMAX CHOCOLE",
				"quantity": {
				  "quantityNumber": 2,
				  "quantityUnit": "C/U"
				},
				"referenceSku": "1231",
				"references": [
				  {
					"type": "offering",
					"value": "242705"
				  }
				],
				"insurance": {
				  "unitValue": 2300,
				  "currency": "CLP"
				},
				"dimensions": {
				  "height": 0.44,
				  "width": 0.44,
				  "depth": 0.44,
				  "unit": "cm"
				},
				"weight": {
				  "value": 3.6,
				  "unit": "kg"
				}
			  }
			],
			"packages": [
			  {
				"lpn": "LPN_12345",
				"packageType": "cartonBox",
				"dimensions": {
				  "height": 0.44,
				  "width": 0.44,
				  "depth": 0.44,
				  "unit": "cm"
				},
				"weight": {
				  "value": 12.3,
				  "unit": "kg"
				},
				"insurance": {
				  "unitValue": 2300,
				  "currency": "CLP"
				},
				"itemReferences": [
				  {
					"itemId": "242705",
					"quantity": {
					  "quantityNumber": 1,
					  "quantityUnit": "C/U"
					}
				  }
				]
			  }
			],
			"promisedDate": {
			  "dateFrom": "1985-12-01",
			  "dateTo": "1985-12-01",
			  "timeRangeFrom": "9:00",
			  "timeRangeTo": "21:00",
			  "serviceType": "HOME_DELIVERY",
			  "serviceCategory": "REGULAR"
			},
			"extraFields": [
			  {
				"type": "simpliroute.example",
				"values": [
				  "SD_98"
				]
			  }
			],
			"skills": [
			  {
				"type": "truck",
				"value": "77031"
			  }
			]
		  }
		]
	  }
	]
  }
`)
