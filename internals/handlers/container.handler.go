package handlers

import "mini-paas/internals/services"

type ContainerHandler struct {
	service services.ContainerService
}

func NewContainerHandler(service services.ContainerService) *ContainerHandler{
	return &ContainerHandler{service: service}
}

