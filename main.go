package main

import (
	"log"
	"os"

	hasApi "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
	"github.com/joho/godotenv"
	api "github.com/ryadavDeqode/dq-vault/api"
)

func main() {
	apiClientMeta := &hasApi.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:]) // Ignore command, strictly parse flags
	err := godotenv.Load()

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := hasApi.VaultPluginTLSProvider(tlsConfig)

	err = plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: api.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
