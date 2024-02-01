package main

import (
	"context"
	"fmt"
	"io"
	"os"

	_ "github.com/alecthomas/kong"
	"github.com/k1nky/gophkeeper/internal/crypto"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
	"github.com/k1nky/gophkeeper/internal/service/keeper"
)

type Context struct {
	Debug  bool
	keeper *keeper.Service
	ctx    context.Context
}

type LsCmd struct{}

type PutCmd struct {
	File string `arg:"" optional:"" name:"file" help:"Path to secret file." type:"path"`
	Line string
}

type ShCmd struct {
	Id string
}

func (c *LsCmd) Run(ctx *Context) error {
	list, err := ctx.keeper.ListSecretsByUser(ctx.ctx, 0)
	fmt.Println(list)
	return err
}

func (c *PutCmd) Run(ctx *Context) error {
	line := vault.NewBytesBuffer([]byte(c.Line))
	enc, _ := crypto.NewEncryptReader("secret", line)
	data := vault.NewDataReader(enc)
	meta, err := ctx.keeper.PutSecret(ctx.ctx, vault.Meta{}, data)
	fmt.Println(meta)
	return err
}

func (c *ShCmd) Run(ctx *Context) error {
	meta, err := ctx.keeper.GetSecretMeta(ctx.ctx, vault.UniqueKey(c.Id))
	if err != nil {
		return err
	}
	fmt.Println(meta)
	data, err := ctx.keeper.GetSecretData(ctx.ctx, vault.UniqueKey(c.Id))
	if err != nil {
		return err
	}
	dec, _ := crypto.NewDecryptReader("secret", data)
	defer data.Close()
	_, err = io.Copy(os.Stdout, dec)
	return err
}

var CLI struct {
	Debug bool   `help:"Enable debug mode."`
	Ls    LsCmd  `cmd:"" help:"List secrects."`
	Put   PutCmd `cmd:"" help:"Put secrect."`
	Sh    ShCmd  `cmd:"" help:"Show secrect."`
}
