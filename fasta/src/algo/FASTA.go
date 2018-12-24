package algo

import (
    . "../db/structs"
    . "../structs"
)

/* --- */

// Data returned by FASTA algorithm
type FastaResult []FastaResultEntry

// Core algorithm for calculating best (by alignment score) sequence matches.
// Given input sequence input.TargetSequence (with DotData) and sequence database,
// the task is to find several best database sequences, i.e. whose alignment has the greatest score.
func FASTA(dbStartOffset int, input *InputBundle, db SequenceDb) FastaResult {
    entryNum        := len(db)
    bestResultsHeap := NewHeapOfSize(input.BestMatchNum)

    // Allocate space for methods of some algorithm steps
    dddAlloc := &DiagonalDotData{}

    for i := 0; i < entryNum; i += 1 {
        alignResult := fastaEntry(db[i], input, dddAlloc)
        bestResultsHeap.Update(alignResult.Score, &FastaResultEntry {
            DbSequenceIndex: dbStartOffset + i,
            Score: alignResult.Score,
        })
    }

    //fmt.Printf("Debug:\n" +
    //   "\tDiagonalDD    %.3f (%d)\n" +
    //   "\tBestDiagonals %.3f (%d)\n" +
    //   "\tMakeSegments  %.3f (%d)\n" +
    //   "\tCutOffFilter  %.3f (%d)\n" +
    //   "\tFloydWarshall %.3f (%d)\n" +
    //   "\tFormStrip     %.3f (%d)\n" +
    //   "\tSmithWaterman %.3f (%d)\n" +
    //   "%v passed to SW\n",
    //   float32(C1 / 1000000) / 1000, C1 / int64(entryNum),
    //   float32(C2 / 1000000) / 1000, C2 / int64(entryNum),
    //   float32(C3 / 1000000) / 1000, C3 / int64(entryNum),
    //   float32(C4 / 1000000) / 1000, C4 / int64(entryNum),
    //   float32(C5 / 1000000) / 1000, C5 / PassedCutOffNum,
    //   float32(C6 / 1000000) / 1000, C6 / PassedCutOffNum,
    //   float32(C7 / 1000000) / 1000, C7 / PassedCutOffNum,
    //   PassedCutOffNum - 1)

    unconvertedValues := bestResultsHeap.ExtractSorted()

    results := make([]FastaResultEntry, len(unconvertedValues))

    for i, value := range unconvertedValues {
        results[i] = *value.(*FastaResultEntry)
    }

    return results
}

/* --- */

//var C1 int64 = 0
//var C2 int64 = 0
//var C3 int64 = 0
//var C4 int64 = 0
//var C5 int64 = 0
//var C6 int64 = 0
//var C7 int64 = 0
//
//var PassedCutOffNum int64 = 1

// Data entry of FastaResult structure.
// Signifies one of the best sequence matches.
type FastaResultEntry struct {
    DbSequenceIndex int
    Score           int
    IsFull          bool
    CorrectedScore  int
    Align           string
}

// FASTA iteration with input sequence input.TargetSequence and DB sequence sDb.
// Can be parallelized: fastaEntry() works independently and can be executed in individual thread.
func fastaEntry(seqEntry SequenceEntry, input *InputBundle, dddAlloc *DiagonalDotData) *AlignResult {
    sDb := seqEntry.Sequence
    seqPair := SeqPair {
        S1: input.TargetSequence.Sequence,
        S2: sDb,
    }

    //t1 := util.CurTime()
    FormDiagonalDotData(dddAlloc, input.TargetSeqDots, seqPair.S2, len(seqPair.S1))
    //t2 := util.CurTime()
    //C1 += t2 - t1
    diags := dddAlloc.SelectBestDiagonals(input.DiagFilterNum, input.DotMatchCutOff)
    //t3 := util.CurTime()
    //C2 += t3 - t2
    segs := TrimToBestSegments(dddAlloc, diags, &seqPair, input.WeightMat, input.DotMatchCutOff)

    //t4 := util.CurTime()
    //C3 += t4 - t3
    filteredSegs := FilterByCutOff(segs, input.CutOff)
    //t5 := util.CurTime()
    //C4 += t5 - t4

    if len(filteredSegs) == 0 {
        return &AlignResult{ Score: 0 }
    }

    //PassedCutOffNum += 1
    _, graphFilteredSegs := FloydWarshall(filteredSegs, input.GapPenalty, input.GraphMaxDistErr)


    //t6 := util.CurTime()
    //C5 += t6 - t5
    strip := GetStripOf(graphFilteredSegs, input.StripExtraWidth, &seqPair)
    //t7 := util.CurTime()
    //C6 += t7 - t6


    //fmt.Printf("%s\n%s\n", seqPair.S2Name, seqPair.S2)
    //for _, seg := range graphFilteredSegs {
    //if (maxScore > 0) {
    //    fmt.Printf("%s\t%d\n", seqEntry.Name[1:6], maxScore)
    //}
    ////}
    //fmt.Println()

    score := SmithWatermanStrip(&seqPair, strip, input.WeightMat, input.GapPenalty)
    //t8 := util.CurTime()
    //C7 += t8 - t7

    return score
}
