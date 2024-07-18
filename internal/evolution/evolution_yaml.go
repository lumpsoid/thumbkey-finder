package evolution

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"tkOptimizer/internal/key"
	"tkOptimizer/internal/keyboard"

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

func ParseLayout(height, width int, l [][]string, w [][]float64) *keyboard.Keyboard {
	rows := len(l)
	if rows == 0 {
		return nil
	}

	k := keyboard.NewEmpty(height, width)
	k.Weights = w

	for r := 0; r < rows; r++ {
		for c := 0; c < len(l[r]); c++ {
			char := l[r][c]
			if char == "" {
				continue
			}
			err := k.InsertKey([]rune(char)[0], key.New(key.Position{X: float64(c), Y: float64(r)}))
			if err != nil {
				fmt.Println(err)
				return nil
			}
		}
	}
	return k
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

	kC := NewKeyboardConfig(
		c.KeyboardConfig.Height,
		c.KeyboardConfig.Width,
		c.KeyboardConfig.Weights,
		[]rune(c.KeyboardConfig.CharSet),
	)
	err = kC.Weights.Check(kC.Height, kC.Width)
	if err != nil {
		return nil, err
	}
	keyboardsPool := make([]*keyboard.Keyboard, 0)
	kP := ParseLayout(
		c.KeyboardConfig.Height,
		c.KeyboardConfig.Width,
		c.KeyboardConfig.Layout,
		c.KeyboardConfig.Weights,
	)
	if kP != nil {
		keyboardsPool = append(keyboardsPool, kP)
	}

	return &Evolution{
		Threads:             c.Threads,
		initPopulation:      c.InitPopulation,
		Population:          c.InitPopulation,
		MinPopulation:       c.MinPopulation,
		Percentile:          c.Percentile,
		MutationProbability: c.MutationProbability,
		KeyboardConfig:      kC,
		Keyboards:           keyboardsPool,
		TestText:            testText,
		MetricHistory:       make([]float64, 0),
	}, nil
}
