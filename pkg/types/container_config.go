package types

type ContainerConfig struct {
	Digest      string
	Cores       int
	Memory      uint64
	Swap        uint64
	SwapPresent bool
	OnBoot      bool
}
