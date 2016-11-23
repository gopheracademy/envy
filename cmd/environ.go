package envy

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
)

func init() {
	cmdEnviron.AddCommand(cmdEnvironRebuild)
	cmdEnviron.AddCommand(cmdEnvironList)
	Cmd.AddCommand(cmdEnviron)
}

var cmdEnviron = &cobra.Command{
	Use: "environ",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var cmdEnvironRebuild = &cobra.Command{
	Short: "rebuild environment image",
	Long:  `Rebuild does a Docker build with your environment Dockerfile.`,

	Use: "rebuild [--force]", // TODO
	Run: func(_ *cobra.Command, args []string) {
		session := GetSession(os.Getenv("ENVY_USER"), os.Getenv("ENVY_SESSION"))
		environ := session.Environ()
		log.Println(session.User.Name, "| rebuilding environ", environ.Name)
		cmd := exec.Command("/bin/docker", "build", "-t", environ.DockerImage(), ".")
		cmd.Dir = environ.Path()
		run(cmd)
		os.Exit(128)
	},
}

var cmdEnvironList = &cobra.Command{
	Short: "list environments",
	Long:  `Lists environments for this user.`,

	Use: "ls",
	Run: func(cmd *cobra.Command, args []string) {
		session := GetSession(os.Getenv("ENVY_USER"), os.Getenv("ENVY_SESSION"))
		for _, environ := range session.User.Environs() {
			fmt.Println(environ)
		}
	},
}

type Environ struct {
	User    *User
	Name    string
	IDEPort int
}

func (e *Environ) Path(parts ...string) string {
	return Envy.Path(append([]string{"users", e.User.Name, "environs", e.Name}, parts...)...)
}

func (e *Environ) DockerImage() string {
	return fmt.Sprintf("%s/%s", e.User.Name, e.Name)
}

func (e *Environ) IDEImage() string {
	return "bketelsen/codebox"
}
func (e *Environ) DockerName() string {
	return fmt.Sprintf("%s.%s", e.User.Name, e.Name)
}

func (e *Environ) IDEName() string {
	return fmt.Sprintf("%s.%s", e.User.Name, "ide")
}
func (e *Environ) IDE(username, password string) int {
	log.Println(e.User.Name, "| creating IDE")
	os.Setenv("ENVY_USER", e.User.Name)
	os.Setenv("ENVY_SESSION", e.IDEName())
	for {
		dockerRemove(e.IDEName())
		args := []string{"run", "-d", "-P",
			fmt.Sprintf("--name=%s", e.IDEName()),

			fmt.Sprintf("--env=ENVY_RELOAD=%v", int32(time.Now().Unix())),
			fmt.Sprintf("--env=ENVY_SESSION=%s", e.IDEName()),
			fmt.Sprintf("--env=ENVY_USER=%s", e.User.Name),
			"--env=DOCKER_HOST=unix:///var/run/docker.sock",
			"--env=ENV=/etc/envyrc",

			fmt.Sprintf("--volume=%s:/var/run/docker.sock", Envy.HostPath(e.Path("run/docker.sock"))),
			fmt.Sprintf("--volume=%s:/var/run/envy.sock:ro", Envy.HostPath(e.Path("run/envy.sock"))),
			fmt.Sprintf("--volume=%s:/etc/envyrc:ro", Envy.HostPath(e.Path("envyrc"))),
			fmt.Sprintf("--volume=%s:/root/environ", Envy.HostPath(e.Path())),
			fmt.Sprintf("--volume=%s:/root", Envy.HostPath(e.User.Path("root"))),
			fmt.Sprintf("--volume=%s:/home/%s", Envy.HostPath(e.User.Path("home")), e.User.Name),
			fmt.Sprintf("--volume=%s:/sbin/envy:ro", Envy.HostPath("bin/envy")),
			fmt.Sprintf("--volume=%s:/sbin/docker:ro", Envy.HostPath("bin/docker")),
		}
		if e.User.Admin() {
			args = append(args, fmt.Sprintf("--volume=%s:/envy", Envy.HostPath()))
		}
		args = append(args, e.IDEImage())
		ws := fmt.Sprintf("/home/%s", e.User.Name)
		args = append(args, ws)
		up := fmt.Sprintf("--users=%s:%s", username, password)
		args = append(args, "--port=80")
		args = append(args, up)
		status := run(exec.Command("/bin/docker", args...))
		if status != 128 {
			return status
		}
	}
}

func GetEnviron(user, name string) *Environ {
	e := &Environ{
		Name: name,
		User: GetUser(user),
	}
	if !exists(e.Path()) {
		copyTree(Envy.DataPath("environ"), e.Path())
	}
	mkdirAll(e.Path("run"))
	if !dockerRunning(e.DockerName()) {
		dockerRemove(e.DockerName())
		log.Println(user, "| starting dind for environ", e.Name)
		portBindings := map[docker.Port][]docker.PortBinding{
			"80/tcp": {{HostIP: "0.0.0.0", HostPort: "18080"}}}
		dockerRunDetached(docker.CreateContainerOptions{
			Name: e.DockerName(),
			Config: &docker.Config{
				Hostname: e.Name,
				Image:    "docker:1.12.1-dind",
			},
			HostConfig: &docker.HostConfig{
				PortBindings:  portBindings,
				Privileged:    true,
				RestartPolicy: docker.RestartPolicy{Name: "always"},
				Binds: []string{
					fmt.Sprintf("%s:/usr/bin/docker", Envy.HostPath("bin/docker")),
					fmt.Sprintf("%s:/var/run", Envy.HostPath(e.Path("run"))),
				},
			},
		})
	}
	/*	if !dockerImage(e.DockerImage()) {
			log.Println(user, "| building environ", e.Name)
			cmd := exec.Command("/bin/docker", "build", "-t", e.DockerImage(), ".")
			cmd.Dir = e.Path()
			assert(cmd.Run())
		}
	*/
	return e
}
