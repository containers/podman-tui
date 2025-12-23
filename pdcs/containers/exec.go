package containers

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v5/pkg/api/handlers"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/env"
	"github.com/rs/zerolog/log"
)

// ExecOption container exec options.
type ExecOption struct {
	Cmd          []string
	Tty          bool
	Detach       bool
	Interactive  bool
	Privileged   bool
	WorkDir      string
	EnvVariables []string
	EnvFile      []string
	User         string
	OutputStream io.Writer
	InputStream  *bufio.Reader
	TtyWidth     int
	TtyHeight    int
	DetachKeys   string
}

// NewExecSession creates a new session and returns its id.
func NewExecSession(id string, opts ExecOption) (string, error) {
	log.Debug().Msgf("pdcs: podman container (%s) exec new session", id)

	conn, err := registry.GetConnection()
	if err != nil {
		return "", err
	}

	// create new exec session
	createConfig, err := genExecCreateConfig(opts)
	if err != nil {
		return "", err
	}

	return containers.ExecCreate(conn, id, createConfig)
}

// ResizeExecTty resizes exec session tty.
func ResizeExecTty(id string, height int, width int) {
	log.Debug().Msgf("pdcs: podman container exec session (%12s) tty resize (height=%d, width=%d)", id, height, width)

	conn, err := registry.GetConnection()
	if err != nil {
		log.Error().Msgf("%v", err)

		return
	}

	for {
		response, err := containers.ExecInspect(conn, id, &containers.ExecInspectOptions{})
		if err != nil {
			log.Error().Msgf("%v", err)

			return
		}

		if response.ExitCode != 0 {
			log.Debug().Msgf("pdcs: podman container cannot resize exec session (%12s) tty, exec already exited", id)

			return
		}

		if response.Running {
			err = containers.ResizeExecTTY(conn, id, new(containers.ResizeExecTTYOptions).WithHeight(height).WithWidth(width))
			if err != nil {
				log.Error().Msgf("%v", err)

				return
			}

			log.Debug().Msgf("pdcs: podman container exec session (%12s) tty resized successfully (height=%d, width=%d)", id, height, width) //nolint:lll

			return
		}
	}
}

// Exec executes command in a given sessionOD.
func Exec(sessionID string, opts ExecOption) { //nolint:cyclop
	log.Debug().Msgf("pdcs: podman container session (%s) exec %v", sessionID, opts)

	conn, err := registry.GetConnection()
	if err != nil {
		_, err := opts.OutputStream.Write([]byte(fmt.Sprintf("%v", err))) //nolint:staticcheck
		if err != nil {
			log.Error().Msgf("%v", err)
		}

		return
	}

	if !opts.Detach {
		attach := !opts.Detach
		execStartAttachOpts := &containers.ExecStartAndAttachOptions{
			AttachOutput: &attach,
			AttachError:  &attach,
			OutputStream: &opts.OutputStream,
			ErrorStream:  &opts.OutputStream,
		}

		if opts.Interactive {
			execStartAttachOpts.AttachInput = &opts.Interactive
			execStartAttachOpts.InputStream = opts.InputStream
		}

		err := containers.ExecStartAndAttach(conn, sessionID, execStartAttachOpts)
		if err != nil {
			log.Error().Msgf("pdcs: podman session (%s) exec error %v", sessionID, err)

			_, err := opts.OutputStream.Write([]byte(fmt.Sprintf("%v", err))) //nolint:staticcheck
			if err != nil {
				log.Error().Msgf("%v", err)
			}
		}

		log.Debug().Msgf("pdcs: podman session (%s) exec finished successfully", sessionID)

		return
	}

	err = containers.ExecStart(conn, sessionID, &containers.ExecStartOptions{})
	if err != nil {
		log.Error().Msgf("pdcs: podman session (%s) exec error %v", sessionID, err)

		_, err := opts.OutputStream.Write([]byte(fmt.Sprintf("%v", err))) //nolint:staticcheck
		if err != nil {
			log.Error().Msgf("%v", err)
		}

		return
	}

	log.Debug().Msgf("pdcs: podman session (%s) exec finished successfully", sessionID)

	_, err = opts.OutputStream.Write([]byte(fmt.Sprintf("session_id ...... : %s\r\n", sessionID))) //nolint:staticcheck
	if err != nil {
		log.Error().Msgf("%v", err)
	}

	_, err = opts.OutputStream.Write([]byte(fmt.Sprintf("exec_mode  ...... : %s\r\n", "detached"))) //nolint:staticcheck
	if err != nil {
		log.Error().Msgf("%v", err)
	}

	_, err = opts.OutputStream.Write([]byte(fmt.Sprintf("exec_command .... : %s\r\n", strings.Join(opts.Cmd, " ")))) //nolint:staticcheck,lll
	if err != nil {
		log.Error().Msgf("%v", err)
	}

	_, err = opts.OutputStream.Write([]byte(fmt.Sprintf("exec_status ..... : %s\r\n", "OK"))) //nolint:staticcheck
	if err != nil {
		log.Error().Msgf("%v", err)
	}
}

func genExecCreateConfig(opts ExecOption) (*handlers.ExecCreateConfig, error) {
	var variables []string

	createCfg := &handlers.ExecCreateConfig{}
	createCfg.Cmd = opts.Cmd
	createCfg.Tty = opts.Tty
	createCfg.Detach = opts.Detach //nolint:staticcheck,nolintlint
	createCfg.WorkingDir = opts.WorkDir
	createCfg.User = opts.User

	if len(opts.EnvVariables) > 0 {
		variables = opts.EnvVariables
	}

	// parse env File
	for _, envFile := range opts.EnvFile {
		envVars, err := env.ParseFile(envFile)
		if err != nil {
			log.Error().Msgf("pdcs: podman container exec create config: %v", err)

			return nil, err
		}

		for index, key := range envVars {
			varString := fmt.Sprintf("%s=%s", key, envVars[index])
			variables = append(variables, varString)
		}
	}

	// add xterm number of LINES (rows) and COLUMES (cols)
	varLines := fmt.Sprintf("LINES=%d", opts.TtyHeight)
	varCols := fmt.Sprintf("COLUMNS=%d", opts.TtyWidth)

	variables = append(variables, varLines)
	variables = append(variables, varCols)
	createCfg.Env = variables

	createCfg.AttachStderr = !opts.Detach
	createCfg.AttachStdout = !opts.Detach

	createCfg.DetachKeys = opts.DetachKeys

	if !opts.Detach {
		createCfg.AttachStdin = opts.Interactive
	}

	return createCfg, nil
}
