package archive

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
)

// ArchiveClient is a struct that contains the configuration for the archive client
type ArchiveClient struct {
	debugMode bool
	writer    *bufio.Writer
	writeChan chan string
}

type Archiver interface {
	Write(line string)
	Close()
	Start()
	Flush()
}

// systemsBlockSize returns the block size of the system
// This is used to optimize the buffer size of the writer
func systemsBlockSize() (int, error) {
	os := runtime.GOOS

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
func NewFileArchiveClient(filePath string, metadata string) (*ArchiveClient, error) {

	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	// Get the block size of the system, and use it to optimize the buffer size
	bs, err := systemsBlockSize()
	if err != nil {
		return nil, fmt.Errorf("failed to get block size: %w", err)
	}

	writer := bufio.NewWriterSize(file, bs)

	ac := &ArchiveClient{
		writer:    writer,
		writeChan: make(chan string),
	}

	_, err = ac.writer.WriteString(metadata + "\n")
	if err != nil {
		return nil, fmt.Errorf("failed to add metadata to file: %w", err)
	}
	err = ac.writer.Flush()
	if err != nil {
		return nil, fmt.Errorf("failed to flush metadata to file: %w", err)
	}

	return ac, nil

}

// Start starts the archive client
//
// Goroutines:
//   - Debug mode: If enabled, starts a goroutine that continuously reads from the
//     write channel, and discards the data.
//   - File writing: Starts a goroutine to handle writing data to a file.
func (ac *ArchiveClient) Start() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		ac.Close()
		os.Exit(0)
	}()

	if ac.debugMode {
		go func() {
			for range ac.writeChan {
			}
		}()
	}
	go ac.writeToFile()
}

// writeToFile writes the data from the write channel to the file using the writer
// from the bufio package.
func (ac *ArchiveClient) writeToFile() {
	for line := range ac.writeChan {
		ac.writer.WriteString(line)
		ac.writer.Flush()
	}
}

func (ac *ArchiveClient) Write(line string) {
	ac.writeChan <- line + "\n"
}

func (ac *ArchiveClient) Close() {
	ac.writer.Flush()
	close(ac.writeChan)
}

func (ac *ArchiveClient) Flush() {
	ac.writer.Flush()
}
