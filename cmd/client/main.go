package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/k1nky/gophkeeper/internal/adapter/store"
	"github.com/k1nky/gophkeeper/internal/logger"
	"github.com/k1nky/gophkeeper/internal/service/keeper"
	"github.com/k1nky/gophkeeper/internal/store/meta/bolt"
	"github.com/k1nky/gophkeeper/internal/store/objects/filestore"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	store := store.New(bolt.New("meta.db"), filestore.New("/tmp/ostore2"))
	err := store.Open(ctx)
	fmt.Println(err)
	defer store.Close()
	keeper := keeper.New(store, &logger.Blackhole{})

	cmd := kong.Parse(&CLI)
	err = cmd.Run(&Context{
		Debug:  CLI.Debug,
		keeper: keeper,
		ctx:    ctx,
	})
	fmt.Println(err)

	// store := store.New(bolt.New("meta.db"), filestore.New("/tmp/ostore2"))
	// keeper := keeper.New(store, &logger.Blackhole{})
	// client := gophkeeper.New("http://localhost:8080", "/")
	// client.Open(ctx)
	// claims, err := client.Login(ctx, "u", "p")
	// fmt.Println(claims, err)
	// list, err := client.ListSecrets(ctx)
	// fmt.Println(list, err)
	// data := vault.NewBytesBuffer([]byte("Hit the lights"))
	// meta, err := client.PutSecret(ctx, vault.Meta{Extra: "first test secret"}, data)
	// fmt.Println(meta, err)

	// <-ctx.Done()
}
