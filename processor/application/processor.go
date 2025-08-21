package application

import (
	"errors"
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
	folder := filepath.Join(os.Getenv("WATCH_DIR"), msg.FileName)
	outDir = filepath.Join(outDir, msg.ID)

	log.Println("input folder:", folder)
	log.Println("output directory:", outDir)

	os.MkdirAll(outDir, os.ModePerm)

	var wg sync.WaitGroup
	var errs = make(chan error, 2)

	processVideos(&wg, &errs, folder, outDir)
	processSubtitles(&wg, &errs, folder, outDir)

	go func() {
		wg.Wait()
		log.Println("All processing done for event:", msg.ID)

		close(errs)

		if len(errs) > 0 {
			os.RemoveAll(outDir)
			return
		}

		p.saveVideo(msg)
	}()
}

func (p *Processor) saveVideo(msg *message.Event) {
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

func processVideos(wg *sync.WaitGroup, errs *chan error, folder string, outDir string) {
	exts := []string{"mkv", "mp4", "avi"}
	wg.Add(1)
	f, err := findFileWithExt(folder, exts)

	if err != nil {
		log.Println("No video files found:", err)
		*errs <- err
		wg.Done()
		return
	}

	log.Println("Found video file:", f)

	go func(f string) {
		defer wg.Done()

		log.Println("processing video file:", f)

		cmd := cmdProcessVideo(f, outDir)

		if err := cmd.Start(); err != nil {
			log.Printf("failed to start ffmpeg on %s: %v", f, err)
			return
		}

		if err := cmd.Wait(); err != nil {
			log.Printf("ffmpeg error on %s: %v", f, err)
		} else {
			log.Printf("finished processing %s", f)
		}
	}(f)
}

func cmdProcessVideo(f string, outDir string) *exec.Cmd {
	segmentPath := filepath.Join(outDir, "video_segment_%03d.m4s")
	playListPath := filepath.Join(outDir, "video.m3u8")

	cmd := exec.Command(
		"ffmpeg",
		"-hide_banner",
		"-nostdin",
		"-sn",
		"-i", f,
		"-map", "0:v:0",
		"-map", "0:a:0",
		"-c:v", "copy",
		"-tag:v", "hvc1",
		"-c:a", "copy",
		"-f", "hls",
		"-hls_time", "6",
		"-hls_playlist_type", "vod",
		"-hls_segment_type", "fmp4",
		"-hls_flags", "independent_segments",
		"-hls_fmp4_init_filename", "init.mp4",
		"-hls_segment_filename", segmentPath,
		playListPath,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func findFileWithExt(folder string, exts []string) (string, error) {
	if len(exts) == 0 {
		return "", errors.New("No extensions provided")
	}

	log.Println("Searching for files in folder:", folder)

	for _, ext := range exts {
		pattern := filepath.Join(folder, "*."+ext)
		files, err := filepath.Glob(pattern)

		log.Printf("Searching for files with pattern: %s", pattern)

		if err != nil {
			log.Printf("glob error for %s: %v", pattern, err)
			continue
		}

		log.Printf("Found files with %s extension: %v", ext, files)

		if len(files) > 0 {
			return files[0], nil
		}
	}

	return "", errors.New("No files found with specified extensions")
}

func processSubtitles(wg *sync.WaitGroup, errs *chan error, folder string, outDir string) {
	exts := []string{"srt"}
	f, err := findFileWithExt(folder, exts)

	wg.Add(1)

	if err != nil {
		log.Println("No subtitle files found:", err)
		*errs <- err
		wg.Done()
		return
	}

	go func(f string) {
		defer wg.Done()

		log.Println("processing subtitle file:", f)
		out := filepath.Join(outDir, "pt-BR.vtt")
		cmd := exec.Command("ffmpeg", "-i", f, out)

		if err := cmd.Start(); err != nil {
			log.Println("ffmpeg error:", err)
			return
		}

		if err := cmd.Wait(); err != nil {
			log.Printf("ffmpeg error on %s: %v", f, err)
		} else {
			log.Printf("finished processing %s", f)
		}
	}(f)
}
