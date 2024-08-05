package utils

import (
	"log"
	"os"
	"os/exec"
)

const (
	fotaImage         = "/tmp/fota.img"
	ubootEnvBinary    = "/tmp/fota/u-boot_env.bin"
	ubootEnvPartition = "/dev/mmcblk0p4"
	ubootEnvCount     = 2047
	tmpFotaFolder     = "/tmp/fota"
)

func FotaExtractFile() {
	if _, err := os.Stat(fotaImage); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", fotaImage)
		return
	}
	log.Printf("File %s exists.\n", fotaImage)

	// Create the /tmp/fota folder if it doesn't exist
	if _, err := os.Stat(tmpFotaFolder); os.IsNotExist(err) {
		err := os.Mkdir(tmpFotaFolder, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", tmpFotaFolder, err)
		}
	}

	// Extract the fota.img to /tmp/fota folder
	cmd := exec.Command("tar", "--strip-components=1", "-xf", fotaImage, "-C", tmpFotaFolder)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error extracting %s: %v", fotaImage, err)
	}
	log.Println("Extraction completed successfully")
}

func FotaUbootEnv() {
	if _, err := os.Stat(ubootEnvBinary); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", ubootEnvBinary)
	} else {
		log.Printf("File %s exists.\n", ubootEnvBinary)
	}

	//cmd1 := exec.Command("dd", "if=/dev/zero", "of="+ubootEnvPartition, "bs=512", "count="+strconv.Itoa(ubootEnvCount))
	cmd1 := exec.Command("echo", "erase ubootEnv")
	err := cmd1.Run()
	if err != nil {
		log.Printf("Error running first dd command: %v", err)
	}
	log.Println("First dd command executed successfully")

	// Second dd command
	//cmd2 := exec.Command("dd", "if="+ubootEnvBinary, "of="+ubootEnvPartition)
	cmd2 := exec.Command("echo", "flash ubootEnv")
	err = cmd2.Run()
	if err != nil {
		log.Printf("Error running second dd command: %v", err)
	}
	log.Println("Second dd command executed successfully")
}
