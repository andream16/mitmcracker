package cracker

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"sync"
	"time"

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
	inserter      repository.Inserter
	plainText     string
	cipherText    string
}

type task struct {
	key  string
	text string
	mode Mode
	fn   func(key, plainText string) (string, error)
}

// New returns a new cracker tuned for a particular system capacity & key length.
func New(keyLength uint, inserter repository.Inserter, plainText, cipherText string) (*Cracker, error) {
	keyNum, err := getKeyNumber(keyLength)
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
	}, nil
}

// Crack returns the matching key pair. False is returned when the pair is not found.
func (c *Cracker) Crack() (*KeyPair, bool, error) {

	var (
		taskNum int
		tasks   = make(chan task, c.keyNum)
		wg      sync.WaitGroup
		keyPair *KeyPair
		found   bool
	)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			log.Println(
				fmt.Sprintf(
					"%d tasks out of %d.",
					taskNum,
					c.keyNum*2,
				),
			)
		}
	}()

	wg.Add(c.goRoutinesNum)

	for i := 0; i <= c.goRoutinesNum; i++ {
		go func(i int, inFound *bool, inKp *KeyPair, wg *sync.WaitGroup) {
			for task := range tasks {
				taskNum++
				defer wg.Done()

				cipherText, err := task.fn(task.key, task.text)
				if err != nil {
					log.Println(fmt.Sprintf("unexpected error while processing task %d: %v", taskNum, err))
					continue
				}

				kp, found, err := c.inserter.Insert(task.key, cipherText, repository.Mode(task.mode))
				if err != nil {
					log.Println(fmt.Sprintf("unexpected error while inserting key for task %d: %v", taskNum, err))
					continue
				}

				if found {
					inFound = &found
					inKp = &KeyPair{
						EncodeKey: kp.EncodeKey,
						DecodeKey: kp.DecodeKey,
					}
				}
			}
		}(i, &found, keyPair, &wg)
	}

	for k := 0; k <= c.keyNum; k++ {
		if found {
			return keyPair, true, nil
		}
		fK := formatKey(k, c.keyLength)
		tasks <- task{
			key:  fK,
			text: c.plainText,
			mode: encodeMode,
			fn:   encode,
		}
		tasks <- task{
			key:  fK,
			text: c.cipherText,
			mode: decodeMode,
			fn:   decode,
		}
	}

	close(tasks)
	wg.Wait()

	return nil, false, nil
}

func encode(key, plainText string) (string, error) {
	out, err := exec.Command("./resources/encrypt", "-s", key, plainText).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func decode(key, cipherText string) (string, error) {
	out, err := exec.Command("./resources/decrypt", "-s", key, cipherText).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func getMaxConcurrency() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func getKeyNumber(keyLength uint) (int, error) {
	if l, ok := map[uint]int{
		24: 16777216,
		28: 268435456,
		32: 4294967296,
	}[keyLength]; ok {
		return l, nil
	}
	return 0, fmt.Errorf("unexpected key length %d. Valid key legths are 24, 28 and 32 bits", keyLength)
}

func formatKey(key int, keyLength uint) string {
	var s string
	switch keyLength {
	case 24:
		s = "%06x"
	case 28:
		s = "%07x"
	case 32:
		s = "%08x"
	}
	return fmt.Sprintf(s, key)
}
