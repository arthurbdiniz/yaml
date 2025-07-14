/*
Copyright 2025 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kyaml

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"sigs.k8s.io/randfill"
	"sigs.k8s.io/yaml"
)

// Return the input string with the longest common leading whitespace removed
// from each line.
func dedent(in string) string {
	lines := strings.Split(in, "\n")
	pfx := ""
	for i, line := range lines {
		if i == 0 {
			pfx = leadingWhitespace(line)
			continue
		}
		if strings.HasPrefix(line, pfx) {
			continue
		}
		pfx = commonPrefix(pfx, line)
	}
	if pfx == "" {
		return in
	}
	for i, line := range lines {
		lines[i] = strings.TrimPrefix(line, pfx)
	}
	return strings.Join(lines, "\n")
}

func leadingWhitespace(s string) string {
	if len(s) == 0 {
		return ""
	}
	for i, r := range s {
		if r != ' ' && r != '\t' {
			return s[:i]
		}
	}
	return s // all whitespace
}

func commonPrefix(a, b string) string {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}
	for i := 0; i < minLen; i++ {
		if a[i] != b[i] {
			return a[:i]
		}
	}
	return a[:minLen]
}

// The following types are intended to provide robust coverage for testing.

type AllTypesStruct struct {
	SimpleStruct1 `json:","`      // embedded unnamed
	SimpleStruct2 `json:"named,"` // embedded named

	// Omitted fields (`json:"-"`) are handled in a different test
	NamedDash SimpleStruct4 `json:"-,"`

	Plain            PlainStruct            `json:"plain"`
	OmitEmpty        OmitEmptyStruct        `json:"omitEmpty"`
	OmitZero         OmitZeroStruct         `json:"omitZero"`
	OmitEmptyAndZero OmitEmptyAndZeroStruct `json:"omitEmptyAndZero"`
}

type SimpleStruct struct {
	String string `json:"string"`
	Bool   bool   `json:"bool"`
	Int    int    `json:"int"`
}
type SimpleStruct1 SimpleStruct
type SimpleStruct2 SimpleStruct
type SimpleStruct3 SimpleStruct
type SimpleStruct4 SimpleStruct

type PlainStruct struct {
	// Basic types
	String  string    `json:"string"`
	Bool    bool      `json:"bool"`
	Int     int       `json:"int"`
	Int8    int8      `json:"int8"`
	Int16   int16     `json:"int16"`
	Int32   int32     `json:"int32"`
	Int64   int64     `json:"int64"`
	Uint    uint      `json:"uint"`
	Uint8   uint8     `json:"uint8"`
	Uint16  uint16    `json:"uint16"`
	Uint32  uint32    `json:"uint32"`
	Uint64  uint64    `json:"uint64"`
	Float32 float32   `json:"float32"`
	Float64 float64   `json:"float64"`
	Time    time.Time `json:"time"`
	Bytes   []byte    `json:"bytes"`

	// Pointers to basic types
	StringPtr  *string    `json:"stringPtr"`
	BoolPtr    *bool      `json:"boolPtr"`
	IntPtr     *int       `json:"intPtr"`
	Int8Ptr    *int8      `json:"int8Ptr"`
	Int16Ptr   *int16     `json:"int16Ptr"`
	Int32Ptr   *int32     `json:"int32Ptr"`
	Int64Ptr   *int64     `json:"int64Ptr"`
	UintPtr    *uint      `json:"uintPtr"`
	Uint8Ptr   *uint8     `json:"uint8Ptr"`
	Uint16Ptr  *uint16    `json:"uint16Ptr"`
	Uint32Ptr  *uint32    `json:"uint32Ptr"`
	Uint64Ptr  *uint64    `json:"uint64Ptr"`
	Float32Ptr *float32   `json:"float32Ptr"`
	Float64Ptr *float64   `json:"float64Ptr"`
	TimePtr    *time.Time `json:"timePtr,omitempty"` // Time implements json.Marshaler, but on a value receiver

	// Slices of basic types
	StringSlice  []string    `json:"stringSlice"`
	BoolSlice    []bool      `json:"boolSlice"`
	IntSlice     []int       `json:"intSlice"`
	Int8Slice    []int8      `json:"int8Slice"`
	Int16Slice   []int16     `json:"int16Slice"`
	Int32Slice   []int32     `json:"int32Slice"`
	Int64Slice   []int64     `json:"int64Slice"`
	UintSlice    []uint      `json:"uintSlice"`
	Uint8Slice   []uint8     `json:"uint8Slice"`
	Uint16Slice  []uint16    `json:"uint16Slice"`
	Uint32Slice  []uint32    `json:"uint32Slice"`
	Uint64Slice  []uint64    `json:"uint64Slice"`
	Float32Slice []float32   `json:"float32Slice"`
	Float64Slice []float64   `json:"float64Slice"`
	TimeSlice    []time.Time `json:"timeSlice"`

	// Maps of string to basic types
	StringMap  map[string]string    `json:"stringMap"`
	BoolMap    map[string]bool      `json:"boolMap"`
	IntMap     map[string]int       `json:"intMap"`
	Int8Map    map[string]int8      `json:"int8Map"`
	Int16Map   map[string]int16     `json:"int16Map"`
	Int32Map   map[string]int32     `json:"int32Map"`
	Int64Map   map[string]int64     `json:"int64Map"`
	UintMap    map[string]uint      `json:"uintMap"`
	Uint8Map   map[string]uint8     `json:"uint8Map"`
	Uint16Map  map[string]uint16    `json:"uint16Map"`
	Uint32Map  map[string]uint32    `json:"uint32Map"`
	Uint64Map  map[string]uint64    `json:"uint64Map"`
	Float32Map map[string]float32   `json:"float32Map"`
	Float64Map map[string]float64   `json:"float64Map"`
	TimeMap    map[string]time.Time `json:"timeMap"`

	// Slice of slices
	StringSliceSlice  [][]string    `json:"stringSliceSlice"`
	BoolSliceSlice    [][]bool      `json:"boolSliceSlice"`
	IntSliceSlice     [][]int       `json:"intSliceSlice"`
	Int8SliceSlice    [][]int8      `json:"int8SliceSlice"`
	Int16SliceSlice   [][]int16     `json:"int16SliceSlice"`
	Int32SliceSlice   [][]int32     `json:"int32SliceSlice"`
	Int64SliceSlice   [][]int64     `json:"int64SliceSlice"`
	UintSliceSlice    [][]uint      `json:"uintSliceSlice"`
	Uint8SliceSlice   [][]uint8     `json:"uint8SliceSlice"`
	Uint16SliceSlice  [][]uint16    `json:"uint16SliceSlice"`
	Uint32SliceSlice  [][]uint32    `json:"uint32SliceSlice"`
	Uint64SliceSlice  [][]uint64    `json:"uint64SliceSlice"`
	Float32SliceSlice [][]float32   `json:"float32SliceSlice"`
	Float64SliceSlice [][]float64   `json:"float64SliceSlice"`
	TimeSliceSlice    [][]time.Time `json:"timeSliceSlice"`

	// Slice of maps
	StringSliceMap  []map[string]string    `json:"stringSliceMap"`
	BoolSliceMap    []map[string]bool      `json:"boolSliceMap"`
	IntSliceMap     []map[string]int       `json:"intSliceMap"`
	Int8SliceMap    []map[string]int8      `json:"int8SliceMap"`
	Int16SliceMap   []map[string]int16     `json:"int16SliceMap"`
	Int32SliceMap   []map[string]int32     `json:"int32SliceMap"`
	Int64SliceMap   []map[string]int64     `json:"int64SliceMap"`
	UintSliceMap    []map[string]uint      `json:"uintSliceMap"`
	Uint8SliceMap   []map[string]uint8     `json:"uint8SliceMap"`
	Uint16SliceMap  []map[string]uint16    `json:"uint16SliceMap"`
	Uint32SliceMap  []map[string]uint32    `json:"uint32SliceMap"`
	Uint64SliceMap  []map[string]uint64    `json:"uint64SliceMap"`
	Float32SliceMap []map[string]float32   `json:"float32SliceMap"`
	Float64SliceMap []map[string]float64   `json:"float64SliceMap"`
	TimeSliceMap    []map[string]time.Time `json:"timeSliceMap"`

	// Map of string to slices
	StringMapSlice  map[string][]string    `json:"stringMapSlice"`
	BoolMapSlice    map[string][]bool      `json:"boolMapSlice"`
	IntMapSlice     map[string][]int       `json:"intMapSlice"`
	Int8MapSlice    map[string][]int8      `json:"int8MapSlice"`
	Int16MapSlice   map[string][]int16     `json:"int16MapSlice"`
	Int32MapSlice   map[string][]int32     `json:"int32MapSlice"`
	Int64MapSlice   map[string][]int64     `json:"int64MapSlice"`
	UintMapSlice    map[string][]uint      `json:"uintMapSlice"`
	Uint8MapSlice   map[string][]uint8     `json:"uint8MapSlice"`
	Uint16MapSlice  map[string][]uint16    `json:"uint16MapSlice"`
	Uint32MapSlice  map[string][]uint32    `json:"uint32MapSlice"`
	Uint64MapSlice  map[string][]uint64    `json:"uint64MapSlice"`
	Float32MapSlice map[string][]float32   `json:"float32MapSlice"`
	Float64MapSlice map[string][]float64   `json:"float64MapSlice"`
	TimeMapSlice    map[string][]time.Time `json:"timeMapSlice"`

	// Map of string to maps
	StringMapMap  map[string]map[string]string    `json:"stringMapMap"`
	BoolMapMap    map[string]map[string]bool      `json:"boolMapMap"`
	IntMapMap     map[string]map[string]int       `json:"intMapMap"`
	Int8MapMap    map[string]map[string]int8      `json:"int8MapMap"`
	Int16MapMap   map[string]map[string]int16     `json:"int16MapMap"`
	Int32MapMap   map[string]map[string]int32     `json:"int32MapMap"`
	Int64MapMap   map[string]map[string]int64     `json:"int64MapMap"`
	UintMapMap    map[string]map[string]uint      `json:"uintMapMap"`
	Uint8MapMap   map[string]map[string]uint8     `json:"uint8MapMap"`
	Uint16MapMap  map[string]map[string]uint16    `json:"uint16MapMap"`
	Uint32MapMap  map[string]map[string]uint32    `json:"uint32MapMap"`
	Uint64MapMap  map[string]map[string]uint64    `json:"uint64MapMap"`
	Float32MapMap map[string]map[string]float32   `json:"float32MapMap"`
	Float64MapMap map[string]map[string]float64   `json:"float64MapMap"`
	TimeMapMap    map[string]map[string]time.Time `json:"timeMapMap"`

	// Recursive types
	Self           *PlainStruct                       `json:"self"`
	SelfSlice      []*PlainStruct                     `json:"selfSlice"`
	SelfMap        map[string]*PlainStruct            `json:"selfMap"`
	SelfSliceSlice [][](*PlainStruct)                 `json:"selfSliceSlice"`
	SelfSliceMap   []map[string]*PlainStruct          `json:"selfSliceMap"`
	SelfMapSlice   map[string][]*PlainStruct          `json:"selfMapSlice"`
	SelfMapMap     map[string]map[string]*PlainStruct `json:"selfMapMap"`
}

type OmitEmptyStruct struct {
	// Basic types
	String  string    `json:"string,omitempty"`
	Bool    bool      `json:"bool,omitempty"`
	Int     int       `json:"int,omitempty"`
	Int8    int8      `json:"int8,omitempty"`
	Int16   int16     `json:"int16,omitempty"`
	Int32   int32     `json:"int32,omitempty"`
	Int64   int64     `json:"int64,omitempty"`
	Uint    uint      `json:"uint,omitempty"`
	Uint8   uint8     `json:"uint8,omitempty"`
	Uint16  uint16    `json:"uint16,omitempty"`
	Uint32  uint32    `json:"uint32,omitempty"`
	Uint64  uint64    `json:"uint64,omitempty"`
	Float32 float32   `json:"float32,omitempty"`
	Float64 float64   `json:"float64,omitempty"`
	Time    time.Time `json:"time,omitempty"`
	Bytes   []byte    `json:"bytes,omitempty"`

	// Pointers to basic types
	StringPtr  *string    `json:"stringPtr,omitempty"`
	BoolPtr    *bool      `json:"boolPtr,omitempty"`
	IntPtr     *int       `json:"intPtr,omitempty"`
	Int8Ptr    *int8      `json:"int8Ptr,omitempty"`
	Int16Ptr   *int16     `json:"int16Ptr,omitempty"`
	Int32Ptr   *int32     `json:"int32Ptr,omitempty"`
	Int64Ptr   *int64     `json:"int64Ptr,omitempty"`
	UintPtr    *uint      `json:"uintPtr,omitempty"`
	Uint8Ptr   *uint8     `json:"uint8Ptr,omitempty"`
	Uint16Ptr  *uint16    `json:"uint16Ptr,omitempty"`
	Uint32Ptr  *uint32    `json:"uint32Ptr,omitempty"`
	Uint64Ptr  *uint64    `json:"uint64Ptr,omitempty"`
	Float32Ptr *float32   `json:"float32Ptr,omitempty"`
	Float64Ptr *float64   `json:"float64Ptr,omitempty"`
	TimePtr    *time.Time `json:"timePtr,omitempty"`

	// Slices of basic types
	StringSlice  []string    `json:"stringSlice,omitempty"`
	BoolSlice    []bool      `json:"boolSlice,omitempty"`
	IntSlice     []int       `json:"intSlice,omitempty"`
	Int8Slice    []int8      `json:"int8Slice,omitempty"`
	Int16Slice   []int16     `json:"int16Slice,omitempty"`
	Int32Slice   []int32     `json:"int32Slice,omitempty"`
	Int64Slice   []int64     `json:"int64Slice,omitempty"`
	UintSlice    []uint      `json:"uintSlice,omitempty"`
	Uint8Slice   []uint8     `json:"uint8Slice,omitempty"`
	Uint16Slice  []uint16    `json:"uint16Slice,omitempty"`
	Uint32Slice  []uint32    `json:"uint32Slice,omitempty"`
	Uint64Slice  []uint64    `json:"uint64Slice,omitempty"`
	Float32Slice []float32   `json:"float32Slice,omitempty"`
	Float64Slice []float64   `json:"float64Slice,omitempty"`
	TimeSlice    []time.Time `json:"timeSlice,omitempty"`

	// Maps of string to basic types
	StringMap  map[string]string    `json:"stringMap,omitempty"`
	BoolMap    map[string]bool      `json:"boolMap,omitempty"`
	IntMap     map[string]int       `json:"intMap,omitempty"`
	Int8Map    map[string]int8      `json:"int8Map,omitempty"`
	Int16Map   map[string]int16     `json:"int16Map,omitempty"`
	Int32Map   map[string]int32     `json:"int32Map,omitempty"`
	Int64Map   map[string]int64     `json:"int64Map,omitempty"`
	UintMap    map[string]uint      `json:"uintMap,omitempty"`
	Uint8Map   map[string]uint8     `json:"uint8Map,omitempty"`
	Uint16Map  map[string]uint16    `json:"uint16Map,omitempty"`
	Uint32Map  map[string]uint32    `json:"uint32Map,omitempty"`
	Uint64Map  map[string]uint64    `json:"uint64Map,omitempty"`
	Float32Map map[string]float32   `json:"float32Map,omitempty"`
	Float64Map map[string]float64   `json:"float64Map,omitempty"`
	TimeMap    map[string]time.Time `json:"timeMap,omitempty"`

	// Slice of slices
	StringSliceSlice  [][]string    `json:"stringSliceSlice,omitempty"`
	BoolSliceSlice    [][]bool      `json:"boolSliceSlice,omitempty"`
	IntSliceSlice     [][]int       `json:"intSliceSlice,omitempty"`
	Int8SliceSlice    [][]int8      `json:"int8SliceSlice,omitempty"`
	Int16SliceSlice   [][]int16     `json:"int16SliceSlice,omitempty"`
	Int32SliceSlice   [][]int32     `json:"int32SliceSlice,omitempty"`
	Int64SliceSlice   [][]int64     `json:"int64SliceSlice,omitempty"`
	UintSliceSlice    [][]uint      `json:"uintSliceSlice,omitempty"`
	Uint8SliceSlice   [][]uint8     `json:"uint8SliceSlice,omitempty"`
	Uint16SliceSlice  [][]uint16    `json:"uint16SliceSlice,omitempty"`
	Uint32SliceSlice  [][]uint32    `json:"uint32SliceSlice,omitempty"`
	Uint64SliceSlice  [][]uint64    `json:"uint64SliceSlice,omitempty"`
	Float32SliceSlice [][]float32   `json:"float32SliceSlice,omitempty"`
	Float64SliceSlice [][]float64   `json:"float64SliceSlice,omitempty"`
	TimeSliceSlice    [][]time.Time `json:"timeSliceSlice,omitempty"`

	// Slice of maps
	StringSliceMap  []map[string]string    `json:"stringSliceMap,omitempty"`
	BoolSliceMap    []map[string]bool      `json:"boolSliceMap,omitempty"`
	IntSliceMap     []map[string]int       `json:"intSliceMap,omitempty"`
	Int8SliceMap    []map[string]int8      `json:"int8SliceMap,omitempty"`
	Int16SliceMap   []map[string]int16     `json:"int16SliceMap,omitempty"`
	Int32SliceMap   []map[string]int32     `json:"int32SliceMap,omitempty"`
	Int64SliceMap   []map[string]int64     `json:"int64SliceMap,omitempty"`
	UintSliceMap    []map[string]uint      `json:"uintSliceMap,omitempty"`
	Uint8SliceMap   []map[string]uint8     `json:"uint8SliceMap,omitempty"`
	Uint16SliceMap  []map[string]uint16    `json:"uint16SliceMap,omitempty"`
	Uint32SliceMap  []map[string]uint32    `json:"uint32SliceMap,omitempty"`
	Uint64SliceMap  []map[string]uint64    `json:"uint64SliceMap,omitempty"`
	Float32SliceMap []map[string]float32   `json:"float32SliceMap,omitempty"`
	Float64SliceMap []map[string]float64   `json:"float64SliceMap,omitempty"`
	TimeSliceMap    []map[string]time.Time `json:"timeSliceMap,omitempty"`

	// Map of string to slices
	StringMapSlice  map[string][]string    `json:"stringMapSlice,omitempty"`
	BoolMapSlice    map[string][]bool      `json:"boolMapSlice,omitempty"`
	IntMapSlice     map[string][]int       `json:"intMapSlice,omitempty"`
	Int8MapSlice    map[string][]int8      `json:"int8MapSlice,omitempty"`
	Int16MapSlice   map[string][]int16     `json:"int16MapSlice,omitempty"`
	Int32MapSlice   map[string][]int32     `json:"int32MapSlice,omitempty"`
	Int64MapSlice   map[string][]int64     `json:"int64MapSlice,omitempty"`
	UintMapSlice    map[string][]uint      `json:"uintMapSlice,omitempty"`
	Uint8MapSlice   map[string][]uint8     `json:"uint8MapSlice,omitempty"`
	Uint16MapSlice  map[string][]uint16    `json:"uint16MapSlice,omitempty"`
	Uint32MapSlice  map[string][]uint32    `json:"uint32MapSlice,omitempty"`
	Uint64MapSlice  map[string][]uint64    `json:"uint64MapSlice,omitempty"`
	Float32MapSlice map[string][]float32   `json:"float32MapSlice,omitempty"`
	Float64MapSlice map[string][]float64   `json:"float64MapSlice,omitempty"`
	TimeMapSlice    map[string][]time.Time `json:"timeMapSlice,omitempty"`

	// Map of string to maps
	StringMapMap  map[string]map[string]string    `json:"stringMapMap,omitempty"`
	BoolMapMap    map[string]map[string]bool      `json:"boolMapMap,omitempty"`
	IntMapMap     map[string]map[string]int       `json:"intMapMap,omitempty"`
	Int8MapMap    map[string]map[string]int8      `json:"int8MapMap,omitempty"`
	Int16MapMap   map[string]map[string]int16     `json:"int16MapMap,omitempty"`
	Int32MapMap   map[string]map[string]int32     `json:"int32MapMap,omitempty"`
	Int64MapMap   map[string]map[string]int64     `json:"int64MapMap,omitempty"`
	UintMapMap    map[string]map[string]uint      `json:"uintMapMap,omitempty"`
	Uint8MapMap   map[string]map[string]uint8     `json:"uint8MapMap,omitempty"`
	Uint16MapMap  map[string]map[string]uint16    `json:"uint16MapMap,omitempty"`
	Uint32MapMap  map[string]map[string]uint32    `json:"uint32MapMap,omitempty"`
	Uint64MapMap  map[string]map[string]uint64    `json:"uint64MapMap,omitempty"`
	Float32MapMap map[string]map[string]float32   `json:"float32MapMap,omitempty"`
	Float64MapMap map[string]map[string]float64   `json:"float64MapMap,omitempty"`
	TimeMapMap    map[string]map[string]time.Time `json:"timeMapMap,omitempty"`

	// Recursive types
	Self           *OmitEmptyStruct                       `json:"self,omitempty"`
	SelfSlice      []*OmitEmptyStruct                     `json:"selfSlice,omitempty"`
	SelfMap        map[string]*OmitEmptyStruct            `json:"selfMap,omitempty"`
	SelfSliceSlice [][](*OmitEmptyStruct)                 `json:"selfSliceSlice,omitempty"`
	SelfSliceMap   []map[string]*OmitEmptyStruct          `json:"selfSliceMap,omitempty"`
	SelfMapSlice   map[string][]*OmitEmptyStruct          `json:"selfMapSlice,omitempty"`
	SelfMapMap     map[string]map[string]*OmitEmptyStruct `json:"selfMapMap,omitempty"`
}

type OmitZeroStruct struct {
	// Basic types
	String  string    `json:"string,omitzero"`
	Bool    bool      `json:"bool,omitzero"`
	Int     int       `json:"int,omitzero"`
	Int8    int8      `json:"int8,omitzero"`
	Int16   int16     `json:"int16,omitzero"`
	Int32   int32     `json:"int32,omitzero"`
	Int64   int64     `json:"int64,omitzero"`
	Uint    uint      `json:"uint,omitzero"`
	Uint8   uint8     `json:"uint8,omitzero"`
	Uint16  uint16    `json:"uint16,omitzero"`
	Uint32  uint32    `json:"uint32,omitzero"`
	Uint64  uint64    `json:"uint64,omitzero"`
	Float32 float32   `json:"float32,omitzero"`
	Float64 float64   `json:"float64,omitzero"`
	Time    time.Time `json:"time,omitzero"`
	Bytes   []byte    `json:"bytes,omitzero"`

	// Pointers to basic types
	StringPtr  *string    `json:"stringPtr,omitzero"`
	BoolPtr    *bool      `json:"boolPtr,omitzero"`
	IntPtr     *int       `json:"intPtr,omitzero"`
	Int8Ptr    *int8      `json:"int8Ptr,omitzero"`
	Int16Ptr   *int16     `json:"int16Ptr,omitzero"`
	Int32Ptr   *int32     `json:"int32Ptr,omitzero"`
	Int64Ptr   *int64     `json:"int64Ptr,omitzero"`
	UintPtr    *uint      `json:"uintPtr,omitzero"`
	Uint8Ptr   *uint8     `json:"uint8Ptr,omitzero"`
	Uint16Ptr  *uint16    `json:"uint16Ptr,omitzero"`
	Uint32Ptr  *uint32    `json:"uint32Ptr,omitzero"`
	Uint64Ptr  *uint64    `json:"uint64Ptr,omitzero"`
	Float32Ptr *float32   `json:"float32Ptr,omitzero"`
	Float64Ptr *float64   `json:"float64Ptr,omitzero"`
	TimePtr    *time.Time `json:"timePtr,omitzero"`

	// Slices of basic types
	StringSlice  []string    `json:"stringSlice,omitzero"`
	BoolSlice    []bool      `json:"boolSlice,omitzero"`
	IntSlice     []int       `json:"intSlice,omitzero"`
	Int8Slice    []int8      `json:"int8Slice,omitzero"`
	Int16Slice   []int16     `json:"int16Slice,omitzero"`
	Int32Slice   []int32     `json:"int32Slice,omitzero"`
	Int64Slice   []int64     `json:"int64Slice,omitzero"`
	UintSlice    []uint      `json:"uintSlice,omitzero"`
	Uint8Slice   []uint8     `json:"uint8Slice,omitzero"`
	Uint16Slice  []uint16    `json:"uint16Slice,omitzero"`
	Uint32Slice  []uint32    `json:"uint32Slice,omitzero"`
	Uint64Slice  []uint64    `json:"uint64Slice,omitzero"`
	Float32Slice []float32   `json:"float32Slice,omitzero"`
	Float64Slice []float64   `json:"float64Slice,omitzero"`
	TimeSlice    []time.Time `json:"timeSlice,omitzero"`

	// Maps of string to basic types
	StringMap  map[string]string    `json:"stringMap,omitzero"`
	BoolMap    map[string]bool      `json:"boolMap,omitzero"`
	IntMap     map[string]int       `json:"intMap,omitzero"`
	Int8Map    map[string]int8      `json:"int8Map,omitzero"`
	Int16Map   map[string]int16     `json:"int16Map,omitzero"`
	Int32Map   map[string]int32     `json:"int32Map,omitzero"`
	Int64Map   map[string]int64     `json:"int64Map,omitzero"`
	UintMap    map[string]uint      `json:"uintMap,omitzero"`
	Uint8Map   map[string]uint8     `json:"uint8Map,omitzero"`
	Uint16Map  map[string]uint16    `json:"uint16Map,omitzero"`
	Uint32Map  map[string]uint32    `json:"uint32Map,omitzero"`
	Uint64Map  map[string]uint64    `json:"uint64Map,omitzero"`
	Float32Map map[string]float32   `json:"float32Map,omitzero"`
	Float64Map map[string]float64   `json:"float64Map,omitzero"`
	TimeMap    map[string]time.Time `json:"timeMap,omitzero"`

	// Slice of slices
	StringSliceSlice  [][]string    `json:"stringSliceSlice,omitzero"`
	BoolSliceSlice    [][]bool      `json:"boolSliceSlice,omitzero"`
	IntSliceSlice     [][]int       `json:"intSliceSlice,omitzero"`
	Int8SliceSlice    [][]int8      `json:"int8SliceSlice,omitzero"`
	Int16SliceSlice   [][]int16     `json:"int16SliceSlice,omitzero"`
	Int32SliceSlice   [][]int32     `json:"int32SliceSlice,omitzero"`
	Int64SliceSlice   [][]int64     `json:"int64SliceSlice,omitzero"`
	UintSliceSlice    [][]uint      `json:"uintSliceSlice,omitzero"`
	Uint8SliceSlice   [][]uint8     `json:"uint8SliceSlice,omitzero"`
	Uint16SliceSlice  [][]uint16    `json:"uint16SliceSlice,omitzero"`
	Uint32SliceSlice  [][]uint32    `json:"uint32SliceSlice,omitzero"`
	Uint64SliceSlice  [][]uint64    `json:"uint64SliceSlice,omitzero"`
	Float32SliceSlice [][]float32   `json:"float32SliceSlice,omitzero"`
	Float64SliceSlice [][]float64   `json:"float64SliceSlice,omitzero"`
	TimeSliceSlice    [][]time.Time `json:"timeSliceSlice,omitzero"`

	// Slice of maps
	StringSliceMap  []map[string]string    `json:"stringSliceMap,omitzero"`
	BoolSliceMap    []map[string]bool      `json:"boolSliceMap,omitzero"`
	IntSliceMap     []map[string]int       `json:"intSliceMap,omitzero"`
	Int8SliceMap    []map[string]int8      `json:"int8SliceMap,omitzero"`
	Int16SliceMap   []map[string]int16     `json:"int16SliceMap,omitzero"`
	Int32SliceMap   []map[string]int32     `json:"int32SliceMap,omitzero"`
	Int64SliceMap   []map[string]int64     `json:"int64SliceMap,omitzero"`
	UintSliceMap    []map[string]uint      `json:"uintSliceMap,omitzero"`
	Uint8SliceMap   []map[string]uint8     `json:"uint8SliceMap,omitzero"`
	Uint16SliceMap  []map[string]uint16    `json:"uint16SliceMap,omitzero"`
	Uint32SliceMap  []map[string]uint32    `json:"uint32SliceMap,omitzero"`
	Uint64SliceMap  []map[string]uint64    `json:"uint64SliceMap,omitzero"`
	Float32SliceMap []map[string]float32   `json:"float32SliceMap,omitzero"`
	Float64SliceMap []map[string]float64   `json:"float64SliceMap,omitzero"`
	TimeSliceMap    []map[string]time.Time `json:"timeSliceMap,omitzero"`

	// Map of string to slices
	StringMapSlice  map[string][]string    `json:"stringMapSlice,omitzero"`
	BoolMapSlice    map[string][]bool      `json:"boolMapSlice,omitzero"`
	IntMapSlice     map[string][]int       `json:"intMapSlice,omitzero"`
	Int8MapSlice    map[string][]int8      `json:"int8MapSlice,omitzero"`
	Int16MapSlice   map[string][]int16     `json:"int16MapSlice,omitzero"`
	Int32MapSlice   map[string][]int32     `json:"int32MapSlice,omitzero"`
	Int64MapSlice   map[string][]int64     `json:"int64MapSlice,omitzero"`
	UintMapSlice    map[string][]uint      `json:"uintMapSlice,omitzero"`
	Uint8MapSlice   map[string][]uint8     `json:"uint8MapSlice,omitzero"`
	Uint16MapSlice  map[string][]uint16    `json:"uint16MapSlice,omitzero"`
	Uint32MapSlice  map[string][]uint32    `json:"uint32MapSlice,omitzero"`
	Uint64MapSlice  map[string][]uint64    `json:"uint64MapSlice,omitzero"`
	Float32MapSlice map[string][]float32   `json:"float32MapSlice,omitzero"`
	Float64MapSlice map[string][]float64   `json:"float64MapSlice,omitzero"`
	TimeMapSlice    map[string][]time.Time `json:"timeMapSlice,omitzero"`

	// Map of string to maps
	StringMapMap  map[string]map[string]string    `json:"stringMapMap,omitzero"`
	BoolMapMap    map[string]map[string]bool      `json:"boolMapMap,omitzero"`
	IntMapMap     map[string]map[string]int       `json:"intMapMap,omitzero"`
	Int8MapMap    map[string]map[string]int8      `json:"int8MapMap,omitzero"`
	Int16MapMap   map[string]map[string]int16     `json:"int16MapMap,omitzero"`
	Int32MapMap   map[string]map[string]int32     `json:"int32MapMap,omitzero"`
	Int64MapMap   map[string]map[string]int64     `json:"int64MapMap,omitzero"`
	UintMapMap    map[string]map[string]uint      `json:"uintMapMap,omitzero"`
	Uint8MapMap   map[string]map[string]uint8     `json:"uint8MapMap,omitzero"`
	Uint16MapMap  map[string]map[string]uint16    `json:"uint16MapMap,omitzero"`
	Uint32MapMap  map[string]map[string]uint32    `json:"uint32MapMap,omitzero"`
	Uint64MapMap  map[string]map[string]uint64    `json:"uint64MapMap,omitzero"`
	Float32MapMap map[string]map[string]float32   `json:"float32MapMap,omitzero"`
	Float64MapMap map[string]map[string]float64   `json:"float64MapMap,omitzero"`
	TimeMapMap    map[string]map[string]time.Time `json:"timeMapMap,omitzero"`

	// Recursive types
	Self           *OmitZeroStruct                       `json:"self,omitzero"`
	SelfSlice      []*OmitZeroStruct                     `json:"selfSlice,omitzero"`
	SelfMap        map[string]*OmitZeroStruct            `json:"selfMap,omitzero"`
	SelfSliceSlice [][](*OmitZeroStruct)                 `json:"selfSliceSlice,omitzero"`
	SelfSliceMap   []map[string]*OmitZeroStruct          `json:"selfSliceMap,omitzero"`
	SelfMapSlice   map[string][]*OmitZeroStruct          `json:"selfMapSlice,omitzero"`
	SelfMapMap     map[string]map[string]*OmitZeroStruct `json:"selfMapMap,omitzero"`
}

type OmitEmptyAndZeroStruct struct {
	// Basic types
	String  string    `json:"string,omitempty,omitzero"`
	Bool    bool      `json:"bool,omitempty,omitzero"`
	Int     int       `json:"int,omitempty,omitzero"`
	Int8    int8      `json:"int8,omitempty,omitzero"`
	Int16   int16     `json:"int16,omitempty,omitzero"`
	Int32   int32     `json:"int32,omitempty,omitzero"`
	Int64   int64     `json:"int64,omitempty,omitzero"`
	Uint    uint      `json:"uint,omitempty,omitzero"`
	Uint8   uint8     `json:"uint8,omitempty,omitzero"`
	Uint16  uint16    `json:"uint16,omitempty,omitzero"`
	Uint32  uint32    `json:"uint32,omitempty,omitzero"`
	Uint64  uint64    `json:"uint64,omitempty,omitzero"`
	Float32 float32   `json:"float32,omitempty,omitzero"`
	Float64 float64   `json:"float64,omitempty,omitzero"`
	Time    time.Time `json:"time,omitempty,omitzero"`
	Bytes   []byte    `json:"bytes,omitempty,omitzero"`

	// Pointers to basic types
	StringPtr  *string    `json:"stringPtr,omitempty,omitzero"`
	BoolPtr    *bool      `json:"boolPtr,omitempty,omitzero"`
	IntPtr     *int       `json:"intPtr,omitempty,omitzero"`
	Int8Ptr    *int8      `json:"int8Ptr,omitempty,omitzero"`
	Int16Ptr   *int16     `json:"int16Ptr,omitempty,omitzero"`
	Int32Ptr   *int32     `json:"int32Ptr,omitempty,omitzero"`
	Int64Ptr   *int64     `json:"int64Ptr,omitempty,omitzero"`
	UintPtr    *uint      `json:"uintPtr,omitempty,omitzero"`
	Uint8Ptr   *uint8     `json:"uint8Ptr,omitempty,omitzero"`
	Uint16Ptr  *uint16    `json:"uint16Ptr,omitempty,omitzero"`
	Uint32Ptr  *uint32    `json:"uint32Ptr,omitempty,omitzero"`
	Uint64Ptr  *uint64    `json:"uint64Ptr,omitempty,omitzero"`
	Float32Ptr *float32   `json:"float32Ptr,omitempty,omitzero"`
	Float64Ptr *float64   `json:"float64Ptr,omitempty,omitzero"`
	TimePtr    *time.Time `json:"timePtr,omitempty,omitzero"`

	// Slices of basic types
	StringSlice  []string    `json:"stringSlice,omitempty,omitzero"`
	BoolSlice    []bool      `json:"boolSlice,omitempty,omitzero"`
	IntSlice     []int       `json:"intSlice,omitempty,omitzero"`
	Int8Slice    []int8      `json:"int8Slice,omitempty,omitzero"`
	Int16Slice   []int16     `json:"int16Slice,omitempty,omitzero"`
	Int32Slice   []int32     `json:"int32Slice,omitempty,omitzero"`
	Int64Slice   []int64     `json:"int64Slice,omitempty,omitzero"`
	UintSlice    []uint      `json:"uintSlice,omitempty,omitzero"`
	Uint8Slice   []uint8     `json:"uint8Slice,omitempty,omitzero"`
	Uint16Slice  []uint16    `json:"uint16Slice,omitempty,omitzero"`
	Uint32Slice  []uint32    `json:"uint32Slice,omitempty,omitzero"`
	Uint64Slice  []uint64    `json:"uint64Slice,omitempty,omitzero"`
	Float32Slice []float32   `json:"float32Slice,omitempty,omitzero"`
	Float64Slice []float64   `json:"float64Slice,omitempty,omitzero"`
	TimeSlice    []time.Time `json:"timeSlice,omitempty,omitzero"`

	// Maps of string to basic types
	StringMap  map[string]string    `json:"stringMap,omitempty,omitzero"`
	BoolMap    map[string]bool      `json:"boolMap,omitempty,omitzero"`
	IntMap     map[string]int       `json:"intMap,omitempty,omitzero"`
	Int8Map    map[string]int8      `json:"int8Map,omitempty,omitzero"`
	Int16Map   map[string]int16     `json:"int16Map,omitempty,omitzero"`
	Int32Map   map[string]int32     `json:"int32Map,omitempty,omitzero"`
	Int64Map   map[string]int64     `json:"int64Map,omitempty,omitzero"`
	UintMap    map[string]uint      `json:"uintMap,omitempty,omitzero"`
	Uint8Map   map[string]uint8     `json:"uint8Map,omitempty,omitzero"`
	Uint16Map  map[string]uint16    `json:"uint16Map,omitempty,omitzero"`
	Uint32Map  map[string]uint32    `json:"uint32Map,omitempty,omitzero"`
	Uint64Map  map[string]uint64    `json:"uint64Map,omitempty,omitzero"`
	Float32Map map[string]float32   `json:"float32Map,omitempty,omitzero"`
	Float64Map map[string]float64   `json:"float64Map,omitempty,omitzero"`
	TimeMap    map[string]time.Time `json:"timeMap,omitempty,omitzero"`

	// Slice of slices
	StringSliceSlice  [][]string    `json:"stringSliceSlice,omitempty,omitzero"`
	BoolSliceSlice    [][]bool      `json:"boolSliceSlice,omitempty,omitzero"`
	IntSliceSlice     [][]int       `json:"intSliceSlice,omitempty,omitzero"`
	Int8SliceSlice    [][]int8      `json:"int8SliceSlice,omitempty,omitzero"`
	Int16SliceSlice   [][]int16     `json:"int16SliceSlice,omitempty,omitzero"`
	Int32SliceSlice   [][]int32     `json:"int32SliceSlice,omitempty,omitzero"`
	Int64SliceSlice   [][]int64     `json:"int64SliceSlice,omitempty,omitzero"`
	UintSliceSlice    [][]uint      `json:"uintSliceSlice,omitempty,omitzero"`
	Uint8SliceSlice   [][]uint8     `json:"uint8SliceSlice,omitempty,omitzero"`
	Uint16SliceSlice  [][]uint16    `json:"uint16SliceSlice,omitempty,omitzero"`
	Uint32SliceSlice  [][]uint32    `json:"uint32SliceSlice,omitempty,omitzero"`
	Uint64SliceSlice  [][]uint64    `json:"uint64SliceSlice,omitempty,omitzero"`
	Float32SliceSlice [][]float32   `json:"float32SliceSlice,omitempty,omitzero"`
	Float64SliceSlice [][]float64   `json:"float64SliceSlice,omitempty,omitzero"`
	TimeSliceSlice    [][]time.Time `json:"timeSliceSlice,omitempty,omitzero"`

	// Slice of maps
	StringSliceMap  []map[string]string    `json:"stringSliceMap,omitempty,omitzero"`
	BoolSliceMap    []map[string]bool      `json:"boolSliceMap,omitempty,omitzero"`
	IntSliceMap     []map[string]int       `json:"intSliceMap,omitempty,omitzero"`
	Int8SliceMap    []map[string]int8      `json:"int8SliceMap,omitempty,omitzero"`
	Int16SliceMap   []map[string]int16     `json:"int16SliceMap,omitempty,omitzero"`
	Int32SliceMap   []map[string]int32     `json:"int32SliceMap,omitempty,omitzero"`
	Int64SliceMap   []map[string]int64     `json:"int64SliceMap,omitempty,omitzero"`
	UintSliceMap    []map[string]uint      `json:"uintSliceMap,omitempty,omitzero"`
	Uint8SliceMap   []map[string]uint8     `json:"uint8SliceMap,omitempty,omitzero"`
	Uint16SliceMap  []map[string]uint16    `json:"uint16SliceMap,omitempty,omitzero"`
	Uint32SliceMap  []map[string]uint32    `json:"uint32SliceMap,omitempty,omitzero"`
	Uint64SliceMap  []map[string]uint64    `json:"uint64SliceMap,omitempty,omitzero"`
	Float32SliceMap []map[string]float32   `json:"float32SliceMap,omitempty,omitzero"`
	Float64SliceMap []map[string]float64   `json:"float64SliceMap,omitempty,omitzero"`
	TimeSliceMap    []map[string]time.Time `json:"timeSliceMap,omitempty,omitzero"`

	// Map of string to slices
	StringMapSlice  map[string][]string    `json:"stringMapSlice,omitempty,omitzero"`
	BoolMapSlice    map[string][]bool      `json:"boolMapSlice,omitempty,omitzero"`
	IntMapSlice     map[string][]int       `json:"intMapSlice,omitempty,omitzero"`
	Int8MapSlice    map[string][]int8      `json:"int8MapSlice,omitempty,omitzero"`
	Int16MapSlice   map[string][]int16     `json:"int16MapSlice,omitempty,omitzero"`
	Int32MapSlice   map[string][]int32     `json:"int32MapSlice,omitempty,omitzero"`
	Int64MapSlice   map[string][]int64     `json:"int64MapSlice,omitempty,omitzero"`
	UintMapSlice    map[string][]uint      `json:"uintMapSlice,omitempty,omitzero"`
	Uint8MapSlice   map[string][]uint8     `json:"uint8MapSlice,omitempty,omitzero"`
	Uint16MapSlice  map[string][]uint16    `json:"uint16MapSlice,omitempty,omitzero"`
	Uint32MapSlice  map[string][]uint32    `json:"uint32MapSlice,omitempty,omitzero"`
	Uint64MapSlice  map[string][]uint64    `json:"uint64MapSlice,omitempty,omitzero"`
	Float32MapSlice map[string][]float32   `json:"float32MapSlice,omitempty,omitzero"`
	Float64MapSlice map[string][]float64   `json:"float64MapSlice,omitempty,omitzero"`
	TimeMapSlice    map[string][]time.Time `json:"timeMapSlice,omitempty,omitzero"`

	// Map of string to maps
	StringMapMap  map[string]map[string]string    `json:"stringMapMap,omitempty,omitzero"`
	BoolMapMap    map[string]map[string]bool      `json:"boolMapMap,omitempty,omitzero"`
	IntMapMap     map[string]map[string]int       `json:"intMapMap,omitempty,omitzero"`
	Int8MapMap    map[string]map[string]int8      `json:"int8MapMap,omitempty,omitzero"`
	Int16MapMap   map[string]map[string]int16     `json:"int16MapMap,omitempty,omitzero"`
	Int32MapMap   map[string]map[string]int32     `json:"int32MapMap,omitempty,omitzero"`
	Int64MapMap   map[string]map[string]int64     `json:"int64MapMap,omitempty,omitzero"`
	UintMapMap    map[string]map[string]uint      `json:"uintMapMap,omitempty,omitzero"`
	Uint8MapMap   map[string]map[string]uint8     `json:"uint8MapMap,omitempty,omitzero"`
	Uint16MapMap  map[string]map[string]uint16    `json:"uint16MapMap,omitempty,omitzero"`
	Uint32MapMap  map[string]map[string]uint32    `json:"uint32MapMap,omitempty,omitzero"`
	Uint64MapMap  map[string]map[string]uint64    `json:"uint64MapMap,omitempty,omitzero"`
	Float32MapMap map[string]map[string]float32   `json:"float32MapMap,omitempty,omitzero"`
	Float64MapMap map[string]map[string]float64   `json:"float64MapMap,omitempty,omitzero"`
	TimeMapMap    map[string]map[string]time.Time `json:"timeMapMap,omitempty,omitzero"`

	// Recursive types
	Self           *OmitEmptyAndZeroStruct                       `json:"self,omitempty,omitzero"`
	SelfSlice      []*OmitEmptyAndZeroStruct                     `json:"selfSlice,omitempty,omitzero"`
	SelfMap        map[string]*OmitEmptyAndZeroStruct            `json:"selfMap,omitempty,omitzero"`
	SelfSliceSlice [][](*OmitEmptyAndZeroStruct)                 `json:"selfSliceSlice,omitempty,omitzero"`
	SelfSliceMap   []map[string]*OmitEmptyAndZeroStruct          `json:"selfSliceMap,omitempty,omitzero"`
	SelfMapSlice   map[string][]*OmitEmptyAndZeroStruct          `json:"selfMapSlice,omitempty,omitzero"`
	SelfMapMap     map[string]map[string]*OmitEmptyAndZeroStruct `json:"selfMapMap,omitempty,omitzero"`
}

func TestKYAMLEncoderFuzzRoundTrip(t *testing.T) {
	for i := 0; i < 1000; i++ {
		t.Run("i="+strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			// Create and fill an instance.
			original := &AllTypesStruct{}
			f := randfill.New().NilChance(0.5).NumElements(1, 3).MaxDepth(3)
			f.Fill(original)

			// Marshal to KYAML.
			ky := &Encoder{}
			yb, err := ky.Marshal(original)
			if err != nil {
				t.Fatalf("iteration %d: failed to render KYAML: %v", i, err)
			}

			// Parse back from KYAML with the standard parser.
			parsed := &AllTypesStruct{}
			if err := yaml.Unmarshal(yb, parsed); err != nil {
				t.Fatalf("iteration %d: failed to parse KYAML: %v\nKYAML:\n%s", i, err, string(yb))
			}

			// Compare.
			if diff := cmp.Diff(original, parsed, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("iteration %d: objects differ after round trip (-original +parsed):\n%s\nKYAML:\n%s", i, diff, string(yb))
			}
		})
	}
}

func TestKYAMLOmittedField(t *testing.T) {
	type Struct struct {
		String string      `json:"string"`
		Bool   bool        `json:"bool"`
		Int    int         `json:"int"`
		Plain  PlainStruct `json:"-"`
	}

	for i := 0; i < 1000; i++ {
		// Create and fill an instance.
		original := &Struct{}
		f := randfill.New().NilChance(0.5).NumElements(1, 3).MaxDepth(3)
		f.Fill(original)

		// Marshal to KYAML.
		ky := &Encoder{}
		yb, err := ky.Marshal(original)
		if err != nil {
			t.Fatalf("iteration %d: failed to render KYAML: %v", i, err)
		}

		// Parse back from KYAML with the standard parser.
		parsed := &Struct{}
		if err := yaml.Unmarshal(yb, parsed); err != nil {
			t.Fatalf("iteration %d: failed to parse KYAML: %v\nKYAML:\n%s", i, err, string(yb))
		}

		// Wipe the state that should not have been produced.
		original.Plain = PlainStruct{}

		// Compare.
		if diff := cmp.Diff(original, parsed, cmpopts.EquateEmpty()); diff != "" {
			t.Fatalf("iteration %d: objects differ after round trip (-original +parsed):\n%s\nKYAML:\n%s", i, diff, string(yb))
		}
	}
}

type SelfMarshalStruct struct {
	String string `json:"string,omitempty"`
}

func (s SelfMarshalStruct) MarshalJSON() ([]byte, error) {
	return []byte(`{"key":"value"}`), nil
}

func TestKYAMLSelfMarshal(t *testing.T) {
	original := &SelfMarshalStruct{String: "string"}
	ky := &Encoder{}
	yb, err := ky.Marshal(original)
	if err != nil {
		t.Fatalf("failed to render KYAML: %v", err)
	}
	expected := "{\n  key: \"value\",\n}\n"
	if s := string(yb); s != expected {
		t.Fatalf("wrong result:\nexpected: %q\n     got: %q", expected, s)
	}
}

func TestKYAMLOutputSyntax(t *testing.T) {
	type testCase struct {
		name     string
		input    any
		expected string
	}

	tests := []testCase{
		// ints
		{"positive int", int(123), `123`},
		{"negative int", int(-123), `-123`},
		{"zero int", int(0), `0`},
		{"positive int8", int8(123), `123`},
		{"negative int8", int8(-123), `-123`},
		{"zero int8", int8(0), `0`},
		{"positive int16", int16(123), `123`},
		{"negative int16", int16(-123), `-123`},
		{"zero int16", int16(0), `0`},
		{"positive int32", int32(123), `123`},
		{"negative int32", int32(-123), `-123`},
		{"zero int32", int32(0), `0`},
		{"positive int64", int64(123), `123`},
		{"negative int64", int64(-123), `-123`},
		{"zero int64", int64(0), `0`},
		// uints
		{"positive uint", uint(123), `123`},
		{"zero uint", uint(0), `0`},
		{"positive uint8", uint8(123), `123`},
		{"zero uint8", uint8(0), `0`},
		{"positive uint16", uint16(123), `123`},
		{"zero uint16", uint16(0), `0`},
		{"positive uint32", uint32(123), `123`},
		{"zero uint32", uint32(0), `0`},
		{"positive uint64", uint64(123), `123`},
		{"zero uint64", uint64(0), `0`},
		// floats
		{"positive float32", float32(3.5), `3.5`},
		{"negative float32", float32(-3.5), `-3.5`},
		{"zero float32", float32(0.0), `0`},
		{"positive float64", float64(3.5), `3.5`},
		{"negative float64", float64(-3.5), `-3.5`},
		{"zero float64", float64(0.0), `0`},
		// bool
		{"true bool", true, `true`},
		{"false bool", false, `false`},
		// string
		{"empty string", "", `""`},
		{"regular string", "abc", `"abc"`},
		// struct
		{"no-init struct", struct {
			I int
			S string
		}{}, "{\n  I: 0,\n  S: \"\",\n}"},
		{"init struct", struct {
			I int
			S string
		}{1, "one"}, "{\n  I: 1,\n  S: \"one\",\n}"},
		{"empty struct", struct{}{}, `{}`},
		{"omitempty struct", struct {
			I int    `json:",omitempty"`
			S string `json:",omitempty"`
			B bool   `json:",omitempty"`
			P *int   `json:",omitempty"`
		}{}, "{}"},
		{"omitempty struct nil slice", struct {
			S []int `json:",omitempty"`
		}{}, "{}"},
		{"omitempty struct empty slice", struct {
			S []int `json:",omitempty"`
		}{S: []int{}}, "{}"},
		{"omitempty struct nil map", struct {
			M map[int]int `json:",omitempty"`
		}{}, "{}"},
		{"omitempty struct empty map", struct {
			M map[int]int `json:",omitempty"`
		}{M: map[int]int{}}, "{}"},
		// slice of primitive
		{"non-empty slice", []int{1, 2, 3}, "[\n  1,\n  2,\n  3,\n]"},
		{"empty slice", []int{}, `[]`},
		{"nil slice", []int(nil), `null`},
		// array of primitive
		{"empty array", [3]int{}, "[\n  0,\n  0,\n  0,\n]"},
		{"non-empty array", [3]int{1, 2, 3}, "[\n  1,\n  2,\n  3,\n]"},
		{"zero-len array", [0]int{}, `[]`},
		// map
		{"non-empty map[string]", map[string]int{"c": 3, "a": 1, "b": 2}, "{\n  a: 1,\n  b: 2,\n  c: 3,\n}"},
		{"non-empty map[int]", map[int]int{3: 3, 1: 1, 2: 2}, "{\n  \"1\": 1,\n  \"2\": 2,\n  \"3\": 3,\n}"},
		{"empty map", map[string]int{}, `{}`},
		{"nil map", map[int]int(nil), `null`},
		{"string map with nulls", map[string]int{
			"unambiguous": 4,
			"null":        3,
			"Null":        2,
			"NULL":        1,
			"~":           5,
		}, "{\n  \"NULL\": 1,\n  \"Null\": 2,\n  \"null\": 3,\n  unambiguous: 4,\n  \"~\": 5,\n}"},
		{"string map with trues", map[string]int{
			"true": 8,
			"True": 4,
			"TRUE": 3,
			"on":   7,
			"On":   2,
			"ON":   1,
			"yes":  9,
			"Yes":  6,
			"YES":  5,
		}, "{\n  \"ON\": 1,\n  \"On\": 2,\n  \"TRUE\": 3,\n  \"True\": 4,\n  \"YES\": 5,\n  \"Yes\": 6,\n  \"on\": 7,\n  \"true\": 8,\n  \"yes\": 9,\n}"},
		{"string map with falses", map[string]int{
			"false": 7,
			"False": 2,
			"FALSE": 1,
			"off":   9,
			"Off":   6,
			"OFF":   5,
			"no":    8,
			"No":    4,
			"NO":    3,
		}, "{\n  \"FALSE\": 1,\n  \"False\": 2,\n  \"NO\": 3,\n  \"No\": 4,\n  \"OFF\": 5,\n  \"Off\": 6,\n  \"false\": 7,\n  \"no\": 8,\n  \"off\": 9,\n}"},
		{"string map with ints", map[string]int{
			"1":        2,
			"-1":       1,
			"_1":       3,
			"__1__2__": 4,
		}, "{\n  \"-1\": 1,\n  \"1\": 2,\n  \"_1\": 3,\n  \"__1__2__\": 4,\n}"},
		{"string map with floats", map[string]int{
			"3.14":  5,
			".inf":  3,
			"-.inf": 2,
			"+.inf": 1,
			".nan":  4,
		}, "{\n  \"+.inf\": 1,\n  \"-.inf\": 2,\n  \".inf\": 3,\n  \".nan\": 4,\n  \"3.14\": 5,\n}"},
		{"string map with unquoted keys", map[string]int{
			"safe":              3,
			"_":                 1,
			"_with_underscore_": 2,
			"with-dash":         4,
			"with.dot":          5,
			"with/slash":        6,
		}, "{\n  _: 1,\n  _with_underscore_: 2,\n  safe: 3,\n  with-dash: 4,\n  with.dot: 5,\n  with/slash: 6,\n}"},
		{"string map with quoted keys", map[string]int{
			"not safe":        1,
			"with\\backslash": 2,
		}, "{\n  \"not safe\": 1,\n  \"with\\\\backslash\": 2,\n}"},
		{"string map with dash keys", map[string]int{
			"-":              1,
			"-leading-dash":  2,
			"trailing-dash-": 3,
		}, "{\n  \"-\": 1,\n  \"-leading-dash\": 2,\n  \"trailing-dash-\": 3,\n}"},
		{"string map with dot keys", map[string]int{
			".":             1,
			".leading.dot":  2,
			"trailing.dot.": 3,
		}, "{\n  \".\": 1,\n  \".leading.dot\": 2,\n  \"trailing.dot.\": 3,\n}"},
		{"string map with slash keys", map[string]int{
			"/":               1,
			"/leading/slash":  2,
			"trailing/slash/": 3,
		}, "{\n  \"/\": 1,\n  \"/leading/slash\": 2,\n  \"trailing/slash/\": 3,\n}"},
		{"string map with dates", map[string]int{
			"2006":                            2,
			"2006-1-2":                        3,
			"2006-1-2T15:4:5.999999999-08:00": 4,
			"11:00":                           1,
		}, "{\n  \"11:00\": 1,\n  \"2006\": 2,\n  \"2006-1-2\": 3,\n  \"2006-1-2T15:4:5.999999999-08:00\": 4,\n}"},
		{"multi-line-string-key map", map[string]int{
			"1\n 2\n  \n3": 123,
			"4\n 5\n  \n6": 456,
		}, "{\n  \"1\\n 2\\n  \\n3\": 123,\n  \"4\\n 5\\n  \\n6\": 456,\n}"},
		// pointer
		{"non-nil pointer", new(int), `0`},
		{"nil pointer", (*int)(nil), `null`},
		// slice of compound
		{"slice of struct", []struct{ I int }{
			{1}, {2}, {3},
		}, "[{\n  I: 1,\n}, {\n  I: 2,\n}, {\n  I: 3,\n}]"},
		{"slice of slice", [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		}, "[[\n  1,\n  2,\n  3,\n], [\n  4,\n  5,\n  6,\n], [\n  7,\n  8,\n  9,\n]]"},
		{"slice of map", []map[string]int{
			{"a": 1, "b": 2},
			{"c": 3, "d": 4},
			{"e": 5, "f": 6},
		}, "[{\n  a: 1,\n  b: 2,\n}, {\n  c: 3,\n  d: 4,\n}, {\n  e: 5,\n  \"f\": 6,\n}]"},
		// interface
		// TODO: figure out how to make a reflect.Value where Kind() ==
		// Interface. As far as I can tell, ValueOf() returns either the
		// underlying type's Kind or Invalid.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ky := &Encoder{}
			yb, err := ky.Marshal(tt.input)
			if err != nil {
				t.Fatalf("failed to render KYAML: %v", err)
			}
			// always has a newline at the end
			if result := strings.TrimRight(string(yb), "\n"); result != tt.expected {
				t.Errorf("Marshal(%v):\nwanted: %q\n   got: %q", tt.input, tt.expected, result)
			}
		})
	}
}

func TestKYAMLFromYAML(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		expected string
	}

	tests := []testCase{
		{ // Without comments
			name:     "empty YAML",
			input:    ``,
			expected: ``,
		}, {
			name:  "int",
			input: `42`,
			expected: `
				---
				42
				`,
		}, {
			name:  "bool",
			input: `true`,
			expected: `
				---
				true
				`,
		}, {
			name:  "naked string",
			input: `a string`,
			expected: `
				---
				"a string"
				`,
		}, {
			name:  "quoted string",
			input: `a string`,
			expected: `
				---
				"a string"
				`,
		}, {
			name: "empty list",
			input: `
				[
				]
				`,
			expected: `
				---
				[]
				`,
		}, {
			name: "list of int",
			input: `
				---
				- 1
				- 2
				- 3
				`,
			expected: `
				---
				[
				  1,
				  2,
				  3,
				]
				`,
		}, {
			name: "list of bool",
			input: `
				---
				- true
				- false
				- true
				`,
			expected: `
				---
				[
				  true,
				  false,
				  true,
				]
				`,
		}, {
			name: "list of string",
			input: `
				---
				- naked
				- "dquoted"
				- 'squoted'
				`,
			expected: `
				---
				[
				  "naked",
				  "dquoted",
				  "squoted",
				]
				`,
		}, {
			name: "empty mapping",
			input: `
				{
				}
				`,
			expected: `
				---
				{}
				`,
		}, {
			name: "mapping of ints",
			input: `
				a: 1
				b: 2
				c: 3
				`,
			expected: `
				---
				{
				  a: 1,
				  b: 2,
				  c: 3,
				}
				`,
		}, {
			name: "mapping of bool",
			input: `
				a: true
				b: false
				c: true
				`,
			expected: `
				---
				{
				  a: true,
				  b: false,
				  c: true,
				}
				`,
		}, {
			name: "mapping of string",
			input: `
				a: naked
				b: "dquoted"
				c: 'squoted'
				`,
			expected: `
				---
				{
				  a: "naked",
				  b: "dquoted",
				  c: "squoted",
				}
				`,
		}, {
			// Note: there's no escaping within single-quotes.
			name: "mapping of strings with quotes",
			input: `
				a: naked "with" 'quotes' embedded
				b: "dquoted \"with\" 'quotes' embedded"
				c: 'squoted "with" quotes embedded'
				`,
			expected: `
				---
				{
				  a: "naked \"with\" 'quotes' embedded",
				  b: "dquoted \"with\" 'quotes' embedded",
				  c: "squoted \"with\" quotes embedded",
				}
				`,
		}, {
			name: "list of list",
			input: `
				-
				  - 1
				  - 2
				  - 3
				-
				  - true
				  - false
				  - true
				-
				  - naked
				  - "dquoted"
				  - 'squoted'
				`,
			expected: `
				---
				[[
				  1,
				  2,
				  3,
				], [
				  true,
				  false,
				  true,
				], [
				  "naked",
				  "dquoted",
				  "squoted",
				]]
				`,
		}, {
			name: "list of mapping",
			input: `
				- a: 1
				  b: 2
				  c: 3
				- a: true
				  b: false
				  c: true
				- a: naked
				  b: "dquoted"
				  c: 'squoted'
				`,
			expected: `
				---
				[{
				  a: 1,
				  b: 2,
				  c: 3,
				}, {
				  a: true,
				  b: false,
				  c: true,
				}, {
				  a: "naked",
				  b: "dquoted",
				  c: "squoted",
				}]
				`,
		}, {
			name: "mapping of list",
			input: `
				a:
				- 1
				- 2
				- 3
				b:
				- true
				- false
				- true
				c:
				- naked
				- "dquoted"
				- 'squoted'
				`,
			expected: `
				---
				{
				  a: [
				    1,
				    2,
				    3,
				  ],
				  b: [
				    true,
				    false,
				    true,
				  ],
				  c: [
				    "naked",
				    "dquoted",
				    "squoted",
				  ],
				}
				`,
		}, {
			name: "mapping with null-string keys",
			input: `
				null: 123
				Null: 123
				NULL: 123
				~: 123
				`,
			expected: `
				---
				{
				  "null": 123,
				  "Null": 123,
				  "NULL": 123,
				  "~": 123,
				}
				`,
		}, {
			name: "mapping with true-string keys",
			input: `
				true: 123
				True: 123
				TRUE: 123
				on: 123
				On: 123
				ON: 123
				yes: 123
				Yes: 123
				YES: 123
				`,
			expected: `
				---
				{
				  "true": 123,
				  "True": 123,
				  "TRUE": 123,
				  "on": 123,
				  "On": 123,
				  "ON": 123,
				  "yes": 123,
				  "Yes": 123,
				  "YES": 123,
				}
				`,
		}, {
			name: "mapping with false-string keys",
			input: `
				false: 123
				False: 123
				FALSE: 123
				off: 123
				Off: 123
				OFF: 123
				no: 123
				No: 123
				NO: 123
				`,
			expected: `
				---
				{
				  "false": 123,
				  "False": 123,
				  "FALSE": 123,
				  "off": 123,
				  "Off": 123,
				  "OFF": 123,
				  "no": 123,
				  "No": 123,
				  "NO": 123,
				}
				`,
		}, {
			name: "mapping with int-string keys",
			input: `
				1: 123
				-1: 123
				+1: 123
				_1: 123
				-_1: 123
				+_1: 123
				__1__2__: 123
				_-_1__2__: 123
				_+_1__2__: 123
				`,
			expected: `
				---
				{
				  "1": 123,
				  "-1": 123,
				  "+1": 123,
				  "_1": 123,
				  "-_1": 123,
				  "+_1": 123,
				  "__1__2__": 123,
				  "_-_1__2__": 123,
				  "_+_1__2__": 123,
				}
				`,
		}, {
			name: "mapping with float-string keys",
			input: `
				3.14: 123
				-3.14: 123
				+3.14: 123
				.inf: 123
				-.inf: 123
				+.inf: 123
				.nan: 123
				`,
			expected: `
				---
				{
				  "3.14": 123,
				  "-3.14": 123,
				  "+3.14": 123,
				  ".inf": 123,
				  "-.inf": 123,
				  "+.inf": 123,
				  ".nan": 123,
				}
				`,
		}, {
			name: "mapping with naked-string keys",
			input: `
				safe: 123
				_: 123
				_with_underscore_: 123
				with-dash: 123
				with.dot: 123
				with/slash: 123
				`,
			expected: `
				---
				{
				  safe: 123,
				  _: 123,
				  _with_underscore_: 123,
				  with-dash: 123,
				  with.dot: 123,
				  with/slash: 123,
				}
				`,
		}, {
			name: "mapping with unsafe-string keys",
			input: `
				not safe: 123
				with\backslash: 123
				`,
			expected: `
				---
				{
				  "not safe": 123,
				  "with\\backslash": 123,
				}
				`,
		}, {
			name: "mapping with dash-string keys",
			input: `
				-: 123
				-leading-dash: 123
				trailing-dash-: 123
				`,
			expected: `
				---
				{
				  "-": 123,
				  "-leading-dash": 123,
				  "trailing-dash-": 123,
				}
				`,
		}, {
			name: "mapping with dot-string keys",
			input: `
				.: 123
				.leading.dot: 123
				trailing.dot.: 123
				`,
			expected: `
				---
				{
				  ".": 123,
				  ".leading.dot": 123,
				  "trailing.dot.": 123,
				}
				`,
		}, {
			name: "mapping with slash-string keys",
			input: `
				/: 123
				/leading/slash: 123
				trailing/slash/: 123
				`,
			expected: `
				---
				{
				  "/": 123,
				  "/leading/slash": 123,
				  "trailing/slash/": 123,
				}
				`,
		}, {
			name: "mapping with date-string keys",
			input: `
				2006: 123
				2006-1-2: 123
				2006-1-2T15:4:5.999999999-08:00: 123
				11:00: 123
				`,
			expected: `
				---
				{
				  "2006": 123,
				  "2006-1-2": 123,
				  "2006-1-2T15:4:5.999999999-08:00": 123,
				  "11:00": 123,
				}
				`,
		}, {
			name: "multi-line strings",
			input: `
				simple: |
				  This is a multi-line string.
				  It has multiple lines.
				leading_space: |
				  This is a multi-line string.
				    It can
				      retain indentation.
				blank_lines: |
				  This is a multi-line string.
				
				  It can retain blank lines.
				`,
			expected: `
				---
				{
				  simple: "\
				     This is a multi-line string.\n\
				     It has multiple lines.\n\
				    ",
				  leading_space: "\
				     This is a multi-line string.\n\
				    \  It can\n\
				    \    retain indentation.\n\
				    ",
				  blank_lines: "\
				     This is a multi-line string.\n\
				     \n\
				     It can retain blank lines.\n\
				    ",
				}
				`,
		},
		{ // With comments
			name: "scalar with comments",
			input: `
				# This is a head comment.
				# It can be multi-line.
				42    # This is a line comment
				# This is a foot comment.
				# It can also be multi-line.
				`,
			expected: `
				---
				# This is a head comment.
				# It can be multi-line.
				42 # This is a line comment
				# This is a foot comment.
				# It can also be multi-line.
				`,
		}, {
			name: "multi-line string with comments",
			input: `
				# This is a head comment.
				# It can be multi-line.
				foo: |    # This is a line comment
				  this is a
				  multi-line
				  comment
				# This is a foot comment.
				# It can also be multi-line.
				`,
			expected: `
				---
				{
				  # This is a head comment.
				  # It can be multi-line.
				  foo: "\
				     this is a\n\
				     multi-line\n\
				     comment\n\
				    ", # This is a line comment
				  # This is a foot comment.
				  # It can also be multi-line.
				}
				`,
		}, {
			name: "block sequence with comments",
			input: `
				# This seems like a head comment.
				# But it will be attributed to the first item.
				- 1    # line1
				- 2    # line2
				- 3    # line3
				# This is a foot comment.
				# It can also be multi-line.
				`,
			expected: `
				---
				[
				  # This seems like a head comment.
				  # But it will be attributed to the first item.
				  1, # line1
				  2, # line2
				  3, # line3
				  # This is a foot comment.
				  # It can also be multi-line.
				]
				`,
		}, {
			name: "short flow sequence with line-comment",
			input: `
				# This is a head comment.
				# It can be multi-line.
				[ 1, 2, 3 ]    # This is a line comment.
				# This is a foot comment.
				# It can also be multi-line.
				`,
			expected: `
				---
				# This is a head comment.
				# It can be multi-line.
				[
				  1,
				  2,
				  3,
				] # This is a line comment.
				# This is a foot comment.
				# It can also be multi-line.
				`,
		}, {
			name: "flow sequence with comments",
			input: `
				# This is a head comment.
				# It can be multi-line.
				[ # This will be lost.
				  1,    # line1
				  2,    # line2
				  3,    # line3
				] # This is a line comment.
				# This is a foot comment.
				# It can also be multi-line.
				`,
			expected: `
				---
				# This is a head comment.
				# It can be multi-line.
				[
				  1, # line1
				  2, # line2
				  3, # line3
				] # This is a line comment.
				# This is a foot comment.
				# It can also be multi-line.
				`,
		}, {}, {
			name: "block mapping with comments",
			input: `
				# This is a head comment.
				# It can be multi-line.
				foo: 1    # line1
				bar: 2    # line2
				# This is a foot comment.
				# It can also be multi-line.
				`,
			expected: `
				---
				{
				  # This is a head comment.
				  # It can be multi-line.
				  foo: 1, # line1
				  bar: 2, # line2
				  # This is a foot comment.
				  # It can also be multi-line.
				}
				`,
		}, {
			name: "short flow mapping with line-comment",
			input: `
				# This is a head comment.
				# It can be multi-line.
				{ foo: 1, bar: 2 }    # This is a line comment.
				# This is a foot comment.
				# It can also be multi-line.
				`,
			expected: `
				---
				# This is a head comment.
				# It can be multi-line.
				{
				  foo: 1,
				  bar: 2,
				} # This is a line comment.
				# This is a foot comment.
				# It can also be multi-line.
				`,
		}, {
			name: "flow mapping with comments",
			input: `
				# This is a head comment.
				# It can be multi-line.
				{    # This will be lost.
				  foo: 1,    # line1
				  bar: 2,    # line2
				}    # This is a line comment.
				# This is a foot comment.
				# It can also be multi-line.
				`,
			expected: `
				---
				# This is a head comment.
				# It can be multi-line.
				{
				  foo: 1, # line1
				  bar: 2, # line2
				} # This is a line comment.
				# This is a foot comment.
				# It can also be multi-line.
				`,
		}, {
			name: "list of list with comments",
			input: `
				# This is a head comment.
				# It can be multi-line.
				[    # This will be lost.
				  [1],
				  [2],
				  # head3
				  [3],
				  [4],
				  # foot4
				  [5],
				  # foot5
				  # head6
				  [6],
				  [7],    # line7
				  [8],
				  [9],    # line9
				]    # This is a line comment.
				# This is a foot comment.
				# It can also be multi-line.
				`,
			expected: `
				---
				# This is a head comment.
				# It can be multi-line.
				[
				  [
				    1,
				  ],
				  [
				    2,
				  ],
				  # head3
				  [
				    3,
				  ],
				  [
				    4,
				  ],
				  # foot4
				  [
				    5,
				  ],
				  # foot5
				  # head6
				  [
				    6,
				  ],
				  [
				    7,
				  ], # line7
				  [
				    8,
				  ],
				  [
				    9,
				  ], # line9
				] # This is a line comment.
				# This is a foot comment.
				# It can also be multi-line.
				`,
		}, {
			name: "list of map with comments",
			input: `
				# This is a head comment.
				# It can be multi-line.
				[    # This will be lost.
				  {fld: 1},
				  {fld: 2},
				  # head3
				  {fld: 3},
				  {fld: 4},
				  # foot4
				  {fld: 5},
				  # foot5
				  # head6
				  {fld: 6},
				  {fld: 7},    # line7
				  {fld: 8},
				  {fld: 9},    # line9
				]    # This is a line comment.
				# This is a foot comment.
				# It can also be multi-line.
				`,
			expected: `
				---
				# This is a head comment.
				# It can be multi-line.
				[
				  {
				    fld: 1,
				  },
				  {
				    fld: 2,
				  },
				  # head3
				  {
				    fld: 3,
				  },
				  {
				    fld: 4,
				  },
				  # foot4
				  {
				    fld: 5,
				  },
				  # foot5
				  # head6
				  {
				    fld: 6,
				  },
				  {
				    fld: 7,
				  }, # line7
				  {
				    fld: 8,
				  },
				  {
				    fld: 9,
				  }, # line9
				] # This is a line comment.
				# This is a foot comment.
				# It can also be multi-line.
				`,
		}, {
			name: "list of mixed types",
			input: `
				[
				  "a string",
				  12345,
				  true,
				  {fld: 1},
				  [a, b, c],
				]
				`,
			expected: `
				---
				[
				  "a string",
				  12345,
				  true,
				  {
				    fld: 1,
				  },
				  [
				    "a",
				    "b",
				    "c",
				  ],
				]
				`,
		}, {
			name: "doc comments",
			input: `
				# This seems like a doc comment.
				# But it will be attributed to the content.
				---
				# This is a head comment.
				42    # This is a line comment
				# This is a foot comment.
				# It can also be multi-line.
				`,
			expected: `
				---
				# This seems like a doc comment.
				# But it will be attributed to the content.
				# This is a head comment.
				42 # This is a line comment
				# This is a foot comment.
				# It can also be multi-line.
				`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ky := &Encoder{}

			// Trimming is just for test output since all test cases have
			// leading and trailing newlines.
			input := strings.TrimPrefix(tt.input, "\n")
			input = strings.TrimSuffix(input, "\n")
			input = dedent(input)

			expected := strings.TrimPrefix(tt.expected, "\n")
			expected = strings.TrimSuffix(expected, "\n")
			expected = dedent(expected)

			buf := &bytes.Buffer{}
			if err := ky.FromYAML(strings.NewReader(input), buf); err != nil {
				t.Fatalf("FromYAML(%q) returned error: %v", input, err)
			}
			if want, got := expected, buf.String(); got != want {
				t.Errorf("FromYAML:\nexpected:\n```\n%s\n```\ngot:\n```\n%s\n```", expected, got)
			}
		})
	}
}

func TestIsTypeAmbiguous(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		expected bool
	}

	tests := []testCase{
		// Regular strings that should not need quotes
		{"regular string", "regular", false},
		{"alphanumeric", "abc123", false},
		{"underscore", "hello_world", false},
		{"hyphen", "hello-world", false},

		// Boolean-like strings
		{"yes", "yes", true},
		{"YES", "YES", true},
		{"y", "y", true},
		{"Y", "Y", true},
		{"no", "no", true},
		{"NO", "NO", true},
		{"n", "n", true},
		{"N", "N", true},
		{"on", "on", true},
		{"ON", "ON", true},
		{"off", "off", true},
		{"OFF", "OFF", true},

		// Numbers should stay strings
		{"decimal", "1234", true},
		{"underscores", "_1_2_3_4_", true},
		{"leading zero", "0123", true},
		{"plus sign", "+123", true},
		{"negative sign", "-123", true},
		{"large decimal", "123456789012345678901234567890", true},
		{"octal 0", "0777", true},
		{"octal 0o", "0o777", true},
		{"hex lowercase", "0xff", true},
		{"hex uppercase", "0xFF", true},

		// Infinity and NaN
		{"infinity", ".inf", true},
		{"negative infinity", "-.inf", true},
		{"positive infinity", "+.inf", true},
		{"not a number", ".nan", true},
		{"uppercase infinity", ".INF", true},
		{"uppercase nan", ".NAN", true},

		// Scientific notation
		{"scientific", "1e10", true},
		{"scientific uppercase", "1E10", true},

		// Timestamp-like strings
		{"year", "2006", true},
		{"date", "2006-1-2", true},
		{"RCF3339Nano with short date", "2006-1-2T15:4:5.999999999-08:00", true},
		{"RCF3339Nano with lowercase t", "2006-1-2t15:4:5.999999999-08:00", true},
		{"space separated", "2006-1-2 14:4:5.999999999", true},

		// Sexagesimal strings
		{"small sexagesimal int", "1:00", true},
		{"large sexagesimal int 59", "12345:59", true},
		{"invalid sexagesimal int", "1:60", false},
		{"multi-part sexagesimal", "1:2:3:4:5:00:00:01", true},
		{"small sexagesimal float 0", "1:00.0", true},
		{"small sexagesimal float 59", "12345:59.59", true},
		{"small non-sexagesimal int", "12345:60.00", false},
		{"multi-part sexagesimal", "1:2:3:4:5:00:00:01.02", true},

		// Null-like strings
		{"null lowercase", "null", true},
		{"null mixed", "Null", true},
		{"null uppercase", "NULL", true},
		{"tilde", "~", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTypeAmbiguous(tt.input)
			if result != tt.expected {
				t.Errorf("isTypeAmbiguous(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
