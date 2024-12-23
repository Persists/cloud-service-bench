package archive

import (
	"bufio"
	"cloud-service-bench/internal/config"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

// ArchiveClient is a struct that contains the configuration for the archive client
type ArchiveClient struct {
	config    *config.Config
	writer    *bufio.Writer
	writeChan chan string
}

type Archiver interface {
	Write(line string)
}

// systemsBlockSize returns the block size of the system
// This is used to optimize the buffer size of the writer
func systemsBlockSize() (int, error) {
	os := runtime.GOOS

	fmt.Println("OS:", os)
	var cmd *exec.Cmd
	if os == "linux" {
		cmd = exec.Command("stat", "-fc", "%s", "/")
	} else if os == "darwin" {
		cmd = exec.Command("stat", "-f", "%k", "/")
	} else {
		fmt.Println("Unsupported OS, using default buffer size")
		return 4096, nil
	}

	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting block size:", err)
		return 0, err
	}
	var blockSize int
	fmt.Sscanf(string(out), "%d", &blockSize)
	return blockSize, nil

}

// NewArchiveClient creates a new ArchiveClient
func NewFileArchiveClient(config *config.Config, outDir string, isGenerator bool) *ArchiveClient {
	var filePath string
	if isGenerator {
		filePath = outDir + fmt.Sprintf("/%s-%s-lps%s.log", config.Generator.Name, config.Experiment.Id, strconv.Itoa(config.Generator.LogsPerSecond))
	} else {
		//Todo: Implement name for the http client
		filePath = outDir + fmt.Sprintf("/%s-%s-lps%s.log", "HTTP", config.Experiment.Id, strconv.Itoa(config.Generator.LogsPerSecond))
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil
	}

	bs, err := systemsBlockSize()
	if err != nil {
		fmt.Println("Error getting block size:", err)
		return nil
	}

	writer := bufio.NewWriterSize(file, bs)

	return &ArchiveClient{
		config:    config,
		writer:    writer,
		writeChan: make(chan string),
	}
}

func (ac *ArchiveClient) Start() {
	go ac.writeToFile()
}

func (ac *ArchiveClient) writeToFile() {
	for {
		select {
		case line := <-ac.writeChan:
			ac.writer.WriteString(line)
		}
	}
}

func (ac *ArchiveClient) Write(line string) {
	// ac.writeChan <- line
}
