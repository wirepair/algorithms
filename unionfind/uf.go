package main

import (
	"flag"
	"fmt"
	"github.com/wirepair/algorithms/helpers"
	"log"
	"os"
)

// UnionFind interface used for implementations of
// quick-find, quick-union and weighted quick-union.
type UnionFinder interface {
	Init(N int64)              // Sets number of sites.
	Count() int64              // returns number of seperate components
	Connected(p, q int64) bool // Searches to see if two points are connected
	Union(p, q int64)          // updates connected sites.
	Find(p int64) int64        // Finds the value of a site in our site list.
}

// An implementation of Quick-Find
type UnionQuickFind struct {
	id    []int64 // site array.
	count int64   // number of components.
}

// Initializes our site container with
// values 0~N-1 sites.
// So id[0] = 0, id[1] = 1, etc
func (UF *UnionQuickFind) Init(N int64) {
	UF.count = N
	UF.id = make([]int64, N)
	for i := int64(0); i < N; i++ {
		UF.id[i] = i
	}
}

// returns the # of components.
// where components is a list of sites that are connected
func (UF *UnionQuickFind) Count() int64 {
	return UF.count
}

// checks if p is connected to q
func (UF *UnionQuickFind) Connected(p, q int64) bool {
	//fmt.Printf("p: %d q: %d connected? %v\n", p, q, UF.Find(p) == UF.Find(q))
	return UF.Find(p) == UF.Find(q)
}

// Union gets the values of p/q respectively from our
// sites (id) then iterates over the site list
// and sets the id to the value of q.
// finally it decrements the count from the total
// possible list of connected components.
func (UF *UnionQuickFind) Union(p, q int64) {
	rootP := UF.Find(p)
	rootQ := UF.Find(q)

	if rootP == rootQ {
		return
	}
	for i := 0; i < len(UF.id); i++ {
		if UF.id[i] == rootP {
			UF.id[i] = rootQ
		}
	}
	//fmt.Printf("%#v\n", UF.id)
	UF.count--
}

// Simply returns the value of where p is in
// our site list.
func (UF *UnionQuickFind) Find(p int64) int64 {
	return UF.id[p]
}

// An implementation of Quick-Union
type QuickUnionFind struct {
	UnionQuickFind // embeds id, count, Init, Connected and Count.
}

// iterates when p does not equal the value of
// the site at p in our site list and updates
// p to the new value, and returns it.
func (UF *QuickUnionFind) Find(p int64) int64 {
	for p != UF.id[p] {
		p = UF.id[p]
	}
	return p
}

// finds p/q and updates the site index at i
// to the value of j. This creates a linked list
// like structure where each site points to
// the next site
func (UF *QuickUnionFind) Union(p, q int64) {
	i := UF.Find(p)
	j := UF.Find(q)
	if i == j {
		return
	}
	UF.id[i] = j
	//fmt.Printf("%#v\n", UF.id)
	UF.count--
}

// An implementation of Weighted Quick Union
type WeightedQuickUnion struct {
	QuickUnionFind         // embeds id, count, Find, Connected
	sz             []int64 // holds the size of the tree.
}

// WeightedQuickUnion adds a second size array
// to gauruntee logarithmic performance.
func (UF *WeightedQuickUnion) Init(N int64) {
	// we could UF.QuickUnionFind.Init(N) but why loop twice...
	UF.id = make([]int64, N)
	UF.sz = make([]int64, N)
	for i := int64(0); i < N; i++ {
		UF.id[i] = i
		UF.sz[i] = 1
	}
	UF.count = N
}

func (UF *WeightedQuickUnion) Union(p, q int64) {
	i := UF.Find(p)
	j := UF.Find(q)
	if i == j {
		return
	}
	// make smaller root point to larger one.
	if UF.sz[i] < UF.sz[j] {
		UF.id[i] = j
		UF.sz[j] += UF.sz[i]
	} else {
		UF.id[j] = i
		UF.sz[i] += UF.sz[j]
	}
	UF.count--
}

var filename string
var ufType string

func init() {
	flag.StringVar(&filename, "f", "stdin", "filename or stdin.")
	flag.StringVar(&ufType, "u", "weighted", "unionfind type: quickfind, quickunion, weighted")
}

// Creates a UnionFinder based on the requested type.
func getFinder() UnionFinder {
	var uf UnionFinder
	switch ufType {
	case "weighted":
		uf = new(WeightedQuickUnion)
	case "quickfind":
		uf = new(UnionQuickFind)
	case "quickunion":
		uf = new(QuickUnionFind)
	default:
		log.Fatal("error must choose a type of quickfind, quickunion or weighted")
	}
	return uf
}

func main() {
	var err error
	var input *os.File
	flag.Parse()
	if filename == "stdin" {
		input = os.Stdin
	} else if input, err = os.Open(filename); err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Opened %s for input.\n", filename)

	helper := helpers.New(input)
	sites, err := helper.GetSites()
	if err != nil {
		log.Fatal(err)
	}

	uf := getFinder()
	fmt.Printf("Using unionfind of type %s.\n", ufType)

	uf.Init(sites)

	intChan := make(chan int64)

	go func() {
		helper.GetInt(intChan)
	}()

	for p := range intChan {
		q := <-intChan
		//fmt.Printf("%d %d\n", p, q)
		if uf.Connected(p, q) {
			//fmt.Printf("%d %d are not connected.\n", p, q)
			continue
		}
		uf.Union(p, q)

	}
	fmt.Printf("%d components.\n", uf.Count())
	helper.Close()
}
