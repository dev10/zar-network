package tests

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestApp_CreateModifyDeleteCDP(t *testing.T) {
	// Setup
	mapp, keeper := setUpMockAppWithoutGenesis()
	genAccs, addrs, _, privKeys := mock.CreateGenAccounts(1, cs(c("xrp", 100)))
	testAddr := addrs[0]
	testPrivKey := privKeys[0]
	mock.SetGenesis(mapp, genAccs)
	// setup pricefeed, TODO can this be shortened a bit?
	header := abci.Header{Height: mapp.LastBlockHeight() + 1}
	mapp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ctx := mapp.BaseApp.NewContext(false, header)
	keeper.pricefeed.AddAsset(ctx, "xrp", "xrp test")
	keeper.pricefeed.SetPrice(
		ctx, sdk.AccAddress{}, "xrp",
		sdk.MustNewDecFromStr("1.00"),
		sdk.NewInt(10))
	keeper.pricefeed.SetCurrentPrices(ctx)
	mapp.EndBlock(abci.RequestEndBlock{})
	mapp.Commit()

	// Create CDP
	msgs := []sdk.Msg{NewMsgCreateOrModifyCDP(testAddr, "xrp", i(10), i(5))}
	mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, abci.Header{Height: mapp.LastBlockHeight() + 1}, msgs, []uint64{0}, []uint64{0}, true, true, testPrivKey)

	mock.CheckBalance(t, mapp, testAddr, cs(c(StableDenom, 5), c("xrp", 90)))

	// Modify CDP
	msgs = []sdk.Msg{NewMsgCreateOrModifyCDP(testAddr, "xrp", i(40), i(5))}
	mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, abci.Header{Height: mapp.LastBlockHeight() + 1}, msgs, []uint64{0}, []uint64{1}, true, true, testPrivKey)

	mock.CheckBalance(t, mapp, testAddr, cs(c(StableDenom, 10), c("xrp", 50)))

	// Delete CDP
	msgs = []sdk.Msg{NewMsgCreateOrModifyCDP(testAddr, "xrp", i(-50), i(-10))}
	mock.SignCheckDeliver(t, mapp.Cdc, mapp.BaseApp, abci.Header{Height: mapp.LastBlockHeight() + 1}, msgs, []uint64{0}, []uint64{2}, true, true, testPrivKey)

	mock.CheckBalance(t, mapp, testAddr, cs(c("xrp", 100)))
}
