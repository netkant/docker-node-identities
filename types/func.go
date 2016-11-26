package types

// ...
import (
    "strings"
    "bytes"
    "os/exec"
    "syscall"
)

// ...
func getIndexValue(parts []string, index int, defaultValue string) string {
    x, v := index, defaultValue
	if len(parts)-1 >= x  {
		v = parts[x]
	}

    return strings.Trim(v, " ")
}

// ...
func commandExists(name string) string {
    stdout, exitCode := commandExecute("which", name)
    if exitCode > 0 {
        return ""
    }

    return strings.TrimSuffix(string(stdout), "\n")
}

// ...
func commandExecute(name string, arg ...string) (stdout string, exitCode int) {
    var outbuf bytes.Buffer

    // ...
    cmd := exec.Command(name, arg...)
    cmd.Stdout = &outbuf

    // ...
    err := cmd.Run()
    stdout = outbuf.String()

    if err != nil {
        // try to get the exit code
        if exitError, ok := err.(*exec.ExitError); ok {
            ws := exitError.Sys().(syscall.WaitStatus)
            exitCode = ws.ExitStatus()
        } else {
            exitCode = 1
        }
    } else {
        // success, exitCode should be 0 if go is ok
        ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
        exitCode = ws.ExitStatus()
    }

    return stdout, exitCode
}
