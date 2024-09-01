package contract

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	clarketmjson "github.com/clarketm/json"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/qri-io/jsonschema"
)

type Contract struct {
	Data          []byte
	Path          string
	HTTPMethod    string
	ContentType   string
	GitToken      string
	RepositoryURL string
}

type APISpec struct {
	contract   *Contract
	bodySchema *jsonschema.Schema
	parameters []*openapi3.ParameterRef
}

func NewAPIAPISpec(contract Contract) (*APISpec, error) {
	obj := &APISpec{
		contract: &Contract{
			Data:          contract.Data,
			Path:          contract.Path,
			HTTPMethod:    contract.HTTPMethod,
			ContentType:   contract.ContentType,
			GitToken:      contract.GitToken,
			RepositoryURL: contract.RepositoryURL,
		},
	}

	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	data, err := loader.LoadFromData(obj.contract.Data)
	if err != nil {
		return nil, err
	}

	if data.Components != nil {
		derefSchemas(data.Components.Schemas)
	}

	s := data.Paths.Find(obj.contract.Path)
	if s == nil {
		return nil, errors.New("contract path: " + obj.contract.Path + " not found")
	}
	var operation *openapi3.Operation
	switch strings.ToUpper(obj.contract.HTTPMethod) {
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
		return nil, errors.New("contract path: " + obj.contract.Path + " " + obj.contract.HTTPMethod + " operation not found")
	}

	obj.parameters = operation.Parameters

	requestBody := operation.RequestBody
	if requestBody == nil {
		return obj, nil
	}

	if obj.contract.ContentType != "application/json" {
		return nil, errors.New("contract path: " + obj.contract.Path + " " + obj.contract.HTTPMethod +
			" current supported media types are: application/json")
	}

	mediaType := requestBody.Value.Content.Get(obj.contract.ContentType)
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

func (v *APISpec) ValidateBodyBytes(json []byte) error {
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

func (v *APISpec) ValidateBodyInterface(json interface{}) error {
	if v.bodySchema == nil {
		return nil
	}

	b, err := clarketmjson.Marshal(json)
	if err != nil {
		return err
	}
	return v.ValidateBodyBytes(b)
}

func (v *APISpec) ValidateRequestHeaders(headers http.Header) error {
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
