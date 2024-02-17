package atlas

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	_ "ariga.io/atlas-go-sdk/recordriver"
	"ariga.io/atlas-provider-gorm/gormschema"
)

// AtlasMain - Call this function in `atlas/main.go` file.
// Create atlas.hcl in a root directory with the content similar to this:
//
//	data "external_schema" "gorm_load_models" {
//	  program = [
//	    "go",
//	    "run",
//	    "./atlas",
//	    "m"
//	  ]
//	}
//
//	data "external" "gorm_load_env" {
//	  program = [
//	    "go",
//	    "run",
//	    "./atlas",
//	    "e"
//	  ]
//	}
//
//	locals {
//	  dbenv = jsondecode(data.external.gorm_load_env)
//	}
//
//	env "gorm" {
//	  src = data.external_schema.gorm_load_models.url
//	  dev = local.dbenv.dev
//	  url = local.dbenv.url
//	  migration {
//	    dir = "file://database/migrations"
//	    revisions_schema = local.dbenv.revisionsSchema
//	  }
//	  schemas = ["dd"]
//	  format {
//	    migrate {
//	      diff = "{{ sql . \"  \" }}"
//	    }
//	  }
//	  excludes = local.dbenv.revisionsSchema
//	}
func AtlasMain(dialect string, initialSql string, config AtlasConfig, models ...any) {
	// argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 0 {
		fmt.Println("No arguments provided")
		os.Exit(1)
	}

	switch argsWithoutProg[0] {
	case "m":
		loadModels(dialect, initialSql, models...)
	case "e":
		loadEnv(config)
	default:
		fmt.Println("Invalid argument, please use 'm' or 'e'")
		os.Exit(1)
	}

}

// To be run with atlas cli defined in atlas.hcl as
//
//	data "external" "gorm_load_models" {
//	  program = [
//	    "go",
//	    "run",
//	    "./atlas",
//	    "m"
//	  ]
//	}
func loadModels(dialect string, sql string, models ...any) {

	stmts, err := gormschema.New(dialect).Load(models)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	_, err = io.WriteString(os.Stdout, sql+stmts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write output: %v\n", err)
		os.Exit(1)
	}
}

// To be run with atlas cli defined in atlas.hcl as
//
//	data "external" "gorm_load_env" {
//	  program = [
//	    "go",
//	    "run",
//	    "./atlas",
//	    "e"
//	  ]
//	}
//
//	locals {
//	  dbenv = jsondecode(data.external.gorm_load_env)
//	}
func loadEnv(config AtlasConfig) {
	jsonConfig, err := json.Marshal(config)
	if err != nil {
		log.Fatal("Can not marshal atlas config: ", err)
	}

	_, err = io.WriteString(os.Stdout, string(jsonConfig))
	if err != nil {
		log.Fatal("Can not write atlas config: ", err)
	}
}
