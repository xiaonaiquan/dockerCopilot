// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task

import (
	"context"

	"dockerCopilot/internal/svc"
	"dockerCopilot/internal/types"
	"dockerCopilot/internal/utiles"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskProgressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaskProgressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskProgressLogic {
	return &TaskProgressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskProgressLogic) TaskProgress(req *types.TaskReq) (resp *types.Resp, err error) {
	// todo: add your logic here and delete this line
	task, exist := utiles.GlobalTaskManager.GetTask(req.Uuid)
	if !exist {
		return &types.Resp{
			Code: 1,
			Msg:  "任务不存在",
			Data: "",
		}, nil
	}
	return &types.Resp{
		Code: 0,
		Msg:  "任务状态获取成功",
		Data: task.Status,
	}, nil
}
