package chiapos

import (
	"crypto/sha256"
	"math/big"
)

// CalculateIterationsQuality implementation calculate_iterations_quality()
func CalculateIterationsQuality(qualityData []byte, size int32, difficulty uint64, ccSpOutputHash []byte) uint64 {
	// https://github.com/Chia-Network/chia-blockchain/blob/1.0rc9/src/consensus/pot_iterations.py#L46
	//def calculate_iterations_quality(
	//  difficulty_constant_factor: uint128,
	//  quality_string: bytes32,
	//  size: int,
	//  difficulty: uint64,
	//  cc_sp_output_hash: bytes32,
	//) -> uint64:
	// """
	// Calculates the number of iterations from the quality. This is derives as the difficulty times the constant factor
	// times a random number between 0 and 1 (based on quality string), divided by plot size.
	// """
	// sp_quality_string: bytes32 = std_hash(quality_string + cc_sp_output_hash)
	//
	// iters = uint64(
	//  int(difficulty)
	//  * int(difficulty_constant_factor)
	//  * int.from_bytes(sp_quality_string, "big", signed=False)
	//  // (int(pow(2, 256)) * int(_expected_plot_size(size)))
	// )
	// return max(iters, uint64(1))

	// quality_str_to_quality
	// https://github.com/Chia-Network/chia-blockchain/blob/1.0rc9/src/consensus/pos_quality.py#L22
	spQualityHash := sha256.Sum256(append(append(make([]byte, 0, 64), qualityData...), ccSpOutputHash...))
	spQuality := new(big.Int).SetBytes(spQualityHash[:])
	// ((2 * k) + 1) * (2 ** (k - 1))
	// https://github.com/Chia-Network/chia-blockchain/blob/1.0rc9/src/consensus/pos_quality.py#L10
	expectedPlotSize := new(big.Int).Mul(big.NewInt(int64(2*size+1)), new(big.Int).Lsh(big.NewInt(1), uint(size-1)))
	// t * _expected_plot_size(k)
	pow2Sqrt256 := new(big.Int).Lsh(big.NewInt(1), 256) // 2 ** 256
	qualityStrToQuality := new(big.Int).Mul(pow2Sqrt256, expectedPlotSize)

	// difficultyFactor
	// https://github.com/Chia-Network/chia-blockchain/blob/1.0rc9/src/consensus/pot_iterations.py#L46
	//difficultyFactor := new(big.Int).Lsh(big.NewInt(1), 65) // 2^65: 36893488147419103232
	difficultyFactor := new(big.Int).Lsh(big.NewInt(1), 67) // 2^67: 147573952589676412928
	// uint128(int(difficulty) * int(difficulty_constant_factor))
	difficultyInt := new(big.Int).SetUint64(difficulty)
	m := new(big.Int).Mul(difficultyInt, difficultyFactor)
	// * spQuality
	m = m.Mul(m, spQuality)
	// //
	iters := new(big.Int).Div(m, qualityStrToQuality).Uint64()
	if iters < 1 {
		return 1
	}
	return iters
}
func ByteAlign(numBits uint32) uint32 {
	return numBits + (8-((numBits)%8))%8
}
