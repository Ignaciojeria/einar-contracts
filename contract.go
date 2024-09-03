package contract

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	clarketmjson "github.com/clarketm/json"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/qri-io/jsonschema"
)

type EndpointDetails struct {
	ContractData []byte
	Path         string
	HTTPMethod   string
	ContentType  string
}

type Endpoint struct {
	endpointDetails *EndpointDetails
	bodySchema      *jsonschema.Schema
	parameters      []*openapi3.ParameterRef
}

func LoadSpecEndpoint(endpointDetails EndpointDetails) (*Endpoint, error) {
	obj := &Endpoint{
		endpointDetails: &EndpointDetails{
			ContractData: endpointDetails.ContractData,
			Path:         endpointDetails.Path,
			HTTPMethod:   endpointDetails.HTTPMethod,
			ContentType:  endpointDetails.ContentType,
		},
	}

	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	data, err := loader.LoadFromData(obj.endpointDetails.ContractData)
	if err != nil {
		return nil, err
	}

	if data.Components != nil {
		derefSchemas(data.Components.Schemas)
	}

	s := data.Paths.Find(obj.endpointDetails.Path)
	if s == nil {
		return nil, errors.New("contract path: " + obj.endpointDetails.Path + " not found")
	}
	var operation *openapi3.Operation
	switch strings.ToUpper(obj.endpointDetails.HTTPMethod) {
	case http.MethodGet:
		operation = s.Get
	case http.MethodPost:
		operation = s.Post
	case http.MethodPut:
		operation = s.Put
	case http.MethodPatch:
		operation = s.Patch
	case http.MethodDelete:
		operation = s.Delete
	default:
		return nil, errors.New("invalid HTTPMethod operation. Available options: " +
			http.MethodGet + " " +
			http.MethodPost + " " +
			http.MethodPut + " " +
			http.MethodPatch + " " +
			http.MethodDelete)
	}

	if operation == nil {
		return nil, errors.New("contract path: " + obj.endpointDetails.Path + " " + obj.endpointDetails.HTTPMethod + " operation not found")
	}

	obj.parameters = operation.Parameters

	requestBody := operation.RequestBody
	if requestBody == nil {
		return obj, nil
	}

	if obj.endpointDetails.ContentType != "application/json" {
		return nil, errors.New("contract path: " + obj.endpointDetails.Path + " " + obj.endpointDetails.HTTPMethod +
			" current supported media types are: application/json")
	}

	mediaType := requestBody.Value.Content.Get(obj.endpointDetails.ContentType)
	if mediaType == nil {
		mediaType = &openapi3.MediaType{}
	}

	schemaRef := mediaType.Schema
	if schemaRef == nil {
		obj.bodySchema = &jsonschema.Schema{}
	} else {
		schemaJSON, err := schemaRef.Value.MarshalJSON()
		if err != nil {
			return nil, err
		}
		jsonSchema := &jsonschema.Schema{}
		if err := json.Unmarshal(schemaJSON, jsonSchema); err != nil {
			return nil, err
		}
		obj.bodySchema = jsonSchema
	}

	return obj, nil
}

func derefSchemas(schemas openapi3.Schemas) {
	for _, schemaRef := range schemas {
		derefSchemaRef(schemaRef)
	}
}

func derefSchemaRefs(schemaRefs openapi3.SchemaRefs) {
	for _, schemaRef := range schemaRefs {
		derefSchemaRef(schemaRef)
	}
}

func derefSchemaRef(schemaRef *openapi3.SchemaRef) {
	if schemaRef == nil {
		return
	}

	schemaRef.Ref = ""

	val := schemaRef.Value

	derefSchemaRefs(val.OneOf)
	derefSchemaRefs(val.AnyOf)
	derefSchemaRefs(val.AllOf)
	derefSchemaRef(val.Not)
	derefSchemaRef(val.Items)
	derefSchemas(val.Properties)
	derefSchemaRef(val.AdditionalProperties.Schema)
}

func (v *Endpoint) ValidateBodyBytes(json []byte) error {
	if v.bodySchema == nil {
		return nil
	}

	errs, err := v.bodySchema.ValidateBytes(context.Background(), json)
	if err != nil {
		return err
	}
	if len(errs) == 0 {
		return nil
	}
	var sb strings.Builder
	for _, e := range errs {
		sb.WriteString(fmt.Sprintf("Validation error: %s\n", e.Error()))
	}
	return errors.New(sb.String())
}

func (v *Endpoint) ValidateBodyInterface(json interface{}) error {
	if v.bodySchema == nil {
		return nil
	}

	b, err := clarketmjson.Marshal(json)
	if err != nil {
		return err
	}
	return v.ValidateBodyBytes(b)
}

func (v *Endpoint) ValidateRequestHeaders(headers http.Header) error {
	var sb strings.Builder

	for _, paramRef := range v.parameters {
		param := paramRef.Value

		if param.In == openapi3.ParameterInHeader && param.Required {
			headerName := param.Name
			headerValue := headers.Get(headerName)

			if headerValue == "" {
				sb.WriteString(fmt.Sprintf("Header '%s' is required but not present\n", headerName))
			}
		}
	}

	if sb.Len() > 0 {
		return errors.New(sb.String())
	}

	return nil
}
