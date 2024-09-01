package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	fmt.Println("\n>>> System and Docker info")
	fmt.Println("[1] Show overall report (system resources, Docker containers, and logs)")
	fmt.Println("[2] Show containers live statistics")
	fmt.Println("[3] Show real-time containers statistics")
	fmt.Println("[4] Show all Docker log file sizes and total")
	fmt.Println("[5] Delete all Docker log files (Warning)")
	fmt.Println("[6] View Swarm cluster info")
	fmt.Println("[x] Exit")
	footer()
	fmt.Print("Enter 1-6 (or x for exit):")

	var ch string
	fmt.Scanln(&ch)
	fmt.Println("Your choice is [" + ch + "] please wait ... \n")

	switch ch {
	case "1":
		showOverallReport()
		footer()
	case "2":
		showLiveStats()
		footer()
	case "3":
		showRealTimeStats()
		footer()
	case "4":
		showDockerLogs()
		footer()
	case "5":
		if confirm() {
			removeDockerLogs()
		} else {
			fmt.Println("Operation Cancelled")
		}
	case "6":
		showSwarmInfo()
		footer()
	case "x":
		fmt.Println("Exited by user \n")
		os.Exit(0)
	default:
		fmt.Println("Try again. (You should enter a number from 1 to 6)")
	}
}

func runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout  // This sends the command's output directly to the terminal, preserving line breaks
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func showOverallReport() {
	clearScreen()
	runCommand("hostnamectl")
	fmt.Println()
	runCommand("hostname", "-i")
	fmt.Println("\nPublic IP:")
	runCommand("curl", "ifconfig.io")
	fmt.Println("\nFirewall status:")
	runCommand("sudo", "ufw", "status")
	fmt.Println("\nRoute:")
	runCommand("ip", "r")
	fmt.Println("\n>>> Docker Total:")
	runCommand("sudo", "docker", "system", "df", "--format", "table {{.Type}}\t{{.TotalCount}}\t{{.Active}}\t{{.Size}}")
	fmt.Println("\n>>> Containers list:")
	runCommand("sudo", "docker", "ps")
	showCPUUsage()
	fmt.Println("\n>>> Storage:")
	runCommand("df", "-h")
	fmt.Println("\n>>> RAM:")
	runCommand("free", "-h")
	fmt.Println("\n>>> Log file size:")
	showDockerLogSize()
        showRealTimeStats()
}

func showLiveStats() {
	clearScreen()
	fmt.Println(">>> Docker Total:")
	runCommand("sudo", "docker", "system", "df", "--format", "table {{.Type}}\t{{.TotalCount}}\t{{.Active}}\t{{.Size}}")
	fmt.Println("\n>>> Containers statistics:")
	runCommand("sudo", "docker", "stats", "--no-stream")
}

func showRealTimeStats() {
	for i := 1; i <= 3; i++ {
		fmt.Printf("#%d Docker containers resource results:\n", i)
		runCommand("sudo", "docker", "stats", "--no-stream")
		fmt.Println()
                time.Sleep(1 * time.Second)
	}
}

func showDockerLogs() {
	clearScreen()
	fmt.Println(">>> Docker log files detail:")
	runCommand("bash", "-c", "sudo du -h $(sudo docker inspect --format='{{.LogPath}}' $(sudo docker ps -qa))")
	runCommand("bash", "-c", "sudo du -ch $(sudo docker inspect --format='{{.LogPath}}' $(sudo docker ps -qa)) | tail -n1")
	fmt.Println("\nDisk Usage:")
	runCommand("df", "-h")
}

func removeDockerLogs() {
	clearScreen()
	fmt.Println("Before:")
	runCommand("bash", "-c", "sudo du -h $(sudo docker inspect --format='{{.LogPath}}' $(sudo docker ps -qa))")
	runCommand("bash", "-c", "sudo du -ch $(sudo docker inspect --format='{{.LogPath}}' $(sudo docker ps -qa)) | tail -n1")
	fmt.Println("\nDisk Usage:")
	runCommand("df", "-h")
	runCommand("sudo", "sh", "-c", "truncate -s 0 /var/lib/docker/containers/*/*-json.log")
	fmt.Println("\n<< All Docker log files successfully deleted >>")
	fmt.Println("After:")
	runCommand("bash", "-c", "sudo du -h $(sudo docker inspect --format='{{.LogPath}}' $(sudo docker ps -qa))")
	runCommand("bash", "-c", "sudo du -ch $(sudo docker inspect --format='{{.LogPath}}' $(sudo docker ps -qa)) | tail -n1")
	fmt.Println("\nDisk Usage:")
	runCommand("sudo", "df", "-h")
}

func showSwarmInfo() {
	clearScreen()
	fmt.Println(">>> Docker & Swarm info:")
	runCommand("sudo", "docker", "info")

	if getSwarmState() == "inactive" {
		fmt.Println("\n>>> Docker swarm mode is disabled, this node is not a part of the swarm cluster.\n")
		os.Exit(0)
	} else {
		fmt.Println("\n>>> Docker swarm mode is enabled, this node is a part of the swarm cluster.\n")
	}
	fmt.Println("\n>>> Swarm cluster nodes (servers):")
	runCommand("sudo", "docker", "node", "ls")
	showSwarmNodesDetail()
}

func getSwarmState() string {
	cmd := exec.Command("sudo", "docker", "info", "--format", "{{.Swarm.LocalNodeState}}")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return "inactive"
	}
	return strings.TrimSpace(string(out))
}

func showSwarmNodesDetail() {
	fmt.Println("\n>>> Swarm nodes detail:")
	runCommand("bash", "-c", "sudo docker node ls -q | xargs sudo docker node inspect -f '{{ .Description.Hostname }} ({{ .ID }}) {{ .Status.Addr }} State={{ .Status.State }} {{ .Spec.Availability }} Role={{ .Spec.Role }} {{ .Description.Platform.Architecture }} OS={{ .Description.Platform.OS }} RAM(B)={{ .Description.Resources.MemoryBytes }} docker_ver={{ .Description.Engine.EngineVersion }} , labels={{ range $k, $v := .Spec.Labels }}{{ $k }}={{ $v }} {{end}}'")
}

func showCPUUsage() {
	fmt.Println("\n>>> CPU usage:")
	runCommand("bash", "-c", "top -bn2 | grep '%Cpu' | tail -1 | grep -P '(....|...) id,'|awk '{print \"CPU Usage: \" 100-$8}'")
}

func showDockerLogSize() {
	runCommand("bash", "-c", "sudo du -h $(sudo docker inspect --format='{{.LogPath}}' $(sudo docker ps -qa))")
	fmt.Println(" -------------")
	runCommand("bash", "-c", "sudo du -ch $(sudo docker inspect --format='{{.LogPath}}' $(sudo docker ps -qa)) | tail -n1")
	fmt.Println("----------------------------------------------------------------------------------------------------")
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func confirm() bool {
	var input string
	fmt.Printf("Do you want to continue with this operation? [y|n]: ")
	_, err := fmt.Scanln(&input)
	if err != nil {
		os.Exit(0)
	}
	return strings.ToLower(input) == "y" || strings.ToLower(input) == "yes"
}

func footer() {
	fmt.Println("App by Elias Naserkhaki. enjoy it ;) \n\n")
}
