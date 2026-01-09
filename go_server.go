// main.go - Enhanced Go Server
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "math"
    "math/big"
    "net/http"
    "runtime"
    "time"
)

type CalculationRequest struct {
	Number int `json:"number"`
}

type EnhancedResponse struct {
	Result            interface{} `json:"result"`
	TimeTaken         int64       `json:"timeTaken"`         // in nanoseconds
	MemoryUsage       uint64      `json:"memoryUsage"`       // in bytes
	ExecutionTime     int64       `json:"executionTime"`     // in nanoseconds
	GoroutinesUsed    int         `json:"goroutinesUsed"`    // Goroutines count
	BinarySize        string      `json:"binarySize"`        // Approx binary size
	DependencySize    string      `json:"dependencySize"`    // Go modules size
	StartupTime       int64       `json:"startupTime"`       // in milliseconds
	Throughput        float64     `json:"throughput"`        // ops/sec
	Error             string      `json:"error,omitempty"`
}

var startupTime time.Time

func init() {
	startupTime = time.Now()
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// 1. Factorial (Big Integer)
func factorial(n int) *big.Int {
	if n < 0 {
		return big.NewInt(0)
	}
	if n == 0 {
		return big.NewInt(1)
	}

	result := big.NewInt(1)
	for i := 1; i <= n; i++ {
		result.Mul(result, big.NewInt(int64(i)))
	}
	return result
}

// 2. Fibonacci (Big Integer)
func fibonacci(n int) *big.Int {
	if n <= 0 {
		return big.NewInt(0)
	}
	if n == 1 {
		return big.NewInt(1)
	}

	a := big.NewInt(0)
	b := big.NewInt(1)
	for i := 2; i <= n; i++ {
		a.Add(a, b)
		a, b = b, a
	}
	return b
}

// 3. Prime Counting
func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func countPrimes(limit int) int {
	count := 0
	for i := 2; i <= limit; i++ {
		if isPrime(i) {
			count++
		}
	}
	return count
}

// 4. Matrix Multiplication (CPU Intensive)
func matrixMultiply(size int) [][]float64 {
	// Create matrices
	a := make([][]float64, size)
	b := make([][]float64, size)
	result := make([][]float64, size)
	
	for i := 0; i < size; i++ {
		a[i] = make([]float64, size)
		b[i] = make([]float64, size)
		result[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			a[i][j] = float64(i + j + 1)
			b[i][j] = float64(i - j + 1)
		}
	}
	
	// Multiply matrices
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			sum := 0.0
			for k := 0; k < size; k++ {
				sum += a[i][k] * b[k][j]
			}
			result[i][j] = sum
		}
	}
	
	return result
}

// 5. Quick Sort (Algorithm Performance)
func quickSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	
	pivot := arr[len(arr)/2]
	var left, right []int
	
	for _, v := range arr {
		if v < pivot {
			left = append(left, v)
		} else if v > pivot {
			right = append(right, v)
		}
	}
	
	left = quickSort(left)
	right = quickSort(right)
	
	return append(append(left, pivot), right...)
}

// 6. String Operations (Memory Intensive)
func stringOperations(count int) string {
	var builder string
	for i := 0; i < count; i++ {
		builder += fmt.Sprintf("String%d: Lorem ipsum dolor sit amet, consectetur adipiscing elit. ", i)
		if i%100 == 0 {
			builder += "\n"
		}
	}
	return builder
}

// 7. Concurrent Fibonacci (Goroutine Demo)
func concurrentFibonacci(n int) *big.Int {
	if n <= 1000 {
		return fibonacci(n)
	}
	
	ch := make(chan *big.Int, 2)
	
	go func() {
		ch <- fibonacci(n/2)
	}()
	
	go func() {
		ch <- fibonacci(n - n/2)
	}()
	
	a := <-ch
	b := <-ch
	return new(big.Int).Add(a, b)
}

// 8. JSON Serialization/Deserialization
func jsonOperations(count int) interface{} {
	type Data struct {
		ID    int      `json:"id"`
		Name  string   `json:"name"`
		Items []string `json:"items"`
		Value float64  `json:"value"`
	}
	
	var data []Data
	for i := 0; i < count; i++ {
		data = append(data, Data{
			ID:    i,
			Name:  fmt.Sprintf("Item%d", i),
			Items: []string{"a", "b", "c", "d", "e"},
			Value: math.Sin(float64(i)),
		})
	}
	
	jsonBytes, _ := json.Marshal(data)
	var decoded []Data
	json.Unmarshal(jsonBytes, &decoded)
	
	return struct {
		OriginalSize  int `json:"originalSize"`
		JsonSize      int `json:"jsonSize"`
		DecodedLength int `json:"decodedLength"`
	}{
		OriginalSize: len(data),
		JsonSize:     len(jsonBytes),
		DecodedLength: len(decoded),
	}
}

// 9. Memory Intensive Allocation
func memoryIntensive(size int) [][]int {
    matrix := make([][]int, size)
    for i := 0; i < size; i++ {
        matrix[i] = make([]int, size)
        for j := 0; j < size; j++ {
            matrix[i][j] = i * j
        }
    }
    return matrix
}

