package bootstrapper

import (
	"testing"
)

func TestProperBinaryStrings(t *testing.T) {
	if Player.String() != "WindowsPlayer" {
		t.Fatal("Player's Binary name is not as specified")
	}
	if Studio.String() != "WindowsStudio" {
		t.Fatal("Studio's Binary name is not as specified")
	}
	if BinaryType(-1).String() != "unknown" {
		t.FailNow()
	}
}
