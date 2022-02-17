package networktypes

import (
	"fmt"

	launchtypes "github.com/tendermint/spn/x/launch/types"

	"github.com/tendermint/starport/starport/pkg/cosmosutil"
)

// VerifyRequest verifies the validity of the request from its content (static check)
func VerifyRequest(request launchtypes.Request) error {
	req, ok := request.Content.Content.(*launchtypes.RequestContent_GenesisValidator)
	if ok {
		err := VerifyAddValidatorRequest(req)
		if err != nil {
			return NewWrappedErrInvalidRequest(request.RequestID, err.Error())
		}
	}

	return nil
}

// VerifyAddValidatorRequest verify the validator request parameters
func VerifyAddValidatorRequest(req *launchtypes.RequestContent_GenesisValidator) error {
	// If this is an add validator request
	var (
		peer           = req.GenesisValidator.Peer
		valAddress     = req.GenesisValidator.Address
		consPubKey     = req.GenesisValidator.ConsPubKey
		selfDelegation = req.GenesisValidator.SelfDelegation
	)

	// Check values inside the gentx are correct
	info, _, err := cosmosutil.ParseGentx(req.GenesisValidator.GenTx)
	if err != nil {
		return fmt.Errorf("cannot parse gentx %s", err.Error())
	}

	// Change the address prefix fetched from the gentx to the one used on SPN
	// Because all on-chain stored address on SPN uses the SPN prefix
	spnFetchedAddress, err := cosmosutil.ChangeAddressPrefix(info.DelegatorAddress, SPN)
	if err != nil {
		return err
	}

	// Check validator address
	if valAddress != spnFetchedAddress {
		return fmt.Errorf(
			"the validator address %s doesn't match the one inside the gentx %s",
			valAddress,
			spnFetchedAddress,
		)
	}

	// Check validator address
	if !info.PubKey.Equal(consPubKey) {
		return fmt.Errorf(
			"the consensus pub key %s doesn't match the one inside the gentx %s",
			string(consPubKey),
			string(info.PubKey),
		)
	}

	// Check self delegation
	if selfDelegation.Denom != info.SelfDelegation.Denom ||
		!selfDelegation.IsEqual(info.SelfDelegation) {
		return fmt.Errorf(
			"the self delegation %s doesn't match the one inside the gentx %s",
			selfDelegation.String(),
			info.SelfDelegation.String(),
		)
	}

	// Check the format of the peer
	if !cosmosutil.VerifyPeerFormat(peer) {
		return fmt.Errorf(
			"the peer address %s doesn't match the peer format <host>:<port>",
			peer.String(),
		)
	}
	return nil
}
