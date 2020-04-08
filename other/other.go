package main

import (
	"flag"
	"fmt"
)

func main() {
	var parallelJobs int64
	flag.Int64Var(&parallelJobs, "parallel", 3, "Number of jobs to run in parallel")
	flag.Parse()

	args := flag.Args()

	fmt.Printf("Parallel %d\n", parallelJobs)
	fmt.Println(args)
	
}
