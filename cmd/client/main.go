package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/k1nky/gophkeeper/internal/adapter/gophkeeper"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	// store := store.New(bolt.New("meta.db"), filestore.New("/tmp/ostore2"))
	// keeper := keeper.New(store, &logger.Blackhole{})
	client := gophkeeper.New("http://localhost:8080", "/")
	client.Open(ctx)
	claims, err := client.Login(ctx, "u", "p")
	fmt.Println(claims, err)
	list, err := client.ListSecrets(ctx)
	fmt.Println(list, err)
	data := vault.NewBytesBuffer([]byte("Hit the lights"))
	meta, err := client.PutSecret(ctx, vault.Meta{Extra: "first test secret"}, data)
	fmt.Println(meta, err)

	<-ctx.Done()
}
