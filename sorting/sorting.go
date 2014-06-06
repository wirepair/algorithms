package main

import (
	"flag"
	"fmt"
	"github.com/wirepair/algorithms/helpers"
	"log"
	"os"
)

// Sorter represents an object which is sortable
// by being able to compare / exchange values.
type Sorter interface {
	Len() int
	Less(i, j int) bool
	Exch(i, j int)
	Get(i int) interface{}
	Set(i int, v interface{})
}

type IntSlice []int

func (s IntSlice) Len() int {
	return len(s)
}

func (s IntSlice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s IntSlice) Exch(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s IntSlice) Set(i int, v interface{}) {
	s[i] = v.(int)
}

func (s IntSlice) Get(i int) interface{} {
	return s[i]
}

type StringSlice []string

func (s StringSlice) Len() int {
	return len(s)
}

func (s StringSlice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s StringSlice) Exch(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s StringSlice) Set(i int, v interface{}) {
	s[i] = v.(string)
}

func (s StringSlice) Get(i int) interface{} {
	return s[i]
}

// A sort function which takes in a Sorter and sorts the data.
type SortFunc func(data Sorter)

// Returns the specified sort function, returning SelectionSort as the default
func GetSortFunc(sortType string) SortFunc {
	switch sortType {
	case "insertion":
		return InsertionSort
	case "selection":
		return SelectionSort
	case "shell":
		return ShellSort
	case "merge":
		return MergeSort
	case "quick":
		return QuickSort
	}
	// default.
	return SelectionSort
}

func SelectionSort(data Sorter) {
	N := data.Len()
	for i := 0; i < N; i++ {
		min := i
		for j := i + 1; j < N; j++ {
			if data.Less(j, min) {
				min = j
			}
		}
		data.Exch(i, min)
	}
}

func InsertionSort(data Sorter) {
	insertionSort(data, 0, data.Len())
}

func insertionSort(data Sorter, lo, hi int) {
	for i := lo; i <= hi; i++ {
		// insert data[i] amoung data[i-1], data[i-2], data[i-3]...
		for j := i; j > lo && data.Less(j, j-1); j-- {
			data.Exch(j, j-1)
			//fmt.Printf("%v\n", data)
		}
	}
}

func ShellSort(data Sorter) {
	N := data.Len()
	h := 1
	// Increment sequence. Why 3*h+1? Because math.
	for h < N/3 {
		h = 3*h + 1 // 1, 4, 13, 40, 121, 364, 1093, ...
		fmt.Printf("values in %s of size: %d make h: %d\n", filename, N, h)
	}
	for h >= 1 {
		// h-sort the aray
		for i := h; i < N; i++ {
			// insert data[i] amoung data[i-h], data[i-2*h], data[i-3*h]...
			for j := i; j >= h && data.Less(j, j-h); j -= h {
				data.Exch(j, j-h)
			}
		}
		h = h / 3
	}
}

// Ported from http://algs4.cs.princeton.edu/22mergesort/MergeX.java.html
func MergeSort(data Sorter) {
	var aux Sorter
	switch t := data.(type) {
	case StringSlice:
		aux = make(StringSlice, data.Len())
		copy(aux.(StringSlice), t)
	case IntSlice:
		aux = make(IntSlice, data.Len())
		copy(aux.(IntSlice), t)
	}

	mergeSort(aux, data, 0, data.Len()-1)
}

func mergeSort(src, dst Sorter, lo, hi int) {
	// 7 cutoff to insertion sort
	if hi <= lo+7 {
		insertionSort(dst, lo, hi)
		return
	}
	mid := lo + (hi-lo)/2
	mergeSort(dst, src, lo, mid)
	fmt.Printf("Src 1st:\n%v\n", src)
	mergeSort(dst, src, mid+1, hi)
	fmt.Printf("Src 2nd:\n%v\n", src)
	if !src.Less(mid+1, mid) {
		for i := lo; i <= hi; i++ {
			dst.Set(i, src.Get(i))
			return
		}

	}
	fmt.Printf("after copy src: %v\n", src)
	merge(src, dst, lo, mid, hi)
}

func merge(src, dst Sorter, lo, mid, hi int) {
	i := lo
	j := mid + 1
	for k := lo; k <= hi; k++ {
		if i > mid {
			dst.Set(k, src.Get(j))
			j = j + 1
		} else if j > hi {
			dst.Set(k, src.Get(i))
			i = i + 1
		} else if src.Less(j, i) {
			dst.Set(k, src.Get(j))
			j = j + 1
		} else {
			dst.Set(k, src.Get(i))
			i = i + 1
		}
	}
	fmt.Printf("Finished merging\n")
}

func QuickSort(data Sorter) {
	quickSort(data, 0, data.Len()-1)
}

func quickSort(data Sorter, lo, hi int) {
	if hi <= lo {
		return
	}
	j := partition(data, lo, hi)
	quickSort(data, lo, j-1)
	quickSort(data, j+1, hi)
}

func partition(data Sorter, lo, hi int) int {
	i := lo
	j := hi + 1
	v := lo
	for {
		for {
			i++
			if !data.Less(i, v) || i == hi {
				break
			}
		}
		for {
			j--
			if !data.Less(v, j) || j == lo {
				break
			}
		}
		if i >= j {
			break
		}
		data.Exch(i, j)
	}
	data.Exch(lo, j)
	return j
}

var filename string
var sortType string

func init() {
	flag.StringVar(&filename, "f", "stdin", "filename or stdin.")
	flag.StringVar(&sortType, "s", "selection", "sort type: selection, insertion, shell, merge.")
}

func main() {
	var err error
	var input *os.File

	flag.Parse()
	if filename == "stdin" {
		input = os.Stdin
	} else if input, err = os.Open(filename); err != nil {
		log.Fatal(err)
	}
	sort := GetSortFunc(sortType)

	helper := helpers.New(input)
	strChan := make(chan string)
	data := make(StringSlice, 0)

	go helper.GetString(strChan)

	for v := range strChan {
		data = append(data, v)
	}
	fmt.Printf("we got our data: %v\n", data)
	sort(data)
	fmt.Printf("...and sorted: %v", data)
}
