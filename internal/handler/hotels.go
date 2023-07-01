package handler

import "context"

func NewDefaultHotelsHandler() *HotelsHandler {
	return &HotelsHandler{}
}

type HotelsHandler struct {
}

func (h *HotelsHandler) Execute(
	ctx context.Context, params map[string][]string,
) (interface{}, error) {
	return nil, nil
}
