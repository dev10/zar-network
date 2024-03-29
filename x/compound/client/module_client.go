package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
	cdpcmd "github.com/zar-network/zar-network/x/compound/client/cli"
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
	cdpQueryCmd := &cobra.Command{
		Use:   "cdp",
		Short: "Querying commands for the cdp module",
	}

	cdpQueryCmd.AddCommand(client.GetCommands(
		cdpcmd.GetCmd_GetCdp(mc.storeKey, mc.cdc),
		cdpcmd.GetCmd_GetCdps(mc.storeKey, mc.cdc),
		cdpcmd.GetCmd_GetUnderCollateralizedCdps(mc.storeKey, mc.cdc),
		cdpcmd.GetCmd_GetParams(mc.storeKey, mc.cdc),
	)...)

	return cdpQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	cdpTxCmd := &cobra.Command{
		Use:   "cdp",
		Short: "cdp transactions subcommands",
	}

	cdpTxCmd.AddCommand(client.PostCommands(
		cdpcmd.GetCmdModifyCdp(mc.cdc),
	)...)

	return cdpTxCmd
}
