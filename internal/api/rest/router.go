package rest

import (
	"snappfood/internal/api/rest/handlers/order"
	"snappfood/internal/api/rest/handlers/vendor"
)

func (s *Server) SetupAPIRoutes(
	orderHandler *order.OrderHandler,
	vendorHandler *vendor.VendorHandler,
) {
	r := s.engine

	r.POST("orders/report", orderHandler.ReportDelay)
	r.GET("orders/proccess", orderHandler.Proccess)
	r.GET("vendors/report", vendorHandler.GetDelayReport)

}
