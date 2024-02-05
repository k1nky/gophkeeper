package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/k1nky/gophkeeper/internal/adapter/gophkeeper"
	"github.com/k1nky/gophkeeper/internal/crypto"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
	"github.com/k1nky/gophkeeper/internal/logger"
	"github.com/k1nky/gophkeeper/internal/service/keeper"
	"github.com/k1nky/gophkeeper/internal/service/sync"
)

type Context struct {
	keeper *keeper.Service
	ctx    context.Context
	sync   *sync.Service
	client *gophkeeper.Adapter
	log    *logger.Logger
	secret string
}

type LsCmd struct {
	Remote bool `optional:"" name:"remote" help:"List secrets from remote storage."`
}

type PutCmd struct {
	Type  string `required:"" name:"type" enum:"text,file,login,card" default:"text"`
	Value string `arg:"" name:"value" help:""`
	Alias string `optional:"" name:"alias" help:"Secret entry alias."`
}

type ShCmd struct {
	Id    string `optional:"" name:"id" help:"Secret entry ID to show."`
	Alias string `optional:"" name:"alias" help:"Secret entry alias to show."`
}

type PushCmd struct {
	Id    string `optional:"" name:"id" help:"Secret entry ID to push."`
	Alias string `optional:"" name:"alias" help:"Secret entry alias to push."`
	All   bool   `optional:"" name:"all" help:"Push all secrets from local storage."`
}

type PullCmd struct {
	Id  string `optional:"" name:"id" help:"Secret entry ID to pull."`
	All bool   `optional:"" name:"all" help:"Pull all secrets from remote storage."`
}

type remoteVaultFlag string

// TODO: delete secret
var cli struct {
	Debug          bool            `optional:"" name:"debug" env:"DEBUG" help:"Enable debug mode."`
	RemoteVault    remoteVaultFlag `optional:"" name:"remote-vault" env:"REMOTE_VAULT"`
	MetaStoreDSN   string          `optional:"" name:"meta-store-dsn" env:"META_STORE_DSN" default:"/tmp/client-meta.db"`
	ObjectStoreDSN string          `optional:"" name:"object-store-dsn" env:"OBJECT_STORE_DSN" default:"/tmp/client-vault"`
	User           string          `optional:"" name:"remote-user" env:"REMOTE_VAULT_USER"`
	Password       string          `optional:"" name:"remote-password" env:"REMOTE_VAULT_PASSWORD"`
	Secret         string          `optional:"" name:"secret" env:"VAULT_SECRET" default:"secret"`
	Ls             LsCmd           `cmd:"" help:"List secrects from local or remote storage."`
	Put            PutCmd          `cmd:"" help:"Put secrect to local storage."`
	Push           PushCmd         `cmd:"" help:"Push secrect to remote storage."`
	Sh             ShCmd           `cmd:"" help:"Show secrect from local storage."`
	Pull           PullCmd         `cmd:"" help:"Pull secrect from remote storage."`
}

func (c *PushCmd) Run(ctx *Context) error {
	if c.All {
		list, err := ctx.sync.PullAll(ctx.ctx)
		fmt.Printf("PULLED:\n %s", list)
		return err
	}
	meta, err := ctx.keeper.GetSecretMeta(ctx.ctx, vault.MetaID(c.Id))
	if err != nil {
		return err
	}
	_, err = ctx.sync.Push(ctx.ctx, *meta)
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
	var value io.ReadCloser
	m := vault.Meta{
		ID:    vault.NewMetaID(),
		Alias: c.Alias,
	}
	switch c.Type {
	case "text":
		value = vault.NewBytesBuffer([]byte(c.Value))
		m.Type = vault.TypeText
	case "login":
		lp := &vault.LoginPassword{}
		lp.Prompt()
		b, err := lp.Bytes()
		if err != nil {
			return err
		}
		m.Type = vault.TypeLoginPassword
		value = vault.NewBytesBuffer(b)
	case "card":
		cc := &vault.CreditCard{}
		cc.Prompt()
		b, err := cc.Bytes()
		if err != nil {
			return err
		}
		m.Type = vault.TypeCreditCard
		value = vault.NewBytesBuffer(b)
	case "file":
		if f, err := os.OpenFile(c.Value, os.O_RDONLY, 0660); err != nil {
			return err
		} else {
			defer f.Close()
			value = f
		}
		m.Type = vault.TypeFile
	}
	// TODO: вектор инициализации можно хранить в мета-данных
	enc, _ := crypto.NewEncryptReader(ctx.secret, value, nil)
	data := vault.NewDataReader(enc)
	meta, err := ctx.keeper.PutSecret(ctx.ctx, m, data)
	fmt.Println(meta.String())
	return err
}

func (c *PullCmd) Run(ctx *Context) error {
	if c.All {
		_, err := ctx.sync.PullAll(ctx.ctx)
		return err
	}
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
	dec, _ := crypto.NewDecryptReader(ctx.secret, data, nil)
	defer data.Close()
	_, err = io.Copy(os.Stdout, dec)
	return err
}
