package application

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/example/go-streaming/processor/domain/message"
	"github.com/example/go-streaming/processor/domain/repository"
	"github.com/example/go-streaming/processor/domain/video"
)

type Processor struct {
	Repo repository.VideoRepository
}

func (p *Processor) HandleEvent(msg *message.Event, outDir string) {
	folder := msg.Path
	outDir = filepath.Join(outDir, msg.ID)
	os.MkdirAll(outDir, os.ModePerm)

	processVideos(folder, outDir)
	processSubtitles(folder, outDir)

	v := video.Video{
		ID:       msg.ID,
		Event:    msg.Event,
		FileName: msg.FileName,
		Size:     msg.Size,
		Path:     msg.Path,
	}
	if err := p.Repo.Save(&v); err != nil {
		log.Println("failed to store video:", err)
	}
}

func processVideos(folder, outDir string) {
    exts := []string{"mkv", "mp4", "avi"}

    var wg sync.WaitGroup

    for _, ext := range exts {
        pattern := filepath.Join(folder, "*."+ext)
        files, err := filepath.Glob(pattern)
        if err != nil {
            log.Printf("glob error for %s: %v", pattern, err)
            continue
        }

        for _, f := range files {
            wg.Add(1)

            go func(f string) {
                defer wg.Done()

                log.Println("processing video file:", f)
                out := filepath.Join(outDir, "video")

                cmd := exec.Command(
                    "ffmpeg",
                    "-nostdin",
                    "-i", f,
                    "-preset", "slow",
                    "-c:v", "libx264",
                    "-crf", "18",
                    "-hls_time", "10",
                    "-hls_list_size", "0",
                    "-hls_segment_filename", filepath.Join(outDir, "video_segment_%03d.ts"),
                    out+".m3u8",
                )

                // Pipe stdout/stderr if you want logs in your container logs
                cmd.Stdout = os.Stdout
                cmd.Stderr = os.Stderr

                // Start the process
                if err := cmd.Start(); err != nil {
                    log.Printf("failed to start ffmpeg on %s: %v", f, err)
                    return
                }

                // Wait for it to finish
                if err := cmd.Wait(); err != nil {
                    log.Printf("ffmpeg error on %s: %v", f, err)
                } else {
                    log.Printf("finished processing %s", f)
                }
            }(f)
        }
    }
}

func processSubtitles(folder, outDir string) {
	exts := []string{"srt"}
	for _, ext := range exts {
		pattern := filepath.Join(folder, "*."+ext)
		files, _ := filepath.Glob(pattern)
		for _, f := range files {
			log.Println("processing subtitle file:", f)
			out := filepath.Join(outDir, "pt-BR.vtt")
			cmd := exec.Command("ffmpeg", "-i", f, out)
			if err := cmd.Run(); err != nil {
				log.Println("ffmpeg error:", err)
			}
		}
	}
}
