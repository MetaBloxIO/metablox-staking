package interest

import (
	"math/big"
	"strconv"

	"github.com/metabloxStaking/dao"
	"github.com/metabloxStaking/models"
	logger "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

const floatErrorTolerance = 1e-6

func TestTruncateToHour(t *testing.T) {
	time1, _ := time.Parse("2006-01-02 15:04:05", "2022-05-27 13:31:00")
	expected1, _ := time.Parse("2006-01-02 15:04:05", "2022-05-27 13:00:00")
	assert.Equal(t, expected1, TruncateToHour(time1))
}

func TestCalculateCurrentAPY(t *testing.T) {
	tests := []struct {
		name           string
		topUpLimit     *big.Int
		lockUpPeriod   int
		defaultAPY     float64
		totalPrincipal *big.Int
		expected       string
	}{
		{
			name:           "total principal 400",
			topUpLimit:     big.NewInt(500000000000),
			lockUpPeriod:   180,
			defaultAPY:     0.2,
			totalPrincipal: big.NewInt(400000000),
			expected:       "250",
		},
		{
			name:           "total principal 900",
			topUpLimit:     big.NewInt(500000000000),
			lockUpPeriod:   180,
			defaultAPY:     0.2,
			totalPrincipal: big.NewInt(900000000),
			expected:       "111.1111111",
		},
		{
			name:           "total principal 1500",
			topUpLimit:     big.NewInt(500000000000),
			lockUpPeriod:   180,
			defaultAPY:     0.2,
			totalPrincipal: big.NewInt(1500000000),
			expected:       "66.66666667",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := &models.StakingProduct{
				TopUpLimit:   tt.topUpLimit,
				LockUpPeriod: tt.lockUpPeriod,
				DefaultAPY:   tt.defaultAPY,
			}
			assert.Equal(t, tt.expected, CalculateCurrentAPY(product, tt.totalPrincipal).String())
		})
	}
}

func TestUpdateOrderInterest(t *testing.T) {
	err := dao.InitTestDB()
	assert.Nil(t, err)
	defer dao.CleanupTestDB()

	product, err := dao.GetProductInfoByID("1")
	assert.Nil(t, err)

	order := &models.Order{
		ProductID:    product.ID,
		UserDID:      "test",
		Type:         "Pending",
		Amount:       big.NewInt(400000000),
		StringAmount: "400000000",
	}
	id, err := dao.CreateOrder(order)
	assert.Nil(t, err)
	orderID := strconv.Itoa(id)

	now := time.Now().UTC()
	err = dao.InsertPrincipalUpdate(product.ID, order.Amount.String())
	assert.Nil(t, err)

	principalUpdates, err := dao.GetPrincipalUpdates(product.ID)
	assert.Nil(t, err)

	txInfo := &models.TXInfo{
		OrderID:         orderID,
		TXCurrencyType:  "",
		TXType:          "BuyIn",
		TXHash:          "",
		Principal:       order.Amount,
		StringPrincipal: order.Amount.String(),
		Interest:        big.NewInt(0),
		StringInterest:  "0",
		UserAddress:     "",
		RedeemableTime:  "2022-01-01 00:00:00",
	}
	err = dao.SubmitBuyin(txInfo)
	assert.Nil(t, err)

	stmt1, err := dao.PrepareGetOrderBuyInPrincipal()
	assert.Nil(t, err)
	stmt2, err := dao.PrepareGetMostRecentOrderInterestUntilDate()
	assert.Nil(t, err)

	until := now.Add(time.Hour)

	err = updateOrderInterest(orderID, product, principalUpdates, until, stmt1, stmt2)
	assert.Nil(t, err)
	result, err := dao.GetSortedOrderInterestListUntilDate(orderID, until.Format(TimeFormat))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "11574074", result[0].InterestGain.String())

	until = now.Add(time.Hour * 10)

	err = updateOrderInterest(orderID, product, principalUpdates, until, stmt1, stmt2)
	assert.Nil(t, err)
	result, err = dao.GetSortedOrderInterestListUntilDate(orderID, until.Format(TimeFormat))
	assert.Nil(t, err)
	assert.Equal(t, 10, len(result))
	assert.Equal(t, "11574074", result[9].InterestGain.String())
	assert.Equal(t, "115740740", result[9].TotalInterestGain.String())
}

