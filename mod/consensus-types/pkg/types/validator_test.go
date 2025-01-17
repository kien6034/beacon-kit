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

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

func TestNewValidatorFromDeposit(t *testing.T) {
	tests := []struct {
		name                      string
		pubkey                    crypto.BLSPubkey
		withdrawalCredentials     types.WithdrawalCredentials
		amount                    math.Gwei
		effectiveBalanceIncrement math.Gwei
		maxEffectiveBalance       math.Gwei
		want                      *types.Validator
	}{
		{
			name:   "normal case",
			pubkey: [48]byte{0x01},
			withdrawalCredentials: types.
				NewCredentialsFromExecutionAddress(
					common.ExecutionAddress{0x01},
				),
			amount:                    32e9,
			effectiveBalanceIncrement: 1e9,
			maxEffectiveBalance:       32e9,
			want: &types.Validator{
				Pubkey: [48]byte{0x01},
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						common.ExecutionAddress{0x01},
					),
				EffectiveBalance: 32e9,
				Slashed:          false,
				ActivationEligibilityEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ActivationEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ExitEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				WithdrawableEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
			},
		},
		{
			name:   "effective balance capped at max",
			pubkey: [48]byte{0x02},
			withdrawalCredentials: types.
				NewCredentialsFromExecutionAddress(
					common.ExecutionAddress{0x02},
				),
			amount:                    40e9,
			effectiveBalanceIncrement: 1e9,
			maxEffectiveBalance:       32e9,
			want: &types.Validator{
				Pubkey: [48]byte{0x02},
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						common.ExecutionAddress{0x02},
					),
				EffectiveBalance: 32e9,
				Slashed:          false,
				ActivationEligibilityEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ActivationEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ExitEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				WithdrawableEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
			},
		},
		{
			name:   "effective balance rounded down",
			pubkey: [48]byte{0x03},
			withdrawalCredentials: types.
				NewCredentialsFromExecutionAddress(
					common.ExecutionAddress{0x03},
				),
			amount:                    32.5e9,
			effectiveBalanceIncrement: 1e9,
			maxEffectiveBalance:       32e9,
			want: &types.Validator{
				Pubkey: [48]byte{0x03},
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						common.ExecutionAddress{0x03},
					),
				EffectiveBalance: 32e9,
				Slashed:          false,
				ActivationEligibilityEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ActivationEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				ExitEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				WithdrawableEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := types.NewValidatorFromDeposit(
				tt.pubkey,
				tt.withdrawalCredentials,
				tt.amount,
				tt.effectiveBalanceIncrement,
				tt.maxEffectiveBalance,
			)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestValidator_IsActive(t *testing.T) {
	tests := []struct {
		name      string
		epoch     math.Epoch
		validator *types.Validator
		want      bool
	}{
		{
			name:  "active",
			epoch: 10,
			validator: &types.Validator{
				ActivationEpoch: 5,
				ExitEpoch:       15,
			},
			want: true,
		},
		{
			name:  "not active, before activation",
			epoch: 4,
			validator: &types.Validator{
				ActivationEpoch: 5,
				ExitEpoch:       15,
			},
			want: false,
		},
		{
			name:  "not active, after exit",
			epoch: 16,
			validator: &types.Validator{
				ActivationEpoch: 5,
				ExitEpoch:       15,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.validator.IsActive(tt.epoch))
		})
	}
}

func TestValidator_IsEligibleForActivation(t *testing.T) {
	tests := []struct {
		name           string
		finalizedEpoch math.Epoch
		validator      *types.Validator
		want           bool
	}{
		{
			name:           "eligible",
			finalizedEpoch: 10,
			validator: &types.Validator{
				ActivationEligibilityEpoch: 5,
				ActivationEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
			},
			want: true,
		},
		{
			name:           "not eligible, activation eligibility in future",
			finalizedEpoch: 4,
			validator: &types.Validator{
				ActivationEligibilityEpoch: 5,
				ActivationEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
			},
			want: false,
		},
		{
			name:           "not eligible, already activated",
			finalizedEpoch: 10,
			validator: &types.Validator{
				ActivationEligibilityEpoch: 5,
				ActivationEpoch:            8,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(
				t,
				tt.want,
				tt.validator.IsEligibleForActivation(tt.finalizedEpoch),
			)
		})
	}
}

func TestValidator_IsEligibleForActivationQueue(t *testing.T) {
	maxEffectiveBalance := math.Gwei(32e9)
	tests := []struct {
		name      string
		validator *types.Validator
		want      bool
	}{
		{
			name: "eligible",
			validator: &types.Validator{
				ActivationEligibilityEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				EffectiveBalance: maxEffectiveBalance,
			},
			want: true,
		},
		{
			name: "not eligible, activation eligibility set",
			validator: &types.Validator{
				ActivationEligibilityEpoch: 5,
				EffectiveBalance:           maxEffectiveBalance,
			},
			want: false,
		},
		{
			name: "not eligible, effective balance too low",
			validator: &types.Validator{
				ActivationEligibilityEpoch: math.Epoch(
					constants.FarFutureEpoch,
				),
				EffectiveBalance: maxEffectiveBalance - 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(
				t,
				tt.want,
				tt.validator.IsEligibleForActivationQueue(maxEffectiveBalance),
			)
		})
	}
}

