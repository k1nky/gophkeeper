package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/k1nky/gophkeeper/internal/adapter/gophkeeper"
	"github.com/k1nky/gophkeeper/internal/crypto"
	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
	"github.com/k1nky/gophkeeper/internal/service/keeper"
)

type Context struct {
	Debug  bool
	keeper *keeper.Service
	ctx    context.Context
	client *gophkeeper.Adapter
}

type LsCmd struct {
	Remote bool
}

type PutCmd struct {
	File  string `arg:"" optional:"" name:"file" help:"Path to secret file." type:"path"`
	Line  string
	Alias string
}

type ShCmd struct {
	Id    string
	Alias string
}

type PushCmd struct {
	Id    string
	Alias string
}

type PullCmd struct {
	Id string
}

func (c *PushCmd) Run(ctx *Context) error {
	meta, err := ctx.keeper.GetSecretMeta(ctx.ctx, vault.MetaID(c.Id))
	if err != nil {
		return err
	}
	data, err := ctx.keeper.GetSecretData(ctx.ctx, vault.MetaID(c.Id))
	if err != nil {
		return err
	}
	defer data.Close()
	m, err := ctx.client.PutSecret(ctx.ctx, *meta, data)
	fmt.Println(m, err)
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
		list, err = ctx.keeper.ListSecretsByUser(ctx.ctx, user.LocalUserID)
	}
	fmt.Println(list.String())
	return err
}

func (c *PutCmd) Run(ctx *Context) error {
	line := vault.NewBytesBuffer([]byte(c.Line))
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
	r, w := io.Pipe()
	data := vault.NewDataReader(r)
	go func() {
		// TODO: error handling
		err = ctx.client.GetSecretData(ctx.ctx, vault.MetaID(c.Id), w)
		w.Close()
	}()
	if err != nil {
		return err
	}
	newMeta, err := ctx.keeper.PutSecret(ctx.ctx, *meta, data)
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

var CLI struct {
	Debug bool    `help:"Enable debug mode."`
	Ls    LsCmd   `cmd:"" help:"List secrects."`
	Put   PutCmd  `cmd:"" help:"Put secrect."`
	Push  PushCmd `cmd:"" help:"Push secrect."`
	Sh    ShCmd   `cmd:"" help:"Show secrect."`
	Pull  PullCmd `cmd:"" help:"Pull secrect."`
}
