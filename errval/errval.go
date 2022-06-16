package errval

import "errors"

var (
	ErrBadDID                       = errors.New("did is improperly formatted")
	ErrDIDNotRegistered             = errors.New("did is not registered in contract")
	ErrNotEnoughInterest            = errors.New("order does not meet minimum interest required to redeem")
	ErrTransactionPending           = errors.New("transaction is still pending")
	ErrTransactionReverted          = errors.New("transaction is reverted")
	ErrAddressComparisonFail        = errors.New("tx address does not match order address")
	ErrAmountComparisonFail         = errors.New("tx amount does not match order amount")
	ErrExistingTXHash               = errors.New("provided tx hash is already recorded in db")
	ErrEarlyOrderRedeem             = errors.New("order can only be redeemed on final day of term")
	ErrVCIDFormat                   = errors.New("VC id is improperly formatted")
	ErrUpdateOrderStatus            = errors.New("failed to update order status; it may not exist, or it may already be holding")
	ErrETHAddress                   = errors.New("not a correct ETH address")
	ErrTokenBalance                 = errors.New("token balance is not enough")
	ErrETHBalance                   = errors.New("ETH balance is not enough")
	ErrContractData                 = errors.New("contract data  format is not correct")
	ErrContractMethod               = errors.New("contract method is not correct")
	ErrContractParam                = errors.New("contract params is not correct")
	ErrWalletAddress                = errors.New("pls send tokens to spec wallet address")
	ErrInvalidSession               = errors.New("invalid session id")
	ErrNonceTimeout                 = errors.New("nonce timeout")
	ErrDIDMismatch                  = errors.New("DID in request response mis-match")
	ErrSignatureVerifyError         = errors.New("signature verify error")
	ErrInvalidValidator             = errors.New("failed to find credential for validator")
	ErrInvalidMiner                 = errors.New("failed to find credential for miner")
	ErrOrderAmountTooLow            = errors.New("order amount is too low")
	ErrOrderExceedsTopUpLimit       = errors.New("order amount exceeds product's pool limit")
	ErrMinerRoleNotFound            = errors.New("miner has no assigned wallet address")
	ErrDIDNotConstant               = errors.New("all seeds must belong to the same DID")
	ErrExchangeRateNotNumber        = errors.New("invalid exchange rate value")
	ErrAmountNotNumber              = errors.New("invalid amount value")
	ErrSeedAlreadyExchanged         = errors.New("this seed has already been exchanged")
	ErrDuplicateRequest             = errors.New("request used multiple times in input")
	ErrDuplicateResult              = errors.New("result used multiple times in input")
	ErrAccumulatedInterestNotNumber = errors.New("failed to parse AccumulatedInterest as valid number")
	ErrTotalInterestGainedNotNumber = errors.New("failed to parse TotalInterestGained as valid number")
	ErrPrincipalNotNumber           = errors.New("failed to parse Principal as valid number")
	ErrInterestNotNumber            = errors.New("failed to parse Interest as valid number")
	ErrInterestGainNotNumber        = errors.New("failed to parse InterestGain as valid number")
	ErrTotalInterestGainNotNumber   = errors.New("failed to parse TotalInterestGain as valid number")
	ErrTotalPrincipalNotNumber      = errors.New("failed to parse TotalPrincipal as valid number")
	ErrTopUpLimitNotNumber          = errors.New("failed to parse TopUpLimit as valid number")
	ErrBurnedInterestNotNumber      = errors.New("failed to parse BurnedInterest as valid number")
)
