package structs

/* --- */

// SeqPair diagonal segment with given diagonal offset and trimming bounds P1 and P2.
// Structure implicitly relates to some SeqPair.
// P1 <= P2;
// Score stores sum of symbols pairs' weights according to some WeightMatrix it was built with.
type Segment struct {
    Diag   Diagonal
    P1, P2 int
    Score  int
}

/* --- */

const EmpiricalFilteredCapacity = 4

// Minimal of two diagonals
func MinOffset(d1, d2 Diagonal) Diagonal {
    if d1 < d2 { return d1 } else { return d2 }
}

// Maximal of two diagonals
func MaxOffset(d1, d2 Diagonal) Diagonal {
    if d1 > d2 { return d1 } else { return d2 }
}

// Filters segments by predicate (Score >= cutOff).
func FilterByCutOff(segs []Segment, cutOff int) []Segment {
    result := make([]Segment, 0, EmpiricalFilteredCapacity)

    for _, seg := range segs {
        if seg.Score >= cutOff {
            result = append(result, seg)
        }
    }

    return result
}

// Calculates and returns bounding strip as pair of left and right bounding diagonals.
// Strip contains all given segments.
// Strip width is equal to maximal difference between any two segments plus extraWidth * 2.
func GetStripOf(segs []Segment, extraWidth int, seqPair *SeqPair) Strip {

    // Note that widening left diagonal means decreasing offset, and
    // widening right diagonal means increasing offset.
    leftDiag, rightDiag := segs[0].Diag, segs[0].Diag

    for _, seg := range segs {
        leftDiag  = MinOffset(seg.Diag, leftDiag)
        rightDiag = MaxOffset(seg.Diag, rightDiag)
    }

    // To insure lack of 'out of bounds' clamps both offsets.
    return Strip {
        MaxOffset(leftDiag  - Diagonal(extraWidth), -Diagonal(len(seqPair.S1) - 1)),
        MinOffset(rightDiag + Diagonal(extraWidth),  Diagonal(len(seqPair.S2) - 1)),
    }
}

// Returns segment start point according to its diagonal offset and P1.
func (s *Segment) GetStartPoint() Point {
    if s.Diag < 0 {
        return Point{ s.P1, -int(s.Diag) + s.P1 }
    } else {
        return Point{ int(s.Diag) + s.P1, s.P1 }
    }
}

// Returns segment end point according to its diagonal offset and P2.
func (s *Segment) GetEndPoint() Point {
    if s.Diag < 0 {
        return Point{ s.P2, -int(s.Diag) + s.P2 }
    } else {
        return Point{ int(s.Diag) + s.P2, s.P2 }
    }
}
