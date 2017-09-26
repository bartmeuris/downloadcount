package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hpcloud/tail"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sync"
	"time"
)

type DownloadSource struct {
	Source string
	Regex  string
	Fields []string
	TimeField struct {
		Field string
		Format string
	} `json:"time,omitempty"`

	tailFile *tail.Tail
	lastTime time.Time
	lastLine uint64
	re       *regexp.Regexp
	timefield int
}

func (ds *DownloadSource) getLocation() *tail.SeekInfo {
	var loc int64
	loc = 0
	return &tail.SeekInfo{
		Offset: loc,
		Whence: os.SEEK_SET,
	}
}

func (ds *DownloadSource) Init() error {
	var err error

	// 1. Precompile regex
	ds.re, err = regexp.Compile(ds.Regex)
	if err != nil {
		return err
	}
	if ds.re.NumSubexp() != len(ds.Fields) {
		return fmt.Errorf("Expected %d fields defined for %s, got %d", ds.re.NumSubexp(), ds.Source, len(ds.Fields))
	}
	ds.timefield = 0
	if ds.TimeField.Field != "" {
		for i := range ds.Fields {
			if ds.Fields[i] == ds.TimeField.Field {
				// Add one to skip global match in substring match list
				ds.timefield = i + 1
				break
			}
		}
	}

	// 2. Initialize file tailing
	ds.tailFile, err = tail.TailFile(ds.Source, tail.Config{
		Location:  ds.getLocation(),
		ReOpen:    true,
		MustExist: true,
		Follow:    true,
		Logger:    tail.DiscardingLogger,
	})
	if err != nil {
		return err
	}

	return nil
}

func (ds *DownloadSource) processLine(l *tail.Line) int {
	m := ds.re.FindStringSubmatch(l.Text)
	if len(m) == (len(ds.Fields) + 1) {
		fmt.Printf("MATCH: %s\n", l.Text)
		for i := range ds.Fields {
			fmt.Printf("    %s: '%s'\n", ds.Fields[i], m[i+1])
		}
		if ds.timefield > 0 {
			tm, err := time.Parse(ds.TimeField.Format, m[ds.timefield])
			if err == nil {
				fmt.Printf("    Parsed time: %s => %s\n", m[ds.timefield], tm.String())
			} else {
				fmt.Printf("    ERROR: Could not parse time %s: %s\n", m[ds.timefield], err)
			}
		}
		return 1
	} else {
		//fmt.Printf("NO MATCH: %s (%q)\n", l.Text, m)
		return 0
	}
}

func (ds *DownloadSource) Process() {
	if ds.tailFile == nil {
		return
	}
	count := 0
	t := time.After(10*time.Second)
	var lastTime time.Time
	for {
		select {
		case l, ok := <-ds.tailFile.Lines:
			if !ok {
				return
			}
			count += ds.processLine(l)
		case <- t:
			fmt.Printf("FLUSH: %d since %s\n", count, lastTime.String())
			t = time.After(10*time.Second)
			//count = 0
			lastTime = time.Now()
		}
	}
}

func (ds *DownloadSource) Cleanup() {
	if ds.tailFile != nil {
		ds.tailFile.Cleanup()
		ds.tailFile = nil
	}
}

func main() {
	jsonfile := flag.String("file", "settings.json", "JSON file with configuration")
	logFile := flag.String("logfile", "", "file to log to")
	flag.Parse()

	// Log to a file if requested
	if *logFile != "" {
		if f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err == nil {
			defer f.Close()
			log.SetOutput(f)
		}
	}

	jsonb, err := ioutil.ReadFile(*jsonfile)
	if err != nil {
		log.Fatalf("Could not read provided file %s", *jsonfile)
	}
	var srcs []*DownloadSource
	if err := json.Unmarshal(jsonb, &srcs); err != nil {
		log.Fatalf("ERROR: could not parse JSON: %s", err)
	}
	wg := &sync.WaitGroup{}
	for _, ds := range srcs {
		err := ds.Init()
		if err != nil {
			log.Fatalf("ERROR: Could not initialize download source: %s", err)
		}
		wg.Add(1)
		go func() {
			ds.Process()
			ds.Cleanup()
			wg.Done()
		}()
	}
	wg.Wait()
	log.Printf("All jobs finished, strange?")
}
