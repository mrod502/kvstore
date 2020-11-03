package kvstore

import (
	"fmt"
	"io/ioutil"
	"math"
	"testing"
	"time"

	"github.com/rs/xid"
)

func TestKV(t *testing.T) {
	testData, err := ioutil.ReadFile("stocktwits_pretty.json")
	if err != nil {
		t.Fatal(err)
	}

	million := 1000000
	iters := 10 * million
	ids := make([]string, iters)
	s := NewStore(true)
	for i := 0; i < iters; i++ {
		id := xid.New().String()
		s.Set(id, FooData{Data: testData}, time.Now().Add(5*time.Second).Unix())
		ids[i] = id
	}
	var tNow time.Time
	var timeArr = make([]time.Duration, iters)

	for i, id := range ids {
		tNow = time.Now()
		s.Get(id)
		timeArr[i] = time.Since(tNow)
	}

	//calculate stats
	stdev, avg, pt95, pt99 := calcTimeStats(timeArr)

	fmt.Printf("Standard Dev in access time:%v\nAverage time to access:%v\n95th percentile (calculated using stDev):%v\n99th percentile (calculated using stDev):%v\n", stdev, avg, pt95, pt99)

	tNow = time.Now()
	for _, i := range ids {
		s.Delete(i)
	}
	//test time taken to delete 1,000,000 entries
	fmt.Println("average time to delete:", time.Since(tNow)/time.Duration(iters))

	//test concurrent access
	go func() {
		for _, i := range ids {
			s.Get(i)
		}
	}()
	for _, i := range ids {
		s.Delete(i)
	}

}

//FooData - a test data struct
type FooData struct {
	Data []byte
}

//GetData - implement the ByteObj interface
func (f FooData) GetData() []byte {
	return f.Data
}

func calcTimeStats(in []time.Duration) (stDev, avg, pt95, pt99 time.Duration) {
	var total time.Duration
	var variance time.Duration
	for _, v := range in {
		total += v
	}

	avg = total / time.Duration(len(in))

	for _, v := range in {
		variance += (v - avg) * (v - avg)
	}
	variance /= time.Duration(len(in))

	stDev = time.Duration(math.Sqrt(float64(variance)))
	pt95 = avg + ((165 * stDev) / 100) //95th percentile about 1.65 standar deviations above mean
	pt99 = avg + ((233 * stDev) / 100) //95th percentile about 1.65 standar deviations above mean
	return
}
