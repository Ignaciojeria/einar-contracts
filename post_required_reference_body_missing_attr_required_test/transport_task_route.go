package post_required_reference_body_missing_attr_required_test

import (
	_ "embed"

	contract "github.com/Ignaciojeria/einar-contracts"
)

type TransportTaskRoute struct {
	RouteCreation *contract.Endpoint
	OrderCreation *contract.Endpoint
}

//go:embed transport_task_route.json
var transport_task_route []byte

func NewTransportTaskRoute() (TransportTaskRoute, error) {
	routeCreation, err := contract.LoadSpecEndpoint(
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

	orderCreation, err := contract.LoadSpecEndpoint(
		contract.Contract{
			Data:        transport_task_route,
			Path:        "/order_creation",
			HTTPMethod:  "POST",
			ContentType: "application/json",
		},
	)
	if err != nil {
		return TransportTaskRoute{}, err
	}
	return TransportTaskRoute{
		RouteCreation: routeCreation,
		OrderCreation: orderCreation}, nil
}
