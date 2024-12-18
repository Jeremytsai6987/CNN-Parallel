package scheduler

import (
	"encoding/json"
	"fmt"
	"os"
)


func ReadEffectsFile() ([]EffectJSON, error) {
	effectsPathFile := fmt.Sprintf("../data/effects.txt")
	effectsFile, err := os.Open(effectsPathFile)
	if err != nil {
		return nil, err
	}
	defer effectsFile.Close()
	var effects []EffectJSON
	reader := json.NewDecoder(effectsFile)
    for reader.More() {
        var imageEffect EffectJSON
        err := reader.Decode(&imageEffect)
        if err != nil {
            return nil, err
        }
        effects = append(effects, imageEffect)
    }

    return effects, nil
}

