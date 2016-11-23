package pkg

type Student struct {
	UID           int       `json:"uid"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	FullName      string    `json:"full_name"`
	Email         string    `json:"email"`
	HomeDirectory string    `json:"home_directory"`
	Courses       []*Course `json:"courses"`
}

type Course struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Instructor      string `json:"instructor"`
	InstructorEmail string `json:"instructor_email"`
}

type StudentStorage interface {
	Student(id int) (*Student, error)
	Exists(*Student) (bool, error)
	AddStudent(*Student) error
	RemoveStudent(*Student) error
	Enroll(*Student, *Course) error
	List() ([]*Student, error)
}

type CourseStorage interface {
	Course(int) (*Course, error)
	AddCourse(*Course) error
	RemoveCourse(*Course) error
	Validator
}

type UserManager interface {
	CheckUsername(*Student) error
	AddUser(*Student) error
	RemoveUser(*Student) error
}

// Validator represents a service for authenticating courses
// by validating the token
type Validator interface {
	Validate(token string) (*Course, error)
	Encode(id int) (string, error)
}

// Client creates a connection to the services.
type Client interface {
	Connect() Session
}

// Session represents authenticable connection to the services.
type Session interface {
	SetAuthToken(token string)
	StudentStorage() StudentStorage
	CourseStorage() CourseStorage
}
