package main

import (
	"github.com/chaosvermittlung/coinslot/db/v100"
	"github.com/chaosvermittlung/coinslot/global"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	db100.Initialisation(&global.Conf.Connection)

}
