package cracker

import (
	"fmt"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/andream16/mitmcracker/internal/repository"

	"gopkg.in/cheggaaa/pb.v1"
)

const (
	encodeTask = "encode"
	decodeTask = "decode"
)

var keyLenghtToKeyNumber = map[uint]int{
	24: 16777216,
	28: 268435456,
	32: 4294967296,
}

// Cracker contains the cracker configuration.
type Cracker struct {
	// plainText is the known plain text.
	plainText string
	// encText is the known encoded text.
	encText string
	// maxParallelism is the maximum number of go routines that
	// can be spawned on runtime.
	maxParallelism int
	// keyLenght is the key lenght to be used.
	keyLenght uint
	// keysNumber is the number of keys to be generated.
	keysNumber int
	// repository is the storage to be used.
	repository repository.Repositorer
}

// New returns a new Cracker.
func New(
	plainText string,
	encText string,
	keyLenght uint,
	repository repository.Repositorer,
) *Cracker {
	return &Cracker{
		plainText:      plainText,
		encText:        encText,
		maxParallelism: getMaxParallelism(),
		keyLenght:      keyLenght,
		keysNumber:     GetKeyNumber(keyLenght),
		repository:     repository,
	}
}

type task struct {
	cmd      *exec.Cmd
	key      string
	taskType string
}

// Crack executes the cracker process.
func (c *Cracker) Crack() (string, error) {

	var (
		bar            = pb.StartNew(c.keysNumber)
		maxConcurrency = getMaxParallelism()
		tasks          = make(chan task)
		tK             = 0
		wg             sync.WaitGroup
	)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			bar.Add(tK)
		}
	}()

	defer bar.Finish()

	wg.Add(maxConcurrency)

	for i := 0; i <= maxConcurrency; i++ {
		go func(i int, wg *sync.WaitGroup) {
			for task := range tasks {
				defer wg.Done()
				b, err := task.cmd.Output()
				if err != nil {
					continue
				}
				switch task.taskType {
				case encodeTask:
					c.repository.InsertEnc(task.key, string(b))
				case decodeTask:
					c.repository.InsertDec(task.key, string(b))
				}
			}
		}(i, &wg)
	}

	for k := 0; k <= c.keysNumber; k++ {
		tK = k
		fK := formatKey(k, c.keyLenght)
		tasks <- task{
			cmd:      encode(fK, c.plainText),
			key:      fK,
			taskType: encodeTask,
		}
		tasks <- task{
			cmd:      decode(fK, c.plainText),
			key:      fK,
			taskType: decodeTask,
		}
	}

	close(tasks)
	wg.Wait()

	return c.repository.FindKey(), nil
}

func encode(key, plainText string) *exec.Cmd {
	return exec.Command("./resources/encrypt", "-s", key, plainText)
}

func decode(key, encText string) *exec.Cmd {
	return exec.Command("./resources/decrypt", "-s", key, encText)
}

func getMaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

// GetKeyNumber returns the number of keys to be generated.
func GetKeyNumber(keyLenght uint) int {
	k, _ := keyLenghtToKeyNumber[keyLenght]
	return k
}

func formatKey(key int, keyLenght uint) string {
	frmt := ""
	switch keyLenght {
	case 24:
		frmt = "%06x"
	case 28:
		frmt = "%07x"
	case 32:
		frmt = "%08x"
	}
	return fmt.Sprintf(frmt, key)
}
