package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/24-host"
	"github.com/cosmos/cosmos-sdk/x/ibc/exported"
	"github.com/cosmos/cosmos-sdk/x/ibc/light-clients/solomachine/types"
	ibctesting "github.com/cosmos/cosmos-sdk/x/ibc/testing"
)

type SoloMachineTestSuite struct {
	suite.Suite

	solomachine *ibctesting.Solomachine
	coordinator *ibctesting.Coordinator

	// testing chain used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain

	store sdk.KVStore
}

func (suite *SoloMachineTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(1))

	suite.solomachine = ibctesting.NewSolomachine(suite.T(), suite.chainA.Codec, "testingsolomachine", "testing")
	suite.store = suite.chainA.App.IBCKeeper.ClientKeeper.ClientStore(suite.chainA.GetContext(), types.SoloMachine)

	bz, err := codec.MarshalAny(suite.chainA.Codec, suite.solomachine.ClientState())
	suite.Require().NoError(err)
	suite.store.Set(host.KeyClientState(), bz)
}

func TestSoloMachineTestSuite(t *testing.T) {
	suite.Run(t, new(SoloMachineTestSuite))
}

func (suite *SoloMachineTestSuite) GetSequenceFromStore() uint64 {
	bz := suite.store.Get(host.KeyClientState())
	suite.Require().NotNil(bz)

	var clientState exported.ClientState
	err := codec.UnmarshalAny(suite.chainA.Codec, &clientState, bz)
	suite.Require().NoError(err)
	return clientState.GetLatestHeight().GetEpochHeight()
}

func (suite *SoloMachineTestSuite) GetInvalidProof() []byte {
	invalidProof, err := suite.chainA.Codec.MarshalBinaryBare(&types.TimestampedSignature{Timestamp: suite.solomachine.Time})
	suite.Require().NoError(err)

	return invalidProof
}

func TestUnpackInterfaces_Header(t *testing.T) {
	registry := testdata.NewTestInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)

	pk := secp256k1.GenPrivKey().PubKey().(cryptotypes.PubKey)
	any, err := codectypes.NewAnyWithValue(pk)
	require.NoError(t, err)

	header := types.Header{
		NewPublicKey: any,
	}
	bz, err := header.Marshal()
	require.NoError(t, err)

	var header2 types.Header
	err = header2.Unmarshal(bz)
	require.NoError(t, err)

	err = codectypes.UnpackInterfaces(header2, registry)
	require.NoError(t, err)

	require.Equal(t, pk, header2.NewPublicKey.GetCachedValue())
}

func TestUnpackInterfaces_HeaderData(t *testing.T) {
	registry := testdata.NewTestInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)

	pk := secp256k1.GenPrivKey().PubKey().(cryptotypes.PubKey)
	any, err := codectypes.NewAnyWithValue(pk)
	require.NoError(t, err)

	hd := types.HeaderData{
		NewPubKey: any,
	}
	bz, err := hd.Marshal()
	require.NoError(t, err)

	var hd2 types.HeaderData
	err = hd2.Unmarshal(bz)
	require.NoError(t, err)

	err = codectypes.UnpackInterfaces(hd2, registry)
	require.NoError(t, err)

	require.Equal(t, pk, hd2.NewPubKey.GetCachedValue())
}
