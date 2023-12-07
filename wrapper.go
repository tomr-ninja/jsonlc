package jsonlc

import (
	"encoding/json"
	"hash/maphash"
	"sync"
	"unsafe"
)

var (
	valueByHash = make(map[uint64]unsafe.Pointer)
	mux         = new(sync.RWMutex)
	hashSeed    = maphash.MakeSeed()
)

func FromValue[T any](v T) LowCardinality[T] {
	data, _ := json.Marshal(v)
	h := hash(data)
	valueByHash[h] = unsafe.Pointer(&v)

	return LowCardinality[T]{
		value: &v,
	}
}

type LowCardinality[T any] struct {
	value *T
}

func (v *LowCardinality[T]) Pointer() *T {
	return v.value
}

func (v *LowCardinality[T]) Value() T {
	return *v.value
}

func (v *LowCardinality[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *LowCardinality[T]) UnmarshalJSON(data []byte) error {
	h := hash(data)

	mux.RLock()
	savedPtr, ok := valueByHash[h]
	mux.RUnlock()
	if ok {
		v.value = (*T)(savedPtr)

		return nil
	}

	var (
		newValue T
		err      error
	)
	if err = json.Unmarshal(data, &newValue); err != nil {
		return err
	}
	v.value = &newValue

	mux.Lock()
	valueByHash[h] = unsafe.Pointer(&newValue)
	mux.Unlock()

	return nil
}

func hash(v []byte) uint64 {
	var h maphash.Hash
	h.SetSeed(hashSeed)
	_, _ = h.Write(v)

	return h.Sum64()
}
