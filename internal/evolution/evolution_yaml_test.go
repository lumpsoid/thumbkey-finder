package evolution

import (
	"testing"
)

func TestFromYaml(t *testing.T) {
  t.Log("FromYaml function testing")
  testFile := "../../test/configABCD.yml"
  e, err := FromYaml(testFile)
  if err != nil {
    t.Error(err)
  }
  rightCharset := []rune("abcd")
  for i, charFromConfig := range e.KeyboardConfig.CharSet {
    rightChar := rightCharset[i]
    if rightChar != charFromConfig {
      t.Errorf("Expected %v, got %v on %d", rightChar, charFromConfig, i)
    }
  }
  weight := e.KeyboardConfig.Weights[2][8]
  if weight != 1 {
    t.Error("weight is not 1")
  }
}
