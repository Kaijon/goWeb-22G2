package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	//dtb
	dtbBinary    = "/tmp/fota/leipzig.dtb"
	dtbPartition = "/dev/mmcblk0p5"
	dtbCount     = 2047
	//Kernel
	kernelBinary    = "/tmp/fota/Image"
	kernelPartition = "/dev/mmcblk0p6"
	kernelCount     = 40959
	//rootfs
	rootfsBinary    = "/tmp/fota/rootfs.ext2"
	rootfsPartition = "/dev/mmcblk0p7"
	rootfsCount     = 262134
	//daemon
	daemonBinary = "/tmp/fota/daemon.tar"
	daemonPath   = "/mnt/getac"
	//flash
	flashBinary = "/tmp/fota/flash.tar"
	flashPath   = "/mnt/flash"

	tmpFotaFolder = "/tmp/fota"
)

func PrehookInstall() error {
	cmd := exec.Command("./web/sh/prehookInstall.sh")
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("Error running PrehookInstall: %v", err)
		return fmt.Errorf("Error running PrehookInstall: %v", err)
	}
	return nil
}

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
	return nil
}

func FotaUboot() error {
	if _, err := os.Stat(ubootBinary); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", ubootBinary)
		return fmt.Errorf("File not found")
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
	return nil
}

func FotaImagePreHook() {
	if err := os.MkdirAll("/tmp/partition5", 0755); err != nil {
		log.Printf("Error creating directory: %v", err)
	}

	cmd := exec.Command("mount", "/dev/mmcblk0p5", "/tmp/partition5")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error mounting /dev/mmcblk0p5: %v", err)
	}
}

func FotaImagePostHook() {
	cmd := exec.Command("umount", "/tmp/partition5")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error unmounting /tmp/partition5: %v", err)
	}
}

func FotaKernel() error {
	if _, err := os.Stat(kernelBinary); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", kernelBinary)
		return fmt.Errorf("File not found")
	} else {
		log.Printf("File %s exists.\n", kernelBinary)
	}

	cmd1 := exec.Command("dd", "if=/dev/zero", "of="+kernelPartition, "bs=512", "count="+strconv.Itoa(kernelCount))
	//cmd1 := exec.Command("echo", "erase kernel")
	cmd1.Env = os.Environ()
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	err := cmd1.Run()
	if err != nil {
		log.Printf("Error running erase command: %v", err)
		return fmt.Errorf("error running erase command: %w", err)
	}

	cmd2 := exec.Command("dd", "if="+kernelBinary, "of="+kernelPartition)
	//cmd2 := exec.Command("echo", "flash kernel")
	cmd2.Env = os.Environ()
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	err = cmd2.Run()
	if err != nil {
		log.Printf("Error running erase command: %v", err)
		return fmt.Errorf("error running second command: %w", err)
	}

	/*
		// Remove /tmp/partition5/Image
		filesToRemove := []string{"/tmp/partition5/Image"}
		for _, file := range filesToRemove {
			if err := os.Remove(file); err != nil {
				log.Printf("Error removing %s: %v", file, err)
			} else {
				log.Printf("Removed %s successfully", file)
			}
		}

		cmd := exec.Command("cp", "/tmp/fota/Image", "/tmp/partition5/")
		//cmd := exec.Command("echo", "flash Kernel")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Error copying /tmp/Image: %v", err)
			return fmt.Errorf("Error copying /tmp/fota/Image: %v", err)
		}
	*/

	return nil
}

func FotaDtb() error {
	if _, err := os.Stat(dtbBinary); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", dtbBinary)
		return fmt.Errorf("File not found")
	} else {
		log.Printf("File %s exists.\n", dtbBinary)
	}

	cmd1 := exec.Command("dd", "if=/dev/zero", "of="+dtbPartition, "bs=512", "count="+strconv.Itoa(dtbCount))
	//cmd1 := exec.Command("echo", "erase ubootenv")
	cmd1.Env = os.Environ()
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	err := cmd1.Run()
	if err != nil {
		log.Printf("Error running erase command: %v", err)
		return fmt.Errorf("error running erase command: %w", err)
	}

	cmd2 := exec.Command("dd", "if="+dtbBinary, "of="+dtbPartition)
	//cmd2 := exec.Command("echo", "flash ubootenv")
	cmd2.Env = os.Environ()
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	err = cmd2.Run()
	if err != nil {
		log.Printf("Error running flash command: %v", err)
		return fmt.Errorf("error running flash command: %w", err)
	}

	/*
		// Remove /tmp/partition5/Image
		filesToRemove := []string{"/tmp/partition5/leipzig.dtb"}
		for _, file := range filesToRemove {
			if err := os.Remove(file); err != nil {
				log.Printf("Error removing %s: %v", file, err)
			} else {
				log.Printf("Removed %s successfully", file)
			}
		}

		cmd := exec.Command("cp", "/tmp/fota/leipzig.dtb", "/tmp/partition5/")
		//cmd := exec.Command("echo", "flash dtb")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Error copying /tmp/leipzig.dtb: %v", err)
			return fmt.Errorf("Error copying /tmp/fota/leipzig.dtb: %v", err)
		}
	*/
	return nil
}

