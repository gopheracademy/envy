package bolt_test

import (
	"reflect"
	"testing"

	"github.com/gophertrain/envy/pkg"
)

// Ensure course can be created and retrieved.
func TestCourseService_Add(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Connect().CourseStorage()

	// Mock authentication.
	c.Authenticator.AuthenticateFn = func(_ string) (*pkg.Student, error) {
		return &pkg.Student{UID: 1000}, nil
	}

	course := pkg.Course{
		ID:              1000,
		Name:            "Go Beyond",
		Description:     "Getting more out of Go",
		Instructor:      "Brian Ketelsen",
		InstructorEmail: "me@brianketelsen.com",
	}

	// Create new course
	if err := s.AddCourse(&course); err != nil {
		t.Fatal(err)
	} else if course.ID != 1000 {
		t.Fatalf("unexpected id: %d", course.ID)
	} else if course.Name != "Go Beyond" {
		t.Fatalf("unexpected course name: %s", course.Name)
	}

	// Retrieve dial and compare.
	other, err := s.Course(1000)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&course, other) {
		t.Fatalf("unexpected course: %#v", other)
	}
}

// Ensure course can be validated from token
func TestCourseService_Validate(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Connect().CourseStorage()

	// Mock authentication.
	c.Authenticator.AuthenticateFn = func(_ string) (*pkg.Student, error) {
		return &pkg.Student{UID: 1000}, nil
	}

	course := pkg.Course{
		ID:              1000,
		Name:            "Go Beyond",
		Description:     "Getting more out of Go",
		Instructor:      "Brian Ketelsen",
		InstructorEmail: "me@brianketelsen.com",
	}

	// Create new course
	if err := s.AddCourse(&course); err != nil {
		t.Fatal(err)
	}
	hash, err := s.Encode(course.ID)
	if err != nil {
		t.Fatal(err)
	}
	other, err := s.Validate(hash)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(&course, other) {
		t.Fatalf("unexpected course: %#v", other)
	}
}
