package structs

import "../db/structs"

// Program input data. Principally stores data for FASTA algorithm.
type InputBundle struct {
    TargetSequence  structs.SequenceEntry
    TargetSeqDots   SeqDotData
    WeightMat       *WeightMatrix
    GapPenalty      int
    DiagFilterNum   int
    DotMatchCutOff  uint
    CutOff          int
    GraphMaxDistErr int
    StripExtraWidth int
    BestMatchNum    int
    DisplayAlign    bool
}
