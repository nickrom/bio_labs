package structs

import (
    "../util"
)

/* --- */

// Pair of sequences S1 and S2 in string form.
// For definiteness, when applicable, S1 and S2 forms a table,
// where S1 goes down and S2 goes right.
type SeqPair struct {
    S1, S2, S2Name string
}

/* --- */

// Returns weight from WeightMatrix for given symbols in s1 and s2.
func (sp *SeqPair) WeightIn(wm *WeightMatrix, s1Ind, s2Ind int) int {
    return (*wm)[util.CombineSymbolPair(sp.S1[s1Ind], sp.S2[s2Ind])]
}

// Returns symbol number in given diagonal related to this SeqPair.
func (sp *SeqPair) GetDiagonalLength(d Diagonal) int {
    if d < 0 {
        return util.MinInt(len(sp.S2), len(sp.S1) + int(d))
    } else {
        return util.MinInt(len(sp.S1), len(sp.S2) - int(d))
    }
}
