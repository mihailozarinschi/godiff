package godiff

import (
	"sort"
)

type DeltaType int

func (d DeltaType) String() string {
	switch d {
	case DeltaTypeRemove:
		return "Rem"
	case DeltaTypeAdd:
		return "Add"
	default:
		return ""
	}
}

const (
	DeltaTypeRemove DeltaType = iota
	DeltaTypeAdd
)

// ChunkDelta contains the information about what happened with any given chunk,
// was it removed/added and from/at which position.
type ChunkDelta struct {
	*Chunk
	Type     DeltaType
	Position int
}

// GetChunksDeltas tries to provide the minimum amount of deltas between any 2 given slices of chunks.
// It navigates through both slices at the same time, when chunks start to diverge, it keeps track of
// possible ("temporary") removals and/or additions, with each new step forward it checks within the
// temporary changes for forward and/or backward shifts that happened, and re-positions the cursors accordingly.
// Finally, it re-orders the deltas in the proper order to be applied/patched on the original data.
// The order presumes all removals first, DESC (from end to start), then all additions ASC (from start to end),
func GetChunksDeltas(original, updated []*Chunk) ([]*ChunkDelta, error) {

	type deltaIndex struct {
		delta *ChunkDelta
		index int
	}

	var (
		// Cursors for the original slice and updated slice
		oc, uc int

		// Temp deltas
		removals       []*ChunkDelta
		removalsIndex  = make(map[string]*deltaIndex) // to speedup lookups
		additions      []*ChunkDelta
		additionsIndex = make(map[string]*deltaIndex) // to speedup lookups

		// Final deltas
		deltas []*ChunkDelta
	)

	// Loop until both slices are exhausted
	for oc < len(original) || uc < len(updated) {

		var o *Chunk
		if oc < len(original) {
			o = original[oc]
		}

		var u *Chunk
		if uc < len(updated) {
			u = updated[uc]
		}

		// Chunks' hashes are the same
		if o != nil && u != nil && o.Hash == u.Hash {
			// We got to a converging point, all temp changes so far need to be persisted
			deltas = append(deltas, removals...) // persist removals
			removals = removals[:0]              // clear temp removals
			removalsIndex = make(map[string]*deltaIndex)

			deltas = append(deltas, additions...) // persist additions
			additions = additions[:0]             // clear temp additions
			additionsIndex = make(map[string]*deltaIndex)

			goto moveCursors
		}

		// Chunks' hashes differ.
		// Lookup new chunk in the list of temporarily removed chunks, maybe it was shifted forward
		if u != nil {
			rem, ok := removalsIndex[u.Hash]
			if ok {
				// New chunk was simply shifted forward by a chunk addition, found it in the list of temp removals.
				deltas = append(deltas, removals[:rem.index]...) // persist temp removals until this point
				removals = removals[:0]                          // clear temp removals cache
				removalsIndex = make(map[string]*deltaIndex)
				oc = rem.delta.Position // Reset original slice cursor to the found one
				continue
			}
		}
		// Lookup original chunk in the list of temporarily added chunks, maybe it was shifted backwards
		if o != nil {
			add, ok := additionsIndex[o.Hash]
			if ok {
				// Original chunk was simply shifted backward by a chunk removal
				deltas = append(deltas, additions[:add.index]...) // persist temp additions until this point
				additions = additions[:0]                         // clear temp additions cache
				additionsIndex = make(map[string]*deltaIndex)
				uc = add.delta.Position
				continue
			}
		}

		// No reoccurrence found in the existing temp deltas, add these too.
		if o != nil {
			chunkDelta := &ChunkDelta{o, DeltaTypeRemove, oc}
			removals = append(removals, chunkDelta)
			// index only first occurrence
			if _, ok := removalsIndex[o.Hash]; !ok {
				removalsIndex[o.Hash] = &deltaIndex{delta: chunkDelta, index: len(removals) - 1}
			}
		}
		if u != nil {
			chunkDelta := &ChunkDelta{u, DeltaTypeAdd, uc}
			additions = append(additions, chunkDelta)
			// index only first occurrence
			if _, ok := additionsIndex[u.Hash]; !ok {
				additionsIndex[u.Hash] = &deltaIndex{delta: chunkDelta, index: len(additions) - 1}
			}
		}

	moveCursors:
		if oc < len(original) {
			oc++
		}
		if uc < len(updated) {
			uc++
		}
	}

	deltas = append(deltas, removals...)
	deltas = append(deltas, additions...)

	// Set the order in which the deltas should be applied
	sort.SliceStable(deltas, func(i, j int) bool {
		di := deltas[i]
		dj := deltas[j]
		return di.Type < dj.Type || // Remove (1) < Add (2). Removals before additions
			(di.Type == DeltaTypeRemove && di.Type == dj.Type && di.Position > dj.Position) || // Sort removals DESC
			(di.Type == DeltaTypeAdd && di.Type == dj.Type && di.Position < dj.Position) // Sort additions ASC
	})

	return deltas, nil
}
