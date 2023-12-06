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

func TestLowCardinality_MarshalJSON(t *testing.T) {
	v := &MyStructOptimized{
		FieldOne: 1,
		FieldTwo: FromValue("golang"),
	}

	data, err := json.Marshal(v)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected := `{"field_one":1,"field_two":"golang"}`
	if string(data) != expected {
		t.Errorf("unexpected value: %v", string(data))
	}
}

func TestLowCardinality_UnmarshalJSON(t *testing.T) {
	t.Run("string", func(t *testing.T) {
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
	})

	t.Run("struct_ref", func(t *testing.T) {
		type innerStruct struct {
			FieldOne int `json:"field_one"`
		}
		type testStruct struct {
			FieldOne int                          `json:"field_one"`
			FieldTwo LowCardinality[*innerStruct] `json:"field_two"`
		}

		data := []byte(`{"field_one": 1, "field_two": {"field_one": 1}}`)
		v := &testStruct{}
		if err := json.Unmarshal(data, v); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if v.FieldOne != 1 {
			t.Errorf("unexpected value: %v", v.FieldOne)
		}
		if v.FieldTwo.Value().FieldOne != 1 {
			t.Errorf("unexpected value: %v", v.FieldTwo)
		}
	})
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	b.Run("string", func(b *testing.B) {
		shortString := "golang"
		longString := strings.Repeat("golang", 300)
		extremelyLongString := strings.Repeat(longString, 100)

		dataWithShortString := []byte(`{"field_one": 1, "field_two": "` + shortString + `"}`)
		dataWithLongString := []byte(`{"field_one": 1, "field_two": "` + longString + `"}`)
		dataWithExtremelyLongString := []byte(`{"field_one": 1, "field_two": "` + extremelyLongString + `"}`)

		benchmarkStructs := func(b *testing.B, data []byte) {
			b.Run("standard", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					_ = json.Unmarshal(data, &MyStruct{})
				}
			})

			b.Run("optimized", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					_ = json.Unmarshal(data, &MyStructOptimized{})
				}
			})
		}

		b.Run("short", func(b *testing.B) {
			benchmarkStructs(b, dataWithShortString)
		})
		b.Run("long", func(b *testing.B) {
			benchmarkStructs(b, dataWithLongString)
		})
		b.Run("extremely_long", func(b *testing.B) {
			benchmarkStructs(b, dataWithExtremelyLongString)
		})
	})

	b.Run("struct", func(b *testing.B) {
		type (
			// all 4 are ints, so relatively small types
			// but there are 4 of them, so optimization should be noticeable
			innerStruct struct {
				FieldOne   int `json:"field_one"`
				FieldTwo   int `json:"field_two"`
				FieldThree int `json:"field_three"`
				FieldFour  int `json:"field_four"`
			}
			nonOptStruct struct {
				FieldOne int         `json:"field_one"`
				FieldTwo innerStruct `json:"field_two"`
			}
			optStruct struct {
				FieldOne int                         `json:"field_one"`
				FieldTwo LowCardinality[innerStruct] `json:"field_two"`
			}
		)

		data := []byte(`{"field_one": 1, "field_two": {"field_one": 1, "field_two": 2, "field_three": 3, "field_four": 4}}`)

		b.Run("standard", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = json.Unmarshal(data, &nonOptStruct{})
			}
		})

		b.Run("optimized", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = json.Unmarshal(data, &optStruct{})
			}
		})
	})
}
