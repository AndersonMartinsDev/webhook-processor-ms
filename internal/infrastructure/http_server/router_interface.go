package http_server

import model "webhook-processor-ms/internal/infrastructure/commons/models/routes"

type RouterInterface interface {
	getRoutersModel() []model.RouteModel
}
