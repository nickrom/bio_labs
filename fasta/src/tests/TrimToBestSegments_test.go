package tests

import "testing"
import (
	. "../structs"
)


func TestTrimToBestSegments1(t *testing.T) {
	seqPair := SeqPair {
		"ACTCXZKLMA", // string
		"ACTCAGKLMB", // string
	}

	var diag1 Diagonal = 0 //номер диагонали
	var start1 int = 0 //начало сегмента
	var end1 int  = 8 //конец сегмента
	var score1 = 39 //скор

	dddAlloc := &DiagonalDotData{}

	diagonals := make([]Diagonal, 1)
	diagonals[0] = Diagonal(0)

	segments := TrimToBestSegments(dddAlloc, diagonals, &seqPair, Blosum62(), 0)


	if diag1 != segments[0].Diag  {
		t.Error("Expected ", diag1, ", got ", segments[0].Diag)
	}

	if start1 != segments[0].P1 {
		t.Error("Expected ", start1, ", got ", segments[0].P1)
	}

	if end1 != segments[0].P2 {
		t.Error("Expected ", end1, ", got ", segments[0].P2)
	}

	if score1 != segments[0].Score {
		t.Error("Expected ", score1, ", got ", segments[0].Score)
	}
}

func TestTrimToBestSegments2(t *testing.T) {
	seqPair := SeqPair {
		"ACGTCATCA", // string
		"TAGTGTCA", // string
	}

	var diag1 Diagonal = 2 //номер диагонали
	var start1 int = 1 //начало сегмента
	var end1 int  = 5 //конец сегмента
	var score1 = 24 //скор

	dddAlloc := &DiagonalDotData{}

	diagonals := make([]Diagonal, 1)
	diagonals[0] = Diagonal(2)

	segments := TrimToBestSegments(dddAlloc, diagonals, &seqPair, Blosum62(), 0)


	if diag1 != segments[0].Diag  {
		t.Error("Expected ", diag1, ", got ", segments[0].Diag)
	}

	if start1 != segments[0].P1 {
		t.Error("Expected ", start1, ", got ", segments[0].P1)
	}

	if end1 != segments[0].P2 {
		t.Error("Expected ", end1, ", got ", segments[0].P2)
	}

	if score1 != segments[0].Score {
		t.Error("Expected ", score1, ", got ", segments[0].Score)
	}
}