package envy

import (
	"os"
	"path/filepath"

	"github.com/gophertrain/envy/pkg"
)

type ExecManager struct {
	ScriptPath string
}

func (em *ExecManager) Exists(username string) (bool, error) {
	path := filepath.Join("/", "home", username)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err

}
func (em *ExecManager) NewStudent(student *pkg.Student) error {

	if !grepFile(Envy.Path("config/users"), student.Username) {
		appendFile(Envy.Path("config/users"), student.Username)
	}

	return nil
}

func (em *ExecManager) NewEnrollment(student *pkg.Student, course *pkg.Course) error {
	return nil
}
