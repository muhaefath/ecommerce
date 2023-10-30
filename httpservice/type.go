package httpservice

import "ecommerce/service"

// Handler is struct standart to handle accounting_journal handler
type Handler struct {
	ecommerceSrv service.EcommerceProvider
}

// HandlerConfig is standart configuration for accounting_journal config
type HandlerConfig struct {
	EcommerceSrv service.EcommerceProvider
}
