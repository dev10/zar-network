package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
	pricefeedcmd "github.com/zar-network/zar-network/x/pricefeed/client/cli"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

// NewModuleClient creates client for the module
func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	// Group nameservice queries under a subcommand
	pricefeedQueryCmd := &cobra.Command{
		Use:   "pricefeed",
		Short: "Querying commands for the pricefeed module",
	}

	pricefeedQueryCmd.AddCommand(client.GetCommands(
		pricefeedcmd.GetCmdCurrentPrice(mc.storeKey, mc.cdc),
		pricefeedcmd.GetCmdRawPrices(mc.storeKey, mc.cdc),
		pricefeedcmd.GetCmdAssets(mc.storeKey, mc.cdc),
	)...)

	return pricefeedQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	pricefeedTxCmd := &cobra.Command{
		Use:   "pricefeed",
		Short: "Pricefeed transactions subcommands",
	}

	pricefeedTxCmd.AddCommand(client.PostCommands(
		pricefeedcmd.GetCmdPostPrice(mc.cdc),
	)...)

	return pricefeedTxCmd
}
