package evolution

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"tkOptimizer/internal/key"
	"tkOptimizer/internal/layout"
	"tkOptimizer/internal/weights"

	"github.com/goccy/go-yaml"
)

type KeyboardConfigYaml struct {
	Height  int         `yaml:"height"`
	Width   int         `yaml:"width"`
	Weights [][]float64 `yaml:"weights"`
	Layout  [][]string  `yaml:"layout,omitempty"`
	CharSet string      `yaml:"characters"`
}

type ConfigYaml struct {
	Path                string
	Threads             int                 `yaml:"threads"`
	CharSet             string              `yaml:"characters"`
	InitPopulation      int                 `yaml:"init_population"`
	MinPopulation       int                 `yaml:"min_population"`
	MutationProbability float64             `yaml:"mutation_probability"`
	PlaceThreshold      float64             `yaml:"place_threshold"`
	StaleThreshold      int                 `yaml:"stale_threshold"`
	ResetThreshold      int                 `yaml:"reset_threshold"`
	KeyboardConfig      *KeyboardConfigYaml `yaml:"keyboard"`
	TextPath            string              `yaml:"text_path,omitempty"`
	Text                string              `yaml:"text,omitempty"`
}

func sanitizeText(text string) string {
	return strings.ReplaceAll(text, "\n", "")
}

func (c *ConfigYaml) GetText() (string, error) {
	var testTextString string
	if c.TextPath != "" {
		isAbs := filepath.IsAbs(c.TextPath)
		// if path is local, then we transform it into absolute
		if !isAbs {
			c.TextPath = filepath.Join(c.Path, c.TextPath)
		}
		// error on not existance
		_, err := os.Stat(c.TextPath)
		if err != nil {
			return "", err
		}

		fileText, err := os.Open(c.TextPath)
		if err != nil {
			return "", err
		}

		testText, err := io.ReadAll(fileText)
		if err != nil {
			return "", err
		}
		testTextString = sanitizeText(string(testText))
	}
	if c.Text != "" {
		testTextString = sanitizeText(c.Text)
	}
	return testTextString, nil
}

func (c *ConfigYaml) Check() error {
	if c.Threads > runtime.NumCPU() {
		return errors.New(
			fmt.Sprintf(
				"`threads` number in the config greater than cpu can handle. You specified = %d. Available = %d",
				c.Threads,
				runtime.NumCPU(),
			),
		)
	}
	if c.InitPopulation == 0 {
		return errors.New("`init_population` must be in the config or greater than 0")
	}
	if c.MutationProbability == 0.0 {
		return errors.New("`mutation_probability` must be in the config or greater than 0")
	}
	if c.PlaceThreshold == 0.0 {
		return errors.New("`place_threshold` must be in the config or greater than 0")
	}
	if c.StaleThreshold == 0 {
		return errors.New("`stale_threshold` must be in the config or greater than 0")
	}
	if c.ResetThreshold == 0 {
		return errors.New("`reset_threshold` must be in the config or greater than 0")
	}
	if c.KeyboardConfig.CharSet == "" {
		return errors.New("`characters` must be in the config or greater than 0")
	}
	if c.KeyboardConfig.Width == 0 {
		return errors.New("`width` must be in the config or greater than 0")
	}
	if c.KeyboardConfig.Height == 0 {
		return errors.New("`height` must be in the config or greater than 0")
	}
	if c.Text == "" && c.TextPath == "" {
		return errors.New("`text` or `text_path` must be in the config")
	}
	return nil
}

func ParseLayout(l [][]string) layout.Layout {
	rows := len(l)
	if rows == 0 {
		return layout.Layout{}
	}

	lN := layout.Layout{}
	for r := 0; r < rows; r++ {
		for c := 0; c < len(l[r]); c++ {
			char := l[r][c]
			if char == "" {
				continue
			}
			kN := key.New(key.Position{X: float64(c), Y: float64(r)})
			lN[[]rune(char)[0]] = kN
		}
	}
	return lN
}

func FromYaml(filePath string) (*Evolution, error) {
	// error on not existance
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("filepath: '%s' %s", filePath, err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	var c ConfigYaml
	err = yaml.NewDecoder(file).Decode(&c)
	if err != nil {
		return nil, err
	}
	err = c.Check()
	if err != nil {
		return nil, err
	}
	c.Path = filepath.Dir(filePath)

	testText, err := c.GetText()

	kL := ParseLayout(c.KeyboardConfig.Layout)

	kC := NewKeyboardConfig(
		c.KeyboardConfig.Height,
		c.KeyboardConfig.Width,
		c.KeyboardConfig.Weights,
		kL,
		[]rune(c.KeyboardConfig.CharSet),
	)
	err = kC.Weights.Check(kC.Height, kC.Width)
	if err != nil {
		kC.Weights = weights.New(kC.Height, kC.Width, 1)
	}

	return &Evolution{
		Threads:             c.Threads,
		initPopulation:      c.InitPopulation,
		MinPopulation:       c.MinPopulation,
		MutationProbability: c.MutationProbability,
		PlaceThreshold:      c.PlaceThreshold,
		StaleThreshold:      c.StaleThreshold,
		ResetThreshold:      c.ResetThreshold,
		KeyboardConfig:      kC,
		TestText:            testText,
		DistanceHistory:     make([]float64, 0),
	}, nil
}
