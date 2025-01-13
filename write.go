package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type WriteOverlays map[string]WriteOverlay
type WriteMetaOverlays map[string]metaOverlay

func TransformOverlays(overlays overlays) WriteOverlays {
	transformedOverlays := make(map[string]WriteOverlay)
	for k, v := range overlays {
		transformedOverlays[k.String()] = TransformOverlay(v)
	}
	return transformedOverlays
}

func TransformMetaOverlays(metaOverlays map[overlayKey]metaOverlay) WriteMetaOverlays {
	transformedOverlays := make(WriteMetaOverlays)
	for k, v := range metaOverlays {
		transformedOverlays[k.String()] = v
	}
	return transformedOverlays
}

type WriteOverlay map[string]note
type WriteMetaOverlay struct {
	pressUp   bool
	pressDown bool
}

func TransformOverlay(grid overlay) WriteOverlay {
	transformedOverlay := make(WriteOverlay)
	for k, v := range grid {
		transformedOverlay[k.String()] = v
	}

	return transformedOverlay
}

type AllWrite struct {
	Overlays        WriteOverlays
	LineDefinitions []lineDefinition
	Beats           uint8
	Tempo           int
	Subdivisions    int
	Keyline         uint8
	MetaOverlays    WriteMetaOverlays
}

func (a AllWrite) UntransformOverlays() overlays {
	untransformedOverlays := make(overlays)
	for k, v := range a.Overlays {
		unMarshalledKey := UnMarshalOverlayKey(k)
		untransformedOverlays[unMarshalledKey] = UntransformOverlay(v)
	}
	return untransformedOverlays
}

func (a AllWrite) UntransformMetaOverlays() map[overlayKey]metaOverlay {
	untransformedOverlays := make(map[overlayKey]metaOverlay)
	for k, v := range a.MetaOverlays {
		unMarshalledKey := UnMarshalOverlayKey(k)
		untransformedOverlays[unMarshalledKey] = v
	}
	return untransformedOverlays
}

func UntransformOverlay(grid WriteOverlay) overlay {
	untransformedOverlays := make(overlay)
	for k, v := range grid {
		unMarshalledKey := UnMarshalGridKey(k)
		untransformedOverlays[unMarshalledKey] = v
	}
	return untransformedOverlays
}

func UnMarshalOverlayKey(k string) overlayKey {
	parts := strings.Split(k, "-")
	keyParts := strings.Split(parts[1], "/")
	num, err := strconv.Atoi(keyParts[0])
	denom, err := strconv.Atoi(keyParts[1])
	if err != nil {
		panic("Could not unmarshal overlaykey")
	}
	return overlayKey{uint8(num), uint8(denom)}
}

func UnMarshalGridKey(k string) gridKey {
	// Grid-00-00
	parts := strings.Split(k, "-")
	line, err := strconv.Atoi(parts[1])
	beat, err := strconv.Atoi(parts[2])
	if err != nil {
		panic("Unable to deserialize gridKey")
	}
	return gridKey{uint8(line), uint8(beat)}
}

func Read() (Definition, bool) {
	all := &AllWrite{}
	toml.DecodeFile("saveme.toml", all)
	return Definition{
		overlays:     (*all).UntransformOverlays(),
		lines:        (*all).LineDefinitions,
		beats:        (*all).Beats,
		tempo:        (*all).Tempo,
		keyline:      (*all).Keyline,
		subdivisions: (*all).Subdivisions,
		metaOverlays: (*all).UntransformMetaOverlays(),
	}, true
}

func Write(definition Definition) {
	f, err := os.Create("saveme.toml")
	if err != nil {
		panic("saveme file not saved")
	}
	encoder := toml.NewEncoder(f)
	all := AllWrite{
		Overlays:        TransformOverlays(definition.overlays),
		LineDefinitions: definition.lines,
		Beats:           definition.beats,
		Tempo:           definition.tempo,
		Subdivisions:    definition.subdivisions,
		Keyline:         definition.keyline,
		MetaOverlays:    TransformMetaOverlays(definition.metaOverlays),
	}
	err = encoder.Encode(all)
	if err != nil {
		panic("could not encode " + err.Error())
	}
}
