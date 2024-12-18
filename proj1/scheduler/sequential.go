package scheduler

import (
	"fmt"
)

type EffectJSON struct {
    InPath  string   `json:"inPath"`
    OutPath string   `json:"outPath"`
    Effects []string `json:"effects"`
}





func RunSequential(config Config) {
	effects, err := ReadEffectsFile()
	if err != nil {
		fmt.Println("Error reading effects file")
		return
	}
	for _, effect := range effects {
		inputPath := fmt.Sprintf("../data/in/%s/%s", config.DataDirs, effect.InPath)
        outputPath := fmt.Sprintf("../data/out/%s_%s", config.DataDirs, effect.OutPath)
		err := ApplyEffects(inputPath, outputPath, effect.Effects)
		if err != nil {
			fmt.Println("Error applying effects")
			return
		}
	}

}
