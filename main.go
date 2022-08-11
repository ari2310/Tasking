package main

import (
    "database/sql"
    "log"
    "net/http"
    "text/template"

    _ "github.com/go-sql-driver/mysql"
)

type Employee struct {
    Id              int
    Task_employee   string
    Id_employee     string
    Date            string
    Status          bool
}

func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := ""
    dbName := "tasking"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    selDB, err := db.Query("SELECT * FROM task ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
    }
    emp := Employee{}
    res := []Employee{}
    for selDB.Next() {
        var id, status int
        var task_employee, date, id_employee string
        err = selDB.Scan(&id, &task_employee, &id_employee, &date, &status)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = id
        emp.Task_employee = task_employee
        emp.Id_employee = id_employee
        emp.Date = date
        if status==1{
            emp.Status = false
        }else{emp.Status = true}
        res = append(res, emp)
    }
    tmpl.ExecuteTemplate(w, "Index", res)
    defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM task WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    emp := Employee{}
    for selDB.Next() {
        var id, status int
        var task_employee, date, id_employee string
        err = selDB.Scan(&id, &task_employee, &id_employee, &date, &status)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = id
        emp.Task_employee = task_employee
        emp.Id_employee = id_employee
        emp.Date = date
        // emp.Status = status
    }
    tmpl.ExecuteTemplate(w, "Edit", emp)
    defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        task_employee := r.FormValue("task_employee")
        id_employee := r.FormValue("id_employee")
        date := r.FormValue("date")
        status := 0
        insForm, err := db.Prepare("INSERT INTO task(task_employee, id_employee, date, status) VALUES(?,?,?,?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(task_employee, id_employee, date, status)
        log.Println("INSERT: Name: " + task_employee + " | City: " + id_employee + " | City: " + date)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        task_employee := r.FormValue("task_employee")
        id_employee := r.FormValue("id_employee")
        date := r.FormValue("date")
        id := r.FormValue("uid")
        insForm, err := db.Prepare("UPDATE task SET task_employee=?, id_employee=?, date=? WHERE id=?")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(task_employee, id_employee, date, id)
        log.Println("UPDATE: Name: " + task_employee + " | City: " + id_employee)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    emp := r.URL.Query().Get("id")
    delForm, err := db.Prepare("DELETE FROM task WHERE id=?")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(emp)
    log.Println("DELETE")
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Done(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    emp := r.URL.Query().Get("id")
    delForm, err := db.Prepare("UPDATE task SET status=1 WHERE id=?")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(emp)
    log.Println("DELETE")
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func main() {
    log.Println("Server started on: http://localhost:9090")
    http.HandleFunc("/", Index)
    http.HandleFunc("/new", New)
    http.HandleFunc("/edit", Edit)
    http.HandleFunc("/insert", Insert)
    http.HandleFunc("/update", Update)
    http.HandleFunc("/delete", Delete)
    http.HandleFunc("/status", Done)
    http.ListenAndServe(":9090", nil)
}