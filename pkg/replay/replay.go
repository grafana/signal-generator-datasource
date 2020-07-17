package replay

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/grafana/waveform-datasource/pkg/common/globpath"
	"github.com/grafana/waveform-datasource/pkg/models"
	"github.com/grafana/waveform-datasource/pkg/parsers"
	"golang.org/x/text/encoding"
)

type Replay struct {
	Files      []string `toml:"files"`
	FileTag    string   `toml:"file_tag"`
	Iterations int      `toml:"iterations"`
	parser     parsers.Parser

	filenames []string
	decoder   *encoding.Decoder
}

func (r *Replay) Start(acc chan *models.InfluxLine) error {
	err := r.refreshFilePaths()
	if err != nil {
		return err
	}
	for _, k := range r.filenames {
		metrics, err := r.readMetrics(k)
		if err != nil {
			return err
		}

		go r.processMetrics(metrics, acc)
	}

	return nil
}

func (r *Replay) Stop() {

}

func (r *Replay) SetParser(p parsers.Parser) {
	r.parser = p
}

func (r *Replay) refreshFilePaths() error {
	var allFiles []string
	for _, file := range r.Files {
		g, err := globpath.Compile(file)
		if err != nil {
			return fmt.Errorf("could not compile glob %v: %v", file, err)
		}
		files := g.Match()
		if len(files) <= 0 {
			return fmt.Errorf("could not find file: %v", file)
		}
		allFiles = append(allFiles, files...)
	}

	r.filenames = allFiles
	return nil
}

func (r *Replay) readMetrics(filename string) ([]*models.InfluxLine, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("E! Error file: %v could not be read, %s", filename, err)
	}

	metrics, err := r.parser.Parse(fileContents)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func (r *Replay) processMetrics(metrics []*models.InfluxLine, acc chan *models.InfluxLine) {
	for i := 0; i != r.Iterations; i++ {
		prevTime := metrics[0].Timestamp
		for _, metric := range metrics {
			currTime := metric.Timestamp
			delay := currTime.Sub(prevTime)
			time.Sleep(delay)
			prevTime = currTime
			// shout, shoult, let it all out
			acc <- metric
		}
	}
}
