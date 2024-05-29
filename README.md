# TTKAgent

TTKAgent is a lightweight agent designed to facilitate communication between an MQTT broker and a local system. It subscribes to MQTT topics to receive commands, executes them on the local system, and publishes the output back to the MQTT broker. It provides flexibility in configuration and ensures basic security through authorization.

## Features

- **Command Execution**: Listens to MQTT commands and executes them on the local system.
- **Authorization**: Requires a key for command authorization, ensuring secure command execution.
- **Output Publication**: Publishes command output (result or error) back to the MQTT broker.
- **Flexible Configuration**: Can be configured via command line flags or a JSON configuration file.
- **Daemon Mode**: Optionally runs as a daemon for background operation.

## Requirements

- Go (Golang) installed on your system.
- Access to an MQTT broker.

## Installation

1. Clone the TTKAgent repository to your local machine:

    ```bash
    git clone https://github.com/TTKLabs/TTKAgent.git
    ```

2. Navigate to the project directory:

    ```bash
    cd TTKAgent
    ```

3. Ensure you have Go modules enabled:

    ```bash
    export GO111MODULE=on
    ```

4. Use the provided builder script to compile the code for multiple platforms. Ensure the script is executable:

    ```bash
    chmod +x build.sh
    ```
5. Install the required packages

```bash
go get github.com/eclipse/paho.mqtt.golang
go get github.com/sevlyar/go-daemon
```

6. Run the builder script to compile the code:

    ```bash
    ./build.sh
    ```

7. Once the compilation process is complete, you'll find the compiled binaries in the `build` directory.

## Usage

You can configure TTKAgent using either command line flags or a JSON configuration file.

### Command Line Flags

Use the following command line flags to configure TTKAgent:

- `-u, --username`: MQTT broker username.
- `-p, --password`: MQTT broker password.
- `-b, --mqttbroker`: MQTT broker address.
- `-m, --maintopic`: Main MQTT topic.
- `-k, --key`: Key for command authorization.
- `-d, --daemon`: Run as daemon (true/false).

Example:

```bash
./TTKAgent -u <username> -p <password> -b <mqtt_broker_address> -m <main_mqtt_topic> -k <key> -d

```
### JSON Configuration File

Alternatively, you can provide a JSON configuration file with the following fields:

- **username**: MQTT broker username.
- **password**: MQTT broker password.
- **mqttbroker**: MQTT broker address.
- **maintopic**: Main MQTT topic.
- **key**: Key for command authorization.
- **daemon**: Run as daemon (true/false).

Example JSON configuration file (`config.json`):

```json
{
  "username": "<username>",
  "password": "<password>",
  "mqttbroker": "<mqtt_broker_address>",
  "maintopic": "<main_mqtt_topic>",
  "key": "<key>",
  "daemon": true
}

```

To use the JSON configuration file, pass it as a command line argument:

```bash
./TTKAgent -c config.json

```

# Disclaimer

Use the TTKAgent with care, especially when interacting with sensitive systems. Ensure that the MQTT broker is properly secured to prevent unauthorized access to execute commands.

# Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

# License

This project is licensed under the BSD License.