func TestValidator_IsSlashable(t *testing.T) {
	tests := []struct {
		name      string
		epoch     math.Epoch
		validator *types.Validator
		want      bool
	}{
		{
			name:  "slashable",
			epoch: 10,
			validator: &types.Validator{
				Slashed:           false,
				ActivationEpoch:   5,
				WithdrawableEpoch: 15,
			},
			want: true,
		},
		{
			name:  "not slashable, already slashed",
			epoch: 10,
			validator: &types.Validator{
				Slashed:           true,
				ActivationEpoch:   5,
				WithdrawableEpoch: 15,
			},
			want: false,
		},
		{
			name:  "not slashable, before activation",
			epoch: 4,
			validator: &types.Validator{
				Slashed:           false,
				ActivationEpoch:   5,
				WithdrawableEpoch: 15,
			},
			want: false,
		},
		{
			name:  "not slashable, after withdrawable",
			epoch: 16,
			validator: &types.Validator{
				Slashed:           false,
				ActivationEpoch:   5,
				WithdrawableEpoch: 15,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.validator.IsSlashable(tt.epoch))
		})
	}
}

func TestValidator_IsFullyWithdrawable(t *testing.T) {
	tests := []struct {
		name      string
		balance   math.Gwei
		epoch     math.Epoch
		validator *types.Validator
		want      bool
	}{
		{
			name:    "fully withdrawable",
			balance: 32e9,
			epoch:   10,
			validator: &types.Validator{
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						common.ExecutionAddress{0x01},
					),
				WithdrawableEpoch: 5,
			},
			want: true,
		},
		{
			name:    "not fully withdrawable, non-eth1 credentials",
			balance: 32e9,
			epoch:   10,
			validator: &types.Validator{
				WithdrawalCredentials: types.
					WithdrawalCredentials{0x00},
				WithdrawableEpoch: 5,
			},
			want: false,
		},
		{
			name:    "not fully withdrawable, before withdrawable epoch",
			balance: 32e9,
			epoch:   4,
			validator: &types.Validator{
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						common.ExecutionAddress{0x01},
					),
				WithdrawableEpoch: 5,
			},
			want: false,
		},
		{
			name:    "not fully withdrawable, zero balance",
			balance: 0,
			epoch:   10,
			validator: &types.Validator{
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						common.ExecutionAddress{0x01},
					),
				WithdrawableEpoch: 5,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(
				t,
				tt.want,
				tt.validator.IsFullyWithdrawable(tt.balance, tt.epoch),
			)
		})
	}
}

func TestValidator_IsPartiallyWithdrawable(t *testing.T) {
	maxEffectiveBalance := math.Gwei(32e9)
	tests := []struct {
		name      string
		balance   math.Gwei
		validator *types.Validator
		want      bool
	}{
		{
			name:    "partially withdrawable",
			balance: 33e9,
			validator: &types.Validator{
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						common.ExecutionAddress{0x01},
					),
				EffectiveBalance: maxEffectiveBalance,
			},
			want: true,
		},
		{
			name:    "not partially withdrawable, non-eth1 credentials",
			balance: 33e9,
			validator: &types.Validator{
				WithdrawalCredentials: types.WithdrawalCredentials{
					0x00,
				},
				EffectiveBalance: maxEffectiveBalance,
			},
			want: false,
		},
		{
			name:    "not partially withdrawable, not at max effective balance",
			balance: 33e9,
			validator: &types.Validator{
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						common.ExecutionAddress{0x01},
					),
				EffectiveBalance: maxEffectiveBalance - 1,
			},
			want: false,
		},
		{
			name:    "not partially withdrawable, no excess balance",
			balance: 32e9,
			validator: &types.Validator{
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						common.ExecutionAddress{0x01},
					),
				EffectiveBalance: maxEffectiveBalance,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(
				t,
				tt.want,
				tt.validator.IsPartiallyWithdrawable(
					tt.balance,
					maxEffectiveBalance,
				),
			)
		})
	}
}

func TestValidator_HasEth1WithdrawalCredentials(t *testing.T) {
	tests := []struct {
		name      string
		validator *types.Validator
		want      bool
	}{
		{
			name: "has eth1 credentials",
			validator: &types.Validator{
				WithdrawalCredentials: types.
					NewCredentialsFromExecutionAddress(
						common.ExecutionAddress{0x01},
					),
			},
			want: true,
		},
		{
			name: "does not have eth1 credentials",
			validator: &types.Validator{
				WithdrawalCredentials: types.WithdrawalCredentials{
					0x00,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(
				t,
				tt.want,
				tt.validator.HasEth1WithdrawalCredentials(),
			)
		})
	}
}

func TestValidator_HasMaxEffectiveBalance(t *testing.T) {
	maxEffectiveBalance := math.Gwei(32e9)
	tests := []struct {
		name      string
		validator *types.Validator
		want      bool
	}{
		{
			name: "has max effective balance",
			validator: &types.Validator{
				EffectiveBalance: maxEffectiveBalance,
			},
			want: true,
		},
		{
			name: "does not have max effective balance",
			validator: &types.Validator{
				EffectiveBalance: maxEffectiveBalance - 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(
				t,
				tt.want,
				tt.validator.HasMaxEffectiveBalance(maxEffectiveBalance),
			)
		})
	}
}
