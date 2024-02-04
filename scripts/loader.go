package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	profile "github.com/thestoicway/user_service/internal/profile/model"
	usr "github.com/thestoicway/user_service/internal/usr/model"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(&usr.UserDB{}, &profile.ProfileDB{})

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
