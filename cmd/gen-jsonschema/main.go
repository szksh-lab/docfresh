package main

import (
	"fmt"
	"log"

	"github.com/suzuki-shunsuke/docfresh/pkg/controller/run"
	"github.com/suzuki-shunsuke/gen-go-jsonschema/jsonschema"
)

func main() {
	if err := core(); err != nil {
		log.Fatal(err)
	}
}

func core() error {
	if err := jsonschema.Write(&run.BlockInput{}, "json-schema/comment.json"); err != nil {
		return fmt.Errorf("create or update a JSON Schema: %w", err)
	}
	return nil
}
