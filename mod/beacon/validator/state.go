// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package validator

import (
	"context"
	"time"

	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// computeAndSetStateRoot computes the state root of an outgoing block
// and sets it in the block.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
]) computeAndSetStateRoot(
	ctx context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
) error {
	s.logger.Info(
		"computing state root for block 🌲",
		"slot", blk.GetSlot(),
	)

	var stateRoot primitives.Root
	stateRoot, err := s.computeStateRoot(ctx, st, blk)
	if err != nil {
		s.logger.Error(
			"failed to compute state root while building block ❗️ ",
			"slot", blk.GetSlot(),
			"error", err,
		)
		return err
	}

	s.logger.Info("state root computed for block 💻 ",
		"slot", blk.GetSlot(),
		"state_root", stateRoot,
	)
	blk.SetStateRoot(stateRoot)
	return nil
}

// computeStateRoot computes the state root of an outgoing block.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
]) computeStateRoot(
	ctx context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
) (primitives.Root, error) {
	startTime := time.Now()
	defer s.metrics.measureStateRootComputationTime(startTime)
	if _, err := s.stateProcessor.Transition(
		// TODO: We should think about how having optimistic
		// engine enabled here would affect the proposer when
		// the payload in their block has come from a remote builder.
		&transition.Context{
			Context:                 ctx,
			OptimisticEngine:        true,
			SkipPayloadVerification: true,
			SkipValidateResult:      true,
			SkipValidateRandao:      true,
		},
		st, blk,
	); err != nil {
		return primitives.Root{}, err
	}

	return st.HashTreeRoot()
}
