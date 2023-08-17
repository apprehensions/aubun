package state

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/vinegarhq/aubun/bootstrapper"
	"github.com/vinegarhq/aubun/internal/dirs"
)

var path = filepath.Join(dirs.Config, "state.toml")

type ApplicationState struct {
	Version  string
	Packages []string
}

type ApplicationStates map[string]ApplicationState

type State struct {
	DxvkInstalled bool
	Applications  ApplicationStates
}

func Save(state *State) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("# State saved by aubun. DO NOT EDIT DIRECTLY.\n\n")
	if err != nil {
		return err
	}

	err = toml.NewEncoder(file).Encode(state)
	if err != nil {
		return err
	}

	return nil
}

func Load() (State, error) {
	var state State

	_, err := toml.DecodeFile(path, &state)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return State{}, err
	}

	return state, nil
}

func SaveManifest(manif *bootstrapper.Manifest) error {
	name := manif.Version.Type.String()

	log.Printf("Saving Manifest State for %s", name)

	state, err := Load()
	if err != nil {
		return err
	}

	app := ApplicationState{
		Version: manif.Version.GUID,
	}
	for _, pkg := range manif.Packages {
		app.Packages = append(app.Packages, pkg.Checksum)
	}

	if state.Applications == nil {
		state.Applications = make(ApplicationStates, 0)
	}

	state.Applications[name] = app

	return Save(&state)
}

func SaveDxvk(installed bool) error {
	log.Printf("Saving installed DXVK State")

	state, err := Load()
	if err != nil {
		return err
	}

	state.DxvkInstalled = installed

	return Save(&state)
}

func Packages() ([]string, error) {
	var packages []string

	states, err := Load()
	if err != nil {
		return []string{}, err
	}

	for _, info := range states.Applications {
		packages = append(packages, info.Packages...)
	}

	return packages, nil
}

func Version(bt bootstrapper.BinaryType) (string, error) {
	states, err := Load()
	if err != nil {
		return "", err
	}

	return states.Applications[bt.String()].Version, nil
}

func Versions() ([]string, error) {
	var versions []string

	states, err := Load()
	if err != nil {
		return []string{}, err
	}

	for _, info := range states.Applications {
		versions = append(versions, info.Version)
	}

	return versions, nil
}

func DxvkInstalled() (bool, error) {
	states, err := Load()
	if err != nil {
		return false, err
	}

	return states.DxvkInstalled, nil
}
