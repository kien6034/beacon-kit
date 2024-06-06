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

package components

import (
	"cosmossdk.io/core/log"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	execution "github.com/berachain/beacon-kit/mod/execution/pkg/engine"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/config"
	payloadbuilder "github.com/berachain/beacon-kit/mod/payload/pkg/builder"
	"github.com/berachain/beacon-kit/mod/payload/pkg/cache"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

func ProvideLocalBuilder(
	cfg *config.Config,
	chainSpec primitives.ChainSpec,
	logger log.Logger,
	executionEngine *execution.Engine[*types.ExecutionPayload],
) *payloadbuilder.PayloadBuilder[
	BeaconState, *types.ExecutionPayload, *types.ExecutionPayloadHeader,
] {
	return payloadbuilder.New[
		BeaconState, *types.ExecutionPayload, *types.ExecutionPayloadHeader,
	](
		&cfg.PayloadBuilder,
		chainSpec,
		logger.With("service", "payload-builder"),
		executionEngine,
		cache.NewPayloadIDCache[engineprimitives.PayloadID, [32]byte, math.Slot](),
	)
}