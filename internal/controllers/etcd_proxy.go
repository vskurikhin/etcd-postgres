package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/victor-skurikhin/etcd-client/v1/internal/controllers/dto"
	"github.com/victor-skurikhin/etcd-client/v1/internal/env"
	"github.com/victor-skurikhin/etcd-client/v1/internal/services"
	"log/slog"
	"sync"

	clientV3 "go.etcd.io/etcd/client/v3"
)

type EtcdProxy interface {
	Delete(*fiber.Ctx) error
	Get(*fiber.Ctx) error
	Put(*fiber.Ctx) error
}

type etcdProxy struct {
	clientConfig     clientV3.Config
	etcdProxyService services.EtcdProxyService
	sLog             *slog.Logger
}

type tContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

type tIdentity struct {
	ID        string
	RequestID uuid.UUID
}

var _ EtcdProxy = (*etcdProxy)(nil)
var (
	onceEtcdProxy = new(sync.Once)
	etcdProxyCont *etcdProxy
)

// GetEtcdProxyController — потокобезопасное (thread-safe) создание
// REST веб-сервиса etcd proxy.
func GetEtcdProxyController(ctx context.Context, cfg env.Config) EtcdProxy {

	onceEtcdProxy.Do(func() {
		etcdProxyCont = new(etcdProxy)
		etcdProxyCont.clientConfig = *cfg.EtcdClientConfig()
		etcdProxyCont.etcdProxyService = services.GetEtcdProxyService(ctx, cfg)
		etcdProxyCont.sLog = cfg.Logger()
	})
	return etcdProxyCont
}

func (f *etcdProxy) Delete(fCtx *fiber.Ctx) error {

	ctxCancel, identity, err := f.contextWithRequestIdentity(fCtx)

	if err != nil {
		return fCtx.
			Status(fiber.StatusBadRequest).
			JSON(dto.StatusMessageInvalidRequestID(identity.ID))
	}
	defer ctxCancel.cancel()
	key := fCtx.Params("name", "default")

	if err = f.etcdProxyService.ApiDelete(ctxCancel.ctx, key); err != nil {
		return fCtx.
			Status(fiber.StatusBadRequest).
			JSON(dto.StatusMessageRequestID{
				Status:    "fail",
				Message:   err.Error(),
				RequestID: identity.RequestID,
			})
	} else {
		return fCtx.
			Status(fiber.StatusOK).
			JSON(dto.StatusRequestID{Status: "success", RequestID: identity.RequestID})
	}
}

func (f *etcdProxy) Get(fCtx *fiber.Ctx) error {

	ctxCancel, identity, err := f.contextWithRequestIdentity(fCtx)

	if err != nil {
		return fCtx.
			Status(fiber.StatusBadRequest).
			JSON(dto.StatusMessageInvalidRequestID(identity.ID))
	}
	defer ctxCancel.cancel()
	key := fCtx.Params("name", "default")

	if result, err := f.etcdProxyService.ApiGet(ctxCancel.ctx, key); err != nil {
		code := fiber.StatusBadRequest

		if err == services.ErrNotFound {
			code = fiber.StatusNotFound
		}
		return fCtx.
			Status(code).
			JSON(dto.StatusMessageRequestID{
				Status:    "fail",
				Message:   err.Error(),
				RequestID: identity.RequestID,
			})
	} else {
		return fCtx.
			Status(fiber.StatusOK).
			JSON(dto.StatusResultRequestID{Status: "success", Result: result, RequestID: identity.RequestID})
	}
}

func (f *etcdProxy) Put(fCtx *fiber.Ctx) error {

	ctxCancel, identity, err := f.contextWithRequestIdentity(fCtx)

	if err != nil {
		return fCtx.
			Status(fiber.StatusBadRequest).
			JSON(dto.StatusMessageInvalidRequestID(identity.ID))
	}
	defer ctxCancel.cancel()
	var payload dto.Result

	if err = fCtx.BodyParser(&payload); err != nil {
		return fCtx.
			Status(fiber.StatusBadRequest).
			JSON(dto.StatusMessageRequestID{
				Status:    "fail",
				Message:   err.Error(),
				RequestID: identity.RequestID,
			})
	}
	errors := dto.ValidateStruct(payload)

	if errors != nil {
		return fCtx.
			Status(fiber.StatusBadRequest).
			JSON(errors)
	}
	if err = fCtx.BodyParser(&payload); err != nil {
		return fCtx.
			Status(fiber.StatusBadRequest).
			JSON(dto.StatusMessageRequestID{
				Status:    "fail",
				Message:   err.Error(),
				RequestID: identity.RequestID,
			})
	}
	key := fCtx.Params("name", "default")

	if err = f.etcdProxyService.ApiPut(ctxCancel.ctx, dto.KeyValue{Key: key, Value: payload.Value}); err != nil {
		return fCtx.
			Status(fiber.StatusBadRequest).
			JSON(dto.StatusMessageRequestID{
				Status:    "fail",
				Message:   err.Error(),
				RequestID: identity.RequestID,
			})
	} else {
		return fCtx.
			Status(fiber.StatusOK).
			JSON(dto.StatusRequestID{Status: "success", RequestID: identity.RequestID})
	}
}

func (f *etcdProxy) contextWithRequestIdentity(fCtx *fiber.Ctx) (tContext, tIdentity, error) {

	var err error
	var requestId uuid.UUID

	if id, ok := fCtx.Locals("requestid").(string); ok {
		if requestId, err = uuid.Parse(id); err != nil {
			return tContext{}, tIdentity{ID: id}, err
		}
	}
	ctx, cancel := context.WithTimeout(
		context.WithValue(fCtx.Context(), "request-id", requestId.String()),
		f.clientConfig.DialTimeout,
	)
	return tContext{ctx: ctx, cancel: cancel}, tIdentity{RequestID: requestId}, nil
}
