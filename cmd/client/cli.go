package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/k1nky/gophkeeper/internal/adapter/gophkeeper"
	"github.com/k1nky/gophkeeper/internal/crypto"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
	"github.com/k1nky/gophkeeper/internal/service/keeper"
	"github.com/k1nky/gophkeeper/internal/service/sync"
)

type Context struct {
	Debug       bool   `optional:"" name:"debug" env:"DEBUG"`
	RemoteVault string `optional:"" name:"remote-vault" env:"REMOTE_VAULT"`
	LocalVault  string `optional:"" name:"local-vault" env:"LOCAL_VAULT"`
	User        string `optional:"" name:"user" env:"user"`
	Password    string `optional:"" name:"password" env:"password"`
	keeper      *keeper.Service
	ctx         context.Context
	sync        *sync.Service
	client      *gophkeeper.Adapter
}

type LsCmd struct {
	Remote bool `optional:"" name:"remote" help:"List secrets from remote storage."`
}

type PutCmd struct {
	File  string `arg:"" optional:"" name:"file" help:"Path to file with the secret to be placed in local storage." type:"path"`
	Text  string `arg:"" optional:"" name:"text" help:"Secret text to save in local storage."`
	Alias string `arg:"" optinal:"" name:"alias" help:"Secret entry alias."`
}

type ShCmd struct {
	Id    string `arg:"" optinal:"" name:"alias" help:"Secret entry ID to show."`
	Alias string `arg:"" optinal:"" name:"alias" help:"Secret entry alias to show."`
}

type PushCmd struct {
	Id    string `arg:"" optinal:"" name:"alias" help:"Secret entry ID to push."`
	Alias string `arg:"" optinal:"" name:"alias" help:"Secret entry alias to push."`
}

type PullCmd struct {
	Id string `arg:"" optinal:"" name:"alias" help:"Secret entry ID to pull."`
}

var CLI struct {
	Debug bool    `help:"Enable debug mode."`
	Ls    LsCmd   `cmd:"" help:"List secrects from local or remote storage."`
	Put   PutCmd  `cmd:"" help:"Put secrect to local storage."`
	Push  PushCmd `cmd:"" help:"Push secrect to remote storage."`
	Sh    ShCmd   `cmd:"" help:"Show secrect from local storage."`
	Pull  PullCmd `cmd:"" help:"Pull secrect from remote storage."`
}

func (c *PushCmd) Run(ctx *Context) error {
	meta, err := ctx.keeper.GetSecretMeta(ctx.ctx, vault.MetaID(c.Id))
	if err != nil {
		return err
	}
	err = ctx.sync.Push(ctx.ctx, *meta)
	return err
}

func (c *LsCmd) Run(ctx *Context) error {
	var (
		list vault.List
		err  error
	)
	if c.Remote {
		list, err = ctx.client.ListSecrets(ctx.ctx)
	} else {
		list, err = ctx.keeper.ListSecretsByUser(ctx.ctx)
	}
	fmt.Println(list.String())
	return err
}

func (c *PutCmd) Run(ctx *Context) error {
	line := vault.NewBytesBuffer([]byte(c.Text))
	// TODO: вектор инициализации можно хранить в мета-данных
	enc, _ := crypto.NewEncryptReader("secret", line, nil)
	data := vault.NewDataReader(enc)
	meta, err := ctx.keeper.PutSecret(ctx.ctx, vault.Meta{
		ID:    vault.NewMetaID(),
		Alias: c.Alias,
	}, data)
	fmt.Println(meta.String())
	return err
}

func (c *PullCmd) Run(ctx *Context) error {
	meta, err := ctx.client.GetSecretMeta(ctx.ctx, vault.MetaID(c.Id))
	if err != nil {
		return err
	}
	newMeta, err := ctx.sync.Pull(ctx.ctx, *meta)
	fmt.Println(newMeta.String())
	return err
}

func (c *ShCmd) Run(ctx *Context) error {
	var (
		meta *vault.Meta
		err  error
	)
	if len(c.Id) != 0 {
		meta, err = ctx.keeper.GetSecretMeta(ctx.ctx, vault.MetaID(c.Id))
	} else {
		meta, err = ctx.keeper.GetSecretMetaByAlias(ctx.ctx, c.Alias)
	}
	if err != nil {
		return err
	}
	fmt.Println(meta)
	data, err := ctx.keeper.GetSecretData(ctx.ctx, vault.MetaID(meta.ID))
	if err != nil {
		return err
	}
	dec, _ := crypto.NewDecryptReader("secret", data, nil)
	defer data.Close()
	_, err = io.Copy(os.Stdout, dec)
	return err
}
