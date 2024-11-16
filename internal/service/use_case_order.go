package service

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/go-worker-order/internal/repository/dynamo"
	"github.com/go-worker-order/internal/repository/storage"
	"github.com/go-worker-order/internal/core"
	"github.com/go-worker-order/internal/erro"
	"github.com/go-worker-order/internal/lib"
)

var childLogger = log.With().Str("service", "service").Logger()

type WorkerService struct {
	workerRepo		*storage.WorkerRepository
	appServer		*core.WorkerAppServer
	workerDynamo	*dynamo.DynamoRepository
}

func NewWorkerService(	workerRepo		*storage.WorkerRepository,
						appServer		*core.WorkerAppServer,
						workerDynamo	*dynamo.DynamoRepository) *WorkerService{
	childLogger.Debug().Msg("NewWorkerService")

	return &WorkerService{
		workerRepo:	workerRepo,
		workerDynamo: 	workerDynamo,
		appServer:	appServer,
	}
}

func (s WorkerService) OrderUpdate(ctx context.Context, order core.Order) (error){
	childLogger.Debug().Msg("OrderUpdate")
	
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

	_, err = s.workerDynamo.Add(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func (s WorkerService) OrderAdd(ctx context.Context, order core.Order) (error){
	childLogger.Debug().Msg("OrderAdd")
	
	span := lib.Span(ctx, "service.OrderAdd")

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
	res, err := s.workerRepo.Add(ctx, tx, &order)
	if err != nil {
		return err
	}
	
	order.ID = res.ID
	order.CreateAt = res.CreateAt

	_, err = s.workerDynamo.Add(ctx, order)
	if err != nil {
		return err
	}

	return nil
}