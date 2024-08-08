package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	fotaImage = "/tmp/fota.img"
	//u-boot settings
	ubootBinary    = "/tmp/fota/u-boot.bin"
	ubootPartition = "/dev/mmcblk0p3"
	ubootCount     = 32767
	//env settings
	ubootEnvBinary    = "/tmp/fota/u-boot_env.bin"
	ubootEnvPartition = "/dev/mmcblk0p4"
	ubootEnvCount     = 2047
	//Kernel & dtb
	kernelBinary    = "/tmp/fota/Image"
	dtbBinary       = "/tmp/fota/leipzig.dtb"
	kernelPartition = "/dev/mmcblk0p5"
	//rootfs
	rootfsBinary    = "/tmp/fota/rootfs.ext2"
	rootfsPartition = "/dev/mmcblk0p6"
	rootfsCount     = 262134
	//storage
	storageBinary = "/tmp/fota/storage"

	tmpFotaFolder = "/tmp/fota"
)

func FotaExtractFile() error {
	if _, err := os.Stat(fotaImage); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", fotaImage)
		return fmt.Errorf("file %s does not exist", fotaImage)
	}
	log.Printf("File %s exists.\n", fotaImage)

	// Create the /tmp/fota folder if it doesn't exist
	if _, err := os.Stat(tmpFotaFolder); os.IsNotExist(err) {
		err := os.Mkdir(tmpFotaFolder, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", tmpFotaFolder, err)
			return fmt.Errorf("failed to create directory %s: %v", tmpFotaFolder, err)
		}
	}

	// Extract the fota.img to /tmp/fota folder
	cmd := exec.Command("tar", "--strip-components=1", "-xf", fotaImage, "-C", tmpFotaFolder)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error extracting %s: %v", fotaImage, err)
		return fmt.Errorf("error extracting %s: %v", fotaImage, err)
	}
	log.Println("Extraction completed successfully")
	return nil
}

func FotaUboot() error {
	if _, err := os.Stat(ubootBinary); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", ubootBinary)
		return fmt.Errorf("file not found")
	} else {
		log.Printf("File %s exists.\n", ubootBinary)
	}

	//cmd1 := exec.Command("dd", "if=/dev/zero", "of="+ubootEnvPartition, "bs=512", "count="+strconv.Itoa(ubootEnvCount))
	cmd1 := exec.Command("echo", "erase uboot")
	err := cmd1.Run()
	if err != nil {
		log.Printf("Error running first dd command: %v", err)
		return fmt.Errorf("error running first command: %w", err)
	}
	log.Println("First dd command executed successfully")

	// Second dd command
	//cmd2 := exec.Command("dd", "if="+ubootEnvBinary, "of="+ubootEnvPartition)
	cmd2 := exec.Command("echo", "flash uboot")
	err = cmd2.Run()
	if err != nil {
		log.Printf("Error running second dd command: %v", err)
		return fmt.Errorf("error running second command: %w", err)
	}
	log.Println("Second dd command executed successfully")
	return nil
}

func FotaUbootEnv() error {
	if _, err := os.Stat(ubootEnvBinary); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", ubootEnvBinary)
		return fmt.Errorf("file not found")
	} else {
		log.Printf("File %s exists.\n", ubootEnvBinary)
	}

	//cmd1 := exec.Command("dd", "if=/dev/zero", "of="+ubootEnvPartition, "bs=512", "count="+strconv.Itoa(ubootEnvCount))
	cmd1 := exec.Command("echo", "erase ubootenv")
	err := cmd1.Run()
	if err != nil {
		log.Printf("Error running first dd command: %v", err)
		return fmt.Errorf("error running first command: %w", err)
	}
	log.Println("First dd command executed successfully")

	// Second dd command
	//cmd2 := exec.Command("dd", "if="+ubootEnvBinary, "of="+ubootEnvPartition)
	cmd2 := exec.Command("echo", "flash ubootenv")
	err = cmd2.Run()
	if err != nil {
		log.Printf("Error running second dd command: %v", err)
		return fmt.Errorf("error running second command: %w", err)
	}
	log.Println("Second dd command executed successfully")
	return nil
}

