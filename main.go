package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "outfile", "[duration]")
		return
	}
	outfile := os.Args[1]
	duration := "3600"
	if len(os.Args) > 2 {
		duration = os.Args[2]
	}
	recordScreen(outfile, duration)
}

func recordScreen(outfile, duration string) {
	cmd := exec.Command("ffmpeg",
		"-video_size", "3840x2160",
		"-framerate", "30",
		"-f", "x11grab",
		"-i", ":0.0",
		"-d", duration,
		"-y",
		outfile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	err = cmd.Process.Signal(sig)
	if err != nil {
		panic(err)
	}
	st, err := cmd.Process.Wait()
	if err != nil {
		panic(err)
	}
	os.Exit(st.ExitCode())
}
