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

	var (
		// Cursors for the original slice and updated slice
		oc, uc int

		// Temp deltas
		removals  []*ChunkDelta
		additions []*ChunkDelta

		// Final deltas
		deltas []*ChunkDelta
	)

mainLoop:
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
			deltas = append(deltas, removals...)  // persist removals
			removals = removals[:0]               // clear temp removals
			deltas = append(deltas, additions...) // persist additions
			additions = additions[:0]             // clear temp additions
			goto moveCursors
		}

		// Chunks' hashes differ.
		// Lookup new chunk in the list of temporarily removed chunks, maybe it was shifted forward
		if u != nil {
			for remPos, remChunk := range removals {
				if remChunk.Hash == u.Hash {
					// New chunk was simply shifted forward by a chunk addition, found it in the list of temp removals.
					deltas = append(deltas, removals[:remPos]...) // persist temp removals until this point
					removals = removals[:0]                       // clear temp removals cache
					oc = remChunk.Position                        // Reset original slice cursor to the found one
					continue mainLoop
				}
			}
		}
		// Lookup original chunk in the list of temporarily added chunks, maybe it was shifted backwards
		if o != nil {
			for addPos, addChunk := range additions {
				if addChunk.Hash == o.Hash {
					// Original chunk was simply shifted backward by a chunk removal
					deltas = append(deltas, additions[:addPos]...) // persist temp additions until this point
					additions = additions[:0]                      // clear temp additions cache
					uc = addChunk.Position                         // Reset updated slice cursor to the found one
					continue mainLoop
				}
			}
		}

		// No reoccurrence found in the existing temp deltas, add these too.
		if o != nil {
			removals = append(removals, &ChunkDelta{o, DeltaTypeRemove, oc})
		}
		if u != nil {
			additions = append(additions, &ChunkDelta{u, DeltaTypeAdd, uc})
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
