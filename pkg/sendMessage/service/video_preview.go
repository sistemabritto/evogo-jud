package send_service

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"time"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

// Patch Cramer (31/10/2026): sem miniatura, duração e dimensões o WhatsApp mostra o vídeo
// como um bloco escuro com o tamanho do arquivo. Gera os três com ffprobe/ffmpeg (já na imagem).
// Em qualquer erro deixa o campo vazio e o envio segue como antes.
func videoPreview(fileData []byte) (thumb []byte, seconds, width, height uint32) {
	tmp, err := os.CreateTemp("", "evgo-video-*.mp4")
	if err != nil {
		return
	}
	defer os.Remove(tmp.Name())
	if _, err = tmp.Write(fileData); err != nil {
		tmp.Close()
		return
	}
	tmp.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if out, err := exec.CommandContext(ctx, "ffprobe", "-v", "error", "-select_streams", "v:0",
		"-show_entries", "stream=width,height:stream_side_data=rotation:format=duration", "-of", "json", tmp.Name()).Output(); err == nil {
		var p struct {
			Streams []struct {
				Width        int `json:"width"`
				Height       int `json:"height"`
				SideDataList []struct {
					Rotation float64 `json:"rotation"`
				} `json:"side_data_list"`
			} `json:"streams"`
			Format struct {
				Duration string `json:"duration"`
			} `json:"format"`
		}
		if json.Unmarshal(out, &p) == nil {
			if len(p.Streams) > 0 {
				w, h := p.Streams[0].Width, p.Streams[0].Height
				for _, sd := range p.Streams[0].SideDataList {
					if int(sd.Rotation)%180 != 0 {
						w, h = h, w
					}
				}
				if w > 0 && h > 0 {
					width, height = uint32(w), uint32(h)
				}
			}
			if d, err := strconv.ParseFloat(p.Format.Duration, 64); err == nil && d > 0 {
				seconds = uint32(d + 0.5)
			}
		}
	}
	at := "1"
	if seconds < 2 {
		at = "0"
	}
	// Um quadro em PNG (rotação já aplicada pelo ffmpeg) e a mesma miniatura das imagens:
	// 72 px de largura, JPEG do Go com qualidade 50 (makeJPEGThumbnail). A de 320 px via mjpeg não aparecia no WhatsApp.
	frame, err := exec.CommandContext(ctx, "ffmpeg", "-v", "error", "-ss", at, "-i", tmp.Name(), "-frames:v", "1",
		"-vf", "scale=240:-2", "-f", "image2", "-c:v", "png", "pipe:1").Output()
	if err == nil {
		thumb = makeJPEGThumbnail(frame, 72)
	}
	return
}

func addVideoPreview(vm *waE2E.VideoMessage, fileData []byte) {
	if vm == nil {
		return
	}
	thumb, seconds, width, height := videoPreview(fileData)
	if len(thumb) > 0 {
		vm.JPEGThumbnail = thumb
	}
	if seconds > 0 {
		vm.Seconds = proto.Uint32(seconds)
	}
	if width > 0 && height > 0 {
		vm.Width = proto.Uint32(width)
		vm.Height = proto.Uint32(height)
	}
}
