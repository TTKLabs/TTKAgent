package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/sevlyar/go-daemon"
)

// Config struct for JSON configuration
type Config struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	MQTTBroker string `json:"mqttbroker"`
	MainTopic  string `json:"maintopic"`
	Key        string `json:"key"`
	Daemon     bool   `json:"daemon"`
}

var version = "1.0.0"

// PrintUsage prints the usage information
func PrintUsage() {
	fmt.Println("Usage: TTKAgent [options]")
	fmt.Println("Options:")
	fmt.Println("  -u, --username       MQTT broker username")
	fmt.Println("  -p, --password       MQTT broker password")
	fmt.Println("  -b, --mqttbroker     MQTT broker address")
	fmt.Println("  -m, --maintopic      Main MQTT topic")
	fmt.Println("  -k, --key            Key for command authorization")
	fmt.Println("  -d, --daemon         Run as daemon (true/false)")
	fmt.Println("  -c, --config         JSON configuration file")
	fmt.Println("  -h, --help           Show this help message")
	fmt.Printf("\nTTKAgent %s [%s/%s]\n", version, runtime.GOOS, runtime.GOARCH)
}

func main() {
	// Command line flags
	username := flag.String("u", "", "MQTT broker username")
	password := flag.String("p", "", "MQTT broker password")
	mqttBroker := flag.String("b", "", "MQTT broker address")
	mainTopic := flag.String("m", "", "Main MQTT topic")
	key := flag.String("k", "", "Key for command authorization")
	daemonFlag := flag.Bool("d", false, "Run as daemon")
	configFile := flag.String("c", "", "JSON configuration file")
	flag.Parse()

	// Print usage if necessary
	if flag.NFlag() == 0 {
		PrintUsage()
		return
	}

	// Read from JSON configuration file if specified
	var config Config
	if *configFile != "" {
		file, err := os.ReadFile(*configFile)
		if err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}
		if err := json.Unmarshal(file, &config); err != nil {
			log.Fatalf("Error parsing config file: %v", err)
		}
	}

	// Override command line flags with configuration file values
	if config.Username != "" {
		username = &config.Username
	}
	if config.Password != "" {
		password = &config.Password
	}
	if config.MQTTBroker != "" {
		mqttBroker = &config.MQTTBroker
	}
	if config.MainTopic != "" {
		mainTopic = &config.MainTopic
	}
	if config.Key != "" {
		key = &config.Key
	}
	if config.Daemon {
		daemonFlag = &config.Daemon
	}

	// Ensure all required parameters are provided
	if *username == "" || *password == "" || *mqttBroker == "" || *mainTopic == "" || *key == "" {
		PrintUsage()
		return
	}

	// Daemonize if specified
	if *daemonFlag {
		cntxt := &daemon.Context{
			PidFileName: "/var/run/mydaemon.pid",
			PidFilePerm: 0644,
			LogFileName: "/var/log/mydaemon.log",
			LogFilePerm: 0640,
			WorkDir:     "/",
			Umask:       027,
			Args:        []string{"[TTKAgent]", version, fmt.Sprintf("[%s/%s]", runtime.GOOS, runtime.GOARCH)},
		}

		d, err := cntxt.Reborn()
		if err != nil {
			log.Fatalf("Unable to run as daemon: %v", err)
		}
		if d != nil {
			return
		}
		defer cntxt.Release()
	}

	// Display banner
	fmt.Printf("TTKAgent %s [%s/%s]\n", version, runtime.GOOS, runtime.GOARCH)

	// MQTT client options
	opts := MQTT.NewClientOptions().AddBroker(*mqttBroker).SetClientID("TTKAgent")
	opts.SetUsername(*username)
	opts.SetPassword(*password)

	// Create MQTT client
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to MQTT broker: %v", token.Error())
	}
	cmdKey := ""
	client.Subscribe(fmt.Sprintf("%s/key", *mainTopic), 0, func(client MQTT.Client, msg MQTT.Message) {
		cmdKey = string(msg.Payload())
	})
	// Handle command messages
	client.Subscribe(fmt.Sprintf("%s/command", *mainTopic), 0, func(client MQTT.Client, msg MQTT.Message) {
		command := string(msg.Payload())

		// Check if the key matches
		if cmdKey != *key {
			log.Printf("Received invalid key: %s", cmdKey)
			return
		}

		// Log the received command
		log.Printf("Received command: %s", command)

		// Execute command received
		output, err := exec.Command("sh", "-c", command).Output()
		if err != nil {
			output = []byte(err.Error())
		}

		// Publish output to output topic
		client.Publish(fmt.Sprintf("%s/%s", *mainTopic, "output"), 1, false, output)
	})

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	client.Disconnect(250)
}
