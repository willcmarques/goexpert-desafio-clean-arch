package database

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/willcmarques/goexpert-desafio-clean-arch/internal/entity"

	// sqlite3
	_ "github.com/mattn/go-sqlite3"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRepositoryTestSuite) setup() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.NoError(err)
	db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	suite.Db = db
}

func (suite *OrderRepositoryTestSuite) SetupSuite() {
	suite.setup()
}

func (suite *OrderRepositoryTestSuite) SetupTest() {
	suite.setup()
}

func (suite *OrderRepositoryTestSuite) TearDownTest() {
	suite.Db.Close()
}

func (suite *OrderRepositoryTestSuite) TearDownSuite() {
	suite.Db.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (suite *OrderRepositoryTestSuite) TestGivenAnOrder_WhenSave_ThenShouldSaveOrder() {
	order, err := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())
	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order)
	suite.NoError(err)

	var orderResult entity.Order
	err = suite.Db.QueryRow("Select id, price, tax, final_price from orders where id = ?", order.ID).
		Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)

	suite.NoError(err)
	suite.Equal(order.ID, orderResult.ID)
	suite.Equal(order.Price, orderResult.Price)
	suite.Equal(order.Tax, orderResult.Tax)
	suite.Equal(order.FinalPrice, orderResult.FinalPrice)
}

func (suite *OrderRepositoryTestSuite) TestShouldReturnAllOrders() {
	stmt, err := suite.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	suite.NoError(err)
	defer stmt.Close()
	stmt.Exec("1", 10.0, 0.1, 10.1)
	stmt.Exec("2", 20.0, 0.2, 20.2)
	stmt.Exec("3", 30.0, 0.3, 30.3)

	repo := NewOrderRepository(suite.Db)
	orders, err := repo.List()
	suite.NoError(err)
	suite.Equal(3, len(orders))
}
