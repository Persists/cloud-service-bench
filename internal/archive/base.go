package archive

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
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
func NewFileArchiveClient(filePath string, metadata string) (*ArchiveClient, error) {

	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	bs, err := systemsBlockSize()
	if err != nil {
		return nil, fmt.Errorf("failed to get block size: %w", err)
	}

	writer := bufio.NewWriterSize(file, bs)

	ac := &ArchiveClient{
		writer:    writer,
		writeChan: make(chan string),
	}

	_, err = ac.writer.WriteString(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to add metadata to file: %w", err)
	}
	err = ac.writer.Flush()
	if err != nil {
		return nil, fmt.Errorf("failed to flush metadata to file: %w", err)
	}

	return ac, nil

}

func (ac *ArchiveClient) Start() {
	if ac.debugMode {
		go func() {
			for line := range ac.writeChan {
				fmt.Println(line)
			}
		}()
	}
	go ac.writeToFile()
}

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
	close(ac.writeChan)
}

func (ac *ArchiveClient) Flush() {
	ac.writer.Flush()
}
