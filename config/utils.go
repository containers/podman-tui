package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"time"

	cntssh "github.com/containers/common/pkg/ssh"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman/v4/libpod/define"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func configPath() (string, error) {
	if configHome := os.Getenv("XDG_CONFIG_HOME"); configHome != "" {
		return filepath.Join(configHome, _configPath), nil
	}

	home, err := utils.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, UserAppConfig), nil
}

// localNodeUnixSocket return local node unix socket file.
func localNodeUnixSocket() string {
	var (
		sockDir string
		socket  string
	)

	currentUser := os.Getenv("USER")
	uid := os.Getenv("UID")

	if currentUser == "root" || uid == "0" {
		sockDir = "/run/"
	} else {
		sockDir = os.Getenv("XDG_RUNTIME_DIR")
	}

	socket = "unix:" + sockDir + "/podman/podman.sock"

	return socket
}

func getUserInfo(uri *url.URL) (*url.Userinfo, error) {
	var (
		usr *user.User
		err error
	)

	if u, found := os.LookupEnv("_CONTAINERS_ROOTLESS_UID"); found {
		usr, err = user.LookupId(u)
		if err != nil {
			return nil, fmt.Errorf("%w failed to lookup rootless user", err)
		}
	} else {
		usr, err = user.Current()
		if err != nil {
			return nil, fmt.Errorf("%w failed to obtain current user", err)
		}
	}

	pw, set := uri.User.Password()
	if set {
		return url.UserPassword(usr.Username, pw), nil
	}

	return url.User(usr.Username), nil
}

// most of the codes are from https://github.com/containers/podman/blob/main/cmd/podman/system/connection/add.go.
func getUDS(uri *url.URL, iden string) (string, error) {
	cfg, err := validateAndConfigure(uri, iden)
	if err != nil {
		return "", fmt.Errorf("%w failed to validate", err)
	}

	dial, err := ssh.Dial("tcp", uri.Host, cfg)
	if err != nil {
		return "", fmt.Errorf("%w failed to connect", err)
	}

	defer dial.Close()

	session, err := dial.NewSession()
	if err != nil {
		return "", fmt.Errorf("%w failed to create new ssh session on %q", err, uri.Host)
	}

	defer session.Close()

	// Override podman binary for testing etc
	podman := "podman"
	if v, found := os.LookupEnv("PODMAN_BINARY"); found {
		podman = v
	}

	infoJSON, err := execRemoteCommand(dial, podman+" info --format=json")
	if err != nil {
		return "", err
	}

	var info define.Info
	if err := json.Unmarshal(infoJSON, &info); err != nil {
		return "", fmt.Errorf("%w failed to parse 'podman info' results", err)
	}

	if info.Host.RemoteSocket == nil || len(info.Host.RemoteSocket.Path) == 0 {
		return "", fmt.Errorf("%w %s", ErrRemotePodmanUDSReport, uri.Host)
	}

	return info.Host.RemoteSocket.Path, nil
}

// validateAndConfigure will take a ssh url and an identity key (rsa and the like)
// and ensure the information given is valid iden can be blank to mean no identity key
// once the function validates the information it creates and returns an ssh.ClientConfig.
func validateAndConfigure(uri *url.URL, iden string) (*ssh.ClientConfig, error) { //nolint:cyclop
	var signers []ssh.Signer

	if iden != "" { // iden might be blank if coming from image scp or if no validation is needed
		value := iden
		passPhrase := ""

		if v, found := os.LookupEnv("CONTAINER_PASSPHRASE"); found {
			passPhrase = v
		}

		if passPhrase == "" {
			passPhrase = "_empty_pass_"
		}

		s, err := cntssh.PublicKey(value, []byte(passPhrase))
		if err != nil {
			infoText := "set 'CONTAINER_PASSPHRASE' variable if password is required"

			return nil, fmt.Errorf("%w failed to read identity %q, %s", err, value, infoText)
		}

		signers = append(signers, s)
		log.Debug().Msgf("config: SSH Ident Key %q %s %s", value, ssh.FingerprintSHA256(s.PublicKey()), s.PublicKey().Type())
	}

	// validate ssh information, specifically the unix file socket used by the ssh agent.
	if sock, found := os.LookupEnv("SSH_AUTH_SOCK"); found {
		log.Debug().Msgf("config: Found SSH_AUTH_SOCK %q, ssh-agent signer enabled", sock)

		c, err := net.Dial("unix", sock)
		if err != nil {
			return nil, err
		}

		agentSigners, err := agent.NewClient(c).Signers()
		if err != nil {
			return nil, err
		}

		signers = append(signers, agentSigners...)

		for _, s := range agentSigners {
			log.Debug().Msgf("config: SSH Agent Key %s %s", ssh.FingerprintSHA256(s.PublicKey()), s.PublicKey().Type())
		}
	}

	// now we validate and check for the authorization methods, most notaibly public key authorization.
	var authMethods []ssh.AuthMethod

	if len(signers) > 0 {
		dedup := make(map[string]ssh.Signer)

		for _, s := range signers {
			fp := ssh.FingerprintSHA256(s.PublicKey())
			if _, found := dedup[fp]; found {
				log.Debug().Msgf("config: Dedup SSH Key %s %s", ssh.FingerprintSHA256(s.PublicKey()), s.PublicKey().Type())
			}

			dedup[fp] = s
		}

		var uniq []ssh.Signer

		for _, s := range dedup {
			uniq = append(uniq, s)
		}

		authMethods = append(authMethods, ssh.PublicKeysCallback(func() ([]ssh.Signer, error) {
			return uniq, nil
		}))
	}

	tick, err := time.ParseDuration("40s")
	if err != nil {
		return nil, err
	}

	cfg := &ssh.ClientConfig{
		User:            uri.User.Username(),
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint:gosec
		Timeout:         tick,
	}

	return cfg, nil
}

// execRemoteCommand takes a ssh client connection and a command to run and executes the
// command on the specified client. The function returns the Stdout from the client or the Stderr.
func execRemoteCommand(dial *ssh.Client, run string) ([]byte, error) {
	sess, err := dial.NewSession() // new ssh client session
	if err != nil {
		return nil, err
	}

	defer sess.Close()

	var (
		buffer    bytes.Buffer
		bufferErr bytes.Buffer
	)

	sess.Stdout = &buffer    // output from client funneled into buffer
	sess.Stderr = &bufferErr // err form client funneled into buffer

	if err := sess.Run(run); err != nil { // run the command on the ssh client
		return nil, fmt.Errorf("%w %s", err, bufferErr.String())
	}

	return buffer.Bytes(), nil
}
