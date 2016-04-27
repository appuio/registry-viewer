package registry

import "fmt"
// simport "bufio"
import "strings"
import "os/exec"
import "bytes"
//import "net/http"

type Proc struct {
  stdout string
  stderr string
  err error
  checkedErrors bool
}

func (proc *Proc) CheckErrors() {
  if proc.checkedErrors && proc.err != nil {
    panic(proc.err.Error() + "\n" + proc.stderr)
  }
}

func (proc *Proc) Stdout() string {
  proc.CheckErrors()
  return proc.stdout
}

func (proc *Proc) StdoutBytes() []byte {
  proc.CheckErrors()
  return []byte(proc.stdout)
}

func (proc *Proc) StdoutLines() []string {
  proc.CheckErrors()
  return strings.Split(proc.stdout, "\n")
}

func (proc *Proc) Err() error {
  proc.checkedErrors = true
  return proc.err
}

func Sh(cmd string, a ...interface{}) *Proc {
//  var proc proc
  var command *exec.Cmd
  if len(a) == 0 {
    command = exec.Command("bash", "-c", cmd)
  } else {
    command = exec.Command("bash", "-c", fmt.Sprintf(cmd, a...))
  }
  var stderr bytes.Buffer
  command.Stderr = &stderr
  out, err := command.Output()
  
  return &Proc { 
    stdout: strings.TrimSuffix(string(out), "\n"),
    stderr: strings.TrimSuffix(stderr.String(), "\n"),
    err: err,
  }
}

/* func sh(cmd string, a ...interface{}) string {
  proc := sh2(cmd, a...)
  if proc.err != nil {
      panic(cmd + " " + proc.err.Error() + "\n" + proc.stderr)
  }

  return proc.out
} */

