// SPDX-License-Identifier: MIT
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
//
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package core

import (
	"fmt"

	"github.com/berachain/beacon-kit/mod/config/params"
	"github.com/berachain/beacon-kit/mod/core/state"
	"github.com/berachain/beacon-kit/mod/core/types"
	enginetypes "github.com/berachain/beacon-kit/mod/execution/types"
	"github.com/ethereum/go-ethereum/common"
)

// PayloadValidator is responsible for validating incoming execution
// payloads to ensure they are valid.
type PayloadValidator struct {
	cfg *params.BeaconChainConfig
}

// NewPayloadValidator creates a new payload validator.
func NewPayloadValidator(cfg *params.BeaconChainConfig) *PayloadValidator {
	return &PayloadValidator{
		cfg: cfg,
	}
}

// ValidatePayload validates the incoming payload.
func (pv *PayloadValidator) ValidatePayload(
	st state.BeaconState,
	body types.BeaconBlockBody,
) error {
	if body == nil || body.IsNil() {
		return types.ErrNilBlkBody
	}

	payload := body.GetExecutionPayload()
	if payload == nil || payload.IsNil() {
		return types.ErrNilPayload
	}

	// TODO: Once deneb genesis data contains execution payload, remove this.
	var safeHash common.Hash

	// Get the current epoch.
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	// Handle genesis case.
	if slot <= 1 {
		safeHash, err = st.GetEth1BlockHash()
		if err != nil {
			return err
		}
	} else {
		var executionPayload enginetypes.ExecutionPayload
		executionPayload, err = st.GetLatestExecutionPayload()
		if err != nil {
			return err
		}
		safeHash = executionPayload.GetBlockHash()
	}

	if safeHash != payload.GetParentHash() {
		return fmt.Errorf(
			"parent block with hash %x is not finalized, expected finalized hash %x",
			payload.GetParentHash(),
			safeHash,
		)
	}

	// When we are validating a payload we expect that it was produced by
	// the proposer for the slot that it is for.
	expectedMix, err := st.GetRandaoMixAtIndex(
		uint64(pv.cfg.SlotToEpoch(slot)) % pv.cfg.EpochsPerHistoricalVector)
	if err != nil {
		return err
	}

	// Ensure the prev randao matches the local state.
	if payload.GetPrevRandao() != expectedMix {
		return fmt.Errorf(
			"prev randao does not match, expected: %x, got: %x",
			expectedMix, payload.GetPrevRandao(),
		)
	}

	// TODO: Verify timestamp data once Clock is done.
	// if expectedTime, err := spec.TimeAtSlot(slot, genesisTime); err != nil {
	// 	return fmt.Errorf("slot or genesis time in state is corrupt, cannot
	// compute time: %v", err)
	// } else if payload.Timestamp != expectedTime {
	// 	return fmt.Errorf("state at slot %d, genesis time %d, expected execution
	// payload time %d, but got %d",
	// 		slot, genesisTime, expectedTime, payload.Timestamp)
	// }

	if uint64(len(body.GetBlobKzgCommitments())) > pv.cfg.MaxBlobsPerBlock {
		return fmt.Errorf(
			"too many blob kzg commitments, expected: %d, got: %d",
			pv.cfg.MaxBlobsPerBlock,
			len(body.GetBlobKzgCommitments()),
		)
	}

	// Verify the number of withdrawals.
	if withdrawals := payload.GetWithdrawals(); uint64(
		len(payload.GetWithdrawals()),
	) > pv.cfg.MaxWithdrawalsPerPayload {
		return fmt.Errorf(
			"too many withdrawals, expected: %d, got: %d",
			pv.cfg.MaxWithdrawalsPerPayload, len(withdrawals),
		)
	}

	return nil
}