func FotaKernel() error {
	if _, err := os.Stat(ubootEnvBinary); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", ubootEnvBinary)
		return fmt.Errorf("file not found")
	} else {
		log.Printf("File %s exists.\n", ubootEnvBinary)
	}

	//cmd1 := exec.Command("dd", "if=/dev/zero", "of="+ubootEnvPartition, "bs=512", "count="+strconv.Itoa(ubootEnvCount))
	cmd1 := exec.Command("echo", "erase ubootenv")
	err := cmd1.Run()
	if err != nil {
		log.Printf("Error running first dd command: %v", err)
		return fmt.Errorf("error running first command: %w", err)
	}
	log.Println("First dd command executed successfully")

	// Second dd command
	//cmd2 := exec.Command("dd", "if="+ubootEnvBinary, "of="+ubootEnvPartition)
	cmd2 := exec.Command("echo", "flash ubootenv")
	err = cmd2.Run()
	if err != nil {
		log.Printf("Error running second dd command: %v", err)
		return fmt.Errorf("error running second command: %w", err)
	}
	log.Println("Second dd command executed successfully")
	return nil
}

func FotaDtb() error {
	if _, err := os.Stat(rootfsBinary); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", rootfsBinary)
		return fmt.Errorf("file not found")
	} else {
		log.Printf("File %s exists.\n", rootfsBinary)
	}

	//cmd1 := exec.Command("dd", "if=/dev/zero", "of="+ubootEnvPartition, "bs=512", "count="+strconv.Itoa(ubootEnvCount))
	cmd1 := exec.Command("echo", "erase rootfs")
	err := cmd1.Run()
	if err != nil {
		log.Printf("Error running first dd command: %v", err)
		return fmt.Errorf("error running first command: %w", err)
	}
	log.Println("First dd command executed successfully")

	// Second dd command
	//cmd2 := exec.Command("dd", "if="+ubootEnvBinary, "of="+ubootEnvPartition)
	cmd2 := exec.Command("echo", "flash rootfs")
	err = cmd2.Run()
	if err != nil {
		log.Printf("Error running second dd command: %v", err)
		return fmt.Errorf("error running second command: %w", err)
	}
	log.Println("Second dd command executed successfully")
	return nil
}

func FotaRootFs() error {
	if _, err := os.Stat(rootfsBinary); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", rootfsBinary)
		return fmt.Errorf("file not found")
	} else {
		log.Printf("File %s exists.\n", rootfsBinary)
	}

	//cmd1 := exec.Command("dd", "if=/dev/zero", "of="+ubootEnvPartition, "bs=512", "count="+strconv.Itoa(ubootEnvCount))
	cmd1 := exec.Command("echo", "erase rootfs")
	err := cmd1.Run()
	if err != nil {
		log.Printf("Error running first dd command: %v", err)
		return fmt.Errorf("error running first command: %w", err)
	}
	log.Println("First dd command executed successfully")

	// Second dd command
	//cmd2 := exec.Command("dd", "if="+ubootEnvBinary, "of="+ubootEnvPartition)
	cmd2 := exec.Command("echo", "flash rootfs")
	err = cmd2.Run()
	if err != nil {
		log.Printf("Error running second dd command: %v", err)
		return fmt.Errorf("error running second command: %w", err)
	}
	log.Println("Second dd command executed successfully")
	return nil
}

func FotaStorage() error {
	if _, err := os.Stat(storageBinary); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", storageBinary)
		return fmt.Errorf("File not found")
	} else {
		log.Printf("File %s exists.\n", storageBinary)
	}

	//cmd1 := exec.Command("dd", "if=/dev/zero", "of="+ubootEnvPartition, "bs=512", "count="+strconv.Itoa(ubootEnvCount))
	cmd1 := exec.Command("echo", "erase storageBinary")
	err := cmd1.Run()
	if err != nil {
		log.Printf("Error running first dd command: %v", err)
		return fmt.Errorf("error running first command: %w", err)
	}
	log.Println("First dd command executed successfully")

	// Second dd command
	//cmd2 := exec.Command("dd", "if="+ubootEnvBinary, "of="+ubootEnvPartition)
	cmd2 := exec.Command("echo", "flash storageBinary")
	err = cmd2.Run()
	if err != nil {
		log.Printf("Error running second dd command: %v", err)
		return fmt.Errorf("error running second command: %w", err)
	}
	log.Println("Second dd command executed successfully")
	return nil
}
