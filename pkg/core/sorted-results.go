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

type SortedResults struct {
	Size    int
	MaxSize int
	Tree    *avl.Tree
}

//=============================================================================

func NewSortedResults(maxSize int, comparator utils.Comparator) *SortedResults {
	return &SortedResults{
		Size   : 0,
		MaxSize: maxSize,
		Tree   : avl.NewWith(comparator),
	}
}

//=============================================================================

func (sr *SortedResults) Add(item any) {
	sr.Tree.Put(item, nil)
	sr.Size++

	if sr.Size > sr.MaxSize {
		sr.Tree.Remove(sr.Tree.Right().Key)
		sr.Size--
	}
}

//=============================================================================

func (sr *SortedResults) ToList() []any {
	if sr.Tree == nil {
		return []any{}
	}

	return sr.Tree.Keys()
}

//=============================================================================
