package handlers

import "skillForce/internal/usecase"

type Handler struct {
	useCase usecase.UsecaseInterface
}

func NewHandler(uc *usecase.Usecase) *Handler {
	return &Handler{
		useCase: uc,
	}
}
