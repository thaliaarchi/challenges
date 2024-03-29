package main

import (
	"math/big"
	"testing"
)

func TestFib(t *testing.T) {
	expect := toBig([]int64{0, 1, 1, 2, 3, 5, 8, 13, 21, 34,
		55, 89, 144, 233, 377, 610, 987, 1597, 2584, 4181})
	num := Fib()
	for i := range expect {
		if n := <-num; n.Cmp(expect[i]) != 0 {
			t.Errorf("%d: got %s, want %s", i+1, n, expect[i])
		}
	}
}

func TestFibCached(t *testing.T) {
	fib := FibCached()
	fib50 := fib(50)
	fib100 := fib(100)
	if len(fib50) != 50 || len(fib100) != 100 {
		t.Fatalf("got len(fib50) = %d and len(fib100) = %d, want 50 and 100", len(fib50), len(fib100))
	}
	for i := range fib50 {
		if fib50[i] != fib100[i] {
			t.Fatalf("%d: result not cached, got %s and %s", i, fib50[i], fib100[i])
		}
	}
}

func TestFibCachedValues(t *testing.T) {
	tests := []struct {
		index int
		num   string
	}{
		{0, "0"},
		{24, "46368"},                // overflows int16
		{25, "75025"},                // overflows uint16
		{37, "24157817"},             // overflows precise float32
		{47, "2971215073"},           // overflows int32
		{48, "4807526976"},           // overflows uint32
		{79, "14472334024676221"},    // overflows precise float64
		{93, "12200160415121876738"}, // overflows int64
		{94, "19740274219868223167"}, // overflows uint64
		{100, "354224848179261915075"},
		{200, "280571172992510140037611932413038677189525"},
		{300, "222232244629420445529739893461909967206666939096499764990979600"},
		{400, "176023680645013966468226945392411250770384383304492191886725992896575345044216019675"},
		{500, "139423224561697880139724382870407283950070256587697307264108962948325571622863290691557658876222521294125"},
		{600, "110433070572952242346432246767718285942590237357555606380008891875277701705731473925618404421867819924194229142447517901959200"},
		{700, "87470814955752846203978413017571327342367240967697381074230432592527501911290377655628227150878427331693193369109193672330777527943718169105124275"},
		{800, "69283081864224717136290077681328518273399124385204820718966040597691435587278383112277161967532530675374170857404743017623467220361778016172106855838975759985190398725"},
		{900, "54877108839480000051413673948383714443800519309123592724494953427039811201064341234954387521525390615504949092187441218246679104731442473022013980160407007017175697317900483275246652938800"},
		{1000, "43466557686937456435688527675040625802564660517371780402481729089536555417949051890403879840079255169295922593080322634775209689623239873322471161642996440906533187938298969649928516003704476137795166849228875"},
	}

	fib := FibCached()
	nums := fib(1001)
	if len(nums) != 1001 {
		t.Fatalf("got len %d, want 1001", len(nums))
	}
	for i, test := range tests {
		num, ok := new(big.Int).SetString(test.num, 10)
		if !ok {
			t.Errorf("test %d: unable to set string", i+1)
			continue
		}
		if nums[test.index].Cmp(num) != 0 {
			t.Errorf("test %d: got %s, want %s", i+1, nums[test.index], num)
		}
	}
}

func TestFibCachedEqualFibNth(t *testing.T) {
	nums := FibCached()(100)
	for i := range nums {
		if num := FibNth(uint64(i)); num.Cmp(nums[i]) != 0 {
			t.Errorf("%d: results not equal, got %s and %s", i, num, nums[i])
		}
	}
}

func BenchmarkFib10000(b *testing.B) {
	num := Fib()
	for i := 0; i < 10000; i++ {
		<-num
	}
}

func BenchmarkFibNth10000(b *testing.B) {
	FibNth(10000)
}

func toBig(nums []int64) []*big.Int {
	b := make([]*big.Int, 0, len(nums))
	for _, num := range nums {
		b = append(b, big.NewInt(num))
	}
	return b
}
