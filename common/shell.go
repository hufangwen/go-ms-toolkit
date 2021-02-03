/*
@Time : 29/1/2021 公元 16:15
@Author : philiphu
@File : shell
@Software: GoLand
*/
package common

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"github.com/go-cmd/cmd"
	"time"
)


// 自带的垃圾 TODO 废弃
func Shell(command string) ([]byte, error) {
	cmd := exec.Command("sh", "-c", command)
	return cmd.Output()
}

// 废弃
func ShellWithTimeout(command string, timeout int) ([]byte, error) {
	if "" == command {
		return nil, nil
	}
	resp := make([]byte, 0)
	var err error
	ch := make(chan string,1)
	go func(ch chan string,resp []byte) {
		resp, err = Shell(command)
		if err != nil {
			ch <- err.Error()
			close(ch)
			return
		}
		ch <- "done"
		close(ch)
	}(ch,resp)
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		return nil, fmt.Errorf("process request is timeout:%d", time.Duration(timeout)*time.Second)
	case done := <-ch:
		if "done" == done {
			return resp, nil
		} else {
			return nil, errors.New(done)
		}
	}
}

func GoShell(shell string) []string {
	c := cmd.NewCmd("bash","-c",shell)
	<-c.Start()
	return c.Status().Stdout
}

func TimeoutGoShell(timeout time.Duration,	shell string) ([]byte,error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	cmdarray := []string{"-c", fmt.Sprintf("%s", shell)}
	cmd := exec.CommandContext(ctx, "bash", cmdarray...)
	out, err := cmd.CombinedOutput()
	return out,err

}
