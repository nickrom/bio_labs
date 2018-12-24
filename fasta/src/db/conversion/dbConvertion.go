package conversion

import (
    "../../util"
    . "../structs"
    "bufio"
    "fmt"
    "os"
    "runtime"
)

const ClusterNumber   = 8
const ClusterFileName = "DbCluster"

// Converts raw database into structured for program needs.
func ConvertRawDatabase(rawDbPath, convertedDirPath string) {
    sequences := ReadSequencesFromFile(rawDbPath)

    // Write converted database

    clNum := ClusterNumber
    sequencesInCluster := len(sequences) / clNum

    if len(sequences) - len(sequences) % clNum > 0 {
        clNum += 1
    }

    cReady := make(chan int, clNum)
    runtime.GOMAXPROCS(clNum)

    for i := 0; i < clNum; i += 1 {
        go func(index int) {
            start := sequencesInCluster * index
            end   := util.MinInt(len(sequences), sequencesInCluster * (index + 1))

            clusterSequences := sequences[start : end]

            if len(clusterSequences) > 0 {
                createCluster(fmt.Sprintf("%s/%s%d.cl", convertedDirPath, ClusterFileName, index), clusterSequences)
            }

            cReady <- 1
        }(i)
    }

    for i := 0; i < clNum; i++ {
        <-cReady
    }
}

// Loads sequences from given input file with sequences of Fasta format.
func ReadSequencesFromFile(filePath string) []SequenceEntry {

    // Configure input file

    rawDbFile, err := os.Open(filePath)

    if err != nil {
        panic(err)
    }

    reader := bufio.NewReader(rawDbFile)

    // Read raw database

    sequences := make([]SequenceEntry, 0, SequenceInitCap)
    counter   := 0

    for {
        seq, ok := ReadFrom(reader, counter)

        if ok {
            sequences = append(sequences, *seq)
        } else {
            break
        }
        counter += 1
    }

    return sequences
}

// Creates one cluster of structured database sequence part.
func createCluster(clusterFilePath string, clusterSequences []SequenceEntry) {
    clusterFile, err := os.Create(clusterFilePath)

    if err != nil {
        panic(err)
    }

    writer := bufio.NewWriter(clusterFile)

    for _, sequence := range clusterSequences {
        sequence.WriteInto(writer)
    }

    err = writer.Flush()

    if err != nil {
        panic(err)
    }
}
