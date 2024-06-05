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

package runtime

import (
	"context"

	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/beacon/validator"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/runtime/middleware/v2"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/service"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
)

// BeaconKitRuntime is a struct that holds the
// service registry.
type BeaconKitRuntime[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockT interface {
		types.RawBeaconBlock[BeaconBlockBodyT]
		NewFromSSZ([]byte, uint32) (BeaconBlockT, error)
		NewWithVersion(
			math.Slot,
			math.ValidatorIndex,
			primitives.Root,
			uint32,
		) (BeaconBlockT, error)
		Empty(uint32) BeaconBlockT
	},
	BeaconBlockBodyT types.BeaconBlockBody,
	BeaconStateT core.BeaconState[
		*types.BeaconBlockHeader,
		*types.ExecutionPayloadHeader,
		*types.Fork,
		*types.Validator,
		*engineprimitives.Withdrawal,
	],
	BlobSidecarsT BlobSidecars,
	DepositStoreT DepositStore,
	StorageBackendT StorageBackend[
		AvailabilityStoreT, BeaconBlockBodyT,
		BeaconStateT, BlobSidecarsT, DepositStoreT,
	],
	TransactionT types.Tx,
] struct {
	// logger is used for logging within the BeaconKitRuntime.
	logger log.Logger[any]
	// services is a registry of services used by the BeaconKitRuntime.
	services *service.Registry
	// storageBackend is the backend storage interface used by the
	// BeaconKitRuntime.
	storageBackend StorageBackendT
	// chainSpec defines the chain specifications for the BeaconKitRuntime.
	chainSpec primitives.ChainSpec
	// abciFinalizeBlockMiddleware handles ABCI interactions for the
	// BeaconKitRuntime.
	abciFinalizeBlockMiddleware *middleware.FinalizeBlockMiddleware[
		BeaconBlockT, BeaconStateT, BlobSidecarsT,
	]
	// abciValidatorMiddleware is responsible for forward ABCI requests to the
	// validator service.
	abciValidatorMiddleware *middleware.ValidatorMiddleware[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT,
		BeaconStateT, BlobSidecarsT, StorageBackendT, TransactionT,
	]
}

// NewBeaconKitRuntime creates a new BeaconKitRuntime
// and applies the provided options.
func NewBeaconKitRuntime[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockT interface {
		types.RawBeaconBlock[BeaconBlockBodyT]
		NewFromSSZ([]byte, uint32) (BeaconBlockT, error)
		NewWithVersion(
			math.Slot,
			math.ValidatorIndex,
			primitives.Root,
			uint32,
		) (BeaconBlockT, error)
		Empty(uint32) BeaconBlockT
	},
	BeaconBlockBodyT types.BeaconBlockBody,
	BeaconStateT core.BeaconState[
		*types.BeaconBlockHeader, *types.ExecutionPayloadHeader, *types.Fork,
		*types.Validator, *engineprimitives.Withdrawal,
	],
	BlobSidecarsT BlobSidecars,
	DepositStoreT DepositStore,
	StorageBackendT blockchain.StorageBackend[
		AvailabilityStoreT,
		BeaconBlockBodyT,
		BeaconStateT,
		BlobSidecarsT,
		*types.Deposit,
		DepositStoreT,
	],
	TransactionT types.Tx,
](
	chainSpec primitives.ChainSpec,
	logger log.Logger[any],
	services *service.Registry,
	storageBackend StorageBackendT,
	telemetrySink middleware.TelemetrySink,
) (*BeaconKitRuntime[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, StorageBackendT, TransactionT,
], error) {
	var (
		chainService *blockchain.Service[
			AvailabilityStoreT,
			BeaconBlockT,
			BeaconBlockBodyT,
			core.BeaconState[
				*types.BeaconBlockHeader,
				*types.ExecutionPayloadHeader,
				*types.Fork,
				*types.Validator,
				*engineprimitives.Withdrawal,
			],
			BlobSidecarsT,
			*types.Deposit,
			DepositStoreT,
		]
		validatorService *validator.Service[
			BeaconBlockT,
			BeaconBlockBodyT,
			core.BeaconState[
				*types.BeaconBlockHeader,
				*types.ExecutionPayloadHeader,
				*types.Fork,
				*types.Validator,
				*engineprimitives.Withdrawal,
			],
			BlobSidecarsT,
			DepositStoreT,
		]
	)

	if err := services.FetchService(&chainService); err != nil {
		panic(err)
	}

	if err := services.FetchService(&validatorService); err != nil {
		panic(err)
	}

	return &BeaconKitRuntime[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
		BlobSidecarsT, DepositStoreT, StorageBackendT, TransactionT,
	]{
		abciFinalizeBlockMiddleware: middleware.
			NewFinalizeBlockMiddleware[
			BeaconBlockT, BeaconStateT, BlobSidecarsT,
		](
			chainSpec,
			chainService,
			telemetrySink,
		),
		abciValidatorMiddleware: middleware.
			NewValidatorMiddleware[
			AvailabilityStoreT,
			BeaconBlockT,
			BeaconBlockBodyT,
			BeaconStateT,
			BlobSidecarsT,
			StorageBackendT,
			TransactionT,
		](
			chainSpec,
			validatorService,
			telemetrySink,
			storageBackend,
		),
		chainSpec:      chainSpec,
		logger:         logger,
		services:       services,
		storageBackend: storageBackend,
	}, nil
}

// StartServices starts the services.
func (r *BeaconKitRuntime[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, StorageBackendT, TransactionT,
]) StartServices(
	ctx context.Context,
) error {
	return r.services.StartAll(ctx)
}

// ABCIHandler returns the ABCI handler.
func (r *BeaconKitRuntime[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, StorageBackendT, TransactionT,
]) ABCIFinalizeBlockMiddleware() *middleware.FinalizeBlockMiddleware[
	BeaconBlockT, BeaconStateT, BlobSidecarsT,
] {
	return r.abciFinalizeBlockMiddleware
}

// ABCIValidatorMiddleware returns the ABCI validator middleware.
func (r *BeaconKitRuntime[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, StorageBackendT, TransactionT,
]) ABCIValidatorMiddleware() *middleware.ValidatorMiddleware[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT,
	BeaconStateT, BlobSidecarsT, StorageBackendT, TransactionT,
] {
	return r.abciValidatorMiddleware
}