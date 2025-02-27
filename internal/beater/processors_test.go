// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package beater

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"

	"github.com/elastic/apm-server/internal/beater/ratelimit"
	"github.com/elastic/apm-server/internal/model"
)

func TestRateLimitBatchProcessor(t *testing.T) {
	limiter := rate.NewLimiter(1, 10)
	ctx := ratelimit.ContextWithLimiter(context.Background(), limiter)

	batch := make(model.Batch, 5)
	for i := range batch {
		batch[i].Transaction = &model.Transaction{}
	}
	for i := 0; i < 2; i++ {
		err := rateLimitBatchProcessor(ctx, &batch)
		require.NoError(t, err)
	}

	// After the second batch, the rate limiter burst has been exhausted,
	// and the limit is not high enough to allow another one.
	err := rateLimitBatchProcessor(ctx, &batch)
	assert.Equal(t, ratelimit.ErrRateLimitExceeded, err)
}
