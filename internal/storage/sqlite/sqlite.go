package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/Asmitshukl/apiitis/internal/config"
	"github.com/Asmitshukl/apiitis/internal/types"

	//we have to import the sqlite3 driver for the sql package to work with sqlite databases
	//also the underscore is used to import the package for its side effects (initialization) without directly using it in the code
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB	
}


func New(cfg *config.Config )(*Sqlite, error){
	db, err := sql.Open("sqlite3" , cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  name TEXT ,
	  email TEXT ,
	  age INTEGER 
	)`)

	if err != nil {
		return nil, err
	}
	
	
	return &Sqlite{Db: db} , nil
}

func (s *Sqlite) CreateStudent(name string , email string ,age int) (int64 , error){

	stmt , err := s.Db.Prepare("INSERT INTO students (name , email , age) VALUES (? , ? , ?)")
    if err != nil {
		return 0 , err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name,email,age)
	if err != nil {
		return 0 , err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0 , err
	}

	return id , nil

}

func (s *Sqlite) GetStudentById(id int64) (types.Student , error){
	stmt ,err := s.Db.Prepare("SELECT * FROM students WHERE id = ? limit 1")
	if err != nil {
		return types.Student{} , err
	}
	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.ID , &student.Name , &student.Email , &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{} , fmt.Errorf("no student of id %d", id)
		}
		return types.Student{} , fmt.Errorf("query error %w",err)
	}

	return student , nil

}