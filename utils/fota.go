package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
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

	cmd1 := exec.Command("dd", "if=/dev/zero", "of="+ubootPartition, "bs=512", "count="+strconv.Itoa(ubootCount))
	//cmd1 := exec.Command("echo", "erase uboot")
	cmd1.Env = os.Environ()
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	err := cmd1.Run()
	if err != nil {
		log.Printf("Error running erase command: %v", err)
		return fmt.Errorf("error running erase command: %w", err)
	}
	log.Println("\tErase executed successfully")

	// Second dd command
	cmd2 := exec.Command("dd", "if="+ubootBinary, "of="+ubootPartition)
	//cmd2 := exec.Command("echo", "flash uboot")
	cmd2.Env = os.Environ()
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	err = cmd2.Run()
	if err != nil {
		log.Printf("Error running flash command: %v", err)
		return fmt.Errorf("error running flash command: %w", err)
	}
	log.Println("\tFlash executed successfully")
	return nil
}

func FotaUbootEnv() error {
	if _, err := os.Stat(ubootEnvBinary); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", ubootEnvBinary)
		return fmt.Errorf("File not found")
	} else {
		log.Printf("File %s exists.\n", ubootEnvBinary)
	}

	cmd1 := exec.Command("dd", "if=/dev/zero", "of="+ubootEnvPartition, "bs=512", "count="+strconv.Itoa(ubootEnvCount))
	//cmd1 := exec.Command("echo", "erase ubootenv")
	cmd1.Env = os.Environ()
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	err := cmd1.Run()
	if err != nil {
		log.Printf("Error running erase command: %v", err)
		return fmt.Errorf("error running erase command: %w", err)
	}
	log.Println("\tErase executed successfully")

	// Second dd command
	cmd2 := exec.Command("dd", "if="+ubootEnvBinary, "of="+ubootEnvPartition)
	//cmd2 := exec.Command("echo", "flash ubootenv")
	cmd2.Env = os.Environ()
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	err = cmd2.Run()
	if err != nil {
		log.Printf("Error running second dd command: %v", err)
		return fmt.Errorf("error running second command: %w", err)
	}
	log.Println("\tFlash executed successfully")
	return nil
}

func FotaImageHook() error {
	// Create the directory /tmp/partition5
	if err := os.MkdirAll("/tmp/partition5", 0755); err != nil {
		log.Printf("Error creating directory: %v", err)
	}
	log.Println("Directory /tmp/partition5 created successfully")

	// Mount /dev/mmcblk0p5 to /tmp/partition5
	cmd := exec.Command("mount", "/dev/mmcblk0p5", "/tmp/partition5")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error mounting /dev/mmcblk0p5: %v", err)
		return fmt.Errorf("Error mounting /dev/mmcblk0p5: %w", err)
	}
	log.Println("/dev/mmcblk0p5 mounted to /tmp/partition5")
	return nil
}

func FotaKernel() error {
	if _, err := os.Stat(kernelBinary); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", kernelBinary)
		return fmt.Errorf("File not found")
	} else {
		log.Printf("File %s exists.\n", kernelBinary)
	}

	// Remove /tmp/partition5/Image
	filesToRemove := []string{"/tmp/partition5/Image"}
	for _, file := range filesToRemove {
		if err := os.Remove(file); err != nil {
			log.Printf("Error removing %s: %v", file, err)
		} else {
			log.Printf("Removed %s successfully", file)
		}
	}
	// Copy /tmp/Image to /tmp/partition5/
	cmd := exec.Command("cp", "/tmp/fota/Image", "/tmp/partition5/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error copying /tmp/Image: %v", err)
		return fmt.Errorf("Error copying /tmp/fota/Image: %v", err)
	}
	log.Println("/tmp/fota/Image copied to /tmp/partition5/")

	return nil
}

func FotaDtb() error {
	if _, err := os.Stat(dtbBinary); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", dtbBinary)
		return fmt.Errorf("File not found")
	} else {
		log.Printf("File %s exists.\n", dtbBinary)
	}

	// Remove /tmp/partition5/Image
	filesToRemove := []string{"/tmp/partition5/leipzig.dtb"}
	for _, file := range filesToRemove {
		if err := os.Remove(file); err != nil {
			log.Printf("Error removing %s: %v", file, err)
		} else {
			log.Printf("Removed %s successfully", file)
		}
	}

	// Copy /tmp/leipzig.dtb to /tmp/partition5/
	cmd := exec.Command("cp", "/tmp/fota/leipzig.dtb", "/tmp/partition5/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error copying /tmp/leipzig.dtb: %v", err)
		return fmt.Errorf("Error copying /tmp/fota/leipzig.dtb: %v", err)
	}
	log.Println("/tmp/fota/leipzig.dtb copied to /tmp/partition5/")

	return nil
}

func FotaRootFs() error {
	if _, err := os.Stat(rootfsBinary); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", rootfsBinary)
		return fmt.Errorf("file not found")
	} else {
		log.Printf("File %s exists.\n", rootfsBinary)
	}

	cmd1 := exec.Command("dd", "if=/dev/zero", "of="+rootfsPartition, "bs=512", "count="+strconv.Itoa(rootfsCount))
	cmd1.Env = os.Environ()
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	err := cmd1.Run()
	if err != nil {
		log.Printf("Error running erase command: %v", err)
		return fmt.Errorf("error running erase command: %w", err)
	}
	log.Println("\tErase executed successfully")

	// Second dd command
	cmd2 := exec.Command("dd", "if="+rootfsBinary, "of="+rootfsPartition)
	cmd2.Env = os.Environ()
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	err = cmd2.Run()
	if err != nil {
		log.Printf("Error running flase command: %v", err)
		return fmt.Errorf("error running flash command: %w", err)
	}
	log.Println("\tFlash executed successfully")
	return nil
}

func FotaStorage() error {
	if _, err := os.Stat(storageBinary); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", storageBinary)
		return fmt.Errorf("File not found")
	} else {
		log.Printf("File %s exists.\n", storageBinary)
	}

	//cmd1 := exec.Command("dd", "if=/dev/zero", "of="+storageBinary, "bs=512", "count="+strconv.Itoa(rootfsCount))
	cmd1 := exec.Command("echo", "erase FotaStorage")
	err := cmd1.Run()
	if err != nil {
		log.Printf("Error running first dd command: %v", err)
		return fmt.Errorf("error running first command: %w", err)
	}
	log.Println("First dd command executed successfully")

	// Second dd command
	//cmd2 := exec.Command("dd", "if="+ubootEnvBinary, "of="+ubootEnvPartition)
	cmd2 := exec.Command("echo", "flash FotaStorage")
	err = cmd2.Run()
	if err != nil {
		log.Printf("Error running second dd command: %v", err)
		return fmt.Errorf("error running second command: %w", err)
	}
	log.Println("Second dd command executed successfully")
	return nil
}
