package main

var cli struct {
	Debug          bool   `optional:"" name:"debug" env:"DEBUG" help:"Enable debug mode."`
	MetaStoreDSN   string `optional:"" name:"meta-store-dsn" env:"META_STORE_DSN" default:"/tmp/server-meta.db"`
	ObjectStoreDSN string `optional:"" name:"object-store-dsn" env:"OBJECT_STORE_DSN" default:"/tmp/server-vault"`
	Secret         string `optional:"" name:"secret" env:"SECRET"`
	Listen         string `optional:"" name:"listen" env:"LISTEN" default:":8080"`
}
