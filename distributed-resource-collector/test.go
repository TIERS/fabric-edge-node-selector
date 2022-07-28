// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/dmonteroh/distributed-resource-collector/internal"
// 	"golang.org/x/crypto/ssh"
// )

// func main() {

// 	target := internal.LatencyTarget{
// 		Hostname:     "172.26.161.241",
// 		Hostport:     "22",
// 		HostUser:     "tiers",
// 		HostPassword: "2022",
// 	}

// 	host := target.Hostname
// 	port := target.Hostport
// 	user := target.HostUser
// 	pass := target.HostPassword

// 	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
// 	fmt.Println(timestamp[len(timestamp)-1:])
// 	cmd := "echo " + timestamp[len(timestamp)-1:]

// 	//cmd := "ps"

// 	// get host public key
// 	//hostKey := getHostKey(host)

// 	// ssh client config
// 	config := &ssh.ClientConfig{
// 		User: user,
// 		Auth: []ssh.AuthMethod{
// 			ssh.Password(pass),
// 		},
// 		// allow any host key to be used (non-prod)
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

// 		// verify host public key
// 		//HostKeyCallback: ssh.FixedHostKey(hostKey),
// 		// optional host key algo list
// 		HostKeyAlgorithms: []string{
// 			ssh.KeyAlgoRSA,
// 			ssh.KeyAlgoDSA,
// 			ssh.KeyAlgoECDSA256,
// 			ssh.KeyAlgoECDSA384,
// 			ssh.KeyAlgoECDSA521,
// 			ssh.KeyAlgoED25519,
// 		},
// 		// optional tcp connect timeout
// 		Timeout: 5 * time.Second,
// 	}

// 	// connect
// 	client, err := ssh.Dial("tcp", host+":"+port, config)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer client.Close()

// 	// start session
// 	sess, err := client.NewSession()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer sess.Close()

// 	// setup standard out and error
// 	// uses writer interface
// 	sess.Stdout = os.Stdout
// 	sess.Stderr = os.Stderr

// 	// run single command
// 	err = sess.Run(cmd)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// }

// func getHostKey(host string) ssh.PublicKey {
// 	// parse OpenSSH known_hosts file
// 	// ssh or use ssh-keyscan to get initial key
// 	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	var hostKey ssh.PublicKey
// 	for scanner.Scan() {
// 		fields := strings.Split(scanner.Text(), " ")
// 		if len(fields) != 3 {
// 			continue
// 		}
// 		if strings.Contains(fields[0], host) {
// 			var err error
// 			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
// 			if err != nil {
// 				log.Fatalf("error parsing %q: %v", fields[2], err)
// 			}
// 			break
// 		}
// 	}

// 	if hostKey == nil {
// 		log.Fatalf("no hostkey found for %s", host)
// 	}

// 	return hostKey
// }

// package main

// import (
// 	"bytes"
// 	"fmt"
// 	"io"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/dmonteroh/distributed-resource-collector/internal"
// 	"golang.org/x/crypto/ssh"
// )

// func RecoverEndpoint() {
// 	if err := recover(); err != nil {
// 		msg := "Error: [Recovered] "
// 		switch errType := err.(type) {
// 		case string:
// 			msg += err.(string)
// 		case error:
// 			msg += errType.Error()
// 		default:
// 		}
// 		fmt.Println(msg)
// 	}
// }

// func main() {

// 	latencyTargetsJson := "{\"browsers\":{\"firefox\":{\"name\":\"Firefox\",\"pref_url\":\"about:config\",\"releases\":{\"1\":{\"release_date\":\"2004-11-09\",\"status\":\"retired\",\"engine\":\"Gecko\",\"engine_version\":\"1.7\"}}}}}"
// 	latencyTargets, err := internal.LatencyJsonToStrcut(string(latencyTargetsJson))
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	latencyResults := internal.LatencyResults{
// 		Source:  "::1",
// 		Results: []internal.LatencyResult{},
// 	}

// 	fmt.Printf("%+v", latencyTargets)

// 	for _, target := range latencyTargets.Targets {
// 		c1 := make(chan internal.LatencyResult)
// 		go func(target internal.LatencyTarget) {
// 			fmt.Printf("%+v", target)
// 			timestamp := strconv.FormatInt(time.Now().Unix(), 10)
// 			//cmd := "touch latency_" + latencyTargets.Hostname + " && echo '" + timestamp[len(timestamp)-1:] + "' > latency_" + latencyTargets.Hostname + " && cat latency_" + latencyTargets.Hostname
// 			cmd := "echo " + timestamp[len(timestamp)-1:]
// 			funcStart := time.Now()
// 			elapsed := int64(0)
// 			result, ok := sshServer(target, cmd, timestamp[len(timestamp)-1:])
// 			if ok {
// 				elapsed = time.Since(funcStart).Milliseconds()
// 				fmt.Println(result)
// 			} else {
// 				fmt.Println(result)
// 				elapsed = int64(-1)
// 			}
// 			latencyResult := internal.LatencyResult{Hostname: target.Hostname, Latency: elapsed}
// 			fmt.Println(latencyResult.String())
// 			c1 <- latencyResult
// 		}(target)
// 		latencyResults.Results = append(latencyResults.Results, <-c1)
// 	}

// 	// Timestamp after operations
// 	tmpTime := time.Now()
// 	latencyResults.Timestamp = internal.LatencyTimestamp{
// 		TimeLocal:   tmpTime,
// 		TimeSeconds: tmpTime.Unix(),
// 		TimeNano:    tmpTime.UnixNano(),
// 	}
// }

// func main() {
// 	defer RecoverEndpoint()
// 	target := internal.LatencyTarget{
// 		Hostname:     "172.26.161.241",
// 		Hostport:     "21",
// 		HostUser:     "tiers",
// 		HostPassword: "2022",
// 	}
// 	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
// 	expected := timestamp[len(timestamp)-1:]
// 	cmd := "echo " + timestamp[len(timestamp)-1:]

// 	fmt.Println("sshServer")
// 	fmt.Printf("%+v", target)
// 	config := &ssh.ClientConfig{
// 		User: target.HostUser,
// 		Auth: []ssh.AuthMethod{
// 			ssh.Password(target.HostPassword),
// 		},
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// 	}

// 	addr := fmt.Sprintf("%s:%s", target.Hostname, target.Hostport)
// 	fmt.Println(addr)
// 	fmt.Println("-----")
// 	fmt.Printf("%+v", config.Auth)
// 	fmt.Println("-----")
// 	for _, sm := range config.Auth {
// 		fmt.Println("-----")
// 		fmt.Printf("%+v", sm)
// 	}

// 	client, err := ssh.Dial("tcp", addr, config)
// 	if err != nil {
// 		fmt.Println("DIAL")
// 		fmt.Println(err.Error())
// 	}
// 	defer client.Close()

// 	session, err := client.NewSession()
// 	if err != nil {
// 		fmt.Println("SESH")
// 		fmt.Println(err.Error())
// 	}
// 	defer session.Close()

// 	stdout, err := session.StdoutPipe()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	err = session.Start(cmd)
// 	if err != nil {
// 		fmt.Printf("unable to execute remote command: %s", err)
// 	}

// 	var buf bytes.Buffer
// 	if _, err := io.Copy(&buf, stdout); err != nil {
// 		fmt.Printf("reading failed: %s", err)
// 	}

// 	if sttyOutput := buf.String(); !strings.Contains(sttyOutput, expected) {
// 		fmt.Printf("FALSE RESULT, expected %s and got %s", expected, sttyOutput)
// 	} else {
// 		fmt.Println(buf.String())
// 	}

// }
