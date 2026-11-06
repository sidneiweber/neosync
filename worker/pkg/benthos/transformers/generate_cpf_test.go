package transformers

import (
	"regexp"
	"testing"

	"github.com/nucleuscloud/neosync/worker/pkg/rng"
	"github.com/redpanda-data/benthos/v4/public/bloblang"
	"github.com/stretchr/testify/require"
)

var (
	formattedCpfRegex = regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$`)
	digitsCpfRegex    = regexp.MustCompile(`^\d{11}$`)
)

func Test_cpfCheckDigit(t *testing.T) {
	// 529.982.247-25 is a commonly used valid CPF example.
	base := []int{5, 2, 9, 9, 8, 2, 2, 4, 7}
	first := cpfCheckDigit(base)
	require.Equal(t, 2, first)
	require.Equal(t, 5, cpfCheckDigit(append(base, first)))
}

func Test_isValidCpf(t *testing.T) {
	require.True(t, isValidCpf("529.982.247-25"))
	require.True(t, isValidCpf("52998224725"))
	require.False(t, isValidCpf("529.982.247-26"))
	require.False(t, isValidCpf("111.111.111-11"))
	require.False(t, isValidCpf("1234567890"))
	require.False(t, isValidCpf("529a982.247-25"))
}

func Test_generateCpf(t *testing.T) {
	randomizer := rng.New(1)
	for i := 0; i < 1000; i++ {
		formatted := generateCpf(randomizer, true)
		require.Regexp(t, formattedCpfRegex, formatted)
		require.True(t, isValidCpf(formatted), formatted)

		digits := generateCpf(randomizer, false)
		require.Regexp(t, digitsCpfRegex, digits)
		require.True(t, isValidCpf(digits), digits)
	}
}

func Test_generateCpf_Deterministic(t *testing.T) {
	require.Equal(t, generateCpf(rng.New(42), true), generateCpf(rng.New(42), true))
}

func Test_GenerateCpf_Bloblang(t *testing.T) {
	ex, err := bloblang.Parse(`root = generate_cpf()`)
	require.NoError(t, err)
	res, err := ex.Query(nil)
	require.NoError(t, err)
	require.Regexp(t, formattedCpfRegex, res)
	require.True(t, isValidCpf(res.(string)))

	ex, err = bloblang.Parse(`root = generate_cpf(formatted: false, seed: 10)`)
	require.NoError(t, err)
	res, err = ex.Query(nil)
	require.NoError(t, err)
	require.Regexp(t, digitsCpfRegex, res)
	require.True(t, isValidCpf(res.(string)))
}

func Test_GenerateCpf_Opts(t *testing.T) {
	opts, err := NewGenerateCpfOptsFromConfig(nil)
	require.NoError(t, err)
	res, err := NewGenerateCpf().Generate(opts)
	require.NoError(t, err)
	require.Regexp(t, formattedCpfRegex, res)
	require.Equal(t, "generate_cpf(formatted:true)", opts.BuildBloblangString())
}
