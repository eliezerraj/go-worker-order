package storage

import (
	"context"
	"errors"
	"time"

	"github.com/go-worker-order/internal/core"
	"github.com/go-worker-order/internal/lib"

	"github.com/go-worker-order/internal/repository/pg"

	"github.com/rs/zerolog/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var childLogger = log.With().Str("repository", "storage").Logger()

//-----------------------------------------------
type WorkerRepository struct {
	databasePG pg.DatabasePG
}

func NewWorkerRepository(databasePG pg.DatabasePG) WorkerRepository {
	childLogger.Debug().Msg("NewWorkerRepository")
	return WorkerRepository{
		databasePG: databasePG,
	}
}

func (w WorkerRepository) StartTx(ctx context.Context) (pgx.Tx, *pgxpool.Conn,error) {
	childLogger.Debug().Msg("StartTx")

	span := lib.Span(ctx, "storage.StartTx")
	defer span.End()

	conn, err := w.databasePG.Acquire(ctx)
	if err != nil {
		childLogger.Error().Err(err).Msg("Erro Acquire")
		return nil, nil, errors.New(err.Error())
	}
	
	tx, err := conn.Begin(ctx)
    if err != nil {
        return nil, nil ,errors.New(err.Error())
    }

	return tx, conn, nil
}

func (w WorkerRepository) ReleaseTx(connection *pgxpool.Conn) {
	childLogger.Debug().Msg("ReleaseTx")

	defer connection.Release()
}

func (w WorkerRepository) Update(ctx context.Context, tx pgx.Tx, order *core.Order) (int64, error){
	childLogger.Debug().Msg("Update")
	childLogger.Debug().Interface("order : ", order).Msg("")

	span := lib.Span(ctx, "storage.Update")	
    defer span.End()

	query := `Update public.order
				set status = $1,
					update_at = $2
				where id = $3`

	childLogger.Debug().Interface("====> order.ID : ", order.ID).Msg("")
	row, err := tx.Exec(ctx, 
						query, 
						order.Status,
						time.Now(),
						order.ID)
	if err != nil {
		childLogger.Error().Err(err).Msg("Exec statement")
		return 0, errors.New(err.Error())
	}

	childLogger.Debug().Interface("rowsAffected : ", row.RowsAffected()).Msg("")

	return int64(row.RowsAffected()) , nil
}