// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
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

package beacondb

import (
	"github.com/berachain/beacon-kit/light/mod/storage/beacondb/keys"
	"github.com/berachain/beacon-kit/mod/primitives"
)

// GetNextWithdrawalIndex returns the next withdrawal index.
func (kv *KVStore) GetNextWithdrawalIndex() (uint64, error) {
	res, err := kv.provider.Query(
		kv.ctx,
		keys.BeaconStoreKey,
		kv.nextWithdrawalIndex.Key(),
		0,
	)
	if err != nil {
		return 0, err
	}

	index, err := kv.nextWithdrawalIndex.Decode(res)
	if err != nil {
		return 0, err
	}

	return index, nil
}

// SetNextWithdrawalIndex sets the next withdrawal index.
func (kv *KVStore) SetNextWithdrawalIndex(
	index uint64,
) error {
	panic(writesNotSupported)
}

// GetNextWithdrawalValidatorIndex returns the next withdrawal validator index.
func (kv *KVStore) GetNextWithdrawalValidatorIndex() (
	primitives.ValidatorIndex, error,
) {
	res, err := kv.provider.Query(
		kv.ctx,
		keys.BeaconStoreKey,
		kv.nextWithdrawalValidatorIndex.Key(),
		0,
	)
	if err != nil {
		return 0, err
	}

	index, err := kv.nextWithdrawalValidatorIndex.Decode(res)
	if err != nil {
		return 0, err
	}

	return primitives.ValidatorIndex(index), nil
}

// SetNextWithdrawalValidatorIndex sets the next withdrawal validator index.
func (kv *KVStore) SetNextWithdrawalValidatorIndex(
	index primitives.ValidatorIndex,
) error {
	panic(writesNotSupported)
}

// ExpectedDeposits returns the first numPeek deposits in the queue.
func (kv *KVStore) ExpectedDeposits(
	numView uint64,
) ([]*primitives.Deposit, error) {
	return kv.depositQueue.PeekMulti(kv.ctx, numView)
}

// EnqueueDeposits pushes the deposits to the queue.
func (kv *KVStore) EnqueueDeposits(
	deposits []*primitives.Deposit,
) error {
	panic(writesNotSupported)
}

// DequeueDeposits returns the first numDequeue deposits in the queue.
func (kv *KVStore) DequeueDeposits(
	numDequeue uint64,
) ([]*primitives.Deposit, error) {
	panic(writesNotSupported)
}