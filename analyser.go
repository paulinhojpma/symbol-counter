package response_analyser

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"maps"
	"math/rand/v2"
	"sync"
	"time"
)

type dataAnalyser struct {
	mu           sync.RWMutex
	codesMap     map[string]uint64
	subscriber   map[time.Duration][]subscriber
	restInterval map[time.Duration]int64
	done         chan bool
	sender       chan []byte
	receiver     chan []byte
	valuesSync   []*valueSync
	numMaps      int
	escalateTo   int
	wg           sync.WaitGroup
}

type valueSync struct {
	mu sync.RWMutex
	vm map[string]uint64
}
type subscriber struct {
	writer io.WriteCloser
	active bool
}

var subPool = sync.Pool{
	New: func() any {
		return make([]subscriber, 0)
	},
}
var poolBit = sync.Pool{
	New: func() any {
		return make([]byte, 1024)
	},
}

func NewDataAnalyser() DataExporter {
	da := &dataAnalyser{
		done:         make(chan bool),
		codesMap:     map[string]uint64{},
		subscriber:   map[time.Duration][]subscriber{},
		restInterval: map[time.Duration]int64{},
		sender:       make(chan []byte, 500),
		receiver:     make(chan []byte),
		numMaps:      30,
		valuesSync:   initializeNewMapSyncs(30),
		escalateTo:   5,
	}
	da.sendwriters()
	da.receiveData()
	return da
}

func (d *dataAnalyser) Analyse(ctx context.Context, content []byte) {

	select {
	case <-d.done:
		return
	default:
		d.sender <- content

	}
}

// Esta função deve retornar os valores atuais de cada código único. Ela não deve redefinir os contadores.
func (d *dataAnalyser) GetCurrentCounts() map[string]uint64 {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, v := range d.valuesSync {
		for k, value := range retrieMap(v) {
			d.codesMap[k] += value
		}
	}

	return d.codesMap
}

// Subscribe adiciona um novo assinante de relatório com o intervalo especificado
func (d *dataAnalyser) Subscribe(writer io.WriteCloser, interval time.Duration) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.subscriber[interval]; ok {
		d.subscriber[interval] = append(d.subscriber[interval], subscriber{
			writer: writer,
			active: true,
		})
		return nil
	}
	d.subscriber[interval] = []subscriber{{writer: writer}}
	d.restInterval[interval] = int64(interval / time.Second)
	return nil
}

// Espera-se que Shutdown seja chamado pelo aplicativo principal quando o servidor receber uma chamada de término com contexto para desligamento correto.
func (d *dataAnalyser) Shutdown(ctx context.Context) error {

	<-ctx.Done()
	d.done <- true
	//d.wg.Wait()
	//time.Sleep(time.Second)
	close(d.done)

	return nil
}

func (d *dataAnalyser) addToCodesMap(bits [][]byte) {

	for _, bit := range bits {
		d.inc(bit)
	}

}

func (d *dataAnalyser) inc(bit []byte) {
	randIndex := d.numMaps
	addToMapSync(string(bit), d.valuesSync[rand.IntN(randIndex)])

}

func addToMapSync(key string, v *valueSync) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.vm[key]++
}

func (d *dataAnalyser) sendwriters() {

	ticker := time.NewTicker(time.Second)

	go func() {
		for {

			select {
			case <-d.done:
				//fmt.Println("closing senders")
				d.mu.Lock()

				for _, subs := range d.subscriber {
					for _, s := range subs {
						s.writer.Close()
					}
				}
				d.mu.Unlock()
				d.GetCurrentCounts()
				return
			case <-ticker.C:

				byteArray := poolBit.Get().([]byte)
				byteArray, _ = json.Marshal(d.GetCurrentCounts())

				for key, subs := range d.subscriber {

					if d.restInterval[key] <= 0 {

						for _, s := range subs {

							_, err := s.writer.Write(byteArray)
							if err != nil {
								s.active = false
							}
						}

						d.restInterval[key] = int64(key / time.Second)
					} else {
						d.restInterval[key]--
					}

				}

				poolBit.Put(byteArray)
				////fmt.Println("size of channel - ", len(d.sender))
				// if len(d.sender) >= cap(d.sender)-500 {
				// 	d.escalateMapSync()
				// } else if len(d.sender) <= cap(d.sender)/2 {
				// 	d.desescalateMapSync()
				// }

			}
		}
	}()

}

func (d *dataAnalyser) receiveData() {

	go func() {

		for {
			content := poolBit.Get().([]byte)
			select {
			case <-d.done:
				d.wg.Wait()
				return
			case content = <-d.sender:
				bits := bytes.Split(content, []byte("\n"))
				d.wg.Go(func() { d.addToCodesMap(bits) })
				//go d.addToCodesMap(bits)

			}
			content = make([]byte, 0)
			poolBit.Put(content)
		}

	}()

}

func initializeNewMapSyncs(num int) []*valueSync {
	values := make([]*valueSync, num)
	for i := 0; i < num; i++ {
		values[i] = &valueSync{
			vm: map[string]uint64{},
		}
	}
	return values
}

func retrieMap(v *valueSync) map[string]uint64 {
	v.mu.RLock()
	defer v.mu.RUnlock()
	m := maps.Clone(v.vm)
	v.vm = map[string]uint64{}
	return m
}

func (d *dataAnalyser) removeInactiveSubscribers(subscribers []subscriber) []subscriber {
	leftSub := subPool.Get().([]subscriber)
	leftSub = make([]subscriber, 0)
	for _, sub := range subscribers {
		if sub.active {
			leftSub = append(leftSub, sub)
		}

	}
	subscribers = leftSub
	subPool.Put(leftSub)
	return subscribers
}
