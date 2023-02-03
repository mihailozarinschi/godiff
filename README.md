# godiff

`godiff` offers a set of functions to help with the generation of diffs/deltas between 2 files (original and updated). It does so by using
a rolling-hash algorithm to generate the signature for both files (list of hashes of all the file's chunks), then compares the 2 lists and
provides the diffs/deltas between the two in a format like: which chunk was removed from where, and which chunk was added where.

# Getting started

```go
import "github.com/mihailozarinschi/godiff"
```

## Usecase #1: Generate diffs between 2 remote files

Example: a client has an updated version of a file also existing on a server.
And it wants to upload only the differences, and not the complete file, to save network resources.  

The steps to achieve something like that would be:
1. The client generates a signature of its local version, the signature being a list of hashes of all the chunks
```go
// Open the updated file
updated, err := os.Open("./updated.txt")
if err != nil {
    return fmt.Errorf("error opening updated file: %s", err)
}
defer updated.Close()

// Define the settings used for the hashing and chunking of the files
hash := sha1.New()
minChunkSize := int64(4) // bytes
divisor := int64(16)
prime := int64(7)

// Chunk the data and get the list of chunks' hashes
chunks, err := godiff.ChunkData(updated, hash, minChunkSize, divisor, prime)
if err != nil {
    return fmt.Errorf("error chunking the updated file: %s", err)
}
```

2. The client sends the signature to the server (via HTTP, gRPC, TCP, whatever)

3. The server generates the signature of its version of the file
```go
// Exactly as in step 1, but on the server.
```

4. The server calculates the deltas between the 2 signatures and asks the client for the missing data
```go
chunksDeltas, err := godiff.GetChunksDeltas(originalChunks, updatedChunks)
if err != nil {
    return fmt.Errorf("error getting original vs updated chunks deltas: %s", err)
}

var deltasNeedingData []*godiff.ChunkDelta

for _, delta := range chunksDeltas {
    // Only deltas that Add data will be asked from the client
    if delta.Type == godiff.DeltaTypeAdd {
        deltasNeedingData = append(deltasNeedingData, delta)
    }
}

// Ask the client for the missing chunks' data
```

5. The server patches the original file using the deltas.
```go
// Not implemented yet
```

## Usecase #2: Generate diffs between 2 local files

Example: a "git-like" client that has both versions, and generates and uploads only the exact diffs to a server

```go
// Open the original file
original, err := os.Open("./original.txt")
if err != nil {
    return fmt.Errorf("error opening original file: %s", err)
}
defer original.Close()

// Open the updated file
updated, err := os.Open("./updated.txt")
if err != nil {
    return fmt.Errorf("error opening updated file: %s", err)
}
defer updated.Close()

// Define the settings used for the hashing and chunking of the files
hashFn := sha1.New
minChunkSize := int64(4) // bytes
divisor := int64(16)
prime := int64(7)

// Generate the diffs between the 2 files
diffs, err := godiff.CalcDiffs(original, updated, hashFn, minChunkSize, divisor, prime)
if err != nil {
    return fmt.Errorf("error generating diffs between original and updated file: %s", err)
}

// Use the diffs however you feel the need
for _, diff := range diffs {
    // Each diff contains the:
    // - Hash of the Data, generated with the given hashFn
    // - Type (Remove/Add)
    // - Data, that was added or removed
    // - DataOffset, where the change should be applied in the original file
    // - DataLen, useful when removing the data, to know how many bytes to remove
}
```
