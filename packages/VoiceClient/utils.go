package voiceclient

import (
	"fmt"
	"io"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/hajimehoshi/go-mp3"
	errorhandling "github.com/jala-R/VideoAutomatorGUI/packages/ErrorHandling"
)

func convertMp3ToWav(mp3File io.ReadCloser, wavFile *os.File) error {
	mp3Decoder, err := mp3.NewDecoder(mp3File)
	if err != nil {
		errorhandling.HandleError(fmt.Errorf("failed to decode MP3: %w", err))
		return err
	}

	// Prepare PCM buffer
	pcmData := &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: 1,     // Stereo
			SampleRate:  44100, // 44.1 kHz
		},
		Data: make([]int, 0),
	}

	// Read MP3 data and convert to PCM
	buf := make([]byte, 1024)
	for {
		n, err := mp3Decoder.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			errorhandling.HandleError(fmt.Errorf("error reading MP3 data: %w", err))
			return err
		}

		// Append decoded data to PCM buffer
		for i := 0; i < n; i += 4 {
			if i+3 < n {
				// Read 16-bit samples for left and right channels
				left := int16(buf[i]) | (int16(buf[i+1]) << 8)
				right := int16(buf[i+2]) | (int16(buf[i+3]) << 8)

				// Average the two channels to create a mono sample
				monoSample := int((int(left) + int(right)) / 2)

				// Append the mono sample to PCM data
				pcmData.Data = append(pcmData.Data, monoSample)
			}
		}
	}

	// Encode PCM to WAV
	wavEncoder := wav.NewEncoder(wavFile, pcmData.Format.SampleRate, 16, pcmData.Format.NumChannels, 1)
	if err := wavEncoder.Write(pcmData); err != nil {
		errorhandling.HandleError(fmt.Errorf("failed to write WAV data: %w", err))
		return err
	}
	if err := wavEncoder.Close(); err != nil {
		errorhandling.HandleError(fmt.Errorf("failed to close WAV encoder: %w", err))
		return err
	}

	return nil

}
