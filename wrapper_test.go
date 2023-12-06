package jsonlc

import (
	"encoding/json"
	"strings"
	"testing"
)

type MyStruct struct {
	FieldOne int    `json:"field_one"`
	FieldTwo string `json:"field_two"`
}

type MyStructOptimized struct {
	FieldOne int                    `json:"field_one"`
	FieldTwo LowCardinality[string] `json:"field_two"`
}

func TestLowCardinality_UnmarshalJSON(t *testing.T) {
	data := []byte(`{"field_one": 1, "field_two": "golang"}`)
	v := &MyStructOptimized{}

	if err := json.Unmarshal(data, v); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if v.FieldOne != 1 {
		t.Errorf("unexpected value: %v", v.FieldOne)
	}
	if v.FieldTwo.Value() != "golang" {
		t.Errorf("unexpected value: %v", v.FieldTwo)
	}

	// test that the value is not unmarshaled again\
	v2 := &MyStructOptimized{}
	if err := json.Unmarshal(data, v2); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if v2.FieldTwo.value != v.FieldTwo.value {
		t.Errorf("pointers should be equal")
	}
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	longString := strings.Repeat("abc", 300)
	data := []byte(`{"field_one": 1, "field_two": "` + longString + `"}`)

	b.Run("standard", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = json.Unmarshal(data, &MyStruct{})
		}
	})

	b.Run("optimized", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = json.Unmarshal(data, &MyStructOptimized{})
		}
	})
}
