package image

import (
	"bytes"
	"encoding/base64"
	"log"
	"testing"
)

func TestWatermark(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	// rd := base64.NewDecoder(base64.StdEncoding, strings.NewReader(jpegData))
	data, err := base64.StdEncoding.DecodeString(jpegData)
	if err != nil {
		t.Errorf("decode err %s", err)
		return
	}

	water, _ := base64.StdEncoding.DecodeString(pngWatermarkData)

	var buf bytes.Buffer

	var wopts = []WaterOption{
		{Pos: TopLeft},
		{Pos: TopRight},
		{Pos: Center},
		{Pos: Golden},
		{Pos: BottomLeft},
		{Pos: BottomRight},
	}

	for _, wopt := range wopts {
		err = Watermark(bytes.NewReader(data), bytes.NewReader(water), &buf, wopt)
		if err != nil {
			t.Fatalf("Watermark error: %s", err)
		}
	}

}

const (
	pngWatermarkData = `iVBORw0KGgoAAAANSUhEUgAAAEAAAAAgCAYAAACinX6EAAAAAXNSR0IArs4c6QAAAAlwSFlzAAAuIwAALiMBeKU/dgAAA6ppVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IlhNUCBDb3JlIDUuNC4wIj4KICAgPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4KICAgICAgPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIKICAgICAgICAgICAgeG1sbnM6eG1wPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvIgogICAgICAgICAgICB4bWxuczp0aWZmPSJodHRwOi8vbnMuYWRvYmUuY29tL3RpZmYvMS4wLyIKICAgICAgICAgICAgeG1sbnM6ZXhpZj0iaHR0cDovL25zLmFkb2JlLmNvbS9leGlmLzEuMC8iPgogICAgICAgICA8eG1wOk1vZGlmeURhdGU+MjAxOS0xMC0yOFQwMDoxMDowNTwveG1wOk1vZGlmeURhdGU+CiAgICAgICAgIDx4bXA6Q3JlYXRvclRvb2w+UGl4ZWxtYXRvciAzLjguNjwveG1wOkNyZWF0b3JUb29sPgogICAgICAgICA8dGlmZjpPcmllbnRhdGlvbj4xPC90aWZmOk9yaWVudGF0aW9uPgogICAgICAgICA8dGlmZjpDb21wcmVzc2lvbj4wPC90aWZmOkNvbXByZXNzaW9uPgogICAgICAgICA8dGlmZjpSZXNvbHV0aW9uVW5pdD4yPC90aWZmOlJlc29sdXRpb25Vbml0PgogICAgICAgICA8dGlmZjpZUmVzb2x1dGlvbj4zMDA8L3RpZmY6WVJlc29sdXRpb24+CiAgICAgICAgIDx0aWZmOlhSZXNvbHV0aW9uPjMwMDwvdGlmZjpYUmVzb2x1dGlvbj4KICAgICAgICAgPGV4aWY6UGl4ZWxYRGltZW5zaW9uPjY0PC9leGlmOlBpeGVsWERpbWVuc2lvbj4KICAgICAgICAgPGV4aWY6Q29sb3JTcGFjZT4xPC9leGlmOkNvbG9yU3Bh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`
)
