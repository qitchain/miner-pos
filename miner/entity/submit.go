package entity

import "encoding/json"

type SubmitProof struct {
	//Quality         uint64
	QualityString    string `json:"quality_string"`
	PlotSize         uint32 `json:"plot_size"`
	PlotId           string `json:"plot_id"`
	PoolPublicKey    string `json:"pool_public_key"`
	FarmerPublicKey  string `json:"farmer_public_key"`
	FarmerPrivateKey string `json:"farmer_private_key"`
	SecurityKey      string `json:"security_key"`
	ResponseNumber   int32  `json:"response_number"`
	ProofXs          string `json:"proof_xs"`
	Challenge        string `json:"challenge"`
	SpHash           string `json:"sp_hash"`
	RequiredIters    uint64 `json:"required_iters"`
	Height           uint32 `json:"height"`
	ScanIterations   int64  `json:"scan_iterations"`
}

func (s *SubmitProof) ToString() string {
	data, _ := json.Marshal(s)
	return string(data)
}
