package main
import (
        "fmt"
        "github.com/ontio/ontology/core/store/leveldbstore"
        "time"
        "math/rand"
)

var dbFile = "./test"
var testLevelDB *leveldbstore.LevelDBStore

func main() {
	var err error
    testLevelDB, err = leveldbstore.NewLevelDBStore(dbFile)
	if err != nil {
		fmt.Printf("NewLevelDBStore error:%s\n", err)
		return
	}

    TestBenchPrepare()
    for i := 0; i < 10; i++ {
            TestBenchPut()
            TestRandomGet()
    }
}

func TestBenchPrepare() {
	startTime := time.Now()
	testPerformanceBench(10*1000*1000, 100)
	dur := time.Since(startTime)
	fmt.Printf("put 1G records: time: %v\n", dur.Round(time.Second))
}

func TestBenchPut() {
	startTime := time.Now()
	testPerformanceBench(100, 100)
	dur := time.Since(startTime)
	fmt.Printf("put 10,000 records:time: %v\n", dur.Round(time.Second))
}

func TestRandomGet() {
	startTime := time.Now()
	for i := 0; i < 100*1000; i++ {
		testLevelDB.Get(genRandBytes(64))
	}
	dur := time.Since(startTime)
	fmt.Printf("get 100,000 records:time: %v\n", dur.Round(time.Second))
}

var letterRunes = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func genRandBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return b
}

func testPerformanceBench(batchCount, txPerBatch int) {
	keySize := 64
	valueSize := 512
	for i := 0; i < batchCount; i++ {
		testLevelDB.NewBatch()
		for j := 0; j < txPerBatch; j++ {
			key := genRandBytes(keySize)
			value := genRandBytes(valueSize)
			testLevelDB.BatchPut(key, value)
		}
		testLevelDB.BatchCommit()
	}
}

