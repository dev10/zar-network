package cli

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zar-network/zar-network/x/issue/internal/types"
)

// QueryParamsCmd implements the query params command.
func QueryParamsCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "params",
		Short:   "Query the parameters of the lock process",
		Long:    "Query all the parameters",
		Example: "$ zarcli lock params",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(types.GetQueryParamsPath(), nil)
			if err != nil {
				return err
			}
			var params types.Params
			cdc.MustUnmarshalJSON(res, &params)
			return cliCtx.PrintOutput(params)
		},
	}
}

// QueryCmd implements the query issue command.
func QueryCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "issue [denom]",
		Args:    cobra.ExactArgs(1),
		Short:   "Query the details of the account coin",
		Long:    "Query the details of the account issue coin",
		Example: "$ zar-networkcli bank issue coin174876e800",
		RunE: func(cmd *cobra.Command, args []string) error {
			return processQuery(cdc, args)
		},
	}
}

// QueryIssueCmd implements the query issue command.
func QueryIssueCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "query [issue-id]",
		Args:    cobra.ExactArgs(1),
		Short:   "Query a single issue",
		Long:    "Query details for a issue. You can find the issue-id by running zar-networkcli issue list-issues",
		Example: "$ zar-networkcli issue query-issue coin174876e800",
		RunE: func(cmd *cobra.Command, args []string) error {
			return processQuery(cdc, args)
		},
	}
}

func processQuery(cdc *codec.Codec, args []string) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	issueID := args[0]
	if err := types.CheckIssueId(issueID); err != nil {
		return types.Errorf(err)
	}
	// Query the issue
	res, _, err := cliCtx.QueryWithData(types.GetQueryIssuePath(issueID), nil)
	if err != nil {
		return err
	}
	var issueInfo types.Issue
	cdc.MustUnmarshalJSON(res, &issueInfo)
	return cliCtx.PrintOutput(issueInfo)
}

// QueryAllowanceCmd implements the query allowance command.
func QueryAllowanceCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "query-allowance [issue-id] [owner-address] [spender-address]",
		Args:    cobra.ExactArgs(3),
		Short:   "Query allowance",
		Long:    "Query the amount of tokens that an owner allowed to a spender",
		Example: "$ zar-networkcli issue query-allowance coin174876e800 gard1zu85q8a7wev675k527y7keyrea7wu7crr9vdrs gard1vud9ptwagudgq7yht53cwuf8qfmgkd0qcej0ah",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			issueID := args[0]
			if err := types.CheckIssueId(issueID); err != nil {
				return types.Errorf(err)
			}
			ownerAddress, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			spenderAddress, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(types.GetQueryIssueAllowancePath(issueID, ownerAddress, spenderAddress), nil)
			if err != nil {
				return err
			}
			var approval types.Approval
			cdc.MustUnmarshalJSON(res, &approval)

			return cliCtx.PrintOutput(approval)
		},
	}
}

// QueryFreezeCmd implements the query freeze command.
func QueryFreezeCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:     "query-freeze [issue-id] [acc-address]",
		Args:    cobra.ExactArgs(2),
		Short:   "Query freeze",
		Long:    "Query freeze the transfer from a address",
		Example: "$ zar-networkcli issue query-freeze coin174876e800 gard15l5yzrq3ff8fl358ng430cc32lzkvxc30n405n",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			issueID := args[0]
			if err := types.CheckIssueId(issueID); err != nil {
				return types.Errorf(err)
			}
			accAddress, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(types.GetQueryIssueFreezePath(issueID, accAddress), nil)
			if err != nil {
				return err
			}
			var freeze types.IssueFreeze
			cdc.MustUnmarshalJSON(res, &freeze)

			return cliCtx.PrintOutput(freeze)
		},
	}
}

// QueryIssuesCmd implements the query issue command.
func QueryIssuesCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Query issue list",
		Long:    "Query all or one of the account issue list, the limit default is 30",
		Example: "$ zar-networkcli issue list-issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			address, err := sdk.AccAddressFromBech32(viper.GetString(flagAddress))
			if err != nil {
				return err
			}
			issueQueryParams := types.IssueQueryParams{
				StartIssueId: viper.GetString(flagStartIssueId),
				Owner:        address,
				Limit:        viper.GetInt(flagLimit),
			}
			// Query the issue
			bz, err := cliCtx.Codec.MarshalJSON(issueQueryParams)
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(types.GetQueryIssuesPath(), bz)
			if err != nil {
				return err
			}

			var issues types.CoinIssues
			cdc.MustUnmarshalJSON(res, &issues)
			return cliCtx.PrintOutput(issues)
		},
	}

	cmd.Flags().String(flagAddress, "", "Token owner address")
	cmd.Flags().String(flagStartIssueId, "", "Start issueId of issues")
	cmd.Flags().Int32(flagLimit, 30, "Query number of issue results per page returned")
	return cmd
}

// QueryFreezesCmd implements the query freezes command.
func QueryFreezesCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list-freeze",
		Short:   "Query freeze list",
		Long:    "Query all or one of the issue freeze list",
		Example: "$ zar-networkcli issue list-freeze",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			issueID := args[0]
			if err := types.CheckIssueId(issueID); err != nil {
				return types.Errorf(err)
			}
			res, _, err := cliCtx.QueryWithData(types.GetQueryIssueFreezesPath(issueID), nil)
			if err != nil {
				return err
			}
			var issueFreeze types.IssueAddressFreezeList
			cdc.MustUnmarshalJSON(res, &issueFreeze)
			return cliCtx.PrintOutput(issueFreeze)
		},
	}
	return cmd
}

// QuerySearchIssuesCmd implements the query issue command.
func QuerySearchIssuesCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "search [symbol]",
		Args:    cobra.ExactArgs(1),
		Short:   "Search issues",
		Long:    "Search issues based on symbol",
		Example: "$ zar-networkcli issue search fo",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			// Query the issue
			res, _, err := cliCtx.QueryWithData(types.GetQueryIssueSearchPath(strings.ToUpper(args[0])), nil)
			if err != nil {
				return err
			}
			var issues types.CoinIssues
			cdc.MustUnmarshalJSON(res, &issues)
			return cliCtx.PrintOutput(issues)
		},
	}
	return cmd
}
