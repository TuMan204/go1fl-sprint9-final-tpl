package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	SIZE   = 100_000_000
	CHUNKS = 8
)

var errSizeZero = errors.New("Argument size is 0")

// generateRandomElements generates random elements.
func generateRandomElements(size int) []int {
	if size <= 0 {
		return []int{}
	}

	slice := make([]int, size)

	src := rand.NewSource(time.Now().Unix())
	for i := 0; i < size; i++ {
		slice[i] = int(src.Int63())
	}

	return slice
}

// maximum returns the maximum number of elements.
func maximum(data []int) int {
	if len(data) == 0 {
		return 0
	}

	max := data[0]
	for _, v := range data[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// maxChunks returns the maximum number of elements in a chunks.
func maxChunks(data []int) int {
	if len(data) == 0 {
		return 0
	}

	chunkSize := len(data) / CHUNKS

	if chunkSize <= 2 {
		return maximum(data)
	}

	var wg sync.WaitGroup
	chunkSlice := make([]int, 8)

	for i := 0; i < CHUNKS; i++ {
		slice := []int{}
		if i == CHUNKS-1 {
			slice = data[i*chunkSize:]
		} else {
			slice = data[i*chunkSize : i*chunkSize+chunkSize]
		}

		wg.Add(1)
		go func(slice []int) {
			defer wg.Done()

			max := slice[0]
			for _, v := range slice[1:] {
				if v > max {
					max = v
				}
			}
			chunkSlice[i] = max
		}(slice)
	}
	wg.Wait()

	return maximum(chunkSlice)
}

func main() {
	fmt.Printf("Генерируем %d целых чисел\n", SIZE)
	randomSlice := generateRandomElements(SIZE)

	fmt.Println("Ищем максимальное значение в один поток")
	start := time.Now().UTC()
	max := maximum(randomSlice)
	stop := time.Now().UTC()
	elapsed := stop.Sub(start).Microseconds()

	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d mks\n", max, elapsed)

	fmt.Printf("Ищем максимальное значение в %d потоков\n", CHUNKS)
	start = time.Now().UTC()
	max = maxChunks(randomSlice)
	stop = time.Now().UTC()
	elapsed = stop.Sub(start).Microseconds()

	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d mks\n", max, elapsed)
}
