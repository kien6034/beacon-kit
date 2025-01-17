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

package middleware

import (
	"context"
	"time"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/genesis"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// BlockchainService defines the interface for interacting with the blockchain
// state and processing blocks.
type BlockchainService[
	BeaconBlockT any, BlobSidecarsT ssz.Marshallable,
] interface {
	// ProcessGenesisData processes the genesis data and initializes the beacon
	// state.
	ProcessGenesisData(
		context.Context,
		*genesis.Genesis[
			*types.Deposit, *types.ExecutionPayloadHeaderDeneb,
		],
	) ([]*transition.ValidatorUpdate, error)
	// ProcessBlockAndBlobs processes the given beacon block and associated
	// blobs sidecars.
	ProcessBlockAndBlobs(
		context.Context,
		BeaconBlockT,
		BlobSidecarsT,
	) ([]*transition.ValidatorUpdate, error)

	// ReceiveBlockAndBlobs receives a beacon block and
	// associated blobs sidecars for processing.
	ReceiveBlockAndBlobs(
		ctx context.Context,
		blk BeaconBlockT,
		blobs BlobSidecarsT,
	) error
}

// ValidatorService is responsible for building beacon blocks.
type ValidatorService[
	BeaconBlockT any,
	BeaconStateT any,
	BlobSidecarsT ssz.Marshallable,
] interface {
	// RequestBlockForProposal requests the best beacon block for a given slot.
	// It returns the beacon block, associated blobs sidecars, and an error if
	// any.
	RequestBlockForProposal(
		context.Context, // The context for the request.
		math.Slot, // The slot for which the best block is requested.
	) (BeaconBlockT, BlobSidecarsT, error)
}

// TelemetrySink is an interface for sending metrics to a telemetry backend.
type TelemetrySink interface {
	// MeasureSince measures the time since the given time.
	MeasureSince(key string, start time.Time, args ...string)
}

// StorageBackend is an interface for accessing the storage backend.
type StorageBackend[BeaconStateT any] interface {
	StateFromContext(ctx context.Context) BeaconStateT
}

// BeaconState is an interface for accessing the beacon state.
type BeaconState interface {
	ValidatorIndexByPubkey(
		pubkey crypto.BLSPubkey,
	) (math.ValidatorIndex, error)

	GetBlockRootAtIndex(
		index uint64,
	) (primitives.Root, error)

	ValidatorIndexByCometBFTAddress(
		cometBFTAddress []byte,
	) (math.ValidatorIndex, error)
}
