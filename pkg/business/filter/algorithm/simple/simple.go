//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package simple

import (
	"github.com/algotiqa/portfolio-trader/pkg/business/filter/algorithm/optimization"
	"github.com/algotiqa/portfolio-trader/pkg/business/filter/model"
)

//=============================================================================

type simpleAlgorithm struct {
	ctx optimization.Context
	fc  *optimization.FilterConfig
}

//=============================================================================

func New() optimization.Algorithm {
	return &simpleAlgorithm{}
}

//=============================================================================
//===
//=== Simple algorithm implementation
//===
//=============================================================================

func (sa *simpleAlgorithm) Init(ctx optimization.Context) {
	sa.ctx = ctx
	sa.fc  = ctx.FilterConfig()
}

//=============================================================================

func (sa *simpleAlgorithm) Steps() int64 {
	actSteps := sa.fc.Activation.Steps()
	deaSteps := sa.fc.Deactivation.Steps()

	//-- Edge cases (no optimizazion is selected)

	if actSteps == 0 {
		actSteps = 1
	}

	if deaSteps == 0 {
		deaSteps = 1
	}

	return actSteps * deaSteps
}

//=============================================================================

func (sa *simpleAlgorithm) Optimize() {
	if sa.fc.Activation.AllOff() {
		tf := sa.ctx.Baseline()

		if !sa.generateDeactivations(&tf) {
			//--- Exit on error
			return
		}
	} else {
		if sa.generateCriteriaA(&sa.fc.Activation.NetReturn.Optimization, netReturnFuncA) {
			if sa.generateCriteriaA(&sa.fc.Activation.AverageTrade.Optimization, averageTradeFuncA) {
				if sa.generateCriteriaA(&sa.fc.Activation.EquityVsAverage.Optimization, equityVsAverageFuncA) {
					if sa.generateCriteriaA(&sa.fc.Activation.WinningPercentage.Optimization, winningPercentageFuncA) {
						if sa.generateCriteriaA(&sa.fc.Activation.MaxDrawdown.Optimization, maxDrawdownFuncA) {
							if sa.generateCriteriaA(&sa.fc.Activation.TrendLine.Optimization, trendLineFuncA) {
								if sa.generateCriteriaA(&sa.fc.Activation.ProfitFactor.Optimization, profitFactorFuncA) {
									if sa.generateCriteriaA(&sa.fc.Activation.SharpeRatio.Optimization, sharpeRatioFuncA) {
										if sa.generateCriteriaA(&sa.fc.Activation.Sqn100.Optimization, sqn100FuncA) {
											if sa.generateCriteriaA(&sa.fc.Activation.ReturnMaxDDRatio.Optimization, returnMaxDDRatioFuncA) {
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	sa.ctx.EndProcess()
}

//=============================================================================
//===
//=== Activations
//===
//=============================================================================

func (sa *simpleAlgorithm) generateCriteriaA(o *optimization.Optimization, f SetFunction) bool {
	if o.Enabled {
		for _, leng := range *o.Length.Steps() {
			for _,val := range *o.Value.Steps() {
				//--- Get a copy of the baseline (passed by value)
				tf := sa.ctx.Baseline()
				f(&tf, leng, val)

				if sa.generateDeactivations(&tf) {
					//--- Check if we have to stop the process

					if sa.ctx.IsStopping() {
						return false
					}
				} else {
					return false
				}
			}
		}
	}

	return true
}

//=============================================================================
//===
//=== Deactivations
//===
//=============================================================================

func (sa *simpleAlgorithm) generateDeactivations(tf *model.TradingFilter) bool {
	if sa.fc.Deactivation.AllOff() {
		_,err := sa.ctx.RunAnalysis(tf)
		if err != nil {
			sa.ctx.ExitWithError(err)
			return false
		}

		return true
	}

	if sa.generateCriteriaD(tf, &sa.fc.Deactivation.NetReturn.Optimization, netReturnFuncD) {
		if sa.generateCriteriaD(tf, &sa.fc.Deactivation.AverageTrade.Optimization, averageTradeFuncD) {
			if sa.generateCriteriaD(tf, &sa.fc.Deactivation.EquityVsAverage.Optimization, equityVsAverageFuncD) {
				if sa.generateCriteriaD(tf, &sa.fc.Deactivation.WinningPercentage.Optimization, winningPercentageFuncD) {
					if sa.generateCriteriaD(tf, &sa.fc.Deactivation.MaxDrawdown.Optimization, maxDrawdownFuncD) {
						if sa.generateCriteriaD(tf, &sa.fc.Deactivation.TrendLine.Optimization, trendLineFuncD) {
							if sa.generateCriteriaD(tf, &sa.fc.Deactivation.ProfitFactor.Optimization, profitFactorFuncD) {
								if sa.generateCriteriaD(tf, &sa.fc.Deactivation.SharpeRatio.Optimization, sharpeRatioFuncD) {
									if sa.generateCriteriaD(tf, &sa.fc.Deactivation.Sqn100.Optimization, sqn100FuncD) {
										if sa.generateCriteriaD(tf, &sa.fc.Deactivation.ReturnMaxDDRatio.Optimization, returnMaxDDRatioFuncD) {
											return true
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return false
}

//=============================================================================

func (sa *simpleAlgorithm) generateCriteriaD(base *model.TradingFilter, o *optimization.Optimization, f SetFunction) bool {
	if o.Enabled {
		for _, leng := range *o.Length.Steps() {
			for _, val := range *o.Value.Steps() {
				tf := *base
				f(&tf, leng, val)

				_,err := sa.ctx.RunAnalysis(&tf)
				if err != nil {
					sa.ctx.ExitWithError(err)
					return false
				}

				//--- Check if we have to stop the process

				if sa.ctx.IsStopping() {
					sa.ctx.LogInfo("generateTrendLineD: Got stop request")
					return false
				}
			}
		}
	}

	return true
}

//=============================================================================
//===
//=== Functions
//===
//=============================================================================

type SetFunction func(tf *model.TradingFilter, length int, value float64)

//=============================================================================
//=== Activations
//=============================================================================

func netReturnFuncA(tf *model.TradingFilter, length int, value float64) {
	tf.Activation.NetReturn.Enabled = true
	tf.Activation.NetReturn.Length  = length
	tf.Activation.NetReturn.Value   = value
}

//=============================================================================

func averageTradeFuncA(tf *model.TradingFilter, length int, value float64) {
	tf.Activation.AverageTrade.Enabled = true
	tf.Activation.AverageTrade.Length  = length
	tf.Activation.AverageTrade.Value   = value
}

//=============================================================================

func equityVsAverageFuncA(tf *model.TradingFilter, length int, value float64) {
	tf.Activation.EquityVsAverage.Enabled = true
	tf.Activation.EquityVsAverage.Length  = length
	tf.Activation.EquityVsAverage.Value   = value
}

//=============================================================================

func winningPercentageFuncA(tf *model.TradingFilter, length int, value float64) {
	tf.Activation.WinningPercentage.Enabled = true
	tf.Activation.WinningPercentage.Length  = length
	tf.Activation.WinningPercentage.Value   = value
}

//=============================================================================

func maxDrawdownFuncA(tf *model.TradingFilter, length int, value float64) {
	tf.Activation.MaxDrawdown.Enabled = true
	tf.Activation.MaxDrawdown.Length  = length
	tf.Activation.MaxDrawdown.Value   = value
}

//=============================================================================

func trendLineFuncA(tf *model.TradingFilter, length int, value float64) {
	tf.Activation.TrendLine.Enabled = true
	tf.Activation.TrendLine.Length  = length
	tf.Activation.TrendLine.Value   = value
}

//=============================================================================

func profitFactorFuncA(tf *model.TradingFilter, length int, value float64) {
	tf.Activation.ProfitFactor.Enabled = true
	tf.Activation.ProfitFactor.Length  = length
	tf.Activation.ProfitFactor.Value   = value
}

//=============================================================================

func sharpeRatioFuncA(tf *model.TradingFilter, length int, value float64) {
	tf.Activation.SharpeRatio.Enabled = true
	tf.Activation.SharpeRatio.Length  = length
	tf.Activation.SharpeRatio.Value   = value
}

//=============================================================================

func sqn100FuncA(tf *model.TradingFilter, length int, value float64) {
	tf.Activation.Sqn100.Enabled = true
	tf.Activation.Sqn100.Length  = length
	tf.Activation.Sqn100.Value   = value
}

//=============================================================================

func returnMaxDDRatioFuncA(tf *model.TradingFilter, length int, value float64) {
	tf.Activation.ReturnMaxDDRatio.Enabled = true
	tf.Activation.ReturnMaxDDRatio.Length  = length
	tf.Activation.ReturnMaxDDRatio.Value   = value
}

//=============================================================================
//=== Deactivations
//=============================================================================

func netReturnFuncD(tf *model.TradingFilter, length int, value float64) {
	tf.Deactivation.NetReturn.Enabled = true
	tf.Deactivation.NetReturn.Length  = length
	tf.Deactivation.NetReturn.Value   = value
}

//=============================================================================

func averageTradeFuncD(tf *model.TradingFilter, length int, value float64) {
	tf.Deactivation.AverageTrade.Enabled = true
	tf.Deactivation.AverageTrade.Length  = length
	tf.Deactivation.AverageTrade.Value   = value
}

//=============================================================================

func equityVsAverageFuncD(tf *model.TradingFilter, length int, value float64) {
	tf.Deactivation.EquityVsAverage.Enabled = true
	tf.Deactivation.EquityVsAverage.Length  = length
	tf.Deactivation.EquityVsAverage.Value   = value
}

//=============================================================================

func winningPercentageFuncD(tf *model.TradingFilter, length int, value float64) {
	tf.Deactivation.WinningPercentage.Enabled = true
	tf.Deactivation.WinningPercentage.Length  = length
	tf.Deactivation.WinningPercentage.Value   = value
}

//=============================================================================

func maxDrawdownFuncD(tf *model.TradingFilter, length int, value float64) {
	tf.Deactivation.MaxDrawdown.Enabled = true
	tf.Deactivation.MaxDrawdown.Length  = length
	tf.Deactivation.MaxDrawdown.Value   = value
}

//=============================================================================

func trendLineFuncD(tf *model.TradingFilter, length int, value float64) {
	tf.Deactivation.TrendLine.Enabled = true
	tf.Deactivation.TrendLine.Length  = length
	tf.Deactivation.TrendLine.Value   = value
}

//=============================================================================

func profitFactorFuncD(tf *model.TradingFilter, length int, value float64) {
	tf.Deactivation.ProfitFactor.Enabled = true
	tf.Deactivation.ProfitFactor.Length  = length
	tf.Deactivation.ProfitFactor.Value   = value
}

//=============================================================================

func sharpeRatioFuncD(tf *model.TradingFilter, length int, value float64) {
	tf.Deactivation.SharpeRatio.Enabled = true
	tf.Deactivation.SharpeRatio.Length  = length
	tf.Deactivation.SharpeRatio.Value   = value
}

//=============================================================================

func sqn100FuncD(tf *model.TradingFilter, length int, value float64) {
	tf.Deactivation.Sqn100.Enabled = true
	tf.Deactivation.Sqn100.Length  = length
	tf.Deactivation.Sqn100.Value   = value
}

//=============================================================================

func returnMaxDDRatioFuncD(tf *model.TradingFilter, length int, value float64) {
	tf.Deactivation.ReturnMaxDDRatio.Enabled = true
	tf.Deactivation.ReturnMaxDDRatio.Length  = length
	tf.Deactivation.ReturnMaxDDRatio.Value   = value
}

//=============================================================================
