package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/zar-network/zar-network/x/auction/internal/keeper"
)

// GetCmdGetAuctions queries the auctions in the store
func GetCmdGetAuctions(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "getauctions",
		Short: "get a list of active auctions",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/getauctions", queryRoute), nil)
			if err != nil {
				fmt.Printf("error when getting auctions - %s", err)
				return nil
			}
			var out keeper.QueryResAuctions
			cdc.MustUnmarshalJSON(res, &out)
			if len(out) == 0 {
				out = append(out, "There are currently no auctions")
			}
			return cliCtx.PrintOutput(out)
		},
	}
}
