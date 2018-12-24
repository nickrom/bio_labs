package main

import (
    "./algo"
    "./db/conversion"
    dbStructs "./db/structs"
    "./structs"
    "./util"
    "bufio"
    "fmt"
    "os"
    "runtime"
    "sort"
)

const GoRoutineNum = 128

func main() {
    if len(os.Args) < 4 {
        fmt.Println(
            "Usage: go run cli.go <converted DB clusters path> <target sequence path> <output path> [options]\n" +
            "\t--align - display aligns for found results")
        return
    }

    // Parse arguments

    dbDirPath  := os.Args[1]
    targetPath := os.Args[2]
    outputPath := os.Args[3]
    args       := os.Args[4:]

    displayAlign := findArgument(args, "--align")

    // Initialize some data

    weightMatrix := structs.Blosum62()

    sequence := conversion.ReadSequencesFromFile(targetPath)[0]

    sequenceDb := dbStructs.FromClusters(dbDirPath)

    fmt.Println("DB sequences were successfully loaded")

    // Fill input bundle

    input := structs.InputBundle {
        TargetSequence: sequence,
        TargetSeqDots: structs.BuildSeqDotDataFor(sequence.Sequence),
        WeightMat: weightMatrix,
        GapPenalty: -2,
        DiagFilterNum: 20,
        DotMatchCutOff: 100,
        CutOff: 300,
        GraphMaxDistErr: 20,
        StripExtraWidth: 32,
        BestMatchNum: 1000,
        DisplayAlign: displayAlign,
    }

    // Call FASTA algorithm

    goRoutineNum := util.MinInt(len(sequenceDb), GoRoutineNum)
    cReady  := make(chan *algo.FastaResult, goRoutineNum)
    results := make([]*algo.FastaResult, 0, goRoutineNum)
    chunk   := len(sequenceDb) / goRoutineNum

    runtime.GOMAXPROCS(goRoutineNum)

    t1 := util.CurTime()

    for i := 0; i < goRoutineNum; i++ {
        go func(start int) {
            end := start + chunk

            if end > len(sequenceDb) {
                end = len(sequenceDb)
            }
            res := algo.FASTA(start, &input, sequenceDb[start : end])

            cReady <- &res
        }(i * chunk)
    }

    for i := 0; i < goRoutineNum; i++ {
        res := <-cReady
        results = append(results, res)
    }

    // Union results from all channels, collecting only best matches

    bestResEntries := make([]algo.FastaResultEntry, 0, input.BestMatchNum * goRoutineNum)

    for _, result := range results {
        for _, entry := range *result {
            bestResEntries = append(bestResEntries, entry)
        }
    }

    sort.Slice(bestResEntries, func (i, j int) bool {
        return bestResEntries[i].Score > bestResEntries[j].Score
    })

    bestResEntries = bestResEntries[:util.MinInt(len(bestResEntries), input.BestMatchNum)]

    nontrivialEnd := 0

    for nontrivialEnd < len(bestResEntries) && bestResEntries[nontrivialEnd].Score > 0 {
        nontrivialEnd += 1
    }

    bestResEntries = bestResEntries[:nontrivialEnd]

    // Correct scores and recover aligns with full Smith-Waterman pass

    if input.DisplayAlign {
        for i, entry := range bestResEntries {
            seqPair := structs.SeqPair {
                S1: input.TargetSequence.Sequence,
                S2: sequenceDb[entry.DbSequenceIndex].Sequence,
            }

            fullResult := algo.SmithWatermanFull(&seqPair, input.WeightMat, input.GapPenalty)

            bestResEntries[i].IsFull         = true
            bestResEntries[i].CorrectedScore = fullResult.Score
            bestResEntries[i].Align          = fullResult.Align
        }
    }

    //

    t2 := util.CurTime()

    timeNano := t2 - t1

    // Print result

    clusterFile, _ := os.Create(outputPath)
    writer := bufio.NewWriter(clusterFile)

    fmt.Println("Input sequence:")
    fmt.Printf(">%s\n", sequence.Name)
    fmt.Println(sequence.Sequence)
    fmt.Println("Converted DB clusters directory path:")
    fmt.Println(dbDirPath)

    fmt.Println()
    fmt.Println("FASTA result:")

    for i, entry := range bestResEntries {
        dbSequence := sequenceDb[entry.DbSequenceIndex]

        //fmt.Println()
        //fmt.Printf(">%s\n", dbSequence.Name)
        //fmt.Println(dbSequence.Sequence)
        //fmt.Printf(util.Colorify("Score: %d\n", util.ColorGreen), entry.Score)
        fmt.Fprintf(writer, "%d %s %d\n", i + 1, dbSequence.Name, entry.Score)

        if entry.IsFull {
            fmt.Printf("Protein: %s\n", dbSequence.Name[1:6])
            fmt.Printf(util.Colorify("Corrected score: %d\n", util.ColorLightGreen), entry.CorrectedScore)
            fmt.Printf("Align:\n%s\n", entry.Align)
        }
    }

    writer.Flush()

    fmt.Printf("\nTotal time: %.3f sec\n", float64(timeNano / 1000000) / 1000)
}


// Finds given command line interface input argument from slice.
func findArgument(args []string, argName string) bool {
    for _, arg := range args {
        if arg == argName {
            return true
        }
    }

    return false
}