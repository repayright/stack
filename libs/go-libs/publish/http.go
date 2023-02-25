package publish

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	wHttp "github.com/ThreeDotsLabs/watermill-http/pkg/http"
	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/fx"
)

func newHTTPPublisher(logger watermill.LoggerAdapter, config wHttp.PublisherConfig) (*wHttp.Publisher, error) {
	return wHttp.NewPublisher(config, logger)
}

func newHTTPSubscriber(logger watermill.LoggerAdapter, addr string, config wHttp.SubscriberConfig) (*wHttp.Subscriber, error) {
	return wHttp.NewSubscriber(addr, config, logger)
}

func newHTTPPublisherConfig(httpClient *http.Client, m wHttp.MarshalMessageFunc) wHttp.PublisherConfig {
	return wHttp.PublisherConfig{
		MarshalMessageFunc: m,
		Client:             httpClient,
	}
}

func newHTTPSubscriberConfig(m wHttp.UnmarshalMessageFunc) wHttp.SubscriberConfig {
	return wHttp.SubscriberConfig{
		Router:               nil,
		UnmarshalMessageFunc: m,
	}
}

func defaultHTTPMarshalMessageFunc() wHttp.MarshalMessageFunc {
	return func(url string, msg *message.Message) (*http.Request, error) {
		req, err := wHttp.DefaultMarshalMessageFunc(url, msg)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		return req, nil
	}
}

func defaultHTTPUnmarshalMessageFunc() wHttp.UnmarshalMessageFunc {
	return wHttp.DefaultUnmarshalMessageFunc
}

func httpModule(addr string) fx.Option {
	options := []fx.Option{
		fx.Provide(newHTTPPublisher),
		fx.Provide(newHTTPPublisherConfig),
		fx.Provide(defaultHTTPMarshalMessageFunc),
		fx.Supply(http.DefaultClient),
		fx.Provide(func(p *wHttp.Publisher) message.Publisher {
			return p
		}),
	}
	if addr != "" {
		options = append(options,
			fx.Provide(newHTTPSubscriberConfig),
			fx.Provide(func(logger watermill.LoggerAdapter, config wHttp.SubscriberConfig, lc fx.Lifecycle) (*wHttp.Subscriber, error) {
				ret, err := newHTTPSubscriber(logger, addr, config)
				if err != nil {
					return nil, err
				}
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						_, port, err := net.SplitHostPort(addr)
						if err != nil {
							return err
						}
						go func() {
							if err := ret.StartHTTPServer(); err != nil && err != http.ErrServerClosed {
								panic(err)
							}
						}()
						for {
							select {
							case <-ctx.Done():
								return ctx.Err()
							case <-time.After(100 * time.Millisecond):
								_, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s", port))
								if err == nil {
									return nil
								}
							}
						}
					},
					OnStop: func(ctx context.Context) error {
						return ret.Close()
					},
				})
				return ret, nil
			}),
			fx.Provide(defaultHTTPUnmarshalMessageFunc),
			fx.Provide(func(p *wHttp.Subscriber) message.Subscriber {
				return p
			}),
		)
	}
	return fx.Options(options...)
}
