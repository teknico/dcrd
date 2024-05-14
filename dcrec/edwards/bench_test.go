// Copyright (c) 2015-2022 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package edwards

import (
	"testing"
)

func BenchmarkRecoverXFieldElement(b *testing.B) {
	curve := Edwards()
	vectors := testPointXRecoveryVectors(b.N)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vector := vectors[i]
		isNegative := vector.bIn[31]>>7 == 1
		yFE := encodedBytesToFieldElement(vector.bIn)
		_ = curve.recoverXFieldElement(isNegative, yFE)
	}
}

func BenchmarkIsOnCurve(b *testing.B) {
	curve := Edwards()
	vectors := testPointXRecoveryVectors(b.N)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vector := vectors[i]
		isNegative := vector.bIn[31]>>7 == 1
		notOnCurve := false
		_, y, err := curve.encodedBytesToBigIntPoint(vector.bIn)
		// The random point wasn't on the curve.
		if err != nil {
			notOnCurve = true
		}
		if notOnCurve {
			y = encodedBytesToBigInt(vector.bIn)
		}
		x2 := curve.recoverXBigInt(isNegative, y)
		_ = curve.IsOnCurve(x2, y)
	}
}

func BenchmarkAdd(b *testing.B) {
	curve := Edwards()
	tpcv := testPointConversionVectors(b.N)

	b.ReportAllocs()
	b.ResetTimer()
	var adds int64
	for i := 1; i < b.N; i++ {
		x1, y1, err := curve.encodedBytesToBigIntPoint(tpcv[i-1].bIn)
		// The random point wasn't on the curve.
		if err != nil {
			continue
		}
		x2, y2, err := curve.encodedBytesToBigIntPoint(tpcv[i].bIn)
		// The random point wasn't on the curve.
		if err != nil {
			continue
		}
		_, _ = curve.Add(x1, y1, x2, y2)
		adds++
	}
	b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(adds), "ns/add")
}

func BenchmarkScalarMult(b *testing.B) {
	curve := Edwards()
	vectors := testVectorsScalarMult()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		for _, vector := range vectors {
			x, y, _ := curve.encodedBytesToBigIntPoint(vector.bIn)
			sBig := encodedBytesToBigInt(vector.s) // We need big endian
			_, _ = curve.ScalarMult(x, y, sBig.Bytes())
		}
	}
	b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(len(vectors)*b.N), "ns/mult")
}
