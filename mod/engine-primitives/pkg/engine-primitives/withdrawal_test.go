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

package engineprimitives_test

import (
	"testing"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

func TestWithdrawal(t *testing.T) {
	withdrawal := engineprimitives.Withdrawal{
		Index:     math.U64(1),
		Validator: math.ValidatorIndex(1),
		Address:   common.ExecutionAddress{1, 2, 3, 4, 5},
		Amount:    math.Gwei(1000),
	}

	require.Equal(t, math.U64(1), withdrawal.GetIndex())
	require.Equal(t, math.ValidatorIndex(1), withdrawal.GetValidatorIndex())
	require.Equal(t,
		common.ExecutionAddress{1, 2, 3, 4, 5},
		withdrawal.GetAddress(),
	)
	require.Equal(t, math.Gwei(1000), withdrawal.GetAmount())
}

func TestWithdrawals(t *testing.T) {
	withdrawals := engineprimitives.Withdrawals{
		&engineprimitives.Withdrawal{
			Index:     math.U64(1),
			Validator: math.ValidatorIndex(1),
			Address:   common.ExecutionAddress{1, 2, 3, 4},
			Amount:    math.Gwei(1000),
		},
	}

	root, err := withdrawals.HashTreeRoot()
	require.NoError(t, err)

	require.NotEqual(t, common.Root{}, root)
}

func TestWithdrawal_Equals(t *testing.T) {
	withdrawal1 := &engineprimitives.Withdrawal{
		Index:     math.U64(1),
		Validator: math.ValidatorIndex(1),
		Address:   common.ExecutionAddress{1, 2, 3, 4, 5},
		Amount:    math.Gwei(1000),
	}

	withdrawal2 := &engineprimitives.Withdrawal{
		Index:     math.U64(1),
		Validator: math.ValidatorIndex(1),
		Address:   common.ExecutionAddress{1, 2, 3, 4, 5},
		Amount:    math.Gwei(1000),
	}

	withdrawal3 := &engineprimitives.Withdrawal{
		Index:     math.U64(2),
		Validator: math.ValidatorIndex(2),
		Address:   common.ExecutionAddress{2, 3, 4, 5, 6},
		Amount:    math.Gwei(2000),
	}

	// Test that Equals returns true for two identical withdrawals
	require.True(t, withdrawal1.Equals(withdrawal2))

	// Test that Equals returns false for two different withdrawals
	require.False(t, withdrawal1.Equals(withdrawal3))
}
