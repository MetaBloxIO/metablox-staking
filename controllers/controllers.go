package controllers

import (
	"github.com/MetaBloxIO/metablox-foundation-services/did"
	"github.com/gin-gonic/gin"
	"github.com/metabloxStaking/contract"
	"github.com/metabloxStaking/dao"
	"github.com/metabloxStaking/errval"
	"github.com/metabloxStaking/foundationdao"
	logger "github.com/sirupsen/logrus"
)

const placeholderExchangeRate = 30.0

func validateDID(userDID string) error {
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
		logger.Error(err)
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}
	product.CurrentAPY = 1234 //todo: get value from Colin's code
	ResponseSuccess(c, product)
}

func GetAllProductInfoHandler(c *gin.Context) {
	products, err := dao.GetAllProductInfo()
	if err != nil {
		logger.Error(err)
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}
	ResponseSuccess(c, products)
}

func CreateOrderHandler(c *gin.Context) {
	output, err := CreateOrder(c)
	if err != nil {
		logger.Error(err)
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, output)
}

func SubmitBuyinHandler(c *gin.Context) {
	output, err := SubmitBuyin(c)
	if err != nil {
		logger.Error(err)
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, output)
}

func GetStakingRecordsHandler(c *gin.Context) {
	records, err := GetStakingRecords(c)
	if err != nil {
		logger.Error(err)
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, records)
}

func GetTransactionsByOrderIDHandler(c *gin.Context) {
	orderID := c.Param("id")
	transactions, err := dao.GetTransactionsByOrderID(orderID)
	if err != nil {
		logger.Error(err)
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, transactions)
}

func GetTransactionsByUserDIDHandler(c *gin.Context) {
	userDID := c.Param("did")
	err := validateDID(userDID)
	if err != nil {
		logger.Error(err)
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	transactions, err := dao.GetTransactionsByUserDID(userDID)
	if err != nil {
		logger.Error(err)
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, transactions)
}

func GetOrderInterestHandler(c *gin.Context) {
	orderID := c.Param("id")
	transactions, err := dao.GetOrderInterestByID(orderID)
	if err != nil {
		logger.Error(err)
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, transactions)
}

func RedeemOrderHandler(c *gin.Context) {
	output, err := RedeemOrder(c)
	if err != nil {
		logger.Error(err)
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, output)
}

func RedeemInterestHandler(c *gin.Context) {
	output, err := RedeemInterest(c)
	if err != nil {
		logger.Error(err)
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, output)
}

func GetMinerListHandler(c *gin.Context) {
	did := c.Query("did")
	err := validateDID(did)
	if err != nil {
		logger.Error(err)
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	latitude := c.Query("latitude")
	longitude := c.Query("longitude")

	if latitude == "" || longitude == "" {
		minerList, err := GetMinerList()
		if err != nil {
			logger.Error(err)
			ResponseErrorWithMsg(c, CodeError, err.Error())
			return
		}
		ResponseSuccess(c, minerList)
		return
	}

	closestMiner, err := GetClosestMiner(latitude, longitude)
	if err != nil {
		logger.Error(err)
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, closestMiner)
}

func GetMinerByIDHandler(c *gin.Context) {
	miner, err := GetMinerByID(c)
	if err != nil {
		logger.Error(err)
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, miner)
}

func GetExchangeRateHandler(c *gin.Context) {
	exchangeRateID := c.Param("id")
	exchangeRate, err := foundationdao.GetExchangeRate(exchangeRateID)
	if err != nil {
		logger.Error(err)
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, exchangeRate)
}

func GetRewardHistoryHandler(c *gin.Context) {
	redeemedToken, err := GetRewardHistory(c)
	if err != nil {
		logger.Error(err)
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, redeemedToken)
}

func ExchangeSeedHandler(c *gin.Context) {
	output, err := ExchangeSeed(c)
	if err != nil {
		logger.Error(err)
		ResponseErrorWithMsg(c, CodeError, err.Error())
		return
	}

	ResponseSuccess(c, output)
}