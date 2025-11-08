package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"go-template/api/h"
	"go-template/utils/system"
)

const systemTag = "system-系统工具"

type openPathInput struct {
	Body struct {
		Path string `json:"path" doc:"目标路径" required:"true"`
	} `json:"body"`
}

func registerSystemRoutes(group *huma.Group) {
	huma.Post(group, "/system/open-explorer", func(ctx context.Context, input *openPathInput) (*h.MessageResponse, error) {
		if err := system.OpenExplorer(input.Body.Path); err != nil {
			return nil, mapSystemError(err)
		}

		resp := h.NewMessageResponse("explorer opened")
		resp.Status = http.StatusOK
		return resp, nil
	}, func(op *huma.Operation) {
		op.OperationID = "system-open-explorer"
		op.Summary = "打开文件管理器"
		op.Tags = []string{systemTag}
	})

	huma.Post(group, "/system/open-terminal", func(ctx context.Context, input *openPathInput) (*h.MessageResponse, error) {
		if err := system.OpenTerminal(input.Body.Path); err != nil {
			return nil, mapSystemError(err)
		}

		resp := h.NewMessageResponse("terminal opened")
		resp.Status = http.StatusOK
		return resp, nil
	}, func(op *huma.Operation) {
		op.OperationID = "system-open-terminal"
		op.Summary = "打开终端"
		op.Tags = []string{systemTag}
	})
}

func mapSystemError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, system.ErrUnsupportedOS):
		return huma.Error501NotImplemented(err.Error())
	case errors.Is(err, system.ErrNoFileManager),
		errors.Is(err, system.ErrNoTerminal):
		return huma.Error503ServiceUnavailable(err.Error())
	default:
		return huma.Error500InternalServerError(err.Error())
	}
}
