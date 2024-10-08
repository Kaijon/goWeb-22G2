package main

import (
	"net/http"
	"sort"
	"time"

	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/mp4f"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

// StreamST struct
type Stream struct {
	URL      string `json:"url"`
	Status   bool   `json:"status"`
	OnDemand bool   `json:"on_demand"`
	RunLock  bool   `json:"-"`
	Codecs   []av.CodecData
	Cl       map[string]viewer
}

type view struct {
	c chan av.Packet
}

func serveHTTP() {
	router := gin.Default()
	gin.SetMode(gin.DebugMode)
	router.LoadHTMLGlob("web/templates/*")
	router.GET("/liveview", func(c *gin.Context) {
		fi, all := Config.list()
		Log.Println(fi, all)
		sort.Strings(all)
		c.HTML(http.StatusOK, "liveview.tmpl", gin.H{
			"port":     Config.Server.HTTPPort,
			"suuid":    fi,
			"suuidMap": all,
			"version":  time.Now().String(),
		})
	})
	router.GET("/player/:suuid", func(c *gin.Context) {
		_, all := Config.list()
		Log.Println(all)
		sort.Strings(all)
		c.HTML(http.StatusOK, "liveview.tmpl", gin.H{
			"port":     Config.Server.HTTPPort,
			"suuid":    c.Param("suuid"),
			"suuidMap": all,
			"version":  time.Now().String(),
		})
	})
	router.GET("/ws/:suuid", func(c *gin.Context) {
		handler := websocket.Handler(ws)
		handler.ServeHTTP(c.Writer, c.Request)
	})
	router.POST("/Upgrade", getacFota)
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"status": "oks",
		})
	})
	router.StaticFS("/static", http.Dir("web/static"))
	err := router.Run(Config.Server.HTTPPort)
	if err != nil {
		Log.Fatalln(err)
	}
}
func ws(ws *websocket.Conn) {
	defer ws.Close()
	suuid := ws.Request().FormValue("suuid")
	Log.Println("Request", suuid)
	//Check if stream exist
	//if !Config.ext(suuid) {
	//	Log.Println("Stream Not Found")
	//	return
	//}
	Config.RunIFNotRun(suuid)
	ws.SetWriteDeadline(time.Now().Add(5 * time.Second))
	cuuid, ch := Config.clAd(suuid)
	Log.Printf("clAd:%v, %v", cuuid, ch)
	defer Config.clDe(suuid, cuuid)
	codecs := Config.coGe(suuid)
	if codecs == nil {
		Log.Println("Codecs Error")
		return
	}
	for i, codec := range codecs {
		Log.Printf("codec : %v", codec)
		if codec.Type().IsAudio() && codec.Type() != av.AAC {
			Log.Println("Track", i, "Audio Codec Work Only AAC")
		}
	}
	muxer := mp4f.NewMuxer(nil)

	err := muxer.WriteHeader(codecs)
	if err != nil {
		Log.Println("muxer.WriteHeader", err)
		return
	}
	meta, init := muxer.GetInit(codecs)
	err = websocket.Message.Send(ws, append([]byte{9}, meta...))
	if err != nil {
		Log.Println("websocket.Message.Send", err)
		return
	}
	err = websocket.Message.Send(ws, init)
	if err != nil {
		return
	}
	var start bool
	go func() {
		for {
			var message string
			err := websocket.Message.Receive(ws, &message)
			if err != nil {
				ws.Close()
				return
			}
		}
	}()
	noVideo := time.NewTimer(10 * time.Second)
	var timeLine = make(map[int8]time.Duration)
	for {
		select {
		case <-noVideo.C:
			Log.Println("noVideo")
			return
		case pck := <-ch:
			if pck.IsKeyFrame {
				noVideo.Reset(10 * time.Second)
				start = true
			}
			if !start {
				continue
			}
			timeLine[pck.Idx] += pck.Duration
			pck.Time = timeLine[pck.Idx]
			ready, buf, _ := muxer.WritePacket(pck, false)
			if ready {
				err = ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err != nil {
					return
				}
				err := websocket.Message.Send(ws, buf)
				if err != nil {
					return
				}
			}
		}
	}
}
