package utils

import (
	"context"
	"fmt"
	"time"

	mymetrics "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/metrics"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	logging "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/utils/log"
)

func scheduledUpdate(ctx context.Context, tx pgx.Tx, metrics *mymetrics.DatabaseMetrics) error {
	funcName := logging.GetOnlyFunctionName()
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	SQLCloseAdvert := `UPDATE public.advert
	SET   is_promoted=false, promotion_start=null, promotion_duration=null
	WHERE promotion_start + promotion_duration < now();  `

	logging.LogInfo(logger, "UPDATE advert")

	var err error

	start := time.Now()

	_, err = tx.Exec(ctx, SQLCloseAdvert)

	metrics.AddDuration(funcName, time.Since(start))

	if err != nil {
		logging.LogError(logger, fmt.Errorf("something went wrong update advert promotion status, err=%w",
			err))
		metrics.IncreaseErrors(funcName)

		return err
	}

	return nil
}

func ScheduledUpdate(ctx context.Context, pool *pgxpool.Pool, metrics *mymetrics.DatabaseMetrics) error {
	logger := logging.GetLoggerFromContext(ctx).With(zap.String("func", logging.GetFunctionName()))

	err := pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
		err := scheduledUpdate(ctx, tx, metrics)
		if err != nil {
			logging.LogError(logger, fmt.Errorf("something went wrong update advert promotion status, err=%w",
				err))

			return err
		}

		return nil
	})

	if err != nil {
		logging.LogError(logger, fmt.Errorf("error while updating advert promotion status, err=%w", err))

		return err
	}

	return nil
}
