package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/chriserin/seq/internal/grid"
	overlaykey "github.com/chriserin/seq/internal/overlaykey"
	"github.com/chriserin/seq/internal/overlays"
)

type WriteOverlays map[string]WriteOverlay

func TransformOverlays(overlays overlay) WriteOverlays {
	transformedOverlays := make(map[string]WriteOverlay)
	// for k, v := range overlays {
	// 	transformedOverlays[k.WriteKey()] = TransformOverlay(v)
	// }
	return transformedOverlays
}

type WriteOverlay map[string]note

func TransformOverlay(grid overlay) WriteOverlay {
	transformedOverlay := make(WriteOverlay)
	for k, v := range grid.Notes {
		transformedOverlay[k.String()] = v
	}

	return transformedOverlay
}

type AllWrite struct {
	Overlays        WriteOverlays
	LineDefinitions []grid.LineDefinition
	Beats           uint8
	Tempo           int
	Subdivisions    int
	Keyline         uint8
	Accents         patternAccents
}

func (a AllWrite) UntransformOverlays() *overlay {
	overlays := overlays.InitOverlay(overlaykey.ROOT, nil)
	for k, _ := range a.Overlays {
		unMarshalledKey := UnMarshalOverlayKey(k)
		newOverlay := overlays.Add(unMarshalledKey)
		overlays = newOverlay
	}
	return overlays
}

// func UntransformOverlay(grid WriteOverlay) overlay {
// 	untransformedOverlays := make(overlay)
// 	for k, v := range grid {
// 		unMarshalledKey := UnMarshalGridKey(k)
// 		untransformedOverlays[unMarshalledKey] = v
// 	}
// 	return untransformedOverlays
// }

func UnMarshalOverlayKey(k string) overlayKey {
	parts := strings.Split(k, "-")
	keyParts := strings.Split(parts[1], "/")
	num, err := strconv.Atoi(keyParts[0])
	if err != nil {
		panic("Could not convert keyParts to str")
	}
	denom, err := strconv.Atoi(keyParts[1])
	if err != nil {
		panic("Could not unmarshal overlaykey")
	}
	return overlayKey{
		Shift:      uint8(num),
		Interval:   uint8(denom),
		Width:      0,
		StartCycle: 0,
	}
}

func UnMarshalGridKey(k string) gridKey {
	// Grid-00-00
	parts := strings.Split(k, "-")
	line, err := strconv.Atoi(parts[1])
	if err != nil {
		panic("Unable to deserialize gridKey")
	}
	beat, err := strconv.Atoi(parts[2])
	if err != nil {
		panic("Unable to deserialize gridKey")
	}
	return GK(uint8(line), uint8(beat))
}

func Read() (Definition, bool) {
	all := &AllWrite{}
	files, _ := getSeqFileNames()
	if len(files) > 0 {
		_, _ = toml.DecodeFile(files[len(files)-1], all)
	} else {
		return Definition{}, false
	}
	return Definition{
		overlays:     (*all).UntransformOverlays(),
		lines:        (*all).LineDefinitions,
		beats:        (*all).Beats,
		tempo:        (*all).Tempo,
		keyline:      (*all).Keyline,
		subdivisions: (*all).Subdivisions,
		accents:      (*all).Accents,
	}, true
}

func Write(definition Definition) {

	fileName := newFilename()
	dirname, _ := CreateSeqDir()
	fullFilePath := filepath.Join(dirname, fileName)
	f, err := os.Create(fullFilePath)
	if err != nil {
		panic("saveme file not saved")
	}
	encoder := toml.NewEncoder(f)
	all := AllWrite{
		Overlays:        TransformOverlays(*definition.overlays),
		LineDefinitions: definition.lines,
		Beats:           definition.beats,
		Tempo:           definition.tempo,
		Subdivisions:    definition.subdivisions,
		Keyline:         definition.keyline,
		Accents:         definition.accents,
	}
	err = encoder.Encode(all)
	if err != nil {
		panic("could not encode " + err.Error())
	}
}

var defaultSeqDir = ".seq"
var extension = ".toml"

func SeqDirName() string {
	workingDir, _ := os.Getwd()
	dirPath := filepath.Join(workingDir, defaultSeqDir)
	return dirPath
}

func CreateSeqDir() (string, error) {
	dirPath := SeqDirName()
	err := os.MkdirAll(dirPath, 0755)
	return dirPath, err
}

var PGEX_DATE_FORMAT = "20060102150405"

func newFilename() string {
	formattedNow := time.Now().Format(PGEX_DATE_FORMAT)
	return fmt.Sprintf("%s_%s%s", formattedNow, "seq", extension)
}

func getSeqFileNames() ([]string, error) {
	dirname := SeqDirName()
	dirEntries, err := os.ReadDir(dirname)
	if err != nil {
		return []string{}, errors.New("_pgex dir does not exist, use the exec command to create a .pgex file in a _pgex dir")
	}

	seqFiles := make([]string, 0, len(dirEntries))
	for _, d := range dirEntries {
		pgexFile := regexp.MustCompile(`[0-9]{14}_seq\.toml`)
		if pgexFile.Match([]byte(d.Name())) {
			result := filepath.Join(dirname, d.Name())
			seqFiles = append(seqFiles, result)
		}
	}

	return seqFiles, nil
}
