package transformers

import (
	"testing"

	"github.com/nucleuscloud/neosync/worker/pkg/rng"
	"github.com/redpanda-data/benthos/v4/public/bloblang"
	"github.com/stretchr/testify/require"
)

func Test_transformCpf(t *testing.T) {
	randomizer := rng.New(1)

	t.Run("formatted input keeps format", func(t *testing.T) {
		res := transformCpf(randomizer, "529.982.247-25")
		require.NotNil(t, res)
		require.Regexp(t, formattedCpfRegex, *res)
		require.True(t, isValidCpf(*res))
		require.NotEqual(t, "529.982.247-25", *res)
	})

	t.Run("digits input keeps format", func(t *testing.T) {
		res := transformCpf(randomizer, "52998224725")
		require.NotNil(t, res)
		require.Regexp(t, digitsCpfRegex, *res)
		require.True(t, isValidCpf(*res))
		require.NotEqual(t, "52998224725", *res)
	})

	t.Run("empty input returns nil", func(t *testing.T) {
		require.Nil(t, transformCpf(randomizer, ""))
	})
}

func Test_TransformCpf_Bloblang(t *testing.T) {
	ex, err := bloblang.Parse(`root = transform_cpf(value: "529.982.247-25")`)
	require.NoError(t, err)
	res, err := ex.Query(nil)
	require.NoError(t, err)
	str, ok := res.(*string)
	require.True(t, ok)
	require.Regexp(t, formattedCpfRegex, *str)
	require.True(t, isValidCpf(*str))
}

func Test_TransformCpf_Transform(t *testing.T) {
	opts, err := NewTransformCpfOptsFromConfig(nil)
	require.NoError(t, err)
	res, err := NewTransformCpf().Transform("52998224725", opts)
	require.NoError(t, err)
	str, ok := res.(*string)
	require.True(t, ok)
	require.Regexp(t, digitsCpfRegex, *str)

	_, err = NewTransformCpf().Transform(123, opts)
	require.Error(t, err)

	require.Equal(t, "transform_cpf(value:this.cpf)", opts.BuildBloblangString("cpf"))
}
