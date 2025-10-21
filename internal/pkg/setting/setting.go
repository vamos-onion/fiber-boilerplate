package setting

import (
	"encoding/json"
	"fmt"

	"fiber-boilerplate/internal/pkg/util"

	"github.com/caarlos0/env/v6"
)

// RuntimeBlock :
type RuntimeBlock struct {
	Port int    `env:"PORT" envDefault:"8080"`
	Env  string `env:"ENV" envDefault:"local"`
}

// Runtime :
var Runtime RuntimeBlock

// Configs :
var Configs map[string]json.RawMessage

// init early
func init() {
	fmt.Println(util.String.Words("*", util.Prev.Func(), "->", util.Current.Func()))

	err := env.Parse(&Runtime)
	if err != nil {
		panic(err)
	}
}
