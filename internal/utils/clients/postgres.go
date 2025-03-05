package clients

import (
	"context"
	apperrors "github.com/DavidMovas/SpeakUp-Server/internal/utils/error"
	"github.com/jackc/pgx/v5/pgxpool"
	"net"
	"time"
)

func NewPostgresClient(ctx context.Context, url string, options *PostgresOptions) (*pgxpool.Pool, error) {
	var opts *PostgresOptions

	if options == nil {
		opts = NewPostgresOptions(url)
	} else {
		opts = options
	}

	pingCtx, pingCtxCancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer pingCtxCancel()

	client, err := pgxpool.NewWithConfig(ctx, opts.Config)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	if err = client.Ping(pingCtx); err != nil {
		return nil, apperrors.Internal(err)
	}

	return client, nil
}

type PostgresOptions struct {
	*pgxpool.Config
}

func NewPostgresOptions(url string) *PostgresOptions {
	opts := &PostgresOptions{}
	opts.Config, _ = pgxpool.ParseConfig(url)

	return opts
}

func (o *PostgresOptions) WithHost(host string) *PostgresOptions {
	o.Config.ConnConfig.Host = host
	return o
}

func (o *PostgresOptions) WithUsername(username string) *PostgresOptions {
	o.Config.ConnConfig.User = username
	return o
}

func (o *PostgresOptions) WithPassword(password string) *PostgresOptions {
	o.Config.ConnConfig.Password = password
	return o
}

func (o *PostgresOptions) WithDatabase(database string) *PostgresOptions {
	o.Config.ConnConfig.Database = database
	return o
}

func (o *PostgresOptions) WithDialFunc(fn func(ctx context.Context, network, addr string) (net.Conn, error)) *PostgresOptions {
	o.Config.ConnConfig.DialFunc = fn
	return o
}

func (o *PostgresOptions) WithConnectTimeout(time time.Duration) *PostgresOptions {
	o.Config.ConnConfig.ConnectTimeout = time
	return o
}

func (o *PostgresOptions) WithHealthCheckPeriod(time time.Duration) *PostgresOptions {
	o.Config.HealthCheckPeriod = time
	return o
}

func (o *PostgresOptions) WithMinCons(amount int32) *PostgresOptions {
	o.Config.MinConns = amount
	return o
}

func (o *PostgresOptions) WithMaxCons(amount int32) *PostgresOptions {
	o.Config.MaxConns = amount
	return o
}

func (o *PostgresOptions) WithConnMaxLifetime(time time.Duration) *PostgresOptions {
	o.Config.MaxConnLifetime = time
	return o
}
