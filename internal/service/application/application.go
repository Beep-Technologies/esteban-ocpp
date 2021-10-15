package application

import (
	"context"

	"gorm.io/gorm"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/models"
	application "github.com/Beep-Technologies/beepbeep3-ocpp/internal/repository/application"
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

func (srv Service) CreateApplication(ctx context.Context, req *rpc.CreateApplicationReq) (*rpc.CreateApplicationResp, error) {
	a, err := srv.application.Create(ctx, models.OcppApplication{
		Name: req.Name,
	})

	if err != nil {
		return nil, err
	}

	res := &rpc.CreateApplicationResp{
		Application: &rpc.Application{
			Id:   a.ID,
			Name: a.Name,
		},
	}

	return res, nil
}

func (srv Service) GetApplicationByUuid(ctx context.Context, req *rpc.GetApplicationByUuidReq) (*rpc.GetApplicationByUuidResp, error) {
	a, err := srv.application.GetApplicationByUUID(ctx, req.ApplicationUuid)
	if err != nil {
		return nil, err
	}

	res := &rpc.GetApplicationByUuidResp{
		Application: &rpc.Application{
			Id:   a.ID,
			Uuid: a.UUID,
			Name: a.Name,
		},
	}

	return res, nil
}

func (srv Service) CreateApplicationCallback(ctx context.Context, req *rpc.CreateApplicationCallbackReq) (*rpc.CreateApplicationCallbackResp, error) {
	ac, err := srv.application.CreateCallback(ctx, models.OcppApplicationCallback{
		ApplicationID: req.ApplicationId,
		CallbackEvent: req.CallbackEvent,
		CallbackURL:   req.CallbackUrl,
	})

	if err != nil {
		return nil, err
	}

	res := &rpc.CreateApplicationCallbackResp{
		ApplicationCallback: &rpc.ApplicationCallback{
			Id:            ac.ID,
			ApplicationId: ac.ApplicationID,
			CallbackEvent: ac.CallbackEvent,
			CallbackUrl:   ac.CallbackURL,
		},
	}

	return res, nil
}

func (srv Service) GetApplicationCallback(ctx context.Context, req *rpc.GetApplicationCallbackReq) (*rpc.GetApplicationCallbackResp, error) {
	ac, err := srv.application.GetApplicationCallback(ctx, req.ApplicationId, req.CallbackEvent)

	if err != nil {
		return nil, err
	}

	res := &rpc.GetApplicationCallbackResp{
		ApplicationCallback: &rpc.ApplicationCallback{
			Id:            ac.ID,
			ApplicationId: ac.ApplicationID,
			CallbackEvent: ac.CallbackEvent,
			CallbackUrl:   ac.CallbackURL,
		},
	}

	return res, nil
}

func (srv Service) GetApplicationCallbacks(ctx context.Context, req *rpc.GetApplicationCallbacksReq) (*rpc.GetApplicationCallbacksResp, error) {
	acs, err := srv.application.GetApplicationCallbacks(ctx, req.ApplicationId)

	if err != nil {
		return nil, err
	}

	r := make([]*rpc.ApplicationCallback, 0)

	for _, ac := range acs {
		r = append(r, &rpc.ApplicationCallback{
			Id:            ac.ID,
			ApplicationId: ac.ApplicationID,
			CallbackEvent: ac.CallbackEvent,
			CallbackUrl:   ac.CallbackURL,
		})
	}

	res := &rpc.GetApplicationCallbacksResp{
		ApplicationCallbacks: r,
	}

	return res, nil
}

func (srv Service) DeleteApplicationCallback(ctx context.Context, req *rpc.DeleteApplicationCallbackReq) (*rpc.DeleteApplicationCallbackResp, error) {
	err := srv.application.DeleteCallback(ctx, req.ApplicationId)

	res := &rpc.DeleteApplicationCallbackResp{}

	return res, err
}
