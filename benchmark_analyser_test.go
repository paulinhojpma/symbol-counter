package response_analyser

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var symbols [][]byte = generateSymbols(100)

func BenchmarkAnalyse(b *testing.B) {
	da := NewDataAnalyser()
	ctx := context.Background()
	for b.Loop() {
		da.Analyse(ctx, symbols[rand.Intn(100)])
	}

}

func TestGetCurrentCounts(t *testing.T) {
	da := NewDataAnalyser()
	//fmt.Printf("symbos - [%s]\n", symbols)
	ma := countSymbols(symbols)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Go(func() {
			da.Analyse(ctx, symbols[i])
		})

	}
	wg.Wait()
	time.Sleep(time.Second)
	go cancel()
	da.Shutdown(ctx)
	res := da.GetCurrentCounts()
	fmt.Printf("expected - %+v\n", ma)
	fmt.Printf("result - %+v\n", res)
	err := checkIfCountIsEqual(ma, res)
	if err != nil {
		t.Error(err.Error())
	}

}

func symbolGenerator() []byte {
	simbols := []string{
		"AU_ieu13",
		"103956",
		"ghgqb",
		"10002",
		"a012ne",
		"ZX_qw91",
		"b77xkp",
		"009381",
		"mn_t45",
		"AB12_cd",
		"qwert9",
		"p0x_11",
		"781245",
		"lkmno",
		"tt_902",
		"X1_yZ9",
		"cde_88",
		"190022",
		"rtyuio",
		"mN_556",
		"AA01_bb",
		"zxc_90",
		"554433",
		"jklmn8",
		"uio_p12",
		"NN_445",
		"abc123",
		"X_0099",
		"we_777",
		"tgbnh",
		"998812",
		"plm_01",
		"qrty_55",
		"fgh_778",
		"zzA_100",
	}
	sizeSimbols := len(simbols)
	simbolsMap := map[string]string{}
	size := rand.Intn(30) + 1
	simStr := ""
	for i := 0; i < size; i++ {
		simbPicked := simbols[rand.Intn(sizeSimbols)]

		if _, ok := simbolsMap[simbPicked]; !ok {
			simbolsMap[simbPicked] = simbPicked
			if simStr != "" {
				simStr = fmt.Sprintf("%s\n%s", simStr, simbPicked)
			} else {
				simStr = simbPicked
			}
		}
	}

	return []byte(simStr)

}

func generateSymbols(num int) [][]byte {
	symbols := make([][]byte, 0)
	for i := 0; i < num; i++ {
		symbols = append(symbols, symbolGenerator())
	}
	return symbols
}

func countSymbols(simb [][]byte) map[string]uint64 {
	counts := map[string]uint64{}
	for _, v := range simb {
		bits := bytes.Split(v, []byte("\n"))
		for _, b := range bits {
			counts[string(b)]++
		}

	}
	return counts
}

func checkIfCountIsEqual(expected, compa map[string]uint64) error {
	if len(expected) != len(compa) {
		return errors.New("diferent sizes of maps")
	}
	for k := range expected {
		if expected[k] != compa[k] {
			return fmt.Errorf("error expected key - [%s], is [%d], not [%d]", k, expected[k], compa[k])
		}
	}
	return nil
}
