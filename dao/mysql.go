package dao

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/metabloxStaking/models"
)

var SqlDB *sqlx.DB
var validate *validator.Validate

func InitSql(validatePtr *validator.Validate) error {
	validate = validatePtr

	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.dbname"),
	)

	SqlDB, err = sqlx.Open("mysql", dsn)
	if err != nil {
		logger.Error("Failed to open database: " + err.Error())
		return err
	}

	//Set the maximum number of database connections

	SqlDB.SetConnMaxLifetime(100)

	//Set the maximum number of idle connections on the database

	SqlDB.SetMaxIdleConns(10)

	//Verify connection

	if err := SqlDB.Ping(); err != nil {
		logger.Error("open database fail: ", err)
		return err
	}
	logger.Info("connect success")
	return nil
}

func GetProductInfoByID(productID string) (*models.StakingProduct, error) {
	product := models.NewStakingProduct()

	sqlStr := "select * from StakingProducts where ID = ?"
	err := SqlDB.Get(product, sqlStr, productID)
	if err != nil {
		return nil, err
	}
	err = validate.Struct(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func GetAllProductInfo() ([]*models.StakingProduct, error) {
	var products []*models.StakingProduct
	sqlStr := "select * from StakingProducts"
	rows, err := SqlDB.Queryx(sqlStr)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		product := models.NewStakingProduct()
		err = rows.StructScan(product)
		if err != nil {
			return nil, err
		}
		err = validate.Struct(product)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, err
}

func CreateOrder(order *models.Order) (int, error) {
	err := validate.Struct(order)
	if err != nil {
		return 0, err
	}

	sqlStr := "insert into Orders (ProductID, UserDID, Type, Term, PaymentAddress, Amount, UserAddress) values (:ProductID, :UserDID, :Type, :Term, :PaymentAddress, :Amount, :UserAddress)"
	result, err := SqlDB.NamedExec(sqlStr, order)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func CheckIfTXExists(txHash string) (bool, error) {
	var count int
	sqlStr := "select count(*) from TXInfo where TXHash = ?"
	err := SqlDB.Get(&count, sqlStr, txHash)
	if err != nil {
		return false, err
	}

	return (count != 0), nil
}

func GetTXCreateDate(txHash string) (string, error) {
	var date string
	sqlStr := "select unix_timestamp(CreateDate) from TXInfo where TXHash = ?"
	err := SqlDB.Get(&date, sqlStr, txHash)
	if err != nil {
		return "", err
	}
	return date, nil
}

func GetStakingRecords(did string) ([]*models.StakingRecord, error) {
	var records []*models.StakingRecord
	sqlStr := "select Orders.OrderID, Orders.ProductID, Orders.Type, Orders.Term, TXInfo.CreateDate, Orders.Amount, TXInfo.TXCurrencyType, TXInfo.RedeemableTime from Orders join TXInfo on TXInfo.OrderID = Orders.OrderID where Orders.UserDID = ? and TXInfo.TXType = 'BuyIn' and Orders.Type != 'Pending'"
	rows, err := SqlDB.Queryx(sqlStr, did)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		record := models.NewStakingRecord()
		err = rows.StructScan(record)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func GetInterestInfoByOrderID(id string) (*models.OrderInterestInfo, error) {
	info := models.NewOrderInterestInfo()
	sqlStr := "select AccumulatedInterest, TotalInterestGained from Orders where OrderID = ?"
	err := SqlDB.Get(info, sqlStr, id)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func PrepareGetInterestByOrderID() (*sqlx.Stmt, error) {
	sqlStr := "select AccumulatedInterest, TotalInterestGained from Orders where OrderID = ?"
	stmt, err := SqlDB.Preparex(sqlStr)
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func ExecuteGetInterestStmt(id string, stmt *sqlx.Stmt) (*models.OrderInterestInfo, error) {
	info := models.NewOrderInterestInfo()
	err := stmt.Get(info, id)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func GetTransactionsByOrderID(orderID string) ([]*models.TXInfo, error) {
	var transactions []*models.TXInfo
	sqlStr := "select PaymentNo, OrderID, TXCurrencyType, TXType, TXHash, Principal, Interest, UserAddress, CreateDate, RedeemableTime from TXInfo where OrderID = ?"
	rows, err := SqlDB.Queryx(sqlStr, orderID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		tx := models.NewTXInfo()
		err = rows.StructScan(tx)
		if err != nil {
			return nil, err
		}
		err = validate.Struct(tx)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, err
}

func GetTransactionsByUserDID(userDID string) ([]*models.TXInfo, error) {
	var transactions []*models.TXInfo
	sqlStr := "select TXInfo.PaymentNo, TXInfo.OrderID, TXInfo.TXCurrencyType, TXInfo.TXType, TXInfo.TXHash, TXInfo.Principal, TXInfo.Interest, TXInfo.UserAddress, TXInfo.CreateDate, TXInfo.RedeemableTime from TXInfo join Orders on Orders.OrderID = TXInfo.OrderID where Orders.UserDID = ?"
	rows, err := SqlDB.Queryx(sqlStr, userDID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		tx := models.NewTXInfo()
		err = rows.StructScan(tx)
		if err != nil {
			return nil, err
		}
		err = validate.Struct(tx)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, err
}

func GetOrderInterestByID(orderID string) ([]*models.OrderInterest, error) {
	var interests []*models.OrderInterest
	sqlStr := "select * from OrderInterest where OrderID = ?"
	rows, err := SqlDB.Queryx(sqlStr, orderID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		interest := models.NewOrderInterest()
		err = rows.StructScan(interest)
		if err != nil {
			return nil, err
		}
		err = validate.Struct(interest)
		if err != nil {
			return nil, err
		}
		interests = append(interests, interest)
	}
	return interests, nil
}

func RedeemInterestByOrderID(orderID string) error {
	sqlStr := "update OrderInterest set TotalInterestGain = 0 where OrderID = ? order by ID desc limit 1"
	_, err := SqlDB.Exec(sqlStr, orderID)
	if err != nil {
		return err
	}
	return nil
}

func GetHoldingOrders() ([]*models.Order, error) {
	var orders []*models.Order
	sqlStr := `select * from Orders where Type = 'Holding'`
	rows, err := SqlDB.Queryx(sqlStr)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		order := models.NewOrder()
		err = rows.StructScan(order)
		if err != nil {
			logger.Warn(err)
			continue
		}
		err = validate.Struct(order)
		if err != nil {
			logger.Warn(err)
			continue
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func GetOrdersByProductID(productID string) ([]*models.Order, error) {
	var orders []*models.Order
	sqlStr := `select * from Orders where ProductID = ? and Type = 'Holding'`
	rows, err := SqlDB.Queryx(sqlStr, productID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		order := models.NewOrder()
		err = rows.StructScan(order)
		if err != nil {
			logger.Warn(err)
			continue
		}
		err = validate.Struct(order)
		if err != nil {
			logger.Warn(err)
			continue
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func GetOrderByID(orderID string) (*models.Order, error) {
	order := models.NewOrder()
	sqlStr := "select * from Orders where OrderID = ?"
	err := SqlDB.Get(order, sqlStr, orderID)
	if err != nil {
		return nil, err
	}
	err = validate.Struct(order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func GetOrderCreateDate(orderID string) (string, error) {
	var createDate string
	sqlStr := "select CreateDate from TXInfo where OrderID = ? and TXType = 'BuyIn'"
	err := SqlDB.Get(&createDate, sqlStr, orderID)
	if err != nil {
		return "", err
	}
	return createDate, nil
}

func GetOrderRedeemableDate(orderID string) (string, error) {
	var redeemableDate string
	sqlStr := "select RedeemableTime from TXInfo where OrderID = ? and TXType = 'BuyIn'"
	err := SqlDB.Get(&redeemableDate, sqlStr, orderID)
	if err != nil {
		return "", err
	}
	return redeemableDate, nil
}

func GetUserAddressByOrderID(orderID string) (string, error) {
	var userAddress string
	sqlStr := "select UserAddress from Orders where OrderID = ?"
	err := SqlDB.Get(&userAddress, sqlStr, orderID)
	if err != nil {
		return "", err
	}
	return userAddress, nil
}

func GetOrderBuyInPrincipal(orderID string) (float64, error) {
	var buyInAmount float64
	sqlStr := "select Principal from TXInfo where OrderID = ? and TXType = 'BuyIn'"
	err := SqlDB.Get(&buyInAmount, sqlStr, orderID)
	if err != nil {
		return 0.0, err
	}
	return buyInAmount, nil
}

func GetMinimumInterestByOrderID(orderID string) (int, error) {
	var minInterest int
	sqlStr := "select StakingProducts.MinRedeemValue from StakingProducts join Orders on StakingProducts.ID = Orders.ProductID where Orders.OrderID = ?"
	err := SqlDB.Get(&minInterest, sqlStr, orderID)
	if err != nil {
		return 0, err
	}

	return minInterest, nil
}

func UploadTransaction(tx *models.TXInfo) error {
	err := validate.Struct(tx)
	if err != nil {
		return err
	}
	sqlStr := "insert into TXInfo (OrderID, TXCurrencyType, TXType, TXHash, Principal, Interest, UserAddress, RedeemableTime) values (:OrderID, :TXCurrencyType, :TXType, :TXHash, :Principal, :Interest, :UserAddress, :RedeemableTime)"
	_, err = SqlDB.NamedExec(sqlStr, tx)
	if err != nil {
		return err
	}
	return nil
}

func SubmitBuyin(tx *models.TXInfo) error {
	err := validate.Struct(tx)
	if err != nil {
		return err
	}
	dbTX, err := SqlDB.Beginx()
	if err != nil {
		return err
	}
	sqlStr := "update Orders set Type = 'Holding' where OrderID = ?"
	result, err := dbTX.Exec(sqlStr, tx.OrderID)
	if err != nil {
		dbTX.Rollback()
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		dbTX.Rollback()
		return err
	}
	if rows == 0 {
		dbTX.Rollback()
		return errors.New("failed to update order status; it may not exist, or it may already be holding")
	}

	sqlStr = "insert into TXInfo (OrderID, TXCurrencyType, TXType, TXHash, Principal, Interest, UserAddress, RedeemableTime) values (:OrderID, :TXCurrencyType, :TXType, :TXHash, :Principal, :Interest, :UserAddress, :RedeemableTime)"
	_, err = dbTX.NamedExec(sqlStr, tx)
	if err != nil {
		dbTX.Rollback()
		return err
	}
	dbTX.Commit()
	return nil
}

func GetTotalInterestGained(id string) (float64, error) {
	var interest float64
	sqlStr := "select TotalInterestGained from Orders where OrderID = ?"
	err := SqlDB.Get(&interest, sqlStr, id)
	if err != nil {
		return 0, err
	}
	return interest, nil
}

func HarvestOrderInterest(id string) error {
	sqlStr := "update Orders set TotalInterestGained = AccumulatedInterest where OrderID = ?"
	_, err := SqlDB.Query(sqlStr, id)
	return err
}

func GetProductNameForOrder(id string) (string, error) {
	var name string
	sqlStr := "select StakingProducts.ProductName from StakingProducts join Orders on StakingProducts.ID = Orders.ProductID where Orders.OrderID = ?"
	err := SqlDB.Get(&name, sqlStr, id)
	if err != nil {
		return "", err
	}
	return name, nil
}

func InsertPrincipalUpdate(productID string, totalPrincipal float64) error {
	sqlStr := `insert into PrincipalUpdates (ProductID, TotalPrincipal) values (?, ?)`
	_, err := SqlDB.Exec(sqlStr, productID, totalPrincipal)
	return err
}

func GetLatestPrincipalUpdate(productID string) (*models.PrincipalUpdate, error) {
	update := models.NewPrincipalUpdate()

	sqlStr := "select * from PrincipalUpdates where ProductID = ? order by Time desc"
	err := SqlDB.Get(update, sqlStr, productID)
	if err != nil {
		return nil, err
	}

	return update, nil
}

func GetPrincipalUpdates(productID string) ([]*models.PrincipalUpdate, error) {
	sqlStr := `select * from PrincipalUpdates where ProductID = ? order by Time asc`
	rows, err := SqlDB.Queryx(sqlStr, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	updates := models.NewPrincipalUpdateList()
	for rows.Next() {
		update := models.NewPrincipalUpdate()
		err := rows.StructScan(update)
		if err != nil {
			return nil, err
		}
		err = validate.Struct(update)
		if err != nil {
			return nil, err
		}
		updates = append(updates, update)
	}
	return updates, nil
}

func InsertOrderInterestList(orderInterestList []*models.OrderInterest) error {
	sqlStr := `insert into OrderInterest (OrderID, Time, APY, InterestGain) values (:OrderID, :Time, :APY, :InterestGain)`
	_, err := SqlDB.NamedExec(sqlStr, orderInterestList)
	return err
}

func GetSortedOrderInterestListUntilDate(orderID string, until string) ([]*models.OrderInterest, error) {
	sqlStr := `select * from OrderInterest where OrderID = ? and Time <= ? order by Time asc`
	rows, err := SqlDB.Queryx(sqlStr, orderID, until)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	interestList := models.NewOrderInterestList()
	for rows.Next() {
		interest := models.NewOrderInterest()
		err := rows.StructScan(interest)
		if err != nil {
			return nil, err
		}
		err = validate.Struct(interest)
		if err != nil {
			return nil, err
		}
		interestList = append(interestList, interest)
	}
	return interestList, nil
}

func UpdateOrderAccumulatedInterest(orderID string, accumulatedInterest float64) error {
	sqlStr := "update Orders set AccumulatedInterest = ? where OrderID = ?"
	_, err := SqlDB.Exec(sqlStr, accumulatedInterest, orderID)
	return err
}

func GetActiveOrdersProductIDs() ([]string, error) {
	var ids []string
	sqlStr := `select distinct ProductID from Orders where Type = 'Holding'`
	rows, err := SqlDB.Queryx(sqlStr)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		id := new(string)
		err := rows.Scan(id)
		if err != nil {
			logger.Warn(err)
			continue
		}
		ids = append(ids, *id)
	}
	return ids, nil
}

func UpdateActiveOrdersProductID(oldProductID string, newProductID string) error {
	sqlStr := `update Orders set ProductID = ? where ProductID = ? and Type = 'Holding'`
	_, err := SqlDB.Exec(sqlStr, newProductID, oldProductID)
	return err
}
