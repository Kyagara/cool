package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
)

func MakePreviewGif(input string, output string) error {
	regex := regexp.MustCompile(`Duration: (\d+):(\d+):(\d+\.\d+)`)
	postDir := filepath.Dir(input)

	cmd := exec.Command("ffprobe", "-i", input)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running ffprobe")
	}

	matches := regex.FindStringSubmatch(string(stdout))
	if len(matches) != 4 {
		return fmt.Errorf("error parsing video duration")
	}

	hours, _ := strconv.Atoi(matches[1])
	minutes, _ := strconv.Atoi(matches[2])
	seconds, _ := strconv.ParseFloat(matches[3], 64)

	// if the video is 10 seconds or less, create a single GIF from the entire video
	if seconds <= 10 {
		cmd = exec.Command("ffmpeg", "-y", "-ss", "0", "-t", "10", "-i", input, "-vf", "fps=10,scale='if(gte(iw,320),320,-1)':-1:flags=lanczos", "-c:v", "gif", output)
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("error creating single GIF from video")
		}

		return nil
	}

	totalDuration := float64(hours*3600) + float64(minutes*60) + (seconds)

	// Split into 3 samples
	interval := float64(totalDuration / 3)
	sampleTimes := []float64{}
	for i := 1; i <= 3; i++ {
		sampleTimes = append(sampleTimes, float64(i)*interval-1)
	}

	clipPaths := []string{}

	// Loop through each sample time and create a gif
	for _, startTime := range sampleTimes {
		clipFile := filepath.Join(postDir, fmt.Sprintf("clip_%d.gif", (int(startTime*1000))))

		cmd = exec.Command("ffmpeg", "-y", "-ss", fmt.Sprintf("%f", startTime), "-t", "2", "-i", input, "-vf", "fps=15,scale=320:-1:flags=lanczos", "-c:v", "gif", clipFile)
		err = cmd.Run()
		if err != nil {
			return err
		}

		clipPaths = append(clipPaths, clipFile)
	}

	// Create a file for ffmpeg containing the list of clips
	concatFile := filepath.Join(postDir, "concat_list.txt")
	concatContent := ""
	for _, clip := range clipPaths {
		concatContent += fmt.Sprintf("file '%s'\n", clip)
	}

	err = os.WriteFile(concatFile, []byte(concatContent), 0644)
	if err != nil {
		return err
	}

	// Concatenate the GIF clips into a single GIF
	cmd = exec.Command("ffmpeg", "-y", "-f", "concat", "-safe", "0", "-i", concatFile, "-c:v", "gif", output)
	err = cmd.Run()
	if err != nil {
		return err
	}

	// Remove intermediate files
	for _, clip := range clipPaths {
		err = os.Remove(clip)
		if err != nil {
			return err
		}
	}

	err = os.Remove(concatFile)
	if err != nil {
		return err
	}

	return nil
}
