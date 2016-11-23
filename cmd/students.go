package envy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"html/template"

	"github.com/gophertrain/envy/pkg"
	"github.com/gorilla/mux"
	"github.com/kr/pretty"
)

type SignupResponse struct {
	Error        bool
	ErrorMessage string
	ShellURL     string
	IDEURL       string
	Username     string
}

type Server struct {
	Router      *mux.Router
	Students    pkg.StudentStorage
	Courses     pkg.CourseStorage
	ExecManager *ExecManager
}

func NewServer(students pkg.StudentStorage, courses pkg.CourseStorage) *Server {
	s := &Server{}
	s.Students = students
	s.Courses = courses
	router := mux.NewRouter()
	s.Router = router

	router.HandleFunc("/manager/enroll", s.ServeEnroll).Methods("GET")
	router.HandleFunc("/manager/enroll", s.Enroll).Methods("POST")

	// TODO: Basic Auth
	router.HandleFunc("/manager/courses/", s.AddCourse).Methods("POST")
	router.HandleFunc("/manager/students/{id}", s.GetStudent).Methods("GET")
	return s
}

func (s *Server) AddCourse(w http.ResponseWriter, req *http.Request) {

	var v pkg.Course
	if err := json.NewDecoder(req.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err := s.Courses.AddCourse(&v)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return

	}

	w.WriteHeader(201)

}

func (s *Server) ServeEnroll(w http.ResponseWriter, req *http.Request) {
	t, err := template.New("body").Parse(form)
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, &SignupResponse{})
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

func (s *Server) Enroll(w http.ResponseWriter, req *http.Request) {

	t, err := template.New("body").Parse(form)
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	wr := &SignupResponse{}

	courseid := req.PostFormValue("coursetoken")
	if courseid == "" {
		w.WriteHeader(400)
		e := t.Execute(w, &SignupResponse{Error: true, ErrorMessage: "Course Token Required"})
		if e != nil {
			log.Print("template executing error: ", err)
		}
		return
	}
	course, err := s.Courses.Validate(courseid)
	if err != nil {
		w.WriteHeader(400)
		e := t.Execute(w, &SignupResponse{Error: true, ErrorMessage: "Bad Course Token"})
		if e != nil {
			log.Print("template executing error: ", err)
		}
		return
	}

	username := req.PostFormValue("username")
	password := req.PostFormValue("password")
	password2 := req.PostFormValue("password2")
	if password != password2 {
		w.WriteHeader(400)
		e := t.Execute(w, &SignupResponse{Error: true, ErrorMessage: "Passwords don't match"})
		if e != nil {
			log.Print("template executing error: ", e)
		}
		return
	}

	name := req.PostFormValue("name")
	email := req.PostFormValue("email")

	student := &pkg.Student{
		Username: username,
		Password: password,
		FullName: name,
		Email:    email,
		Courses:  []*pkg.Course{course},
	}

	exists, _ := s.Students.Exists(student)
	if exists {
		log.Print("enrolling existing student in course")
		err := s.Students.Enroll(student, course)
		if err != nil {
			w.WriteHeader(400)
			log.Print("enrolling existing student: ", err)
			e := t.Execute(w, &SignupResponse{Error: true, ErrorMessage: "Error enrolling existing student in course."})
			if e != nil {
				log.Print("template executing error: ", e)
			}
			return
		}
	} else {

		log.Print("enrolling new student in course")
		err = s.Students.AddStudent(student)
		if err != nil {
			w.WriteHeader(400)
			e := t.Execute(w, &SignupResponse{Error: true, ErrorMessage: err.Error()})
			if e != nil {
				log.Print("template executing error: ", e)
			}
			return
		}
		err = s.ExecManager.NewStudent(student)
		if err != nil {

			w.WriteHeader(400)
			e := t.Execute(w, &SignupResponse{Error: true, ErrorMessage: err.Error()})
			if e != nil {
				log.Print("template executing error: ", err)
			}
			return
		}
	}
	err = s.ExecManager.NewEnrollment(student, course)
	if err != nil {

		w.WriteHeader(400)
		e := t.Execute(w, &SignupResponse{Error: true, ErrorMessage: err.Error()})
		if e != nil {
			log.Print("template executing error: ", err)
		}
		return
	}
	environ := GetEnviron(student.Username, student.Username)
	pretty.Println(environ)
	rc := environ.IDE(student.Username, student.Password)
	log.Println("return code", rc)
	w.WriteHeader(201)
	wr.Username = student.Username
	wr.ShellURL = fmt.Sprintf("https://students.brianketelsen.com/u/%s", student.Username)
	wr.IDEURL = "https://ide.brianketelsen.com/"
	e := t.Execute(w, wr)
	if e != nil {
		log.Print("template executing error: ", err)
	}

}

func (s *Server) GetStudent(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("error"))

	}
	student, err := s.Students.Student(id)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("error"))
	}

	w.Write([]byte(fmt.Sprintf("%v", student)))
}

const form = `<!DOCTYPE html>
<html>
<link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/font-awesome/4.3.0/css/font-awesome.min.css">
<link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">
<link href='http://fonts.googleapis.com/css?family=Varela+Round' rel='stylesheet' type='text/css'>
<script   src="http://code.jquery.com/jquery-3.1.1.min.js"   integrity="sha256-hVVnYaiADRTO2PzUGmuLJr8BLUSjGIZsDYGmIJLv2b8="   crossorigin="anonymous"></script>
<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1" />
<body>

<div class="container">
{{ if .Error }}
<div class="alert alert-danger">
<strong>Error:<strong> {{.ErrorMessage}}
</div>
{{ end  }}
{{ if .Username}}
<div class="alert alert-success">
  <strong>Created<strong> {{.Username}} was created<br/>
  <strong>Web IDE:<strong> <a target=new href="{{.IDEURL}}">Click Here</a> 
</div>
{{ end  }}
        <div class="row centered-form">
        <div class="col-xs-12 col-sm-8 col-md-4 col-sm-offset-2 col-md-offset-4">
        	<div class="panel panel-default">
        		<div class="panel-heading">
				<h3 class="panel-title">Enter your user information to join:</h3>
			 			</div>
			 			<div class="panel-body">
			    		<form role="form" action="/manager/enroll" method="POST">
			    					<div class="form-group">
										<input type="text" name="username" id="username" class="form-control input-sm" placeholder="Enter a Username">
			    					</div>
			    					<div class="form-group">
			    						<input type="text" name="name" id="name" class="form-control input-sm" placeholder="Full Name">
			    					</div>

			    			<div class="form-group">
			    				<input type="email" name="email" id="email" class="form-control input-sm" placeholder="Email Address">
			    			</div>

			    					<div class="form-group">
			    						<input type="password" name="password" id="password" class="form-control input-sm" placeholder="Enter a Password">
			    					</div>

			    					<div class="form-group">
			    						<input type="password" name="password2" id="password2" class="form-control input-sm" placeholder="Repeat Password">
			    					</div>
			    			<div class="form-group">
			    				<input type="text" name="coursetoken" id="coursetoken" class="form-control input-sm" placeholder="CourseToken">
			    			</div>
			    			
			    			<input type="submit" value="Register" class="btn btn-info btn-block">
			    		
			    		</form>
			    	</div>
	    		</div>
    		</div>
    	</div>
    </div>
</body>
</html>`
