package container

import (
	"dockerCopilot/internal/logic/container"
	"dockerCopilot/internal/svc"
	"dockerCopilot/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RenameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ContainerRenameReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := container.NewRenameLogic(r.Context(), svcCtx)
		resp, err := l.Rename(&req)
		if err != nil {
			httpx.WriteJson(w, resp.Code, resp)
		} else {
			httpx.WriteJson(w, resp.Code, resp)
		}
	}
}