func TestOrderInterest_MultipleUsers(t *testing.T) {
	events := []struct {
		orderID  string
		amount   *big.Int
		hour     int  // number of hours since product start
		isRedeem bool // false = buy-in, true = redeem
	}{
		{"1", big.NewInt(200000000000), 271, false},
		{"2", big.NewInt(125000000000), 1525, false},
		{"3", big.NewInt(125000000000), 2681, false},
		{"4", big.NewInt(50000000000), 3634, false},
		{"1", big.NewInt(0), 4591, true},
		{"2", big.NewInt(0), 5845, true},
		{"3", big.NewInt(0), 7001, true},
	}

	err := dao.InitTestDB()
	assert.Nil(t, err)
	defer dao.CleanupTestDB()

	product, err := dao.GetProductInfoByID("1")
	assert.Nil(t, err)
	productStart, err := time.Parse(TimeFormat, product.StartDate)
	assert.Nil(t, err)

	prevUpdate := productStart
	for _, event := range events {
		currTime := productStart.Add(time.Hour * time.Duration(event.hour))
		if currTime.After(prevUpdate) {
			UpdateAllOrderInterest(currTime) // update before each new order buyin/redeem
			prevUpdate = currTime
		}
		assert.Nil(t, err)
		if !event.isRedeem { // order buy-in
			order := &models.Order{
				ProductID:    "1",
				UserDID:      "test",
				Type:         "Pending",
				Amount:       event.amount,
				StringAmount: event.amount.String(),
			}
			id, err := dao.BuyinTestOrderWithDate(order, currTime.Format(TimeFormat))
			assert.Nil(t, err)
			assert.Equal(t, event.orderID, id)
		} else {
			err := dao.RedeemTestOrderWithDate(event.orderID, currTime.Format(TimeFormat))
			assert.Nil(t, err)
		}
	}

	// Each test checks the number of entries and the interest gain for all active orders one hour after each new order or redemption.
	tests := []struct {
		orderID              string
		hour                 int
		expectedLen          int
		expectedInterestGain string
	}{
		{orderID: "1", hour: 272, expectedLen: 1, expectedInterestGain: "11574074"},

		{orderID: "1", hour: 1526, expectedLen: 1255, expectedInterestGain: "7122507"},
		{orderID: "2", hour: 1526, expectedLen: 1, expectedInterestGain: "4451566"},

		{orderID: "1", hour: 2682, expectedLen: 2411, expectedInterestGain: "5144032"},
		{orderID: "2", hour: 2682, expectedLen: 1157, expectedInterestGain: "3215020"},
		{orderID: "3", hour: 2682, expectedLen: 1, expectedInterestGain: "3215020"},

		{orderID: "1", hour: 3635, expectedLen: 3364, expectedInterestGain: "4629629"},
		{orderID: "2", hour: 3635, expectedLen: 2110, expectedInterestGain: "2893518"},
		{orderID: "3", hour: 3635, expectedLen: 954, expectedInterestGain: "2893518"},
		{orderID: "4", hour: 3635, expectedLen: 1, expectedInterestGain: "1157407"},

		// order 1 redeemed, should stay at length 4320 from now on
		{orderID: "1", hour: 4592, expectedLen: 4320, expectedInterestGain: "4629629"},
		{orderID: "2", hour: 4592, expectedLen: 3067, expectedInterestGain: "4822530"},
		{orderID: "3", hour: 4592, expectedLen: 1911, expectedInterestGain: "4822530"},
		{orderID: "4", hour: 4592, expectedLen: 958, expectedInterestGain: "1929012"},

		// all redeemed orders should stay at length 4320
		{orderID: "1", hour: 5846, expectedLen: 4320, expectedInterestGain: "4629629"},
		{orderID: "2", hour: 5846, expectedLen: 4320, expectedInterestGain: "4822530"},
		{orderID: "3", hour: 5846, expectedLen: 3165, expectedInterestGain: "8267195"},
		{orderID: "4", hour: 5846, expectedLen: 2212, expectedInterestGain: "3306878"},

		{orderID: "1", hour: 7002, expectedLen: 4320, expectedInterestGain: "4629629"},
		{orderID: "2", hour: 7002, expectedLen: 4320, expectedInterestGain: "4822530"},
		{orderID: "3", hour: 7002, expectedLen: 4320, expectedInterestGain: "8267195"},
		{orderID: "4", hour: 7002, expectedLen: 3368, expectedInterestGain: "11574074"},

		// order 4 was never redeemed, so it should continue gaining interest in the next product term
		{orderID: "1", hour: 7002, expectedLen: 4320, expectedInterestGain: "4629629"},
		{orderID: "2", hour: 7002, expectedLen: 4320, expectedInterestGain: "4822530"},
		{orderID: "3", hour: 7002, expectedLen: 4320, expectedInterestGain: "8267195"},
		{orderID: "4", hour: 12274, expectedLen: 8640, expectedInterestGain: "11574074"},
	}

	// not using tt.Run because we don't want to set up the purchases multiple times
	for _, tt := range tests {
		currTime := productStart.Add(time.Hour * time.Duration(tt.hour))
		if currTime.After(prevUpdate) {
			UpdateAllOrderInterest(currTime) // update before each new order buyin/redeem
			prevUpdate = currTime
		}
		result, err := dao.GetSortedOrderInterestListUntilDate(tt.orderID, currTime.Format(TimeFormat))
		assert.Nil(t, err)
		if assert.Equal(t, tt.expectedLen, len(result)) {
			assert.Equal(t, tt.expectedInterestGain, result[len(result)-1].InterestGain.String())
		} else {
			logger.Warnf("test order %s at hour %d failed", tt.orderID, tt.hour)
		}
	}
}
