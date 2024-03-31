package internal

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/baalimago/clai/internal/glob"
	"github.com/baalimago/clai/internal/models"
	"github.com/baalimago/clai/internal/photo"
	"github.com/baalimago/clai/internal/text"
	"github.com/baalimago/clai/internal/tools"
)

type PromptConfig struct {
	Photo string `yaml:"photo"`
	Query string `yaml:"query"`
}

type Mode int

const (
	HELP Mode = iota
	QUERY
	CHAT
	GLOB
	PHOTO
	VERSION
)

var defaultFlags = Configurations{
	ChatModel:    "gpt-4-turbo-preview",
	PhotoModel:   "dall-e-3",
	PhotoPrefix:  "clai",
	PhotoDir:     fmt.Sprintf("%v/Pictures", os.Getenv("HOME")),
	StdinReplace: "",
	PrintRaw:     false,
	ReplyMode:    false,
}

func getModeFromArgs(cmd string) (Mode, error) {
	switch cmd {
	case "photo", "p":
		return PHOTO, nil
	case "chat", "c":
		return CHAT, nil
	case "query", "q":
		return QUERY, nil
	case "glob", "g":
		return GLOB, nil
	case "help", "h":
		return HELP, nil
	case "version", "v":
		return VERSION, nil
	default:
		return HELP, fmt.Errorf("unknown command: '%s'\n", os.Args[1])
	}
}

func Setup(usage string) (models.Querier, error) {
	flagSet := setupFlags(defaultFlags)
	args := flag.Args()
	mode, err := getModeFromArgs(args[0])
	if err != nil {
		return nil, err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to find home dir: %v", err)
	}
	switch mode {
	case CHAT, QUERY, GLOB:
		tConf, err := tools.LoadConfigFromFile[text.Configurations](homeDir, "textConfig.json", migrateOldChatConfig, &text.DEFAULT)
		if err != nil {
			return nil, fmt.Errorf("failed to load configs: %err", err)
		}
		if mode == CHAT {
			tConf.ChatMode = true
		}
		applyFlagOverridesForText(&tConf, flagSet, defaultFlags)
		if mode == GLOB {
			globStr, err := glob.Setup()
			if err != nil {
				return nil, fmt.Errorf("failed to setup glob: %w", err)
			}
			tConf.Glob = globStr
		}
		err = tConf.SetupPrompts()
		if err != nil {
			return nil, fmt.Errorf("failed to setup prompt: %v", err)
		}
		cq, err := CreateTextQuerier(tConf)
		if err != nil {
			return nil, fmt.Errorf("failed to create text querier: %v\n", err)
		}
		return cq, nil
	case PHOTO:
		pConf, err := tools.LoadConfigFromFile(homeDir, "photoConfig.json", migrateOldPhotoConfig, &photo.DEFAULT)
		if err != nil {
			return nil, fmt.Errorf("failed to load configs: %w", err)
		}
		applyFlagOverridesForPhoto(&pConf, flagSet, defaultFlags)
		err = pConf.SetupPrompts()
		if err != nil {
			return nil, fmt.Errorf("failed to setup prompt: %v", err)
		}
		pq, err := NewPhotoQuerier(pConf)
		if err != nil {
			return nil, fmt.Errorf("failed to create photo querier: %v\n", err)
		}
		return pq, nil
	case HELP:
		fmt.Print(usage)
		os.Exit(0)
	case VERSION:
		bi, ok := debug.ReadBuildInfo()
		if !ok {
			return nil, errors.New("failed to read build info")
		}
		fmt.Printf("version: %v, go version: %v, checksum: %v\n", bi.Main.Version, bi.GoVersion, bi.Main.Sum)
		os.Exit(0)
	default:
		return nil, fmt.Errorf("unknown mode: %v\n", mode)
	}
	return nil, errors.New("unexpected conditional: how did you end up here?")
}
