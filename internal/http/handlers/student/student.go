package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	// "github.com/Asmitshukl/apiitis/internal/http/handlers/student"
	"github.com/Asmitshukl/apiitis/internal/storage"
	"github.com/Asmitshukl/apiitis/internal/types"
	"github.com/Asmitshukl/apiitis/internal/utils/response"
	"github.com/go-playground/validator/v10"
)


func New(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter , r *http.Request) {

		var student types.Student

		slog.Info("creating a student")

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err , io.EOF){
			response.WriteJson(w , http.StatusBadRequest , response.GeneralError(fmt.Errorf("empty body")))
			return 
		}

		if err != nil {
			response.WriteJson(w , http.StatusBadRequest , response.GeneralError(err))
			return 
		}

		//request validataion

		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w , http.StatusBadRequest , response.ValidationError(validateErrs))
			return
		}

		id , err :=storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		if err != nil {
			response.WriteJson(w , http.StatusInternalServerError , err)
			return
		}
		slog.Info("User created successfully",slog.String("userid",fmt.Sprint(id)))

		response.WriteJson(w , http.StatusCreated , map[string]int64{"id": id})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter , r *http.Request) {
		//implement get by id logic here
		id := r.PathValue("id")
		slog.Info("getting a student",slog.String("id",id))

		intid , err := strconv.ParseInt(id , 10 ,64)
		if err != nil {
			response.WriteJson(w,http.StatusBadRequest,err)
		}
		student , err := storage.GetStudentById(intid)
		if err != nil {
			slog.Error("error getting user" ,slog.String("id",id))
			response.WriteJson(w, http.StatusInternalServerError , response.GeneralError(err))
			return 
		}

		response.WriteJson(w , http.StatusOK , student)
	}
}