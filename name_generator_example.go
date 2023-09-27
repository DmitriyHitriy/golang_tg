package main

import (
	"strings"
	"time"

	"github.com/goombaio/namegenerator"
)

func generate_name() string {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	name := strings.ReplaceAll(nameGenerator.Generate()+"_bratkov", "-", "_")

	return name
}
