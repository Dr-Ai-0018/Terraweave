package modalgpu

type GPU struct {
	ID             string  `json:"id"`
	Label          string  `json:"label"`
	Kind           string  `json:"kind"`
	MemoryGB       int     `json:"memory_gb,omitempty"`
	PricePerSecond float64 `json:"price_per_second_usd"`
	PricePerHour   float64 `json:"price_per_hour_usd"`
	MaxCount       int     `json:"max_count"`
	Notes          string  `json:"notes,omitempty"`
}

// List mirrors Modal's public GPU types and pricing checked on 2026-06-04.
func List() []GPU {
	return []GPU{
		gpu("T4", "Nvidia T4", "gpu", 16, 0.000164, 8, ""),
		gpu("L4", "Nvidia L4", "gpu", 24, 0.000222, 8, "TerraWeave default GPU."),
		gpu("A10", "Nvidia A10", "gpu", 24, 0.000306, 4, ""),
		gpu("L40S", "Nvidia L40S", "gpu", 48, 0.000542, 8, "Modal recommends starting here for many inference workloads."),
		gpu("A100", "Nvidia A100", "gpu", 40, 0.000583, 8, "Requests A100 40 GB, but Modal may upgrade to A100 80 GB without changing GPU cost."),
		gpu("A100-40GB", "Nvidia A100, 40 GB", "gpu", 40, 0.000583, 8, ""),
		gpu("A100-80GB", "Nvidia A100, 80 GB", "gpu", 80, 0.000694, 8, ""),
		gpu("RTX-PRO-6000", "Nvidia RTX PRO 6000 Blackwell Server Edition", "gpu", 96, 0.000842, 1, "Modal pricing lists this GPU; the checked Modal GPU guide does not state multi-GPU count support."),
		gpu("H100", "Nvidia H100 SXM", "gpu", 80, 0.001097, 8, "Modal may upgrade H100 requests to H200 without changing GPU cost."),
		gpu("H100!", "Nvidia H100 SXM without automatic H200 upgrade", "scheduling", 80, 0.001097, 8, "Strict H100 request for benchmarking or hardware-specific runs."),
		gpu("H200", "Nvidia H200", "gpu", 141, 0.001261, 8, ""),
		gpu("B200", "Nvidia B200", "gpu", 180, 0.001736, 8, ""),
		gpu("B200+", "Nvidia B200 or B300 compatible pool", "scheduling", 180, 0.001736, 8, "Billed as B200; this is Modal's documented B300 opt-in path and requires CUDA 13+ compatibility when upgraded."),
	}
}

func gpu(id, label, kind string, memoryGB int, perSecond float64, maxCount int, notes string) GPU {
	return GPU{
		ID:             id,
		Label:          label,
		Kind:           kind,
		MemoryGB:       memoryGB,
		PricePerSecond: perSecond,
		PricePerHour:   perSecond * 3600,
		MaxCount:       maxCount,
		Notes:          notes,
	}
}
