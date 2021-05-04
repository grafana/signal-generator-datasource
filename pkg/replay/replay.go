package replay

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Replayer func(messageData []byte) error

func DoReplay(fpath string, interval time.Duration, player Replayer) {
	file, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}

	defer func() { _ = file.Close() }()

	firstTimeNS := int64(0)
	startTime := time.Now()
	offsetTime := int64(0)
	lastTime := int64(0)
	isNano := false
	timestampScale := int64(1)

	// All lines in the path
	batch := ""
	batchCount := 0
	row := 0

	fmt.Printf("STREAMING: %s\n\n", fpath)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		idx := strings.LastIndex(line, " ")
		if idx < 1 || strings.HasPrefix("#", line) || strings.HasPrefix("//", line) {
			continue
		}

		row++
		timestamp, err := strconv.ParseInt(line[(idx+1):], 10, 64)
		if err != nil {
			fmt.Printf("error parsing: %s [%s]\n", line, err.Error())
			continue
		}

		// first line
		if offsetTime < 1 {
			startTime = time.Now()
			isNano = timestamp > 1000066771122836000
			if isNano {
				offsetTime = startTime.UnixNano() - timestamp
			} else {
				timestampScale = int64(time.Millisecond)
				nowMS := startTime.UnixNano() / timestampScale
				offsetTime = nowMS - timestamp
			}
			firstTimeNS = timestamp * timestampScale
		}

		delta := timestamp - lastTime
		if delta > 0 && len(batch) > 0 {
			shitedTime := (timestamp + offsetTime) * timestampScale

			delta := shitedTime - time.Now().UnixNano()
			if delta > 0 {
				sleepTime := time.Nanosecond * time.Duration(delta)
				if sleepTime < interval {
					sleepTime = interval
				}
				elapsed := time.Duration((timestamp * timestampScale) - firstTimeNS)

				fmt.Printf("WRITE [%d] sleep: %s // %s\n", row, sleepTime, elapsed)
				err = player([]byte(batch))
				if err != nil {
					fmt.Printf("error writing: %s\n", err.Error())
				}

				time.Sleep(time.Nanosecond * time.Duration(delta))
			}

			lastTime = timestamp
			batch = ""
			batchCount = 0
		}
		batch += fmt.Sprintf("%s %d\n", line[:idx], timestamp+offsetTime)
		batchCount++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
