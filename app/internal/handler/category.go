package handler

import (
	"github.com/alfattd/category-service/internal/platform/rabbitmq"
	"github.com/alfattd/category-service/internal/service"
)

type CategoryHandler struct {
	svc       *service.CategoryService
	publisher *rabbitmq.Publisher
}

func NewCategoryHandler(svc *service.CategoryService, publisher *rabbitmq.Publisher) *CategoryHandler {
	return &CategoryHandler{
		svc:       svc,
		publisher: publisher,
	}
}
