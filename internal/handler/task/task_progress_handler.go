// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"net/http"

	"dockerCopilot/internal/logic/task"
	"dockerCopilot/internal/svc"
	"dockerCopilot/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func TaskProgressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TaskReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := task.NewTaskProgressLogic(r.Context(), svcCtx)
		resp, err := l.TaskProgress(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
