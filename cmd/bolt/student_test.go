package bolt_test

import (
	"reflect"
	"testing"

	"github.com/gophertrain/envy/studentmgr"
)

// Ensure course can be created and retrieved.
func TestStudentService_AddStudent(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Connect().StudentStorage()

	// Mock authentication.
	c.Authenticator.AuthenticateFn = func(_ string) (*studentmgr.Student, error) {
		return &studentmgr.Student{UID: 1000}, nil
	}

	student := studentmgr.Student{
		UID:           1000,
		Username:      "bketelsen",
		Password:      "Password",
		FullName:      "Brian Ketelsen",
		Email:         "bketelsen@gmail.com",
		HomeDirectory: "/home/bketelsen",
	}

	// Create new student
	if err := s.AddStudent(&student); err != nil {
		t.Fatal(err)
	} else if student.UID != 1000 {
		t.Fatalf("unexpected id: %d", student.UID)
	} else if student.Username != "bketelsen" {
		t.Fatalf("unexpected username: %s", student.Username)
	}

	// Retrieve dial and compare.
	other, err := s.Student(1000)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&student, other) {
		t.Fatalf("unexpected course: %#v", other)
	}
}

func TestStudentService_WithCourse(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Connect().StudentStorage()

	// Mock authentication.
	c.Authenticator.AuthenticateFn = func(_ string) (*studentmgr.Student, error) {
		return &studentmgr.Student{UID: 1000}, nil
	}

	student := studentmgr.Student{
		UID:           1000,
		Username:      "bketelsen",
		Password:      "Password",
		FullName:      "Brian Ketelsen",
		Email:         "bketelsen@gmail.com",
		HomeDirectory: "/home/bketelsen",
	}

	course := studentmgr.Course{
		ID:              1000,
		Name:            "Go Beyond",
		Description:     "Getting more out of Go",
		Instructor:      "Brian Ketelsen",
		InstructorEmail: "me@brianketelsen.com",
	}

	student.Courses = append(student.Courses, &course)

	// Create new student
	if err := s.AddStudent(&student); err != nil {
		t.Fatal(err)
	}
	// Retrieve dial and compare.
	other, err := s.Student(1000)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&student, other) {
		t.Fatalf("unexpected course: %#v", other)
	}
}

func TestStudentService_Enroll(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Connect().StudentStorage()

	// Mock authentication.
	c.Authenticator.AuthenticateFn = func(_ string) (*studentmgr.Student, error) {
		return &studentmgr.Student{UID: 1000}, nil
	}

	student := studentmgr.Student{
		UID:           1000,
		Username:      "bketelsen",
		Password:      "Password",
		FullName:      "Brian Ketelsen",
		Email:         "bketelsen@gmail.com",
		HomeDirectory: "/home/bketelsen",
	}

	course := studentmgr.Course{
		ID:              1000,
		Name:            "Go Beyond",
		Description:     "Getting more out of Go",
		Instructor:      "Brian Ketelsen",
		InstructorEmail: "me@brianketelsen.com",
	}

	// Create new student
	if err := s.AddStudent(&student); err != nil {
		t.Fatal(err)
	}

	if err := s.Enroll(&student, &course); err != nil {
		t.Fatal(err)
	}
	other, err := s.Student(1000)
	if err != nil {
		t.Fatal(err)
	}
	if len(other.Courses) != 1 {
		t.Fatal("Expected student to have a course")
	}
	if other.Courses[0].ID != 1000 {
		t.Fatal("Expected enrolled course to be correct")
	}
}

func TestStudentService_EnrollDup(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Connect().StudentStorage()

	// Mock authentication.
	c.Authenticator.AuthenticateFn = func(_ string) (*studentmgr.Student, error) {
		return &studentmgr.Student{UID: 1000}, nil
	}

	student := studentmgr.Student{
		UID:           1000,
		Username:      "bketelsen",
		Password:      "Password",
		FullName:      "Brian Ketelsen",
		Email:         "bketelsen@gmail.com",
		HomeDirectory: "/home/bketelsen",
	}

	course := studentmgr.Course{
		ID:              1000,
		Name:            "Go Beyond",
		Description:     "Getting more out of Go",
		Instructor:      "Brian Ketelsen",
		InstructorEmail: "me@brianketelsen.com",
	}

	// Create new student
	if err := s.AddStudent(&student); err != nil {
		t.Fatal(err)
	}

	if err := s.Enroll(&student, &course); err != nil {
		t.Fatal(err)
	}

	if err := s.Enroll(&student, &course); err != nil {
		t.Fatal(err)
	}
	other, err := s.Student(1000)
	if err != nil {
		t.Fatal(err)
	}
	if len(other.Courses) != 1 {
		t.Fatal("Expected student to have only one course")
	}
	if other.Courses[0].ID != 1000 {
		t.Fatal("Expected enrolled course to be correct")
	}
}
