package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"

	"github.com/zar-network/zar-network/x/liquidator/internal/types"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/liquidator/outstandingdebt", queryDebtHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/liquidator/seize", seizeCdpHandlerFn(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/liquidator/mint", debtAuctionHandlerFn(cdc, cliCtx)).Methods("POST")
	// r.HandleFunc("liquidator/burn", surplusAuctionHandlerFn(cdc, cliCtx).Methods("POST"))
}

func queryDebtHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/liquidator/%s", types.QueryGetOutstandingDebt), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res) // write JSON to response writer
	}
}

type SeizeAndStartCollateralAuctionRequest struct {
	BaseReq         rest.BaseReq   `json:"base_req"`
	Sender          sdk.AccAddress `json:"sender"`
	CdpOwner        sdk.AccAddress `json:"cdp_owner"`
	CollateralDenom string         `json:"collateral_denom"`
}

func seizeCdpHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get args from post body
		var req SeizeAndStartCollateralAuctionRequest
		if !rest.ReadRESTReq(w, r, cdc, &req) { // This function writes a response on error
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) { // This function writes a response on error
			return
		}

		// Create msg
		msg := types.MsgSeizeAndStartCollateralAuction{
			req.Sender,
			req.CdpOwner,
			req.CollateralDenom,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Generate tx and write response
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

type StartDebtAuctionRequest struct {
	BaseReq rest.BaseReq   `json:"base_req"`
	Sender  sdk.AccAddress `json:"sender"` // TODO use baseReq.From instead?
}

func debtAuctionHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get args from post body
		var req StartDebtAuctionRequest
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			return
		}
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// Create msg
		msg := types.MsgStartDebtAuction{
			req.Sender,
		}
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Generate tx and write response
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
