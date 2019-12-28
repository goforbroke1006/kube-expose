package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/goforbroke1006/kube-expose/pkg/domain"
)
import "gopkg.in/yaml.v2"

func main() {
	data, err := ioutil.ReadFile("./kube-compose.yml")
	if err != nil {
		panic(err)
	}

	cc := domain.ComposeConfig{}

	err = yaml.Unmarshal(data, &cc)
	if err != nil {
		panic(err)
	}

	var pids []int
	var pidsMx sync.Mutex

	app := "/bin/bash"

	for name, svc := range cc.Services {
		go func(svc domain.Service) {
			// TODO:

			var outbuf, errbuf bytes.Buffer

			getPodNameCmd := fmt.Sprintf("kubectl get pods -l name=%s | grep '%s' | awk '{print $1}'", name, name)
			cmd := exec.Command(app, "-c", getPodNameCmd)
			cmd.Stdout = &outbuf
			cmd.Stderr = &errbuf
			err = cmd.Run()
			if err != nil {
				panic(err)
			}
			podName := strings.TrimSpace(outbuf.String())

			for _, portForward := range svc.Ports {
				forwardPortCmd := fmt.Sprintf("kubectl port-forward %s %s", podName, portForward)
				cmd = exec.Command(app, "-c", forwardPortCmd)
				cmd.Stdout = &outbuf
				cmd.Stderr = &errbuf
				err = cmd.Start()
				if err != nil {
					panic(err)
				}

				fmt.Println("	", podName, "up", "ports="+portForward, fmt.Sprintf("pid=%d", cmd.Process.Pid))

				pidsMx.Lock()
				pids = append(pids, cmd.Process.Pid)
				pidsMx.Unlock()
			}
		}(svc)
	}

	done := make(chan bool)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("down forwarded ports...")
		done <- true
	}()

	<-done

	for _, pid := range pids {
		cmd := exec.Command(app, "-c", fmt.Sprintf("kill %d", pid))
		_ = cmd.Start()
	}
}
