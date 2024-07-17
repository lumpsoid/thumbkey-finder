package keyboard

import "testing"

func TestIsExist(t *testing.T) {
  t.Log("Test AlreadyExist function")
  e := AlreadyExist{Message: "yes it is"}
  check := IsAlreadyExist(e)
  if check != true {
    t.Error("Error should be true")
  }
}
