package types

type ContainerConfig struct {
	Digest      string
	Cores       int
	Memory      int64
	Swap        int64
	SwapPresent bool
	OnBoot      bool
}
