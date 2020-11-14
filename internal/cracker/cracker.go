package cracker

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/andream16/mitmcracker/internal/keycalculator"
	"gopkg.in/cheggaaa/pb.v1"

	"github.com/andream16/mitmcracker/internal/decrypter"
	"github.com/andream16/mitmcracker/internal/encrypter"
	"github.com/andream16/mitmcracker/internal/formatter"
	"github.com/andream16/mitmcracker/internal/repository"
)

const (
	encodeMode Mode = "encode"
	decodeMode Mode = "decode"
)

// Mode represents the operation mode.
type Mode string

// Cracker represent a MITM cracker.
type Cracker struct {
	keyNum        int
	keyLength     uint
	goRoutinesNum int
	plainText     string
	cipherText    string
	inserter      repository.Inserter
	encFn         encrypter.Encrypter
	decFn         decrypter.Decrypter
	formatterFn   formatter.Formatter
	keyCalcFn     keycalculator.KeyCalculator
}

type task struct {
	key  string
	text string
	mode Mode
	fn   func(key, plainText string) (string, error)
}

// New returns a new cracker tuned for a particular system capacity & key length.
func New(
	keyLength uint,
	cipherText string,
	plainText string,
	inserter repository.Inserter,
	encFn encrypter.Encrypter,
	decFn decrypter.Decrypter,
	formatterFn formatter.Formatter,
	keyCalcFn keycalculator.KeyCalculator,
) (*Cracker, error) {
	keyNum, err := keyCalcFn(keyLength)
	if err != nil {
		return nil, err
	}

	return &Cracker{
		keyNum:        keyNum,
		keyLength:     keyLength,
		goRoutinesNum: getMaxConcurrency(),
		inserter:      inserter,
		plainText:     plainText,
		cipherText:    cipherText,
		encFn:         encFn,
		decFn:         decFn,
		formatterFn:   formatterFn,
		keyCalcFn:     keyCalcFn,
	}, nil
}

// Crack returns the matching key pair. False is returned when the pair is not found.
func (c *Cracker) Crack() (*KeyPair, bool, error) {

	var (
		taskCtr int
		tasks   = make(chan task)
		taskBar = pb.StartNew(c.keyNum * 2)
		wg      sync.WaitGroup
		keyPair KeyPair
		found   bool
	)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			taskBar.Set(taskCtr)
		}
	}()

	for i := 0; i <= c.goRoutinesNum; i++ {
		go func(i int, refFound *bool, refKp *KeyPair, refTaskCounter *int, wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()
			for task := range tasks {
				*refTaskCounter++

				cipherText, err := task.fn(task.key, task.text)
				if err != nil {
					log.Println(
						fmt.Sprintf(
							"unexpected error while processing task %d: %v",
							refTaskCounter,
							err,
						),
					)
					continue
				}

				//log.Println(
				//	fmt.Sprintf(
				//		"run %s with key %s on text %s. Got common ciphertext %s",
				//		task.mode,
				//		task.key,
				//		task.text,
				//		cipherText,
				//	),
				//)

				kp, wasFound, err := c.inserter.Insert(task.key, cipherText, repository.Mode(task.mode))
				if err != nil {
					log.Println(
						fmt.Sprintf(
							"unexpected error while inserting key for task %d: %v",
							refTaskCounter,
							err,
						),
					)
					continue
				}

				if wasFound {
					*refFound = true
					*refKp = KeyPair{
						EncodeKey: kp.EncodeKey,
						DecodeKey: kp.DecodeKey,
					}
				}
			}
		}(i, &found, &keyPair, &taskCtr, &wg)
	}

	for k := 0; k < c.keyNum; k++ {
		if found {
			return &keyPair, true, nil
		}
		fK := c.formatterFn(k, c.keyLength)
		tasks <- task{
			key:  fK,
			text: c.plainText,
			mode: encodeMode,
			fn:   c.encFn,
		}
		tasks <- task{
			key:  fK,
			text: c.cipherText,
			mode: decodeMode,
			fn:   c.decFn,
		}
	}

	close(tasks)
	wg.Wait()

	return &keyPair, found, nil
}

func getMaxConcurrency() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}
