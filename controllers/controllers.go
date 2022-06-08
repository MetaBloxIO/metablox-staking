package controllers

import (
	"github.com/metabloxStaking/interest"
	"time"

	"github.com/MetaBloxIO/metablox-foundation-services/did"
	"github.com/gin-gonic/gin"
	"github.com/metabloxStaking/contract"
	"github.com/metabloxStaking/dao"
	"github.com/metabloxStaking/errval"
	"github.com/metabloxStaking/foundationdao"
)

const placeholderExchangeRate = 30.0

func ValidateDID(userDID string) error {
	splitDID := did.SplitDIDString(userDID)
	valid := did.IsDIDValid(splitDID)
	if !valid {
		return errval.ErrBadDID
	}
	err := contract.CheckForRegisteredDID(splitDID[2])
	if err != nil {
		return err
	}
	return nil
}

func GetProductInfoByIDHandler(c *gin.Context) {
	productID := c.Param("id")
	product, err := dao.GetProductInfoByID(productID)
	if err != nil {
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}
	principalUpdate, err := dao.GetLatestPrincipalUpdate(product.ID)
	if err != nil {
		product.CurrentAPY = product.DefaultAPY
	} else {
		product.CurrentAPY = interest.CalculateCurrentAPY(product, principalUpdate.TotalPrincipal)
	}
	ResponseSuccess(c, product)
}

func GetAllProductInfoHandler(c *gin.Context) {
	products, err := dao.GetAllProductInfo()
	if err != nil {
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}
	for _, product := range products {
		principalUpdate, err := dao.GetLatestPrincipalUpdate(product.ID)
		if err != nil {
			product.CurrentAPY = product.DefaultAPY
		} else {
			product.CurrentAPY = interest.CalculateCurrentAPY(product, principalUpdate.TotalPrincipal)
		}
	}
	ResponseSuccess(c, products)
}

func CreateOrderHandler(c *gin.Context) {
	output, err := CreateOrder(c)
	if err != nil {
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, output)
}

func SubmitBuyinHandler(c *gin.Context) {
	output, err := SubmitBuyin(c)
	if err != nil {
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, output)
}

func GetStakingRecordsHandler(c *gin.Context) {
	records, err := GetStakingRecords(c)
	if err != nil {
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, records)
}

func GetTransactionsByOrderIDHandler(c *gin.Context) {
	orderID := c.Param("id")
	transactions, err := dao.GetTransactionsByOrderID(orderID)
	if err != nil {
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, transactions)
}

func GetTransactionsByUserDIDHandler(c *gin.Context) {
	userDID := c.Param("did")

	transactions, err := dao.GetTransactionsByUserDID(userDID)
	if err != nil {
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, transactions)
}

func GetOrderInterestHandler(c *gin.Context) {
	orderID := c.Param("id")
	transactions, err := dao.GetSortedOrderInterestListUntilDate(orderID, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, transactions)
}

func RedeemOrderHandler(c *gin.Context) {
	output, err := RedeemOrder(c)
	if err != nil {
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, output)
}

func RedeemInterestHandler(c *gin.Context) {
	output, err := RedeemInterest(c)
	if err != nil {
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, output)
}

func GetMinerListHandler(c *gin.Context) {
	minerList, err := GetMinerList(c)
	if err != nil {
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, minerList)
}

func GetMinerByIDHandler(c *gin.Context) {
	miner, err := GetMinerByID(c)
	if err != nil {
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, miner)
}

func GetExchangeRateHandler(c *gin.Context) {
	exchangeRateID := c.Param("id")
	exchangeRate, err := foundationdao.GetExchangeRate(exchangeRateID)
	if err != nil {
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, exchangeRate)
}

func GetRewardHistoryHandler(c *gin.Context) {
	redeemedToken, err := GetRewardHistory(c)
	if err != nil {
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, redeemedToken)
}

func ExchangeSeedHandler(c *gin.Context) {
	output, err := ExchangeSeed(c)
	if err != nil {
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, output)
}
