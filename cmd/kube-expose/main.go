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

	"gopkg.in/yaml.v2"

	"github.com/goforbroke1006/kube-expose/pkg/domain"
)

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

	for _, namespace := range cc.Namespaces {
		for name, resourceDescr := range namespace.Resources {
			go func(namespaceName string, resName string, rsc domain.Resource) {
				// TODO:

				var outbuf, errbuf bytes.Buffer

				getResourceNameCmd := fmt.Sprintf("kubectl get %s -n %s | grep -m1 %s | awk '{print $1}'", rsc.Type, namespaceName, resName)
				cmd := exec.Command(app, "-c", getResourceNameCmd)
				cmd.Stdout = &outbuf
				cmd.Stderr = &errbuf

				if err := cmd.Run(); err != nil {
					fmt.Println("ERROR:", err.Error())
					return
				}

				resourceID := strings.TrimSpace(outbuf.String())

				for _, pp := range rsc.Ports {
					go func(namespaceName string, resType, portsPair string) {
						forwardPortCmd := fmt.Sprintf("kubectl port-forward -n %s %s/%s %s", namespaceName, resType, resourceID, portsPair)
						cmd := exec.Command(app, "-c", forwardPortCmd)
						cmd.Stdout = &outbuf
						cmd.Stderr = &errbuf

						fmt.Println("$", forwardPortCmd)

						if err = cmd.Start(); err != nil {
							fmt.Println("ERROR:", err.Error())
							return
						}

						fmt.Println("	", resourceID, "up", "ports="+pp, fmt.Sprintf("pid=%d", cmd.Process.Pid))

						pidsMx.Lock()
						pids = append(pids, cmd.Process.Pid)
						pidsMx.Unlock()
					}(namespaceName, rsc.Type, pp)

				}
			}(namespace.Name, name, resourceDescr)
		}
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
		go func(pid int) {
			cmd := exec.Command(app, "-c", fmt.Sprintf("kill %d", pid))
			_ = cmd.Start()
		}(pid)
	}
}
