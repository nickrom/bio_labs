package structs

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "os"
    "strings"
)

const ClusterFileExtension = ".cl"
const SequenceInitCap      = 1000000

/* --- */

// Adapter for database of sequences and their dots data.
type SequenceDb []SequenceEntry

/* --- */

// Loads sequence database from precalculated clusters.
// Cluster directory path must be provided.
func FromClusters(clustersDirPath string) SequenceDb {
    files, err := ioutil.ReadDir(clustersDirPath)

    if err != nil {
        panic(err)
    }

    sequences := make(SequenceDb, 0, SequenceInitCap)

    for _, f := range files {
        if strings.HasSuffix(f.Name(), ClusterFileExtension) {
            fmt.Printf("Processing cluster %v ...\n", f.Name())
            file, err := os.Open(clustersDirPath + "/" + f.Name())

            if err != nil {
                panic(err)
            }

            scanner := bufio.NewScanner(file)
            buf := make([]byte, 0, 1024 * 1024)
            scanner.Buffer(buf, 100 * 1024 * 1024)

            for scanner.Scan() {
                seqName  := scanner.Text()
                scanner.Scan()
                seqValue := scanner.Text()

                sequences = append(sequences, SequenceEntry{
                    Name:     seqName,
                    Sequence: seqValue,
                })
            }

            err = file.Close()

            if err != nil {
                panic(err)
            }
        }
    }

    return sequences
}

// Converts sequences to database entries.
// Sequence names are set to "?".
func sequencesToEntries(seqs []string) []SequenceEntry {
    entries := make([]SequenceEntry, len(seqs))

    for i, seq := range seqs {
        entries[i] = SequenceEntry{
            Name: "?",
            Sequence: seq,
        }
    }

    return entries
}

// Debug function. Creates sequence database by given array of sequences.
func DbBySequences(seqs []string) SequenceDb {
    entries := SequenceDb(sequencesToEntries(seqs))
    return entries
}
