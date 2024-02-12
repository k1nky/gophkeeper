package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/k1nky/gophkeeper/internal/adapter/gophkeeper"
	"github.com/k1nky/gophkeeper/internal/adapter/store"
	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/logger"
	"github.com/k1nky/gophkeeper/internal/service/keeper"
	"github.com/k1nky/gophkeeper/internal/service/sync"
	"github.com/k1nky/gophkeeper/internal/store/meta/bolt"
	"github.com/k1nky/gophkeeper/internal/store/objects/filestore"
)

func newClient(ctx context.Context, url string, u user.User, token string, l *logger.Logger) (*gophkeeper.Adapter, error) {
	if len(url) == 0 {
		return nil, nil
	}
	client := gophkeeper.New(url, "")
	if len(token) > 0 {
		client.SetToken(token)
	} else {
		if cur, err := client.Login(ctx, u.Login, u.Password); err != nil {
			return nil, err
		} else {
			l.Debugf("log on as %s", cur.Login)
		}
	}
	if err := client.Open(ctx); err != nil {
		return nil, err
	}
	return client, nil
}

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	log := logger.New()

	cmd := kong.Parse(&cli)
	if cli.Debug {
		log.SetLevel("debug")
	}
	client, err := newClient(ctx, string(cli.RemoteVault), user.User{
		Login:    cli.User,
		Password: cli.Password,
	}, cli.Token, log)
	if err != nil {
		log.Errorf("connect to %s: %v", cli.RemoteVault, err)
		os.Exit(1)
	}
	store := store.New(bolt.New(cli.MetaStoreDSN), filestore.New(cli.ObjectStoreDSN))
	if err := store.Open(ctx); err != nil {
		log.Errorf("store: %v", err)
		os.Exit(1)
	}
	defer store.Close()
	keeper := keeper.New(store, log)
	sync := sync.New(client, keeper, log)

	if err = cmd.Run(&Context{
		keeper: keeper,
		ctx:    ctx,
		client: client,
		sync:   sync,
		log:    log,
		secret: cli.Secret,
	}); err != nil {
		log.Errorf("command: %s", err)
	}
}
