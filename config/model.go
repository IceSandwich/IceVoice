package config

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type ModelConfig struct {
	// 模型的名字
	Name string `json:"name"`

	// 模型可以完成的任务，比如"tts", "sft"等
	Task string `json:"task"`

	// 模型的架构，比如"vits"等
	Architecture string `json:"architecture"`

	// 模型的类型，比如"14b", "0.5b"等
	ModelType string `json:"model_type"`

	// 量化类型，比如"Q4_K_M"等
	Quantization string `json:"quantization"`

	// 模型按层保存，每层都有一个名字作为key，value为层的哈希值
	Rootfs map[string]string `json:"rootfs"`
}

func WalkManifest(callback func(path string, model ModelConfig) error) error {
	return filepath.Walk(GetManifestsDir(), func(path string, info os.FileInfo, ex error) error {
		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return errors.Wrapf(err, "failed to open manifest file %s", path)
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			return errors.Wrapf(err, "failed to read manifest file %s", path)
		}

		var model ModelConfig
		if err := json.Unmarshal(data, &model); err != nil {
			return errors.Wrapf(err, "failed to parse manifest file %s", path)
		}

		return callback(path, model)
	})
}

func IsManifestEmpty(model string) bool {
	filename := filepath.Join(GetManifestsDir(), model)

	count := 0
	filepath.Walk(filename, func(path string, info fs.FileInfo, err error) error {
		count++ //not correct
		return nil
	})

	return count == 0
}
