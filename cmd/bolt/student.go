package bolt

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/boltdb/bolt"
	"github.com/gophertrain/envy/cmd/bolt/internal"
	envy "github.com/gophertrain/envy/pkg"
)

var _ envy.StudentStorage = &StudentStorage{}

// StudentStorage represents a service for managing students
type StudentStorage struct {
	session *Session
}

func (s *StudentStorage) Student(id int) (*envy.Student, error) {
	// Start read-only transaction.
	tx, err := s.session.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Find and unmarshal student
	var d envy.Student
	if v := tx.Bucket([]byte("Students")).Get(itob(int(id))); v == nil {
		return nil, nil
	} else if err := internal.UnmarshalStudent(v, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

func (s *StudentStorage) Exists(st *envy.Student) (bool, error) {

	err := s.session.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Students"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
			var d envy.Student
			if err := internal.UnmarshalStudent(v, &d); err != nil {
				return err
			}
			if d.UID == st.UID {
				return errors.New("UserID Exists in system.")
			}
			if d.Username == st.Username {
				st.UID = d.UID
				return errors.New("Username Exists in system")
			}

			if d.Email == st.Email {
				st.UID = d.UID
				return errors.New("Email Exists in system")
			}
		}

		return nil
	})
	if err != nil {
		return true, err
	}
	return false, nil
}

func (s *StudentStorage) AddStudent(st *envy.Student) error {

	// Start read-write transaction.
	tx, err := s.session.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// Create new ID.

	b := tx.Bucket([]byte("Students"))

	id, _ := b.NextSequence()
	st.UID = int(id)

	pwbytes, err := hashPassword(st.Password)
	if err != nil {
		return err
	}
	st.Password = string(pwbytes)
	// Marshal and insert record.
	if v, err := internal.MarshalStudent(st); err != nil {
		return err
	} else if err := b.Put(itob(int(st.UID)), v); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *StudentStorage) RemoveStudent(*envy.Student) error {
	panic("not implemented")
}

func (s *StudentStorage) List() ([]*envy.Student, error) {

	tx, err := s.session.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var students []*envy.Student

	s.session.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("Students"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
			var d envy.Student
			if err := internal.UnmarshalStudent(v, &d); err != nil {
				return err
			}
			students = append(students, &d)
		}

		return nil
	})

	return students, err

}

func (s *StudentStorage) Enroll(st *envy.Student, c *envy.Course) error {
	tx, err := s.session.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Find and unmarshal student
	var d envy.Student
	if v := tx.Bucket([]byte("Students")).Get(itob(int(st.UID))); v == nil {
		return errors.New("Student Not Found")
	} else if err := internal.UnmarshalStudent(v, &d); err != nil {
		return err
	}

	for _, course := range d.Courses {
		if course.ID == c.ID {
			return nil // already enrolled
		}
	}
	d.Courses = append(d.Courses, c)

	b := tx.Bucket([]byte("Students"))
	// Marshal and insert record.
	if v, err := internal.MarshalStudent(&d); err != nil {
		return err
	} else if err := b.Put(itob(int(d.UID)), v); err != nil {
		return err
	}

	return tx.Commit()

}

func hashPassword(password string) ([]byte, error) {
	// Hashing the password with the default cost of 10
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

}

func checkPassword(hash, password []byte) error {
	// Comparing the password with the hash
	return bcrypt.CompareHashAndPassword(hash, password)

}
