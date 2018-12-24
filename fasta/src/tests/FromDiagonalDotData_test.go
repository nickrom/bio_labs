package tests

import "testing"

import (
. "../structs"
)


func TestFormDiagonalDotData1(t *testing.T) {
	inputSeq := "AAAAA"
	baseSeq := "AAAAA"
	var n1 uint = 4
	var n2 uint = 0
	var diag1  = 4 //центральная диагональ
	var diag2  = 0 //самая нижняя диагональ
	dddAlloc := &DiagonalDotData{}

	s1Dots := BuildSeqDotDataFor(inputSeq)

	FormDiagonalDotData(dddAlloc, s1Dots, baseSeq, len(inputSeq))

	if n1 != dddAlloc.Data[diag1] {
		t.Error("Expected ", n1, ", got ", dddAlloc.Data[diag1])
	}

	if n2 != dddAlloc.Data[diag2] {
		t.Error("Expected ", n2, ", got ", dddAlloc.Data[diag2])
	}
}

func TestFormDiagonalDotData2(t *testing.T) {
	inputSeq := "TAGTGTCA"
	baseSeq := "ACGTCATCA"
	var n1 uint = 1
	var n2 uint = 3
	var n3 uint = 2
	var diag1  = 7 //центральная диагональ
	var diag2  = 5 //
	var diag3  = 8 //

	dddAlloc := &DiagonalDotData{}

	s1Dots := BuildSeqDotDataFor(inputSeq)

	FormDiagonalDotData(dddAlloc, s1Dots, baseSeq, len(inputSeq))

	if n1 != dddAlloc.Data[diag1] {
		t.Error("Expected ", n1, ", got ", dddAlloc.Data[diag1])
	}

	if n2 != dddAlloc.Data[diag2] {
		t.Error("Expected ", n2, ", got ", dddAlloc.Data[diag2])
	}

	if n3 != dddAlloc.Data[diag3] {
		t.Error("Expected ", n3, ", got ", dddAlloc.Data[diag3])
	}
}




