package cracker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/andream16/mitmcracker/internal/decrypter"
	"github.com/andream16/mitmcracker/internal/encrypter"
	"github.com/andream16/mitmcracker/internal/formatter"
	"github.com/andream16/mitmcracker/internal/keycalculator"
	"github.com/andream16/mitmcracker/internal/perf"
	"github.com/andream16/mitmcracker/internal/repository"

	"golang.org/x/sync/errgroup"
	"gopkg.in/cheggaaa/pb.v1"
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
	perfNumGoRoutines perf.MaxGoRoutineNumber,
) (*Cracker, error) {
	keyNum, err := keyCalcFn(keyLength)
	if err != nil {
		return nil, err
	}

	return &Cracker{
		keyNum:        keyNum,
		keyLength:     keyLength,
		goRoutinesNum: perfNumGoRoutines(),
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
func (c *Cracker) Crack(ctx context.Context) (*KeyPair, bool, error) {

	var (
		taskCtr     int
		tasks       = make(chan task)
		result      = make(chan KeyPair, 1)
		taskCtrChan = make(chan struct{})
		taskBar     = pb.StartNew(c.keyNum * 2)
		keyPair     KeyPair
		found       bool
	)

	defer taskBar.Finish()

	ctx, cancel := context.WithCancel(ctx)
	g, ctx := errgroup.WithContext(ctx)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			taskBar.Set(taskCtr)
		}
	}()

	// insight on progress
	g.Go(func() error {
		for {
			select {
			case <-taskCtrChan:
				taskCtr++
			case <-ctx.Done():
				return ctx.Err()
			default:
				break
			}
		}
	})

	// check for task progress and end
	g.Go(func() error {
		for {
			select {
			case r := <-result:
				found = true
				keyPair = r
				cancel()
				break
			}
			return nil
		}
	})

	// consume
	for i := 0; i < c.goRoutinesNum; i++ {
		g.Go(func() error {
			for {
				select {
				case t := <-tasks:
					taskCtrChan <- struct{}{}

					kp, wasFound, err := c.handleTask(t)
					if err != nil {
						log.Println(fmt.Sprintf("error while processing task %v", err))
						break
					}
					if wasFound {
						result <- kp
					}
				case <-ctx.Done():
					return ctx.Err()
				default:
					break
				}
			}
		})
	}

	// produce
	g.Go(func() error {
		for k := 0; k < c.keyNum; k++ {
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

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				break
			}
		}

		return nil
	})

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return nil, false, err
	}

	close(tasks)
	close(result)
	close(taskCtrChan)

	return &keyPair, found, nil
}

func (c *Cracker) handleTask(t task) (KeyPair, bool, error) {
	var res KeyPair

	cipherText, err := t.fn(t.key, t.text)
	if err != nil {
		return res, false, err
	}

	kp, wasFound, err := c.inserter.Insert(t.key, cipherText, repository.Mode(t.mode))
	if err != nil {
		return res, false, err
	}

	if kp != nil && kp.DecodeKey != "" && kp.EncodeKey != "" {
		res = KeyPair{
			EncodeKey: kp.EncodeKey,
			DecodeKey: kp.DecodeKey,
		}
	}

	return res, wasFound, err
}
