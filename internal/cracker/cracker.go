package cracker

import (
	"context"
	"errors"
	"fmt"
	"log"

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
	encodeMode = "encode"
	decodeMode = "decode"
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
	bufferSize    int
}

type task struct {
	key  string
	text string
	mode string
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
	bufferSize int,
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
		bufferSize:    bufferSize,
	}, nil
}

// Crack returns the matching key pair. False is returned when the pair is not found.
func (c *Cracker) Crack(ctx context.Context) (*KeyPair, bool, error) {

	// 4 go routines are reserved for: waiting for a result, producing and waiting for the other go routines.
	const reservedGoRoutinesNum = 4

	var (
		tasks      = make(chan task, c.goRoutinesNum-reservedGoRoutinesNum)
		insertReqs = make(chan repository.InsertBulkRequest, c.bufferSize)
		result     = make(chan KeyPair, 1)
		taskBar    = pb.StartNew(c.keyNum * 2)
		keyPair    KeyPair
		found      bool
	)

	defer taskBar.Finish()

	ctx, cancel := context.WithCancel(ctx)
	g, ctx := errgroup.WithContext(ctx)

	// result
	g.Go(func() error {
		for r := range result {
			found = true
			keyPair = r
			cancel()
			break
		}
		return nil
	})

	// insert and check
	g.Go(func() error {
		var (
			insertReqBuf = make([]repository.InsertBulkRequest, c.bufferSize)
			idx          int
		)
		for {
			select {
			case r := <-insertReqs:
				insertReqBuf[idx] = r
				if idx == c.bufferSize-1 {
					kp, wasFound, err := c.inserter.InsertBulk(insertReqBuf...)
					if err != nil {
						log.Println(fmt.Sprintf("error while processing task %v", err))
					}
					if wasFound {
						result <- KeyPair{
							EncodeKey: kp.EncodeKey,
							DecodeKey: kp.DecodeKey,
						}
					}
					insertReqBuf, idx = make([]repository.InsertBulkRequest, c.bufferSize), 0
					taskBar.Add(c.bufferSize)
					continue
				}
				idx++
			case <-ctx.Done():
				return ctx.Err()
			default:
				break
			}
		}
	})

	// consumer
	for i := 0; i < c.goRoutinesNum-reservedGoRoutinesNum/2; i++ {
		g.Go(func() error {
			for {
				select {
				case t := <-tasks:
					cipherText, err := t.fn(t.key, t.text)
					if err != nil {
						log.Println(fmt.Sprintf("error while processing task %v", err))
						break
					}

					insertReqs <- repository.InsertBulkRequest{
						Key:        t.key,
						CipherText: cipherText,
						Mode:       t.mode,
					}

				case <-ctx.Done():
					return ctx.Err()
				default:
					break
				}
			}
		})
	}

	// producer
	g.Go(func() error {
		for k := 0; k < c.keyNum; k++ {
			fk := c.formatterFn(k, c.keyLength)
			tasks <- task{
				key:  fk,
				text: c.cipherText,
				mode: decodeMode,
				fn:   c.decFn,
			}
			tasks <- task{
				key:  fk,
				text: c.plainText,
				mode: encodeMode,
				fn:   c.encFn,
			}
		}
		return nil
	})

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return nil, false, err
	}

	close(tasks)
	close(insertReqs)
	close(result)

	return &keyPair, found, nil
}
