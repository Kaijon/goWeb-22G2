package main

import (
	. "getac/goWeb/utils"
	"log"
	"net/http"

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

	err = FotaImageHook()
	if err != nil {
		MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/3", "{\"status\":\"failed\"}")
		MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/3", "{\"status\":\"failed\"}")
		isPass = false
	} else {
		err = FotaKernel()
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
	}

	err = FotaRootFs()
	if err != nil {
		if err.Error() == "File not found" {
			MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/4", "{\"status\":\"skip\"}")
		} else {
			MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/4", "{\"status\":\"failed\"}")
			isPass = false
		}
	} else {
		MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/4", "{\"status\":\"success\", \"percentage\":\"123\"}")
	}

	err = FotaStorage()
	if err != nil {
		if err.Error() == "File not found" {
			MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/5", "{\"status\":\"skip\"}")
		} else {
			MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/5", "{\"status\":\"failed\"}")
		}
	} else {
		MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "status/fota/5", "{\"status\":\"success\"}")
	}

	if isPass == true {
		MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "fota/info", "{\"status\":\"success\"}")
	} else {
		MqttClient.Publish(MQTT_INTERNAL_CLIENT_ID, "fota/info", "{\"status\":\"failed\"}")
	}
}
