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
	shortString := "golang"
	longString := strings.Repeat("golang", 300)
	extremelyLongString := strings.Repeat(longString, 100)

	dataWithShortString := []byte(`{"field_one": 1, "field_two": "` + shortString + `"}`)
	dataWithLongString := []byte(`{"field_one": 1, "field_two": "` + longString + `"}`)
	dataWithExtremelyLongString := []byte(`{"field_one": 1, "field_two": "` + extremelyLongString + `"}`)

	benchmarkStructs := func(b *testing.B, data []byte) {
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

	b.Run("short_string", func(b *testing.B) {
		benchmarkStructs(b, dataWithShortString)
	})
	b.Run("long_string", func(b *testing.B) {
		benchmarkStructs(b, dataWithLongString)
	})
	b.Run("extremely_long_string", func(b *testing.B) {
		benchmarkStructs(b, dataWithExtremelyLongString)
	})
}
