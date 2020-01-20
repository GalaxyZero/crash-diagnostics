// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/crash-diagnostics/script"
	"github.com/vmware-tanzu/crash-diagnostics/ssh"
)

// cmdExec executes script on remote machines
func cmdExec(asCmd *script.AsCommand, authCmd *script.AuthConfigCommand, action script.Command, machine *script.Node, workdir string) error {

	user := asCmd.GetUserId()
	if authCmd.GetUsername() != "" {
		user = authCmd.GetUsername()
	}

	privKey := authCmd.GetPrivateKey()
	if privKey == "" {
		return fmt.Errorf("missing private key file")
	}

	//for _, action := range src.Actions {
	switch cmd := action.(type) {
	case *script.CopyCommand:
		if err := execCopy(user, privKey, machine, asCmd, cmd, workdir); err != nil {
			return err
		}
	case *script.CaptureCommand:
		// capture command output
		if err := execCapture(user, privKey, machine.Address(), cmd, workdir); err != nil {
			return err
		}
	case *script.RunCommand:
		if err := execRun(user, privKey, machine.Address(), cmd, workdir); err != nil {
			return err
		}
	default:
		logrus.Errorf("Unsupported command %T", cmd)
	}
	//}

	return nil
}

func execCapture(user, privKey, hostAddr string, cmdCap *script.CaptureCommand, workdir string) error {
	sshc := ssh.New(user, privKey)
	if err := sshc.Dial(hostAddr); err != nil {
		return err
	}
	defer sshc.Hangup()

	cmdStr, err := cmdCap.GetEffectiveCmdStr()
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s.txt", sanitizeStr(cmdStr))
	filePath := filepath.Join(workdir, fileName)
	logrus.Debugf("CAPTURE command [%s] -into-> %s", cmdStr, filePath)

	cmdReader, err := sshc.SSHRun(cmdStr)
	if err != nil {
		sshErr := fmt.Errorf("CAPTURE remote command %s failed: %s", cmdStr, err)
		logrus.Warn(sshErr)
		return writeCmdError(sshErr, filePath, cmdStr)
	}

	echo := false
	switch cmdCap.GetEcho() {
	case "true", "yes", "on":
		echo = true
	}

	if err := writeCmdOutput(cmdReader, filePath, echo, cmdStr); err != nil {
		return err
	}

	return nil
}

func execRun(user, privKey, hostAddr string, cmdRun *script.RunCommand, workdir string) error {
	sshc := ssh.New(user, privKey)
	if err := sshc.Dial(hostAddr); err != nil {
		return err
	}
	defer sshc.Hangup()

	cmdStr, err := cmdRun.GetEffectiveCmdStr()
	if err != nil {
		return err
	}

	cmdReader, err := sshc.SSHRun(cmdStr)
	if err != nil {
		sshErr := fmt.Errorf("RUN remote command failed: %s: %s", cmdStr, err)
		logrus.Error(sshErr)
		return nil
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, cmdReader); err != nil {
		return fmt.Errorf("RUN: result: %s", err)
	}

	// save result
	result := strings.TrimSpace(buf.String())
	if len(result) < 1 {
		if err := os.Unsetenv("CMD_RESULT"); err != nil {
			return fmt.Errorf("RUN: unset CMD_RESULT: %s", err)
		}
		return nil
	}

	if err := os.Setenv("CMD_RESULT", result); err != nil {
		return fmt.Errorf("RUN: set CMD_RESULT: %s: %s", result, err)
	}

	switch cmdRun.GetEcho() {
	case "true", "yes", "on":
		fmt.Printf("%s\n%s\n", cmdRun.GetCmdString(), result)
	}

	return nil
}

var (
	cliScpName = "scp"
	cliScpArgs = "-rpq"
)

// execCopy uses rsync and requires both rsync and ssh to be installed
func execCopy(user, privKey string, machine *script.Node, asCmd *script.AsCommand, cmd *script.CopyCommand, dest string) error {
	if _, err := exec.LookPath(cliScpName); err != nil {
		return fmt.Errorf("remote copy: %s", err)
	}

	logrus.Debugf("Entering remote COPY command: %s", cmd.Args())

	host, err := machine.Host()
	if err != nil {
		return fmt.Errorf("COPY: %s", err)
	}
	port, err := machine.Port()
	if err != nil {
		return fmt.Errorf("COPY: %s", err)
	}

	asUid, asGid, err := asCmd.GetCredentials()
	if err != nil {
		return err
	}

	for _, path := range cmd.Paths() {

		remotePath := fmt.Sprintf("%s@%s:%s", user, host, path)

		// if path contains file pattern, adjust target
		pathDir, pathFile := filepath.Split(path)
		targetPath := filepath.Join(dest, path)
		targetDir := filepath.Dir(targetPath)
		if strings.Index(pathFile, "*") != -1 {
			targetPath = filepath.Join(dest, pathDir)
			targetDir = targetPath
		}

		if _, err := os.Stat(targetDir); err != nil {
			if !os.IsNotExist(err) {
				return err
			}

			if err := os.MkdirAll(targetDir, 0744); err != nil && !os.IsExist(err) {
				return err
			}
			logrus.Debugf("Created dir %s", targetDir)
		}

		logrus.Debugf("Copying %s to %s", path, targetPath)

		args := []string{cliScpArgs, "-o StrictHostKeyChecking=no", "-P", port, "-i", privKey, remotePath, targetPath}
		output, err := CliRun(uint32(asUid), uint32(asGid), cliScpName, args...)
		if err != nil {
			msgBytes, _ := ioutil.ReadAll(output)
			cliErr := fmt.Errorf("scp command failed: %s: %s", err, string(msgBytes))
			logrus.Warn(cliErr)
			return writeCmdError(cliErr, targetPath, fmt.Sprintf("%s %s", cliScpName, strings.Join(args, " ")))
		}
		logrus.Debug("Remote copy succeeded:", remotePath)
	}

	return nil
}
