package ibm

import (
	"fmt"

	"github.com/infracost/infracost/internal/schema"
	"github.com/shopspring/decimal"
)

// instructLabTokensPerUnit is the number of tokens in one usage unit. Usage is
// expressed in millions of tokens, while the pricing API prices are per token.
var instructLabTokensPerUnit = decimal.NewFromInt(1000000)

const instructLabPricingPlan = "instructlab-pricing-plan"

/*
 * https://cloud.ibm.com/catalog/services/instructlab
 * Standard = instructlab-pricing-plan (only available in us-east)
 */
func GetInstructLabCostComponents(r *ResourceInstance) []*schema.CostComponent {
	if r.Plan == instructLabPricingPlan {
		return []*schema.CostComponent{
			// Synthetic data generation
			instructLabTokensCostComponent(r, "Synthetic data tokens", "SYN_DATA_TOKEN", r.InstructLab_SynDataTokens),
			instructLabTokensCostComponent(r, "Synthetic data tokens (2x expedited)", "SYN_DATA_TOKEN_TWOE", r.InstructLab_SynDataTokens2e),
			instructLabTokensCostComponent(r, "Synthetic data tokens (4x expedited)", "SYN_DATA_TOKEN_FOURE", r.InstructLab_SynDataTokens4e),
			// Model training
			instructLabTokensCostComponent(r, "Model train tokens", "MODEL_TRAIN_TOKEN", r.InstructLab_ModelTrainTokens),
			instructLabTokensCostComponent(r, "Model train tokens (2x expedited)", "MODEL_TRAIN_TOKEN_TWOE", r.InstructLab_ModelTrainTokens2e),
			instructLabTokensCostComponent(r, "Model train tokens (4x expedited)", "MODEL_TRAIN_TOKEN_FOURE", r.InstructLab_ModelTrainTokens4e),
			// Model testing
			instructLabTokensCostComponent(r, "Model test input tokens", "MODEL_TEST_INPUT_TOKEN", r.InstructLab_ModelTestInputTokens),
			instructLabTokensCostComponent(r, "Model test output tokens", "MODEL_TEST_OUTPUT_TOKEN", r.InstructLab_ModelTestOutputTokens),
			// Model inference
			instructLabTokensCostComponent(r, "Qwen3-8B input tokens", "QWEN3-8B-INPUT-TOKENS", r.InstructLab_Qwen3_8BInputTokens),
			instructLabTokensCostComponent(r, "Qwen3-8B output tokens", "QWEN3-8B-OUTPUT-TOKENS", r.InstructLab_Qwen3_8BOutputTokens),
			instructLabTokensCostComponent(r, "GPT-OSS-120B input tokens", "GPT-OSS-120B-INPUT-TOKENS", r.InstructLab_GptOss120BInputTokens),
			instructLabTokensCostComponent(r, "GPT-OSS-120B output tokens", "GPT-OSS-120B-OUTPUT-TOKENS", r.InstructLab_GptOss120BOutputTokens),
			instructLabTokensCostComponent(r, "Gemma 4 31B IT input tokens", "GEMMA_4_31B_IT_INPUT_TOKENS", r.InstructLab_Gemma431BInputTokens),
			instructLabTokensCostComponent(r, "Gemma 4 31B IT output tokens", "GEMMA_4_31B_IT_OUTPUT_TOKENS", r.InstructLab_Gemma431BOutputTokens),
			instructLabTokensCostComponent(r, "Granite 4.0 Small input tokens", "GRANITE-4.0-SMALL-INPUT-TOKEN", r.InstructLab_Granite40SmallInputTokens),
			instructLabTokensCostComponent(r, "Granite 4.0 Small output tokens", "GRANITE-4.0-SMALL-OUTPUT-TOKEN", r.InstructLab_Granite40SmallOutTokens),
			instructLabTokensCostComponent(r, "Granite 4.1 30B input tokens", "GRANITE_4_1_30B_INPUT_TOKENS", r.InstructLab_Granite4130BInputTokens),
			instructLabTokensCostComponent(r, "Granite 4.1 30B output tokens", "GRANITE_4_1_30B_OUTPUT_TOKENS", r.InstructLab_Granite4130BOutputTokens),
			instructLabTokensCostComponent(r, "Granite Embedding 278M Multilingual input tokens", "GRANITE-EMBEDDING-278M-MULTILINGUAL-INPUT-TOKENS", r.InstructLab_GraniteEmbeddingInTokens),
			instructLabTokensCostComponent(r, "Llama 3.3 70B Instruct input tokens", "LLAMA-3.3-70B-INSTRUCT-INPUT-TOKENS", r.InstructLab_Llama3370BInputTokens),
			instructLabTokensCostComponent(r, "Llama 3.3 70B Instruct output tokens", "LLAMA-3.3-70B-INSTRUCT-OUTPUT-TOKENS", r.InstructLab_Llama3370BOutputTokens),
			instructLabTokensCostComponent(r, "Mistral Small 3.2 24B Instruct input tokens", "MISTRAL-SMALL-3.2-24B-INSTRUCT-INPUT-TOKENS", r.InstructLab_MistralSmallInputTokens),
			instructLabTokensCostComponent(r, "Mistral Small 3.2 24B Instruct output tokens", "MISTRAL-SMALL-3.2-24B-INSTRUCT-OUTPUT-TOKENS", r.InstructLab_MistralSmallOutputTokens),
			instructLabTokensCostComponent(r, "Nemotron 3 Nano 30B A3B FP8 input tokens", "NEMOTRON3_NANO30B_A3B_FP8_INPUT_TOKENS", r.InstructLab_NemotronNanoInputTokens),
			instructLabTokensCostComponent(r, "Nemotron 3 Nano 30B A3B FP8 output tokens", "NEMOTRON3_NANO30B_A3B_FP8_OUTPUT_TOKENS", r.InstructLab_NemotronNanoOutputTokens),
		}
	}

	costComponent := schema.CostComponent{
		Name:            fmt.Sprintf("Plan %s not found", r.Plan),
		UnitMultiplier:  decimal.NewFromInt(1),
		MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
	}
	costComponent.SetCustomPrice(decimalPtr(decimal.NewFromInt(0)))
	return []*schema.CostComponent{
		&costComponent,
	}
}

// instructLabTokensCostComponent builds a per-token cost component. The usage is in millions of
// tokens; the quantity is converted to tokens and the unit multiplier displays it per 1M tokens.
func instructLabTokensCostComponent(r *ResourceInstance, name, priceUnit string, usage *float64) *schema.CostComponent {
	var q *decimal.Decimal
	if usage != nil {
		q = decimalPtr(decimal.NewFromFloat(*usage).Mul(instructLabTokensPerUnit))
	}

	return &schema.CostComponent{
		Name:            name,
		Unit:            "1M tokens",
		UnitMultiplier:  instructLabTokensPerUnit,
		MonthlyQuantity: q,
		ProductFilter: &schema.ProductFilter{
			VendorName: strPtr("ibm"),
			Region:     strPtr(r.Location),
			Service:    &r.Service,
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "planName", Value: &r.Plan},
			},
		},
		PriceFilter: &schema.PriceFilter{
			Unit: strPtr(priceUnit),
		},
	}
}
