package bolt

import (
	"fmt"
	"log"

	"github.com/gophertrain/envy/cmd/bolt/internal"
	envy "github.com/gophertrain/envy/pkg"
	hashids "github.com/speps/go-hashids"
)

var _ envy.CourseStorage = &CourseStorage{}

var h *hashids.HashID

func init() {

	hd := hashids.NewData()
	hd.Salt = "learning for everyone"
	hd.MinLength = 5
	h = hashids.NewWithData(hd)

}

type CourseStorage struct {
	session *Session
}

func (c *CourseStorage) Course(id int) (*envy.Course, error) {
	// Start read-only transaction.
	tx, err := c.session.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Find and unmarshal course
	var d envy.Course
	if v := tx.Bucket([]byte("Courses")).Get(itob(int(id))); v == nil {
		return nil, nil
	} else if err := internal.UnmarshalCourse(v, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

func (c *CourseStorage) AddCourse(d *envy.Course) error {

	// Start read-write transaction.
	tx, err := c.session.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// Create new ID.
	b := tx.Bucket([]byte("Courses"))
	// Marshal and insert record.
	if v, err := internal.MarshalCourse(d); err != nil {
		log.Println(err)
		return err
	} else if err := b.Put(itob(int(d.ID)), v); err != nil {
		log.Println(err)
		return err
	}

	return tx.Commit()
}

func (c *CourseStorage) RemoveCourse(*envy.Course) error {
	panic("not implemented")
}

func (c *CourseStorage) Validate(token string) (*envy.Course, error) {
	fmt.Println(token)
	d, err := h.DecodeWithError(token)
	if err != nil {
		return nil, err
	}
	fmt.Println(d)
	return c.Course(d[0])
}

func (c *CourseStorage) Encode(id int) (string, error) {
	return h.Encode([]int{id})
}
