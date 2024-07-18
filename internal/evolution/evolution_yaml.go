package evolution

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"tkOptimizer/internal/key"
	"tkOptimizer/internal/layout"

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
	Percentile          float64             `yaml:"pass_percentile"`
	MutationProbability float64             `yaml:"mutation_probability"`
	KeyboardConfig      *KeyboardConfigYaml `yaml:"keyboard"`
	TextPath            string              `yaml:"text_path,omitempty"`
	Text                string              `yaml:"text,omitempty"`
}

func sanitizeText(text string) string {
	return strings.ReplaceAll(text, "\n", " ")
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
	if c.InitPopulation == 0 {
		return errors.New("`init_population` must be in the config or greater than 0")
	}
	if c.Percentile == 0.0 {
		return errors.New("`pass_percentile` must be in the config or greater than 0")
	}
	if c.MutationProbability == 0.0 {
		return errors.New("`mutation_probability` must be in the config or greater than 0")
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
		return nil, err
	}

	return &Evolution{
		Threads:             c.Threads,
		initPopulation:      c.InitPopulation,
		MinPopulation:       c.MinPopulation,
		Percentile:          c.Percentile,
		MutationProbability: c.MutationProbability,
		KeyboardConfig:      kC,
		TestText:            testText,
		MetricHistory:       make([]float64, 0),
	}, nil
}
