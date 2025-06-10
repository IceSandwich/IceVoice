package config

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
