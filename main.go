package main

import (
	"fmt"
	"log"
	"os"

	"github.com/frostschutz/go-fibmap"
)

func getExtents(filename string) ([]fibmap.Extent, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	fm := fibmap.NewFibmapFile(fd)

	bsz, errno := fm.Figetbsz()
	if errno != 0 {
		return nil, fmt.Errorf("figetbsz: %v", errno)
	}

	stat, err := fd.Stat()
	if err != nil {
		return nil, fmt.Errorf("fstat: %v", err)
	}
	size := stat.Size()

	blocks := uint32((size-1)/int64(bsz)) + 1

	extents, errno := fm.Fiemap(blocks)
	if errno != 0 {
		return nil, fmt.Errorf("fiemap: %v", errno)
	}
	return extents, nil
}

func sharedExtents(a, b string) (shared, total uint64, err error) {
	extentsA, err := getExtents(a)
	if err != nil {
		return 0, 0, err
	}
	extentsB, err := getExtents(b)
	if err != nil {
		return 0, 0, err
	}

	type Extent struct{ Start, Length uint64 }

	type ExtentSet map[Extent]struct{}

	totalLength := func(es ExtentSet) uint64 {
		var total uint64
		for e := range es {
			total += e.Length
		}
		return total
	}

	newExtentSet := func(es []fibmap.Extent) ExtentSet {
		extentSet := ExtentSet{}
		for _, e := range es {
			extent := Extent{e.Physical, e.Length}
			extentSet[extent] = struct{}{}
		}
		return extentSet
	}

	intersect := func(as, bs ExtentSet) ExtentSet {
		extentSet := ExtentSet{}
		for a := range as {
			if _, ok := bs[a]; ok {
				extentSet[a] = struct{}{}
			}
		}
		return extentSet
	}

	extentSetA := newExtentSet(extentsA)
	extentSetB := newExtentSet(extentsB)

	totalLengthA := totalLength(extentSetA)
	totalLengthB := totalLength(extentSetB)
	if totalLengthA != totalLengthB {
		return 0, 0, fmt.Errorf(
			"Files %q (%d) and %q (%d) differ in block length",
			a, totalLengthA, b, totalLengthB)
	}

	totalLengthBoth := totalLengthA

	intersection := intersect(extentSetA, extentSetB)
	totalLengthI := totalLength(intersection)

	return totalLengthI, totalLengthBoth, nil
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("sharedextents: ")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: sharedextents <a> <b>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Reports number bytes between <a> and <b> sharing physical extents.")
		os.Exit(1)
	}

	a := os.Args[1]
	b := os.Args[2]

	shared, total, err := sharedExtents(a, b)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d / %d bytes (%.2f%%)\n", shared, total,
		100*float64(shared)/float64(total))

	if shared == 0 {
		os.Exit(1)
	}
}
