package internal

type OperationType int

const (
	OperationTypeUnspecified OperationType = iota
	OperationTypeInput
	OperationTypeBondTax
	OperationTypeOutputSecurities
	OperationTypeOvernight
	OperationTypeTax
	OperationTypeBondRepaymentFull
	OperationTypeSellCard
	OperationTypeDividendTax
	OperationTypeOutput
	OperationTypeBondRepayment
	OperationTypeTaxCorrection
	OperationTypeServiceFee
	OperationTypeBenefitTax
	OperationTypeMarginFee
	OperationTypeBuy
	OperationTypeBuyCard
	OperationTypeInputSecurities
	OperationTypeSellMargin
	OperationTypeBrokerFee
	OperationTypeBuyMargin
	OperationTypeDividend
	OperationTypeSell
	OperationTypeCoupon
	OperationTypeSuccessFee
	OperationTypeDividendTransfer
	OperationTypeAccruingVarmargin
	OperationTypeWritingOffVarmargin
	OperationTypeDeliveryBuy
	OperationTypeDeliverySell
	OperationTypeTrackMfee
	OperationTypeTrackPfee
	OperationTypeTaxProgressive
	OperationTypeBondTaxProgressive
	OperationTypeDividendTaxProgressive
	OperationTypeBenefitTaxProgressive
	OperationTypeTaxCorrectionProgressive
	OperationTypeTaxRepoProgressive
	OperationTypeTaxRepo
	OperationTypeTaxRepoHold
	OperationTypeTaxRepoRefund
	OperationTypeTaxRepoHoldProgressive
	OperationTypeTaxRepoRefundProgressive
	OperationTypeDivExt
	OperationTypeTaxCorrectionCoupon
	OperationTypeCashFee
	OperationTypeOutFee
	OperationTypeOutStampDuty
	// Пропущенные значения для сохранения оригинальных номеров
	_ = iota + 1 // пропускаем 48 и 49
	OperationTypeOutputSwift
	OperationTypeInputSwift
	_ // пропускаем 52
	OperationTypeOutputAcquiring
	OperationTypeInputAcquiring
	OperationTypeOutputPenalty
	OperationTypeAdviceFee
	OperationTypeTransIisBs
	OperationTypeTransBsBs
	OperationTypeOutMulti
	OperationTypeInpMulti
	OperationTypeOverPlacement
	OperationTypeOverCom
	OperationTypeOverIncome
	OperationTypeOptionExpiration
	OperationTypeFutureExpiration
)

