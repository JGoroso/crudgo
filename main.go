package main

import (
	"database/sql"

	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

var templates = template.Must(template.ParseGlob("templates/*"))

func main() {
	//handlefunc -> when visit "", some function will be activated.
	http.HandleFunc("/", Init)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)

	log.Println("Servidor corriendo...")
	http.ListenAndServe(":8080", nil)

}

type Employee struct {
	Id    int
	Name  string
	Email string
}

func Init(w http.ResponseWriter, r *http.Request) {
	connectCorrect := connectBD()

	data, err := connectCorrect.Query("SELECT * FROM employees")
	if err != nil {
		panic(err.Error())
	}

	employee := Employee{}
	arrayEmployee := []Employee{}

	for data.Next() {
		var id int
		var name, email string

		err = data.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		employee.Id = id
		employee.Name = name
		employee.Email = email

		arrayEmployee = append(arrayEmployee, employee)
	}

	templates.ExecuteTemplate(w, "init", arrayEmployee)
}

func Create(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "create", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("inputName")
		email := r.FormValue("inputEmail")

		connectCorrect := connectBD()

		insertData, err := connectCorrect.Prepare("INSERT INTO employees(name, email) VALUES (?,?)")

		if err != nil {
			panic(err.Error())
		}

		insertData.Exec(name, email)

		http.Redirect(w, r, "/", 301)

	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idEmployee := r.URL.Query().Get("id")

	connectCorrect := connectBD()

	deleteData, err := connectCorrect.Prepare("DELETE FROM employees WHERE id=?")

	if err != nil {
		panic(err.Error())
	}

	deleteData.Exec(idEmployee)

	http.Redirect(w, r, "/", 301)
}

// Edit & Actualizar
func Edit(w http.ResponseWriter, r *http.Request) {
	idEmployee := r.URL.Query().Get("id")

	connectCorrect := connectBD()

	register, err := connectCorrect.Query("SELECT * FROM employees WHERE id =?", idEmployee)
	if err != nil {
		panic(err.Error())
	}

	employee := Employee{}

	for register.Next() {
		var id int
		var name, email string

		err = register.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		employee.Id = id
		employee.Name = name
		employee.Email = email

	}

	templates.ExecuteTemplate(w, "edit", employee)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		id := r.FormValue("id")
		name := r.FormValue("inputName")
		email := r.FormValue("inputEmail")

		connectCorrect := connectBD()

		updateData, err := connectCorrect.Prepare("UPDATE employees SET name=?, email=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}

		updateData.Exec(name, email, id)

		http.Redirect(w, r, "/", 301)
	}
}

// Connect to BD
func connectBD() (connect *sql.DB) {

	Driver := "mysql"
	User := "root"
	Password := ""
	Name := "crudgobd"

	connect, err := sql.Open(Driver, User+":"+Password+"@tcp(127.0.0.1)/"+Name)
	if err != nil {
		panic(err.Error())
	}
	return connect
}
