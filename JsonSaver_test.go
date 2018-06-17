package main

import (
	"testing"
)

type Teststruct struct {
	A string
	B int
	C bool
}

func TestJsonSaver_SaveMem(t *testing.T) {
	js := NewJsonSaver()
	cases := []struct {
		Testin  Teststruct
		Testout string
	}{
		{Teststruct{"233", 1, true}, `{"A":"233","B":1,"C":true}`},
		{Teststruct{"", 0, true}, `{"A":"","B":0,"C":true}`},
		{Teststruct{"abcd", -233, false}, `{"A":"abcd","B":-233,"C":false}`},
		{Teststruct{"p", -10000, false}, `{"A":"p","B":-10000,"C":false}`},
	}

	for _, v := range cases {
		if calcout := js.SaveMem(v.Testin); v.Testout != calcout {
			t.Errorf("Error on: %v != %v", v.Testout, calcout)
		}
	}
	// extra test

	var tin1 []Teststruct
	tin1 = make([]Teststruct, 0)
	tin1 = append(tin1, Teststruct{"2", 1, true})
	tin1 = append(tin1, Teststruct{"2", 1, true})
	tout1 := `[{"A":"2","B":1,"C":true},{"A":"2","B":1,"C":true}]`
	if calcout := js.SaveMem(tin1); tout1 != calcout {
		t.Errorf("Error on: %v != %v", tout1, calcout)
	}

	type SpecialForJson struct {
		ABC int `json:"hello"`
	}

	tin2 := SpecialForJson{1}
	tout2 := `{"hello":1}`
	if calcout := js.SaveMem(tin2); tout2 != calcout {
		t.Errorf("Error on: %v != %v", tout2, calcout)
	}
}

func TestJsonSaver_SaveLoadMem(t *testing.T) {
	js := NewJsonSaver()
	cases := []struct {
		Testin  Teststruct
		Testout Teststruct
	}{
		{Teststruct{"233", 1, true}, Teststruct{"233", 1, true}},
		{Teststruct{"", 0, true}, Teststruct{"", 0, true}},
		{Teststruct{"abcd", -233, false}, Teststruct{"abcd", -233, false}},
		{Teststruct{"p", -10000, false}, Teststruct{"p", -10000, false}},
	}
	for _, v := range cases {
		js.SaveMem(v.Testin)
		var calcout Teststruct
		js.LoadMem(&calcout)
		if v.Testout != calcout {
			t.Errorf("Error on: %v != %v", v.Testout, calcout)
		}
	}
}

func TestJsonSaver_SaveLoadFile(t *testing.T) {
	js := NewJsonSaver()
	const testfilename = "testdir/TFN"
	cases := []struct {
		Testin  string
		Testout string
	}{
		{"aabbcc", "aabbcc"},
		{"", ""},
		{"a b", "a b"},
		{" n d e", " n d e"},
		{"abcd\nabcd", "abcd\nabcd"},
	}
	for _, v := range cases {
		js.SaveFile(testfilename, v.Testin)
		if calcout := js.LoadFile(testfilename); v.Testout != calcout {
			t.Errorf("Error on: %v != %v", v.Testout, calcout)
		}
	}
}

func TestJsonSaver_SaveLoad(t *testing.T) {
	js := NewJsonSaver()
	const testfilename = "testdir/TFN"
	cases := []struct {
		Testin  Teststruct
		Testout Teststruct
	}{
		{Teststruct{"233", 1, true}, Teststruct{"233", 1, true}},
		{Teststruct{"", 0, true}, Teststruct{"", 0, true}},
		{Teststruct{"abcd", -233, false}, Teststruct{"abcd", -233, false}},
		{Teststruct{"p", -10000, false}, Teststruct{"p", -10000, false}},
	}
	for _, v := range cases {
		js.Save(testfilename, v.Testin)
		var calcout Teststruct
		js.Load(testfilename, &calcout)
		if v.Testout != calcout {
			t.Errorf("Error on: %v != %v", v.Testout, calcout)
		}
	}
}
