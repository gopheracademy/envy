package bolt

import (
	"encoding/binary"
	"time"

	"github.com/boltdb/bolt"
	envy "github.com/gophertrain/envy/pkg"
)

// Bolt Implmentation ideas stolen from github.com/benbjohnson/wtf

// Session represents an authenticable connection to the database.
type Session struct {
	db  *bolt.DB
	now time.Time

	courseToken string
	authToken   string
	Validator   envy.Validator
	student     *envy.Student
	course      *envy.Course

	// Services
	studentStore StudentStorage
	courseStore  CourseStorage
}

// newSession returns a new instance of Session attached to db.
func newSession(db *bolt.DB) *Session {
	s := &Session{db: db}
	s.studentStore.session = s
	s.courseStore.session = s
	return s
}

func (s *Session) SetCourseToken(token string) {
	s.courseToken = token
}

// Validate checks the course token
func (s *Session) Validate() (*envy.Course, error) {
	// Return course if already validated
	if s.course != nil {
		return s.course, nil
	}

	// Authenticate using token.
	c, err := s.Validator.Validate(s.courseToken)
	if err != nil {
		return nil, err
	}

	// Cache authenticated course
	s.course = c

	return c, nil
}

// Validate checks the course token
func (s *Session) Encode(i int) (string, error) {
	return s.Validator.Encode(i)
}

// StudentStorage returns the service associated with this session
func (s *Session) StudentStorage() envy.StudentStorage { return &s.studentStore }

// CourseStorage returns the service associated with this session
func (s *Session) CourseStorage() envy.CourseStorage { return &s.courseStore }

// itob returns an 8-byte big-endian encoded byte slice of v.
//
// This function is typically used for encoding integer IDs to byte slices
// so that they can be used as BoltDB keys.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