func measurePerformance(calcFunc func(int) interface{}, n int) EnhancedResponse {
	var m1, m2 runtime.MemStats
	
	// Force GC and get initial stats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	goroutinesBefore := runtime.NumGoroutine()
	
	start := time.Now()
	result := calcFunc(n)
	executionTime := time.Since(start).Nanoseconds()
	
	runtime.ReadMemStats(&m2)
	goroutinesAfter := runtime.NumGoroutine()
	
	memoryUsage := m2.Alloc - m1.Alloc
	
	// Calculate throughput (ops/sec)
	throughput := float64(0)
	if executionTime > 0 {
		throughput = 1e9 / float64(executionTime) // 1 second / execution time
	}
	
	// Approximate binary size (Go produces static binaries)
	binarySize := "8-15 MB"
	dependencySize := "Varies (go.mod based)"
	
	// Calculate startup time
	startupDuration := time.Since(startupTime).Milliseconds()
	
	return EnhancedResponse{
		Result:         result,
		TimeTaken:      executionTime / 1e6, // Convert to milliseconds
		MemoryUsage:    memoryUsage,
		ExecutionTime:  executionTime,
		GoroutinesUsed: goroutinesAfter - goroutinesBefore,
		BinarySize:     binarySize,
		DependencySize: dependencySize,
		StartupTime:    startupDuration,
		Throughput:     throughput,
	}
}

func createHandler(endpoint string, calcFunc func(int) interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		
		if r.Method == "OPTIONS" {
			return
		}
		
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		var req CalculationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		
		response := measurePerformance(func(n int) interface{} {
			return calcFunc(req.Number)
		}, req.Number)
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	// Register all endpoints
	http.HandleFunc("/api/factorial", createHandler("factorial", func(n int) interface{} {
		return factorial(n).String()
	}))
	
	http.HandleFunc("/api/fibonacci", createHandler("fibonacci", func(n int) interface{} {
		return fibonacci(n).String()
	}))
	
	http.HandleFunc("/api/primes", createHandler("primes", func(n int) interface{} {
		return countPrimes(n)
	}))
	
	http.HandleFunc("/api/matrix", createHandler("matrix", func(n int) interface{} {
		size := n
		if size > 200 {
			size = 200 // Limit for safety
		}
		result := matrixMultiply(size)
		return struct {
			MatrixSize int `json:"matrixSize"`
			FirstValue float64 `json:"firstValue"`
			LastValue float64 `json:"lastValue"`
		}{
			MatrixSize: size,
			FirstValue: result[0][0],
			LastValue: result[size-1][size-1],
		}
	}))
	
	http.HandleFunc("/api/sort", createHandler("sort", func(n int) interface{} {
		// Create large array to sort
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = n - i
		}
		sorted := quickSort(arr)
		return struct {
			ArraySize int `json:"arraySize"`
			FirstElement int `json:"firstElement"`
			LastElement int `json:"lastElement"`
		}{
			ArraySize: len(sorted),
			FirstElement: sorted[0],
			LastElement: sorted[len(sorted)-1],
		}
	}))
	
	http.HandleFunc("/api/string", createHandler("string", func(n int) interface{} {
		str := stringOperations(n)
		return struct {
			Length int `json:"length"`
			First100 string `json:"first100"`
		}{
			Length: len(str),
			First100: str[:min(100, len(str))],
		}
	}))
	
	http.HandleFunc("/api/concurrent", createHandler("concurrent", func(n int) interface{} {
		result := concurrentFibonacci(n)
		return result.String()
	}))
	
	http.HandleFunc("/api/json", createHandler("json", func(n int) interface{} {
		return jsonOperations(n)
	}))
	
	http.HandleFunc("/api/memory", createHandler("memory", func(n int) interface{} {
		size := int(math.Sqrt(float64(n)))
		if size > 100 {
			size = 100
		}
		// matrix := memoryIntensive(size)
		return struct {
			MatrixSize int `json:"matrixSize"`
			TotalElements int `json:"totalElements"`
			MemoryEstimate string `json:"memoryEstimate"`
		}{
			MatrixSize: size,
			TotalElements: size * size,
			MemoryEstimate: fmt.Sprintf("%d MB", (size*size*8)/(1024*1024)),
		}
	}))
	
	// Stats endpoint
	http.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		
		stats := struct {
			Goroutines   int     `json:"goroutines"`
			HeapAlloc    uint64  `json:"heapAlloc"`
			StackInUse   uint64  `json:"stackInUse"`
			NumCPU       int     `json:"numCPU"`
			NumCgoCall   int64   `json:"numCgoCall"`
			Uptime       int64   `json:"uptime"`
		}{
			Goroutines:   runtime.NumGoroutine(),
			HeapAlloc:    m.HeapAlloc,
			StackInUse:   m.StackInuse,
			NumCPU:       runtime.NumCPU(),
			NumCgoCall:   runtime.NumCgoCall(),
			Uptime:       time.Since(startupTime).Milliseconds(),
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	})
	
	fmt.Println("🚀 Go server running on http://localhost:8081")
	fmt.Println("📊 Available endpoints:")
	fmt.Println("  /api/factorial   - Big integer factorial")
	fmt.Println("  /api/fibonacci   - Fibonacci sequence")
	fmt.Println("  /api/primes      - Prime counting")
	fmt.Println("  /api/matrix      - Matrix multiplication")
	fmt.Println("  /api/sort        - QuickSort algorithm")
	fmt.Println("  /api/string      - String operations")
	fmt.Println("  /api/concurrent  - Concurrent Fibonacci")
	fmt.Println("  /api/json        - JSON serialization")
	fmt.Println("  /api/memory      - Memory intensive")
	fmt.Println("  /api/stats       - Runtime statistics")
	
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}