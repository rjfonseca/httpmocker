package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	httpAddr = ":8080"

	defaultStatus  = 200
	defaultMessage = "ok"
	defaultHeaders = "Server: mock & Content-Type: application/text"

	datadir    = ""
	outputdir  = time.Now().Format("output-20060102150405")
	strictMode = false

	maxLoops = 0

	sleep       = ""
	sleepJitter = ""
)

func main() {
	flag.StringVar(&httpAddr, "http-addr", getEnvOrString("HTTP_ADDR", httpAddr), "HTTP address and port to bind to. Examples: ':8080', '127.0.0.1:8080'")

	flag.IntVar(&defaultStatus, "default-status", getEnvOrInt("DEFAULT_STATUS", defaultStatus), "Default HTTP status to return when the request doesn't mach any handler")
	flag.StringVar(&defaultMessage, "default-message", getEnvOrString("DEFAULT_MESSAGE", defaultMessage), "Default HTTP body to return when the request doesn't mach any handler")
	flag.StringVar(&defaultHeaders, "default-header", getEnvOrString("DEFAULT_HEADERS", defaultHeaders), "Default HTTP Headers (sepparated by '&') to return when the request doesn't mach any handler")

	flag.StringVar(&outputdir, "output", getEnvOrString("OUTOUT", outputdir), "Save requests and responses to this directory")

	flag.StringVar(&datadir, "data", getEnvOrString("DATA", datadir), "Directory containing all the mock files")
	flag.IntVar(&maxLoops, "max-loops", getEnvOrInt("MAX_LOOPS", maxLoops), "Maximum number of loops before sending only the default response")

	flag.StringVar(&sleep, "sleep", getEnvOrString("SLEEP", sleep), "Sleep duration for all requests, Example: '5s'")
	flag.StringVar(&sleepJitter, "sleep-jitter", getEnvOrString("SLEEP_JITTER", sleepJitter), "Adds a randem extra duration to the sleep tim. Example: '3s'")

	flag.Parse()

	_, err := os.Stat(outputdir)
	if err == nil {
		log.Fatalln("Output directory already exists:", outputdir)
	} else if !os.IsNotExist(err) {
		log.Fatalln("Unable to check if output directory exists:", err)
	}

	err = os.MkdirAll(outputdir, 0777)
	if err != nil {
		log.Fatalln("Failed to reate output directory:", err)
	}

	dh := defaultHandler(defaultStatus,
		parseHeaders(defaultHeaders),
		[]byte(defaultMessage),
	)

	var h http.Handler
	if datadir != "" {
		h, err = makeStrictHandlerFromDir(maxLoops, dh, datadir)

		if err != nil {
			log.Fatalln("Failed to create handlers for data directory:", err)
		}
	} else {
		h = dh
	}

	if sleep != "" {
		s, err := time.ParseDuration(sleep)
		if err != nil {
			log.Fatalln("Failed to parse sleep duration:", err)
		}
		var j time.Duration
		if sleepJitter != "" {
			j, err = time.ParseDuration(sleepJitter)
			if err != nil {
				log.Fatalln("Failed to parse sleep jiiter duration:", err)
			}
		}
		h = withSleep(s, j, h)
	}
	h = registerResponse(outputdir, h)
	h = registerRequest(outputdir, h)
	h = requestCounter(h)

	log.Println("Service started")
	err = http.ListenAndServe(httpAddr, h)
	log.Println("Service terminated:", err)
}

func parseHeaders(s string) map[string]string {
	h := make(map[string]string)
	parts := strings.Split(s, "&")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kv := strings.Split(part, ":")
		if len(kv) != 2 {
			panic("invalid header map string")
		}
		h[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	return h
}

func newHandler(fpath string) http.Handler {
	switch {
	case strings.HasSuffix(fpath, ".mock.json"):
		return mockHandler(fpath)
	default:
		return txtHandler(fpath)
	}
}

func newMultiHandler(maxLoops int, defaultHandler http.Handler, dir string) (http.Handler, error) {
	h, err := makeHandlersFromDir(dir)
	if err != nil {
		return nil, err
	}
	return multiHandler(maxLoops, defaultHandler, h...), nil
}

func makeAllHandlersFromDir(dir string) ([]http.Handler, error) {
	handlers := []http.Handler{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fmt.Println("Loaded txt handler:", path)
		handlers = append(handlers, newHandler(path))
		return nil
	})

	return handlers, err
}

func makeStrictHandlerFromDir(maxLoops int, defaultHandler http.Handler, dir string) (http.Handler, error) {
	m := http.NewServeMux()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		resource := filepath.ToSlash(strings.TrimPrefix(path, dir))
		if resource == "" {
			resource = "/"
		} else {
			resource = strings.TrimSuffix(resource, filepath.Ext(resource))
		}
		var h http.Handler
		if info.IsDir() {
			h, err = newMultiHandler(maxLoops, defaultHandler, path)
			if err != nil {
				return err
			}
		} else {
			h = newHandler(path)
		}
		log.Printf("Loaded '%s' into '%s'\n", path, resource)
		m.Handle(resource, h)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return m, nil
}

func makeHandlersFromDir(dir string) ([]http.Handler, error) {
	d, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	s, err := d.Stat()
	if err != nil {
		return nil, err
	}
	if !s.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", dir)
	}
	handlers := []http.Handler{}

	for {
		infos, err := d.Readdir(100)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		for _, info := range infos {
			if info.IsDir() {
				log.Println("Skiping", info.Name())
				continue
			}
			fpath := filepath.Join(dir, info.Name())
			log.Println("New handler to", fpath)

			handlers = append(handlers, newHandler(fpath))
		}
	}
	return handlers, nil
}
