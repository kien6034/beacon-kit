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

package deneb

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

//go:generate go run github.com/ferranbt/fastssz/sszgen -path deneb.go -objs BeaconState -include ../../../../primitives/pkg/crypto,../../../../primitives/pkg/common,../../../../primitives/pkg/bytes,../../../../primitives/mod.go,../../../../consensus-types/pkg/types,../../../../engine-primitives/pkg/engine-primitives,../../../../primitives/mod.go,../../../../primitives/pkg/math,$GETH_PKG_INCLUDE/common,$GETH_PKG_INCLUDE/common/hexutil -output deneb.ssz.go
//nolint:lll // various json tags.
type BeaconState struct {
	// Versioning
	//
	//nolint:lll
	GenesisValidatorsRoot primitives.Root `json:"genesisValidatorsRoot" ssz-size:"32"`
	Slot                  math.Slot       `json:"slot"`
	Fork                  *types.Fork     `json:"fork"`

	// History
	LatestBlockHeader *types.BeaconBlockHeader `json:"latestBlockHeader"`
	BlockRoots        []primitives.Root        `json:"blockRoots"        ssz-size:"?,32" ssz-max:"8192"`
	StateRoots        []primitives.Root        `json:"stateRoots"        ssz-size:"?,32" ssz-max:"8192"`

	// Eth1
	Eth1Data                     *types.Eth1Data                    `json:"eth1Data"`
	Eth1DepositIndex             uint64                             `json:"eth1DepositIndex"`
	LatestExecutionPayloadHeader *types.ExecutionPayloadHeaderDeneb `json:"latestExecutionPayloadHeader"`

	// Registry
	Validators []*types.Validator `json:"validators" ssz-max:"1099511627776"`
	Balances   []uint64           `json:"balances"   ssz-max:"1099511627776"`

	// Randomness
	RandaoMixes []primitives.Bytes32 `json:"randaoMixes" ssz-size:"?,32" ssz-max:"65536"`

	// Withdrawals
	NextWithdrawalIndex          uint64              `json:"nextWithdrawalIndex"`
	NextWithdrawalValidatorIndex math.ValidatorIndex `json:"nextWithdrawalValidatorIndex"`

	// Slashing
	Slashings     []uint64  `json:"slashings"     ssz-max:"1099511627776"`
	TotalSlashing math.Gwei `json:"totalSlashing"`
}
