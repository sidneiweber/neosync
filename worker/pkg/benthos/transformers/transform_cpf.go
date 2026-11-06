package transformers

import (
	"errors"
	"fmt"
	"strings"

	mgmtv1alpha1 "github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1"
	transformer_utils "github.com/nucleuscloud/neosync/worker/pkg/benthos/transformers/utils"
	"github.com/nucleuscloud/neosync/worker/pkg/rng"
	"github.com/redpanda-data/benthos/v4/public/bloblang"
)

// +neosyncTransformerBuilder:transform:transformCpf

func init() {
	spec := bloblang.NewPluginSpec().
		Description("Transforms an existing Brazilian CPF into a different, valid CPF. Keeps the input format: if the input contains dots or a hyphen the output is formatted as xxx.xxx.xxx-xx, otherwise only the 11 digits are returned.").
		Category("string").
		Param(bloblang.NewAnyParam("value").Optional()).
		Param(bloblang.NewInt64Param("seed").Optional().Description("An optional seed value used to generate deterministic outputs."))

	err := bloblang.RegisterFunctionV2(
		"transform_cpf",
		spec,
		func(args *bloblang.ParsedParams) (bloblang.Function, error) {
			valuePtr, err := args.GetOptionalString("value")
			if err != nil {
				return nil, err
			}

			var value string
			if valuePtr != nil {
				value = *valuePtr
			}

			seedArg, err := args.GetOptionalInt64("seed")
			if err != nil {
				return nil, err
			}

			seed, err := transformer_utils.GetSeedOrDefault(seedArg)
			if err != nil {
				return nil, err
			}
			randomizer := rng.New(seed)

			return func() (any, error) {
				return transformCpf(randomizer, value), nil
			}, nil
		},
	)
	if err != nil {
		panic(err)
	}
}

func NewTransformCpfOptsFromConfig(config *mgmtv1alpha1.TransformCpf) (*TransformCpfOpts, error) {
	return NewTransformCpfOpts(nil)
}

func (t *TransformCpf) Transform(value, opts any) (any, error) {
	parsedOpts, ok := opts.(*TransformCpfOpts)
	if !ok {
		return nil, fmt.Errorf("invalid parsed opts: %T", opts)
	}

	valueStr, ok := value.(string)
	if !ok {
		return nil, errors.New("value is not a string")
	}

	return transformCpf(parsedOpts.randomizer, valueStr), nil
}

// transformCpf returns a new valid CPF that differs from the input and keeps its format.
// Empty input returns nil so NULL-like values stay empty.
func transformCpf(randomizer rng.Rand, value string) *string {
	if value == "" {
		return nil
	}

	formatted := strings.ContainsAny(value, ".-")
	inputDigits := onlyDigits(value)

	for {
		res := generateCpf(randomizer, formatted)
		if onlyDigits(res) != inputDigits {
			return &res
		}
	}
}

func onlyDigits(value string) string {
	var b strings.Builder
	b.Grow(len(value))
	for _, r := range value {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
