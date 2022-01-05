package application

import (
	"context"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	application "github.com/Beep-Technologies/beepbeep3-ocpp/internal/repository/application"
	"github.com/Beep-Technologies/beepbeep3-ocpp/pkg/util"
)

type Service struct {
	db          *gorm.DB
	application application.BaseRepo
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:          db,
		application: application.NewBaseRepo(db),
	}
}

func (srv Service) GetApplicationCallbacks(ctx context.Context, req *rpc.GetApplicationCallbacksReq) (*rpc.GetApplicationCallbacksResp, error) {
	callbacks, err := srv.application.GetApplicationCallbacks(ctx, req.EntityCode, req.CallbackEvent)

	if err != nil {
		return nil, err
	}

	resCallbacks := make([]*rpc.ApplicationCallback, 0)
	for _, callback := range callbacks {
		resCallback := rpc.ApplicationCallback{}
		util.ConvertCopyStruct(&resCallback, &callback, map[string]util.ConverterFunc{})
		resCallbacks = append(resCallbacks, &resCallback)
	}

	res := &rpc.GetApplicationCallbacksResp{
		ApplicationCallbacks: resCallbacks,
	}

	return res, nil
}
