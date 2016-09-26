package db100

import (
	"fmt"
	"log"
	"os"

	"github.com/chaosvermittlung/coinslot/global"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func Initialisation(dbc *global.DBConnection) {
	var err error
	db, err = sqlx.Open(dbc.Driver, dbc.Connection)
	if err != nil {
		log.Fatal(err)
	}
	initDB(dbc)
}

func initDB(dbc *global.DBConnection) {
	switch dbc.Driver {
	case "sqlite3":
		cont, err := global.Exists(dbc.Connection)
		if err != nil {
			log.Fatal(err)
		}
		if cont {
			fmt.Println("cont")
			return
		}
		_, err = os.Create(dbc.Connection)
		if err != nil {
			log.Fatal("Could not create file "+dbc.Connection, err)
		}
		_, err = db.Exec(createSQLlitestmt)
		if err != nil {
			log.Printf("%q: %s\n", err, createSQLlitestmt)
			return
		}
		var u User
		u.Username = "admin"
		u.Password = "admin"
		u.Email = "admin@local"
		u.Right = USERRIGHT_ADMIN
		s, err := global.GenerateSalt()
		if err != nil {
			log.Println(err)
		}
		u.Salt = s

		pw, err := global.GeneratePasswordHash(u.Password, u.Salt)
		if err != nil {
			log.Println(err)
		}
		u.Password = pw
		err = u.Insert()
		if err != nil {
			log.Println(err)
		}
	default:
		log.Fatal("DB Driver unkown. Stopping Server")
	}
}

type UserRight int

const (
	USERRIGHT_MEMBER UserRight = 1 + iota
	USERRIGHT_ADMIN
)

type User struct {
	User_ID  int       `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Salt     string    `json:"-"`
	Email    string    `json:"email"`
	Right    UserRight `json:"userright"`
}

func copyifnotempty(str1, str2 string) string {
	if str2 != "" {
		return str2
	} else {
		return str1
	}
}

func DoesUserExist(username string) (bool, error) {
	var id int
	err := db.Get(&id, "Select Count(*) from Users Where Username = ?", username)
	b := (id > 0)
	return b, err
}

func GetUsers() ([]User, error) {
	var u []User
	err := db.Select(&u, "Select * from Users")
	return u, err
}

func (u *User) GetDetailstoUsername() error {
	err := db.Get(u, "SELECT * from Users Where Username = ? Limit 1", u.Username)
	return err
}

func (u *User) GetDetails() error {
	err := db.Get(u, "SELECT * from Users Where User_ID = ? Limit 1", u.User_ID)
	return err
}

func (u *User) Patch(ou User) error {
	u.Username = copyifnotempty(u.Username, ou.Username)
	if ou.Password != "" {
		p, err := global.GeneratePasswordHash(ou.Password, u.Salt)
		if err != nil {
			return err
		}
		u.Password = p
	}
	u.Email = copyifnotempty(u.Email, ou.Email)
	if ou.Right != 0 {
		u.Right = ou.Right
	}
	return nil
}

func (u *User) Update() error {
	_, err := db.Exec("UPDATE Users SET username = ?, password = ?, email = ?, right = ? WHERE id = ?", u.Username, u.Password, u.Email, u.Right, u.User_ID)
	return err
}

func (u *User) Insert() error {
	res, err := db.Exec("INSERT INTO Users (username, password, salt, email, right) VALUES(?,?,?,?,?)", u.Username, u.Password, u.Salt, u.Email, u.Right)
	if err != nil {
		log.Println(err)
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}
	u.User_ID = int(id)

	return nil
}

func DeleteUser(id int) error {
	_, err := db.Exec("DELETE FROM Users Where User_ID = ?", id)
	return err
}

func (u *User) GetProjects() ([]Project, error) {
	var p []Project
	err := db.Select(&p, "Select * from Project Where User_ID = ?", u.User_ID)
	return p, err
}

type Project struct {
	Project_ID  int
	Name        string
	Goal        float64
	Initiator   int
	Description string
}

func GetAllProjects() ([]Project, error) {
	var p []Project
	err := db.Select(&p, "Select * from Projects")
	return p, err
}

func (p *Project) Insert() error {
	res, err := db.Exec("INSERT INTO Projects (name, goal, initiator, description) VALUES(?,?,?,?)", p.Name, p.Goal, p.Initiator, p.Description)
	if err != nil {
		log.Println(err)
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}
	p.Project_ID = int(id)

	return nil
}

func (p *Project) GetDetails() error {
	err := db.Get(p, "SELECT * from Projects Where Project_ID = ? Limit 1", p.Project_ID)
	return err
}

func (p *Project) GetInitiatorName() (string, error) {
	var n string
	err := db.Get(&n, "Select Username From Users Where User_ID = ?", p.Initiator)
	return n, err
}

func (p *Project) GetFundings() ([]Funding, error) {
	var f []Funding
	err := db.Select(&f, "Select * from fundings Where project_id = ?", p.Project_ID)
	return f, err
}

type Funding struct {
	Project_ID int
	User_ID    int
	Amount     float64
	Confirmed  bool
}

func (f *Funding) Insert() error {
	_, err := db.Exec("INSERT INTO Fundings (Project_ID ,User_ID, Amount, Confirmed) VALUES(?,?,?,?)", f.Project_ID, f.User_ID, f.Amount, f.Confirmed)
	return err
}

func GetFundingAmounts(ff []Funding) (float64, float64) {
	var pro float64
	var got float64
	for _, f := range ff {
		if !f.Confirmed {
			pro = pro + f.Amount
		} else {
			got = got + f.Amount
		}
	}

	return pro, got
}
