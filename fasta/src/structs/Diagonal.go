package structs

import "../util"

/* --- */

// SeqPair diagonal with given offset.
// Structure implicitly relates to some SeqPair.
// offset > 0 => offset goes down (by S1);
// offset < 0 => offset goes right (by S2).
type Diagonal int

/* --- */

// Trims diagonal by finding best (with greatest Score) diagonal segment.
// Returns Segment and its Score.
func (d Diagonal) TrimToBestSegment(seqPair *SeqPair, wm *WeightMatrix) Segment {

    // Initialize minimal and maximal bounds of segment
    p1, p2 := 0, 0
    // Initialize Score accumulation and max Score variables
    score, maxScore := 0, 0
    // Initialize extra variables
    lastZero := 0
    diagLen  := seqPair.GetDiagonalLength(d)
    startRow := util.MaxInt(0, -int(d))
    startCol := util.MaxInt(0,  int(d))

    for i := 0; i < diagLen; i += 1 {
        score += (*wm)[(uint16(seqPair.S1[startRow + i]) << 8) | uint16(seqPair.S2[startCol + i])]
        //score += seqPair.WeightIn(wm, startRow + i, startCol + i) /* Slower */

        if score <= 0 {
            score = 0
            lastZero = i
        }

        if maxScore < score {
            maxScore = score
            p1 = lastZero
            p2 = i
        }
    }

    return Segment { d, p1, p2 , maxScore }
}

func TrimToBestSegments(
    ddd *DiagonalDotData, diags []Diagonal, seqPair *SeqPair, wm *WeightMatrix, dotMatchCutOff uint) []Segment {

    diagNum := len(diags)
    result  := make([]Segment, diagNum)

    for i := 0; i < diagNum; i += 1 {
        if ddd.Data[int(diags[i]) - ddd.StartOffset] >= dotMatchCutOff {
            result[i] = diags[i].TrimToBestSegment(seqPair, wm)
        } else {
            result[i] = Segment { diags[i], 0, 0, 0 }
        }
    }

    return result
}
