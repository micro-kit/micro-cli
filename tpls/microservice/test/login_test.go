package test

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"testing"
	"time"

	"{{ .MicroKitClientRoot }}/client/{{ .BaseServiceNameNotLine }}"
	"{{ .MicroKitClientRoot }}/proto/{{ .BaseServiceNameNotLine }}pb"
)

var (
	cl {{ .BaseServiceNameNotLine }}pb.{{ .BaseServiceNameHump }}Client
)

func TestMain(m *testing.M) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var err error
	cl, err = {{ .BaseServiceNameNotLine }}.NewClient()
	if err != nil {
		log.Panicln(err)
	}
	m.Run()
}
