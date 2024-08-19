package main

import (
	. "getac/goWeb/utils"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func getacFota(c *gin.Context) {
	log.Println("RunFOTA")
	file, err := c.FormFile("image_file")
	if err != nil {
		return
	}

	// Save the uploaded file
	err = c.SaveUploadedFile(file, "/tmp/fota.img")
	if err != nil {
		return
	}

	go fota()

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})

	log.Println("getacFota done")
	//c.HTML(http.StatusOK, "fota.page.gohtml", nil)
}

func fota() {
	var isPass = true
	err := FotaExtractFile()
	if err != nil {
		MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/0", "{\"status\":\"failed\"}")
		log.Println("Wrong IMAGE")
		return
	} else {
		MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/0", "{\"status\":\"success\"}")
	}

	log.Println("========== Run FotaUboot ==========")
	err = FotaUboot()
	if err != nil {
		if err.Error() == "File not found" {
			MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/1", "{\"status\":\"skip\"}")
		} else {
			MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/1", "{\"status\":\"failed\"}")
			isPass = false
		}
	} else {
		MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/1", "{\"status\":\"success\"}")
	}

	log.Println("========== Run FotaUbootEnv ==========")
	err = FotaUbootEnv()
	if err != nil {
		if err.Error() == "File not found" {
			MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/2", "{\"status\":\"skip\"}")
		} else {
			MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/2", "{\"status\":\"failed\"}")
			isPass = false
		}
	} else {
		MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/2", "{\"status\":\"success\"}")
	}

	//log.Println("========== Run FotaImagePreHook ==========")
	//FotaImagePreHook()

	log.Println("========== Run FotaDtb ==========")
	err = FotaDtb()
	if err != nil {
		if err.Error() == "File not found" {
			MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/3", "{\"status\":\"skip\"}")
		} else {
			MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/3", "{\"status\":\"failed\"}")
			isPass = false
		}
	} else {
		MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/3", "{\"status\":\"success\"}")
	}

	log.Println("========== Run FotaKernel ==========")
	err = FotaKernel()
	if err != nil {
		if err.Error() == "File not found" {
			MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/4", "{\"status\":\"skip\"}")
		} else {
			MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/4", "{\"status\":\"failed\"}")
			isPass = false
		}
	} else {
		MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/4", "{\"status\":\"success\"}")
	}

	//log.Println("========== Run FotaImagePostHook ==========")
	//FotaImagePostHook()

	log.Println("========== Run FotaRootFs ==========")
	err = FotaRootFs()
	if err != nil {
		if err.Error() == "File not found" {
			MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/5", "{\"status\":\"skip\"}")
		} else {
			MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/5", "{\"status\":\"failed\"}")
			isPass = false
		}
	} else {
		MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/5", "{\"status\":\"success\"}")
	}

	log.Println("========== Run FotaWeb ==========")
	err = FotaDaemon()
	if err != nil {
		if err.Error() == "File not found" {
			MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/6", "{\"status\":\"skip\"}")
		} else {
			MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/6", "{\"status\":\"failed\"}")
			isPass = false
		}
	} else {
		MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/6", "{\"status\":\"success\"}")
	}

	log.Println("========== Run FotaFlash ==========")
	err = FotaFlash()
	if err != nil {
		if err.Error() == "File not found" {
			MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/7", "{\"status\":\"skip\"}")
		} else {
			MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/7", "{\"status\":\"failed\"}")
			isPass = false
		}
	} else {
		MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/7", "{\"status\":\"success\"}")
	}

	log.Println("========== Update FOTA Status ==========")
	if isPass == true {
		MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "fota/info", "{\"status\":\"success\"}")
		cmd := exec.Command("rm", "-rf", "/tmp/fota")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("remove /tmp/fota: %v", err)
		}
		cmd = exec.Command("rm", "-rf", "/tmp/fota.img")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("remove /tmp/fota.img: %v", err)
		}
	} else {
		MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "fota/info", "{\"status\":\"failed\"}")
	}
}
