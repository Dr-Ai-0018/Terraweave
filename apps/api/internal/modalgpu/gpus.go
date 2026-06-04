package modalgpu

type GPU struct {
	ID             string  `json:"id"`
	Label          string  `json:"label"`
	PricePerSecond float64 `json:"price_per_second_usd"`
	PricePerHour   float64 `json:"price_per_hour_usd"`
	MaxCount       int     `json:"max_count"`
	Notes          string  `json:"notes,omitempty"`
}

// List mirrors Modal's public GPU types and pricing checked on 2026-06-04.
func List() []GPU {
	return []GPU{
		gpu("T4", "Nvidia T4", 0.000164, 8, ""),
		gpu("L4", "Nvidia L4", 0.000222, 8, "TerraWeave default GPU."),
		gpu("A10", "Nvidia A10", 0.000306, 4, ""),
		gpu("A10G", "Nvidia A10G", 0.000306, 4, "Official Modal examples and A10G pricing article use gpu=\"A10G\"; priced the same as A10."),
		gpu("L40S", "Nvidia L40S", 0.000542, 8, "Modal recommends starting here for many inference workloads."),
		gpu("A100", "Nvidia A100", 0.000583, 8, "Alias for A100 40 GB; Modal may upgrade to A100 80 GB without changing GPU cost."),
		gpu("A100-40GB", "Nvidia A100, 40 GB", 0.000583, 8, ""),
		gpu("A100-80GB", "Nvidia A100, 80 GB", 0.000694, 8, ""),
		gpu("RTX-PRO-6000", "Nvidia RTX PRO 6000", 0.000842, 1, "Count support was not specified in the checked Modal GPU guide."),
		gpu("H100", "Nvidia H100", 0.001097, 8, "Modal may upgrade H100 requests to H200 without changing GPU cost."),
		gpu("H100!", "Nvidia H100 without automatic H200 upgrade", 0.001097, 8, ""),
		gpu("H200", "Nvidia H200", 0.001261, 8, ""),
		gpu("B200", "Nvidia B200", 0.001736, 8, ""),
		gpu("B200+", "Nvidia B200 or B300 compatible pool", 0.001736, 8, "Billed as B200; this is Modal's documented B300 opt-in path and requires CUDA 13+ compatibility when upgraded."),
	}
}

func gpu(id, label string, perSecond float64, maxCount int, notes string) GPU {
	return GPU{
		ID:             id,
		Label:          label,
		PricePerSecond: perSecond,
		PricePerHour:   perSecond * 3600,
		MaxCount:       maxCount,
		Notes:          notes,
	}
}
