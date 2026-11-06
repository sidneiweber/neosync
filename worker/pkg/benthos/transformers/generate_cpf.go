package transformers

import (
	"fmt"

	mgmtv1alpha1 "github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1"
	transformer_utils "github.com/nucleuscloud/neosync/worker/pkg/benthos/transformers/utils"
	"github.com/nucleuscloud/neosync/worker/pkg/rng"
	"github.com/redpanda-data/benthos/v4/public/bloblang"
)

// +neosyncTransformerBuilder:generate:generateCpf

func init() {
	spec := bloblang.NewPluginSpec().
		Description("Generates a random, valid Brazilian CPF (Cadastro de Pessoas Físicas) number with correct check digits.").
		Category("string").
		Param(bloblang.NewBoolParam("formatted").Default(true).Description("Whether to format the CPF as xxx.xxx.xxx-xx. If false, returns only the 11 digits.")).
		Param(bloblang.NewInt64Param("seed").Optional().Description("An optional seed value used to generate deterministic outputs."))

	err := bloblang.RegisterFunctionV2(
		"generate_cpf",
		spec,
		func(args *bloblang.ParsedParams) (bloblang.Function, error) {
			formatted, err := args.GetBool("formatted")
			if err != nil {
				return nil, err
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
				return generateCpf(randomizer, formatted), nil
			}, nil
		},
	)
	if err != nil {
		panic(err)
	}
}

func NewGenerateCpfOptsFromConfig(config *mgmtv1alpha1.GenerateCpf) (*GenerateCpfOpts, error) {
	if config == nil {
		return NewGenerateCpfOpts(nil, nil)
	}
	return NewGenerateCpfOpts(config.Formatted, nil)
}

func (t *GenerateCpf) Generate(opts any) (any, error) {
	parsedOpts, ok := opts.(*GenerateCpfOpts)
	if !ok {
		return nil, fmt.Errorf("invalid parsed opts: %T", opts)
	}

	return generateCpf(parsedOpts.randomizer, parsedOpts.formatted), nil
}

// generateCpf returns a valid CPF. The first 9 digits are random and the last 2 are
// check digits computed with the official mod 11 algorithm.
// Sequences of a single repeated digit (e.g. 111.111.111-11) pass the check digit math
// but are rejected by the Receita Federal, so they are never generated.
func generateCpf(randomizer rng.Rand, formatted bool) string {
	digits := make([]int, 11)
	for {
		for i := 0; i < 9; i++ {
			digits[i] = randomizer.Intn(10)
		}
		if !allSameDigit(digits[:9]) {
			break
		}
	}
	digits[9] = cpfCheckDigit(digits[:9])
	digits[10] = cpfCheckDigit(digits[:10])

	return formatCpfDigits(digits, formatted)
}

// cpfCheckDigit computes the next CPF check digit for the given digits.
// Weights start at len(digits)+1 and decrease to 2.
func cpfCheckDigit(digits []int) int {
	sum := 0
	weight := len(digits) + 1
	for _, d := range digits {
		sum += d * weight
		weight--
	}
	rest := sum % 11
	if rest < 2 {
		return 0
	}
	return 11 - rest
}

func allSameDigit(digits []int) bool {
	for _, d := range digits[1:] {
		if d != digits[0] {
			return false
		}
	}
	return true
}

func formatCpfDigits(digits []int, formatted bool) string {
	buf := make([]byte, 0, 14)
	for i, d := range digits {
		if formatted {
			switch i {
			case 3, 6:
				buf = append(buf, '.')
			case 9:
				buf = append(buf, '-')
			}
		}
		buf = append(buf, byte('0'+d))
	}
	return string(buf)
}

// isValidCpf reports whether the input (formatted or not) is a valid CPF.
func isValidCpf(value string) bool {
	digits := make([]int, 0, 11)
	for _, r := range value {
		switch {
		case r >= '0' && r <= '9':
			digits = append(digits, int(r-'0'))
		case r == '.' || r == '-':
		default:
			return false
		}
	}
	if len(digits) != 11 || allSameDigit(digits) {
		return false
	}
	return cpfCheckDigit(digits[:9]) == digits[9] && cpfCheckDigit(digits[:10]) == digits[10]
}