func FotaRootFs() error {
	if _, err := os.Stat(rootfsBinary); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", rootfsBinary)
		return fmt.Errorf("File not found")
	} else {
		log.Printf("File %s exists.\n", rootfsBinary)
	}

	//cmd1 := exec.Command("dd", "if=/dev/zero", "of="+rootfsPartition, "bs=512", "count="+strconv.Itoa(rootfsCount))
	//cmd1 := exec.Command("echo", "erase rootfs")
	//cmd1.Env = os.Environ()
	//cmd1.Stdout = os.Stdout
	//cmd1.Stderr = os.Stderr
	//err := cmd1.Run()
	//if err != nil {
	//	log.Printf("Error running erase command: %v", err)
	//	return fmt.Errorf("error running erase command: %w", err)
	//}

	cmd2 := exec.Command("dd", "if="+rootfsBinary, "of="+rootfsPartition)
	//cmd2 := exec.Command("echo", "flash rootfs")
	cmd2.Env = os.Environ()
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	err := cmd2.Run()
	if err != nil {
		log.Printf("Error running flase command: %v", err)
		return fmt.Errorf("error running flash command: %w", err)
	}
	return nil
}

func FotaDaemon() error {
	if _, err := os.Stat(daemonBinary); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", daemonBinary)
		return fmt.Errorf("File not found")
	} else {
		log.Printf("File %s exists.\n", daemonBinary)
	}

	extractDirPath := "/tmp/tmp_daemon"
	localPath := "/mnt/getac"

	err := os.MkdirAll(extractDirPath, 0755)
	if err != nil {
		log.Printf("Error creating directory: %v", err)
	} else {
		log.Printf("Directory created successfully: %s", extractDirPath)
	}

	if _, err := os.Stat(extractDirPath); os.IsNotExist(err) {
		log.Printf("Directory does not exist: %s", extractDirPath)
	} else {
		log.Printf("Directory exists: %s", extractDirPath)
	}

	//Get Folder list
	cmd := exec.Command("tar", "--strip-components=1", "-xf", daemonBinary, "-C", extractDirPath)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("Error running tar command: %v", err)
	}

	//Get the list of files/folders extracted
	contents, err := os.ReadDir(extractDirPath)
	if err != nil {
		log.Printf("Error reading directory contents: %v", err)
	}
	log.Printf("%v\n", contents)

	//Remove files or directories based on their type
	for _, item := range contents {
		itemPath := filepath.Join(localPath, item.Name())
		if item.IsDir() {
			// If it's a directory, remove it
			err = os.RemoveAll(itemPath)
			if err != nil {
				log.Printf("Error removing directory: %v", err)
			} else {
				log.Printf("Removed directory: %s", itemPath)
			}
		} else {
			// If it's a file, remove it
			err = os.Remove(itemPath)
			if err != nil {
				log.Printf("Error removing file: %v", err)
			} else {
				log.Printf("Removed file: %s", itemPath)
			}
		}
	}

	cmd = exec.Command("tar", "--strip-components=1", "-xf", daemonBinary, "-C", daemonPath)
	err = cmd.Run()
	if err != nil {
		log.Printf("Error : %v", err)
		return fmt.Errorf("error running tar: %w", err)
	}

	cmd = exec.Command("./web/sh/cleanup.sh", "getac")
	err = cmd.Run()
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err != nil {
		log.Printf("Error running cleanup.sh: %v", err)
	}

	cmd = exec.Command("sync")
	err = cmd.Run()
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err != nil {
		log.Printf("sync: %v", err)
	}
	return nil
}

func FotaFlash() error {
	if _, err := os.Stat(flashBinary); os.IsNotExist(err) {
		log.Printf("File %s does not exist.\n", flashBinary)
		return fmt.Errorf("File not found")
	} else {
		log.Printf("File %s exists.\n", flashBinary)
	}

	extractDirPath := "/tmp/tmp_flash"
	localPath := "/mnt/flash"

	err := os.MkdirAll(extractDirPath, 0755)
	if err != nil {
		log.Printf("Error creating directory: %v", err)
	} else {
		log.Printf("Directory created successfully: %s", extractDirPath)
	}

	if _, err := os.Stat(extractDirPath); os.IsNotExist(err) {
		log.Printf("Directory does not exist: %s", extractDirPath)
	} else {
		log.Printf("Directory exists: %s", extractDirPath)
	}

	//Get Folder list
	cmd := exec.Command("tar", "--strip-components=1", "-xf", flashBinary, "-C", extractDirPath)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("Error running tar command: %v", err)
	}

	//Get the list of files/folders extracted
	contents, err := os.ReadDir(extractDirPath)
	if err != nil {
		log.Printf("Error reading directory contents: %v", err)
	}
	log.Printf("%v\n", contents)

	//Remove files or directories based on their type
	for _, item := range contents {
		itemPath := filepath.Join(localPath, item.Name())
		if item.IsDir() {
			// If it's a directory, remove it
			err = os.RemoveAll(itemPath)
			if err != nil {
				log.Printf("Error removing directory: %v", err)
			} else {
				log.Printf("Removed directory: %s", itemPath)
			}
		} else {
			// If it's a file, remove it
			err = os.Remove(itemPath)
			if err != nil {
				log.Printf("Error removing file: %v", err)
			} else {
				log.Printf("Removed file: %s", itemPath)
			}
		}
	}

	cmd = exec.Command("tar", "--strip-components=1", "-xf", flashBinary, "-C", flashPath)
	err = cmd.Run()
	if err != nil {
		log.Printf("Error running extract command: %v", err)
		return fmt.Errorf("Error running extract command: %w", err)
	}

	cmd = exec.Command("./web/sh/cleanup.sh", "flash")
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("Error running cleanup.sh: %v", err)
	}

	cmd = exec.Command("sync")
	err = cmd.Run()
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err != nil {
		log.Printf("sync: %v", err)
	}
	return nil
}
