package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/chaosvermittlung/coinslot/db/v100"
	"github.com/chaosvermittlung/coinslot/global"
)

type frontpage struct {
	Projects []project
	Message  template.HTML
	Username string
}

type project struct {
	Name        string
	Goal        float64
	Promised    float64
	Confirmed   float64
	PromisedP   float64
	ConfirmedP  float64
	Initiator   string
	Description string
	UserIsAdmin bool
	Fundings    []funding
}

type funding struct {
	Funder    string
	Amount    float64
	Confirmed bool
}

const errormessage = `<div class="alert alert-danger" role="alert">
  <span class="glyphicon glyphicon-exclamation-sign" aria-hidden="true"></span>
  <span class="sr-only">Error:</span>
  $MESSAGE$
</div>`

func mainhandler(w http.ResponseWriter, r *http.Request) {

	var fronterr string

	u := db100.User{}
	u.User_ID = 1

	/*u.Username, err = GetCookie(r)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
	}*/

	err := u.GetDetails()
	if err != nil {
		fronterr = global.BuildMessage(errormessage, err.Error())
		log.Println(err)
		return
	}

	fp := frontpage{}

	pp, err := db100.GetAllProjects()

	if err != nil {
		fronterr = global.BuildMessage(errormessage, err.Error())
		log.Println(err)
		return
	}

	var pro []project

	for _, op := range pp {
		var p project
		p.Name = op.Name
		p.Description = op.Description
		p.Goal = op.Goal
		p.UserIsAdmin = (u.User_ID == op.Initiator)
		n, err := op.GetInitiatorName()
		if err != nil {
			log.Println(err)
			fronterr = global.BuildMessage(errormessage, err.Error())
			break
		}
		p.Initiator = n
		ff, err := op.GetFundings()
		if err != nil {
			log.Println(err)
			fronterr = global.BuildMessage(errormessage, err.Error())
			break
		}

		p.Promised, p.Confirmed = db100.GetFundingAmounts(ff)

		p.PromisedP = p.Goal / 100 * p.Promised
		p.ConfirmedP = p.Goal / 100 * p.Confirmed
		if p.UserIsAdmin {
			for _, of := range ff {
				if of.Project_ID == op.Project_ID {
					var f funding
					f.Amount = of.Amount
					var u db100.User
					u.User_ID = of.User_ID
					err := u.GetDetails()
					if err != nil {
						log.Println(err)
						fronterr = global.BuildMessage(errormessage, err.Error())
						break
					}
					f.Funder = u.Username
					f.Confirmed = of.Confirmed
					p.Fundings = append(p.Fundings, f)
				}

			}
		}

		pro = append(pro, p)

	}

	fp.Projects = pro
	fp.Username = u.Username
	fp.Message = template.HTML(fronterr)

	t, err := template.ParseFiles("templates/main.html")

	err = t.Execute(w, &fp)
	if err != nil {
		log.Println(err)
	}
}

func statichandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	db100.Initialisation(&global.Conf.Connection)

	http.HandleFunc("/", mainhandler)
	http.HandleFunc("/static/", statichandler)
	port := ":" + strconv.Itoa(global.Conf.Port)
	http.ListenAndServe(port, nil)
}
