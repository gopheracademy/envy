package internal

import (
	"github.com/gogo/protobuf/proto"
	envy "github.com/gophertrain/envy/pkg"
)

//go:generate protoc --gogo_out=. internal.proto

func MarshalStudent(s *envy.Student) ([]byte, error) {
	c := make([]*Course, len(s.Courses))
	for x, s := range s.Courses {
		c[x] = &Course{
			ID:              proto.Int64(int64(s.ID)),
			Name:            proto.String(s.Name),
			Description:     proto.String(s.Description),
			Instructor:      proto.String(s.Instructor),
			InstructorEmail: proto.String(s.InstructorEmail),
		}
	}
	return proto.Marshal(&Student{
		UID:           proto.Int64(int64(s.UID)),
		Username:      proto.String(s.Username),
		Password:      proto.String(s.Password),
		FullName:      proto.String(s.FullName),
		Email:         proto.String(s.Email),
		HomeDirectory: proto.String(s.HomeDirectory),
		Courses:       c,
	})

}

func UnmarshalStudent(data []byte, d *envy.Student) error {
	var pb Student
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}

	d.UID = int(pb.GetUID())
	d.Username = pb.GetUsername()
	d.Password = pb.GetPassword()
	d.FullName = pb.GetFullName()
	d.Email = pb.GetEmail()
	d.HomeDirectory = pb.GetHomeDirectory()
	for _, course := range pb.GetCourses() {
		c := &envy.Course{}
		c.ID = int(course.GetID())
		c.Name = course.GetName()
		c.Description = course.GetDescription()
		c.Instructor = course.GetInstructor()
		c.InstructorEmail = course.GetInstructorEmail()
		d.Courses = append(d.Courses, c)
	}
	return nil
}

func MarshalCourse(s *envy.Course) ([]byte, error) {

	return proto.Marshal(&Course{
		ID:              proto.Int64(int64(s.ID)),
		Name:            proto.String(s.Name),
		Description:     proto.String(s.Description),
		Instructor:      proto.String(s.Instructor),
		InstructorEmail: proto.String(s.InstructorEmail),
	})

}

func UnmarshalCourse(data []byte, d *envy.Course) error {
	var pb Course
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}

	d.ID = int(pb.GetID())
	d.Name = pb.GetName()
	d.Description = pb.GetDescription()
	d.Instructor = pb.GetInstructor()
	d.InstructorEmail = pb.GetInstructorEmail()
	return nil
}
