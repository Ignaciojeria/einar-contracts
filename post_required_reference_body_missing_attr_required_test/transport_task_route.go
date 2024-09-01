package post_required_reference_body_missing_attr_required_test

import (
	_ "embed"

	"github.com/Ignaciojeria/contract"
)

type TransportTaskRoute struct {
	*contract.APISpec
}

//go:embed transport_task_route.json
var transport_task_route []byte

func NewTransportTaskRoute() (TransportTaskRoute, error) {
	spec, err := contract.NewAPIAPISpec(
		contract.Contract{
			Data:        transport_task_route,
			Path:        "/route_creation",
			HTTPMethod:  "POST",
			ContentType: "application/json",
		},
	)
	if err != nil {
		return TransportTaskRoute{}, err
	}
	return TransportTaskRoute{spec}, nil
}
