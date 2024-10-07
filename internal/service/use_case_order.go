package service

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/go-worker-order/internal/repository/storage"
	"github.com/go-worker-order/internal/core"
	"github.com/go-worker-order/internal/erro"
	"github.com/go-worker-order/internal/lib"
)

var childLogger = log.With().Str("service", "service").Logger()

type WorkerService struct {
	workerRepo		*storage.WorkerRepository
	appServer		*core.WorkerAppServer
}

func NewWorkerService(	workerRepo		*storage.WorkerRepository,
						appServer		*core.WorkerAppServer) *WorkerService{
	childLogger.Debug().Msg("NewWorkerService")

	return &WorkerService{
		workerRepo:	workerRepo,
		appServer:	appServer,
	}
}

func (s WorkerService) OrderUpdate(ctx context.Context, order core.Order) (error){
	childLogger.Debug().Msg("OrderUpdate")
	childLogger.Debug().Interface("===>order: ",order).Msg("")
	
	span := lib.Span(ctx, "service.OrderUpdate")

	tx, conn, err := s.workerRepo.StartTx(ctx)
	if err != nil {
		return err
	}
	
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
		s.workerRepo.ReleaseTx(conn)
		span.End()
	}()

	order.Status = "DONE"
	res_update, err := s.workerRepo.Update(ctx, tx, &order)
	if err != nil {
		return err
	}
	if res_update == 0 {
		err = erro.ErrUpdate
		return err
	}

	return nil
}