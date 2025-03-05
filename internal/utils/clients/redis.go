package clients

import (
	"context"
	"crypto/tls"
	apperrors "github.com/DavidMovas/SpeakUp-Server/internal/utils/error"
	"github.com/redis/go-redis/v9"
	"net"
	"time"
)

var defaultOptions = &redis.Options{
	Addr:            "127.0.0.1:6379",
	MaxRetries:      5,
	MinRetryBackoff: 10 * time.Millisecond,
	MaxRetryBackoff: time.Second,
	DialTimeout:     10 * time.Second,
}

func NewRedisClient(url string, options *RedisOptions) (*redis.Client, error) {
	var opts *RedisOptions

	if options == nil {
		opts = NewRedisOptions(url)
	} else {
		opts = options
	}

	pingCtx, pingCtxCancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer pingCtxCancel()

	client := redis.NewClient(opts.Options)
	if err := client.Ping(pingCtx).Err(); err != nil {
		return nil, apperrors.Internal(err)
	}

	return client, nil
}

type RedisOptions struct {
	*redis.Options
}

func NewRedisOptions(url string) *RedisOptions {
	defaultOptions.Addr = url
	return &RedisOptions{
		Options: defaultOptions,
	}
}

func (o *RedisOptions) WithAddr(addr string) *RedisOptions {
	o.Addr = addr
	return o
}

func (o *RedisOptions) WithUsername(username string) *RedisOptions {
	o.Username = username
	return o
}

func (o *RedisOptions) WithPassword(pw string) *RedisOptions {
	o.Password = pw
	return o
}

func (o *RedisOptions) WithDealer(fn func(ctx context.Context, network, addr string) (net.Conn, error)) *RedisOptions {
	o.Dialer = fn
	return o
}

func (o *RedisOptions) WithOnConnect(fn func(ctx context.Context, conn *redis.Conn) error) *RedisOptions {
	o.OnConnect = fn
	return o
}

func (o *RedisOptions) WithMaxRetries(maxRetries int) *RedisOptions {
	o.MaxRetries = maxRetries
	return o
}

func (o *RedisOptions) WithMinRetryBackoff(time time.Duration) *RedisOptions {
	o.MinRetryBackoff = time
	return o
}

func (o *RedisOptions) WithMaxRetryBackoff(time time.Duration) *RedisOptions {
	o.MaxRetryBackoff = time
	return o
}

func (o *RedisOptions) WithDialTimeout(timeout time.Duration) *RedisOptions {
	o.DialTimeout = timeout
	return o
}

func (o *RedisOptions) WithLimiter(limiter redis.Limiter) *RedisOptions {
	o.Limiter = limiter
	return o
}

func (o *RedisOptions) WithTLSConfig(tlsConfig *tls.Config) *RedisOptions {
	o.TLSConfig = tlsConfig
	return o
}
