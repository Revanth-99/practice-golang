package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Job struct {
	id  int
	num int
}

type Result struct {
	job Job
	sum int
}

var jobs = make(chan Job, 10)
var results = make(chan Result, 10)

func sumOfDigits(num int) (result int) {
	for num != 0 {
		result = result + (num % 10)
		num = num / 10
	}
	return
}

func worker(wg *sync.WaitGroup) {
	for job := range jobs {
		res := Result{job, sumOfDigits(job.num)}
		results <- res
	}
	wg.Done()
}

func createWorkerPool(numOfWorkers int) {
	var wg sync.WaitGroup
	for index := 0; index < numOfWorkers; index++ {
		wg.Add(1)
		go worker(&wg)
	}
	wg.Wait()
	close(results)
}

func getResults(done chan bool) {
	for result := range results {
		fmt.Printf("Job id %d, input random no %d , sum of digits %d\n", result.job.id, result.job.num, result.sum)
	}
	done <- true
}

func allocateJobs(numOfJobs int) {
	for index := 1; index <= numOfJobs; index++ {
		randomNum := rand.Intn(999)
		jobs <- Job{index, randomNum}
	}
	close(jobs)
}

func main() {
	numOfWorkers := 10
	numOfJobs := 100
	done := make(chan bool)
	go createWorkerPool(numOfWorkers)
	go allocateJobs(numOfJobs)
	go getResults(done)
	<-done
}
