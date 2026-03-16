// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package objstore

import (
	"context"
	"strings"
	"testing"

	"github.com/efficientgo/core/testutil"
)

func TestInMem_StartAfter(t *testing.T) {
	ctx := context.Background()
	b := NewInMemBucket()

	// Upload objects with known keys.
	testutil.Ok(t, b.Upload(ctx, "a/file1.txt", strings.NewReader("data1")))
	testutil.Ok(t, b.Upload(ctx, "b/file2.txt", strings.NewReader("data2")))
	testutil.Ok(t, b.Upload(ctx, "c/file3.txt", strings.NewReader("data3")))
	testutil.Ok(t, b.Upload(ctx, "d/file4.txt", strings.NewReader("data4")))

	t.Run("recursive", func(t *testing.T) {
		var items []string
		testutil.Ok(t, b.Iter(ctx, "", func(name string) error {
			items = append(items, name)
			return nil
		}, WithRecursiveIter(), WithStartAfter("b/file2.txt")))

		testutil.Equals(t, []string{"c/file3.txt", "d/file4.txt"}, items)
	})

	t.Run("non-recursive", func(t *testing.T) {
		var items []string
		testutil.Ok(t, b.Iter(ctx, "", func(name string) error {
			items = append(items, name)
			return nil
		}, WithStartAfter("b/")))

		testutil.Equals(t, []string{"c/", "d/"}, items)
	})

	t.Run("start_after_last_key", func(t *testing.T) {
		var items []string
		testutil.Ok(t, b.Iter(ctx, "", func(name string) error {
			items = append(items, name)
			return nil
		}, WithRecursiveIter(), WithStartAfter("d/file4.txt")))

		testutil.Equals(t, 0, len(items))
	})

	t.Run("start_after_empty_string", func(t *testing.T) {
		var items []string
		testutil.Ok(t, b.Iter(ctx, "", func(name string) error {
			items = append(items, name)
			return nil
		}, WithRecursiveIter(), WithStartAfter("")))

		testutil.Equals(t, []string{"a/file1.txt", "b/file2.txt", "c/file3.txt", "d/file4.txt"}, items)
	})
}

func TestInMem_StartAfter_UnsupportedProvider(t *testing.T) {
	// Validate that passing StartAfter to a provider that doesn't support it returns an error.
	err := ValidateIterOptions(
		[]IterOptionType{Recursive},
		WithStartAfter("foo"),
	)
	testutil.NotOk(t, err)
	testutil.Assert(t, strings.Contains(err.Error(), "iter option is not supported"), "expected ErrOptionNotSupported")
}
