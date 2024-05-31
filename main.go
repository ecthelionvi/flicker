package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/pelletier/go-toml"
)

type Config struct {
	Watch struct {
		Directories []string
		Device      string
	}
}

var (
	currentCmd   *exec.Cmd
	currentStdin io.WriteCloser
	wg           sync.WaitGroup
)

func main() {
	generateConfig := flag.Bool("generate-config", false, "Generate a standard flicker.toml file")
	flag.Parse()

	if *generateConfig {
		generateStandardConfig()
		fmt.Println("Standard flicker.toml file generated.")
		return
	}

	config := loadConfig("flicker.toml")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("Event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Modified file:", event.Name)
					hotReloadFlutterApp()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	for _, dir := range config.Watch.Directories {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			log.Printf("Directory %s does not exist, skipping...\n", dir)
			continue
		}
		err = watcher.Add(dir)
		if err != nil {
			log.Fatal(err)
		}
	}

	startFlutterApp(config.Watch.Device)
	<-done
}

func generateStandardConfig() {
	config := Config{
		Watch: struct {
			Directories []string
			Device      string
		}{
			Directories: []string{"lib"},
			Device:      "chrome",
		},
	}

	configData, err := toml.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("flicker.toml", configData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func loadConfig(path string) Config {
	var config Config
	configFile, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	err = toml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func startFlutterApp(device string) {
	cmd := exec.Command("flutter", "run", "-d", device)
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	currentCmd = cmd
	currentStdin = stdin

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal("Error starting flutter command:", err)
	}

	// Wait for the app to start before attaching
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			if strings.Contains(line, "To hot reload") {
				go attachToFlutterApp(device)
			}
			if strings.Contains(line, "Application finished.") {
				log.Println("Application finished. Exiting.")
				cleanup()
				os.Exit(0)
			}
		}
	}()
}

func attachToFlutterApp(device string) {
	cmd := exec.Command("flutter", "attach", "-d", device)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal("Error starting attach command:", err)
	}

	// Reading the output to determine when the app is ready
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// Print stderr
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if err != nil {
				break
			}
			os.Stderr.Write(buf[:n])
		}
	}()

	if err := cmd.Wait(); err != nil {
		log.Println("Error attaching to flutter command:", err)
	}
}

func hotReloadFlutterApp() {
	if currentCmd != nil && currentStdin != nil {
		log.Println("Sending hot reload command")
		if err := sendHotReloadCommand(); err != nil {
			log.Println("Failed to send hot reload command:", err)
		}
	}
}

func sendHotReloadCommand() error {
	if currentStdin == nil {
		return fmt.Errorf("no flutter command running")
	}

	_, err := currentStdin.Write([]byte("r\n"))
	if err != nil {
		return err
	}

	return nil
}

func cleanup() {
	if currentCmd != nil {
		log.Println("Cleaning up...")
		currentCmd.Process.Kill()
		wg.Wait()
	}
}
