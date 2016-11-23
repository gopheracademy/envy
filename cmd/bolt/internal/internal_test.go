package internal_test

import (
	"reflect"
	"testing"

	"github.com/gophertrain/envy/bolt/internal"
	envy "github.com/gophertrain/envy/pkg"
)

// Ensure student can be marshaled and unmarshaled.
func TestMarshalStudent(t *testing.T) {
	c := envy.Course{
		ID:              1000,
		Name:            "Go Beyond",
		Description:     "Getting more out of Go",
		Instructor:      "Brian Ketelsen",
		InstructorEmail: "me@brianketelsen.com",
	}

	v := envy.Student{
		UID:           1001,
		Username:      "bketelsen",
		Password:      "XYX123",
		FullName:      "Brian Ketelsen",
		Email:         "me@brianketelsen.com",
		HomeDirectory: "/home/bketelsen",
		Courses:       []*envy.Course{&c},
	}

	var other envy.Student
	if buf, err := internal.MarshalStudent(&v); err != nil {
		t.Fatal(err)
	} else if err := internal.UnmarshalStudent(buf, &other); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, other) {
		t.Fatalf("unexpected copy: %#v", other)
	}
}

// Ensure course can be marshaled and unmarshaled.
func TestMarshalCourse(t *testing.T) {
	v := envy.Course{
		ID:              1000,
		Name:            "Go Beyond",
		Description:     "Getting more out of Go",
		Instructor:      "Brian Ketelsen",
		InstructorEmail: "me@brianketelsen.com",
	}

	var other envy.Course
	if buf, err := internal.MarshalCourse(&v); err != nil {
		t.Fatal(err)
	} else if err := internal.UnmarshalCourse(buf, &other); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, other) {
		t.Fatalf("unexpected copy: %#v", other)
	}
}
