package matrices

import (
	"math"
)

const epsilon = 1e-9

// gcd remains the same
func gcd(a, b int64) int64 {
	a = int64(math.Abs(float64(a)))
	b = int64(math.Abs(float64(b)))
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// lcm remains the same
func lcm(a, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	res := (a / gcd(a, b)) * b
	return int64(math.Abs(float64(res)))
}

// scaleToSmallestIntegers remains the same, using the more robust LCM approach
func scaleToSmallestIntegers(A [][]float64) [][]int64 {
	r := len(A)
	if r == 0 {
		return nil
	}
	c := len(A[0])
	result := make([][]int64, r)

	for i := 0; i < r; i++ {
		currentLCM := int64(1)
		const maxDenominator = 1000000

		for j := 0; j < c; j++ {
			val := A[i][j]
			if math.Abs(val) < epsilon {
				continue
			}

			f_val := math.Abs(val)
			foundDenominator := int64(1)

			for d := int64(1); d <= maxDenominator; d++ {
				num := f_val * float64(d)
				if math.Abs(num-math.Round(num)) < epsilon {
					foundDenominator = d
					break
				}
			}
			currentLCM = lcm(currentLCM, foundDenominator)
		}

		result[i] = make([]int64, c)
		sign := int64(1)
		for j := 0; j < c; j++ {
			if math.Abs(A[i][j]) >= epsilon {
				if A[i][j] < 0 {
					sign = int64(-1)
				}
				break
			}
		}

		for j := 0; j < c; j++ {
			scaledVal := math.Round(A[i][j] * float64(currentLCM) * float64(sign))
			result[i][j] = int64(scaledVal)
		}
	}
	return result
}

// ToFloatReducedEchelonForm performs Gauss-Jordan elimination on a matrix of float64.
// FIX APPLIED: Improved elimination step to prevent floating-point accumulation errors.
func ToFloatReducedEchelonForm(A [][]float64) [][]float64 {
	r := len(A)
	if r == 0 {
		return nil
	}
	c := len(A[0])
	pivotRow := 0

	for j := 0; j < c && pivotRow < r; j++ {

		// 1. Find pivot row
		maxRow := pivotRow
		for i := pivotRow + 1; i < r; i++ {
			if math.Abs(A[i][j]) > math.Abs(A[maxRow][j]) {
				maxRow = i
			}
		}

		if math.Abs(A[maxRow][j]) < epsilon {
			continue // Column is all zeros below pivotRow
		}

		// 2. Swap rows
		A[pivotRow], A[maxRow] = A[maxRow], A[pivotRow]

		// 3. Normalize: R_pivot = R_pivot / pivotVal
		pivotVal := A[pivotRow][j]
		if math.Abs(pivotVal) < epsilon {
			continue
		}

		for k := j; k < c; k++ {
			A[pivotRow][k] /= pivotVal
		}

		// 4. Eliminate above and below: R_i = R_i - A[i][j] * R_pivot
		for i := 0; i < r; i++ {
			if i != pivotRow {
				factor := A[i][j]

				// FIX: Start inner loop from j+1 and explicitly zero out A[i][j]
				// This is crucial for avoiding float errors in the zeroing step.
				A[i][j] = 0.0 // Set the pivot column entry to zero explicitly

				for k := j + 1; k < c; k++ {
					A[i][k] -= factor * A[pivotRow][k]
				}
			}
		}

		pivotRow++
	}

	// Clean up near-zero entries resulting from float inaccuracies
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if math.Abs(A[i][j]) < epsilon {
				A[i][j] = 0.0
			}
		}
	}

	return A
}

// ToIntegerReducedEchelonForm is the final wrapper function
func ToIntegerReducedEchelonForm(A [][]int64) [][]int64 {

	// 1. Convert initial integer matrix to float64
	r := len(A)
	c := len(A[0])
	floatA := make([][]float64, r)
	for i := 0; i < r; i++ {
		floatA[i] = make([]float64, c)
		for j := 0; j < c; j++ {
			floatA[i][j] = float64(A[i][j])
		}
	}

	// 2. Reduce to RREF using float arithmetic
	rrefFloat := ToFloatReducedEchelonForm(floatA)

	// 3. Scale back to the smallest integers
	finalIntegerRREF := scaleToSmallestIntegers(rrefFloat)

	return finalIntegerRREF
}

func FindFreeVariables(A [][]int64) []int {
	// A is the augmented matrix (m rows x n columns)
	if len(A) == 0 || len(A[0]) == 0 {
		return nil // Handle empty matrix
	}

	N_cols := len(A[0])
	// The variable columns are from index 0 up to N_cols - 2
	N_var_cols := N_cols - 1

	// Use a map to quickly track basic variable column indices
	// Key: 0-indexed column index, Value: true (it's a basic variable)
	basicVariableColumns := make(map[int]bool)

	// 1. Iterate through each row to find the pivot (leading entry)
	for i := 0; i < len(A); i++ {
		// Find the index of the first non-zero entry (the pivot) in row i
		// We only look up to the variable columns (N_var_cols)
		for j := 0; j < N_var_cols; j++ {
			if A[i][j] != 0 {
				// We found a pivot. The 0-indexed column index is j.
				basicVariableColumns[j] = true
				break // Move to the next row
			}
		}
	}

	// 2. Identify Free Variables
	var freeVariables []int

	// Iterate through all possible variable column indices (0 to N_var_cols - 1)
	for k := 0; k < N_var_cols; k++ {
		// If the 0-indexed column k is NOT in the map, it's a free variable column
		if !basicVariableColumns[k] {
			freeVariables = append(freeVariables, k)
		}
	}

	return freeVariables
}