var OperationTypeDescs = map[OperationType]string{
	OperationTypeUnspecified:              "Тип операции не определен.",
	OperationTypeInput:                    "Пополнение брокерского счета.",
	OperationTypeBondTax:                  "Удержание НДФЛ по купонам.",
	OperationTypeOutputSecurities:         "Вывод ЦБ.",
	OperationTypeOvernight:                "Доход по сделке РЕПО овернайт.",
	OperationTypeTax:                      "Удержание налога.",
	OperationTypeBondRepaymentFull:        "Полное погашение облигаций.",
	OperationTypeSellCard:                 "Продажа ЦБ с карты.",
	OperationTypeDividendTax:              "Удержание налога по дивидендам.",
	OperationTypeOutput:                   "Вывод денежных средств.",
	OperationTypeBondRepayment:            "Частичное погашение облигаций.",
	OperationTypeTaxCorrection:            "Корректировка налога.",
	OperationTypeServiceFee:               "Удержание комиссии за обслуживание брокерского счета.",
	OperationTypeBenefitTax:               "Удержание налога за материальную выгоду.",
	OperationTypeMarginFee:                "Удержание комиссии за непокрытую позицию.",
	OperationTypeBuy:                      "Покупка ЦБ.",
	OperationTypeBuyCard:                  "Покупка ЦБ с карты.",
	OperationTypeInputSecurities:          "Перевод ценных бумаг из другого депозитария.",
	OperationTypeSellMargin:               "Продажа в результате Margin-call.",
	OperationTypeBrokerFee:                "Удержание комиссии за операцию.",
	OperationTypeBuyMargin:                "Покупка в результате Margin-call.",
	OperationTypeDividend:                 "Выплата дивидендов.",
	OperationTypeSell:                     "Продажа ЦБ.",
	OperationTypeCoupon:                   "Выплата купонов.",
	OperationTypeSuccessFee:               "Удержание комиссии SuccessFee.",
	OperationTypeDividendTransfer:         "Передача дивидендного дохода.",
	OperationTypeAccruingVarmargin:        "Зачисление вариационной маржи.",
	OperationTypeWritingOffVarmargin:      "Списание вариационной маржи.",
	OperationTypeDeliveryBuy:              "Покупка в рамках экспирации фьючерсного контракта.",
	OperationTypeDeliverySell:             "Продажа в рамках экспирации фьючерсного контракта.",
	OperationTypeTrackMfee:                "Комиссия за управление по счету автоследования.",
	OperationTypeTrackPfee:                "Комиссия за результат по счету автоследования.",
	OperationTypeTaxProgressive:           "Удержание налога по ставке 15%.",
	OperationTypeBondTaxProgressive:       "Удержание налога по купонам по ставке 15%.",
	OperationTypeDividendTaxProgressive:   "Удержание налога по дивидендам по ставке 15%.",
	OperationTypeBenefitTaxProgressive:    "Удержание налога за материальную выгоду по ставке 15%.",
	OperationTypeTaxCorrectionProgressive: "Корректировка налога по ставке 15%.",
	OperationTypeTaxRepoProgressive:       "Удержание налога за возмещение по сделкам РЕПО по ставке 15%.",
	OperationTypeTaxRepo:                  "Удержание налога за возмещение по сделкам РЕПО.",
	OperationTypeTaxRepoHold:              "Удержание налога по сделкам РЕПО.",
	OperationTypeTaxRepoRefund:            "Возврат налога по сделкам РЕПО.",
	OperationTypeTaxRepoHoldProgressive:   "Удержание налога по сделкам РЕПО по ставке 15%.",
	OperationTypeTaxRepoRefundProgressive: "Возврат налога по сделкам РЕПО по ставке 15%.",
	OperationTypeDivExt:                   "Выплата дивидендов на карту.",
	OperationTypeTaxCorrectionCoupon:      "Корректировка налога по купонам.",
	OperationTypeCashFee:                  "Комиссия за валютный остаток.",
	OperationTypeOutFee:                   "Комиссия за вывод валюты с брокерского счета.",
	OperationTypeOutStampDuty:             "Гербовый сбор.",
	OperationTypeOutputSwift:              "SWIFT-перевод.",
	OperationTypeInputSwift:               "SWIFT-перевод.",
	OperationTypeOutputAcquiring:          "Перевод на карту.",
	OperationTypeInputAcquiring:           "Перевод с карты.",
	OperationTypeOutputPenalty:            "Комиссия за вывод средств.",
	OperationTypeAdviceFee:                "Списание оплаты за сервис Советов.",
	OperationTypeTransIisBs:               "Перевод ценных бумаг с ИИС на брокерский счет.",
	OperationTypeTransBsBs:                "Перевод ценных бумаг с одного брокерского счета на другой.",
	OperationTypeOutMulti:                 "Вывод денежных средств со счета.",
	OperationTypeInpMulti:                 "Пополнение денежных средств со счета.",
	OperationTypeOverPlacement:            "Размещение биржевого овернайта.",
	OperationTypeOverCom:                  "Списание комиссии.",
	OperationTypeOverIncome:               "Доход от оверанайта.",
	OperationTypeOptionExpiration:         "Экспирация опциона.",
	OperationTypeFutureExpiration:         "Экспирация фьючерса.",
}
