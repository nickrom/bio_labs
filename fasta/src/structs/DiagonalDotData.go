package structs

import (
    "../util"
)

const MaxSequenceLength = 2 << 17

/* --- */

// Data from DotMatrix about matched diagonal dots of SeqPair.
// Structure implicitly relates to some SeqPair and their SeqDotData's.
// For every diagonal (determined by offset) number of matches is provided;
// StartOffset keeps offset of array, i.e. Data indices start from 0,
// but diagonal offsets starts from StartOffset.
type DiagonalDotData struct {
    Data        [MaxSequenceLength]uint
    StartOffset int
    Length      int
}

/* --- */

// Calculates DiagonalDotData by given SeqDotData of each sequence S1 and S2.
// S1 Length are used for determining array size and start diagonal offset.
// Result is written into given dddRef DiagonalDotData reference for memory issues.
func FormDiagonalDotData(dddRef *DiagonalDotData, s1Dots SeqDotData, s2 string, s1Len int) {

    // Initialize Data

    s2Len := len(s2)
    startOffset := -(s1Len - 1)

    dddRef.Length = s1Len + s2Len - 1
    dddRef.StartOffset = startOffset

    for i := 0; i < dddRef.Length; i += 1 {
        dddRef.Data[i] = 0
    }

    // Fill data

    for s2Ind := 0; s2Ind < s2Len - 1; s2Ind += 1 {
        key    := util.CombineSymbolPair(s2[s2Ind], s2[s2Ind + 1])
        value  := s1Dots[key]
        length := len(value)

        if length > 0 {
            dddRef.Data[s2Ind - value[0] - startOffset] += 1
            if length > 1 {
                dddRef.Data[s2Ind - value[1] - startOffset] += 1
                if length > 2 {
                    dddRef.Data[s2Ind - value[2] - startOffset] += 1
                    if length > 3 {
                        dddRef.Data[s2Ind - value[3] - startOffset] += 1
                        if length > 4 {
                            dddRef.Data[s2Ind - value[4] - startOffset] += 1
                            if length > 5 {
                                for s1Ind := 5; s1Ind < len(value); s1Ind += 1 {
                                   dddRef.Data[s2Ind - s1Ind - startOffset] += 1
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}

// Selects <amount> best (by dot match number) diagonals.
// Uses BinaryHeap for keeping best entries.
func (ddd *DiagonalDotData) SelectBestDiagonals(amount int, dotMatchCutOff uint) []Diagonal {

    // Store array of best values (and indices) with naive traverse and update.

    bestValues := make([]int, amount + 1)
    bestIndices := make([]Diagonal, amount)

    for i := range bestValues {
        bestValues[i] = -1
    }

    bestValues[amount] = 1000000

    for i, j := 0, 0; i < ddd.Length; i += 1 {
        if ddd.Data[i] < dotMatchCutOff {
            continue
        }

        j = 0
        for ; bestValues[j] < int(ddd.Data[i]); j += 1 {}
        j -= 1

        if j >= 0 {
            for k := 0; k < j; k += 1 {
                bestValues[k]  = bestValues[k + 1]
                bestIndices[k] = bestIndices[k + 1]
            }

            bestValues[j]  = int(ddd.Data[i])
            bestIndices[j] = Diagonal(i + ddd.StartOffset)
        }
    }

    firstReal := 0

    for ; firstReal < amount && bestValues[firstReal] == -1; firstReal += 1 {}

    // Return result

    return bestIndices[firstReal:]
}
