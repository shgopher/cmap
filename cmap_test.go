package cmap

import (
	"fmt"
	"github.com/spaolacci/murmur3"
	"strconv"
	"sync"
	"testing"
)


// 简单的对结果进行一个测试
func TestCmap(t *testing.T) {
	cp := NewCmap()
	cp.Set("a", 1)
	cp.Set("b", 2)
	cp.Set("c", 3)
	fmt.Println(cp.Get("a"), cp.Get("c"), cp.Get("b"))
}

// 对设置key - value进行一个测试，
// 模拟100个goroutine进行写
func TestCmap_Set_Get(t *testing.T) {
	cp := NewCmap()
	wg := new(sync.WaitGroup)
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			cp.Set(strconv.FormatInt(int64(i),10),i)
		}(i)
	}
	wg.Wait()
	wg.Add(100)
	count := 0
	for i:= 0;i < 100;i++ {
		go func(i int) {
			defer wg.Done()
			count++
			fmt.Println(cp.Get(strconv.FormatInt(int64(i),10)))
		}(i)
	}
	wg.Wait()
	fmt.Println(count)
}

// 对获取一个 key - value 进行测试
// 模拟 100个goroutine进行读
func TestCmap_Get(t *testing.T) {
	cp := NewCmap()
	wg := new(sync.WaitGroup)
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			cp.Get(strconv.FormatInt(int64(i),10))
		}(i)
	}
	wg.Wait()

}


// 使用sync.mutex的一般大锁进行测试。
func TestCompare(t *testing.T) {
	wg := new(sync.WaitGroup)
	mu := new(sync.Mutex)
	wg.Add(100)
	testMap := make(map[string]interface{})
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			mu.Lock()
			testMap[strconv.FormatInt(int64(i),10)]=i
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	wg.Add(1000)
	count := 0
	for i:= 0;i < 1000;i++ {
		go func(i int) {
			defer wg.Done()
			mu.Lock()
			count++
			fmt.Println(testMap[strconv.FormatInt(int64(i),10)])
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	fmt.Println(count)
}
//goos: darwin
//goarch: amd64
//pkg: github.com/shgopher/cmap
//BenchmarkCamp
//BenchmarkCamp-4   	 1581716	       673 ns/op
//PASS
func BenchmarkCamp(b *testing.B) {
	cp := NewCmap()
	for i := 0; i < b.N; i++ {
		t := i
		go cp.Set(strconv.FormatInt(int64(t),10),t)
	}
}
//goos: darwin
//goarch: amd64
//pkg: github.com/shgopher/cmap
//BenchmarkSyncMutex
//BenchmarkSyncMutex-4   	  712914	      4213 ns/op
//PASS
func BenchmarkSyncMutex(b *testing.B) {
	wg := new(sync.WaitGroup)
	mu := new(sync.Mutex)
	testMap := make(map[string]interface{})
	for i := 0; i < b.N; i++ {
		wg.Add(1)
			go func(i int) {
				defer wg.Done()
				mu.Lock()
				testMap[strconv.FormatInt(int64(i),10)]=i
				mu.Unlock()
			}(i)
	}
	wg.Wait()
}

func BenchmarkFnvHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FnvHash(strconv.FormatInt(int64(i),10))
	}
}

func BenchmarkMurMur(b *testing.B) {
	for i := 0; i < b.N; i++ {
		murmur3.Sum32([]byte(strconv.FormatInt(int64(i),10)))
	}
}
