package cracker

import (
	"fmt"
	"os/exec"
	"runtime"
	"sync"

	"github.com/andream16/mitmcracker/internal/repository"
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
		keysNumber:     getKeyNumber(keyLenght),
		repository:     repository,
	}
}

type encDecTask struct {
	cmd      *exec.Cmd
	key      string
	taskType string
}

type repoTask struct {
	output   []byte
	key      string
	taskType string
}

// Crack executes the cracker process.
func (c *Cracker) Crack() (string, error) {

	var (
		wgTasks     sync.WaitGroup
		wgRepo      sync.WaitGroup
		encDecTasks = make(chan encDecTask, c.maxParallelism)
		repoTasks   = make(chan repoTask, c.maxParallelism)
	)

	for k := 0; k <= c.keysNumber; k++ {
		fK := formatKey(k, c.keyLenght)
		encDecTasks <- encDecTask{
			cmd:      encode(fK, c.plainText),
			key:      fK,
			taskType: encodeTask,
		}
		encDecTasks <- encDecTask{
			cmd:      decode(fK, c.plainText),
			key:      fK,
			taskType: decodeTask,
		}
	}

	for t := range encDecTasks {
		wgTasks.Add(1)
		var (
			b   []byte
			err error
		)
		go func(cmd *exec.Cmd, wg *sync.WaitGroup) {
			defer wg.Done()
			b, err = cmd.Output()
		}(t.cmd, &wgTasks)
		if err != nil {
			return "", err
		}
		repoTasks <- repoTask{
			output:   b,
			key:      t.key,
			taskType: t.taskType,
		}
	}

	for t := range repoTasks {
		wgRepo.Add(1)
		go func(task repoTask, wg *sync.WaitGroup) {
			defer wg.Done()
			switch task.taskType {
			case encodeTask:
				c.repository.InsertEnc(task.key, string(task.output))
			case decodeTask:
				c.repository.InsertDec(task.key, string(task.output))
			}
		}(t, &wgRepo)
	}

	close(encDecTasks)
	close(repoTasks)
	wgTasks.Wait()
	wgRepo.Wait()

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

func getKeyNumber(keyLenght uint) int {
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
