//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package core

import (
	avl "github.com/emirpasic/gods/trees/avltree"
	"github.com/emirpasic/gods/utils"
)

//=============================================================================

type SortedResults[T any] struct {
	Size    int
	MaxSize int
	Tree    *avl.Tree
}

//=============================================================================

func NewSortedResults[T any](maxSize int, comparator utils.Comparator) *SortedResults[T] {
	return &SortedResults[T]{
		Size   : 0,
		MaxSize: maxSize,
		Tree   : avl.NewWith(comparator),
	}
}

//=============================================================================

func (sr *SortedResults[T]) Add(item T) {
	sr.Tree.Put(item, nil)
	sr.Size++

	if sr.Size > sr.MaxSize {
		node := sr.Tree.Right()
		if node != nil {
			sr.Tree.Remove(node.Key)
			sr.Size--
		}
	}
}

//=============================================================================

func (sr *SortedResults[T]) ToList() []T {
	if sr.Tree == nil {
		return []T{}
	}

	keys := sr.Tree.Keys()
	list := make([]T, len(keys))

	for i, k := range keys {
		list[i] = k.(T)
	}
	return list
}

//=============================================================================
