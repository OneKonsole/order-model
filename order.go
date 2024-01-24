package order

import (
	"database/sql"
)

type OrderFullInfos struct {
	AppOrder    Order              `json:"app_order"`
	PaypalOrder PaypalOrderDetails `json:"paypal_order"`
}

type Order struct {
	ID                int    `json:"id"`
	PaypalID          string `json:"paypal_id" validate:"required"`
	UserID            string `json:"user_id" validate:"required,uuid"`
	ClusterName       string `json:"cluster_name" validate:"required,min=1,max=63,isvalidclustername"`
	HasControlPlane   bool   `json:"has_control_plane"`
	HasMonitoring     bool   `json:"has_monitoring"`
	HasAlerting       bool   `json:"has_alerting"`
	ImageStorage      int    `json:"images_storage" validate:"required"`
	MonitoringStorage int    `json:"monitoring_storage" validate:"required_with=has_monitoring"`
}

type PaypalOrderDetails struct {
	ID            string         `json:"id"`
	Intent        string         `json:"intent"`
	Status        string         `json:"status"`
	PurchaseUnits []PurchaseUnit `json:"purchase_units"`
	CreateTime    string         `json:"create_time"`
}

type PurchaseUnit struct {
	Reference string `json:"reference_id"`
	Amount    Amount `json:"amount"`
	Payee     Payee  `json:"payee"`
}

type Amount struct {
	Currency string `json:"currency_code"`
	Value    string `json:"value"`
}

type Payee struct {
	Email    string `json:"email_address"`
	Merchant string `json:"merchant_id"`
}

func (o *Order) GetOrder(db *sql.DB) error {
	columns := "paypal_id, user_id, cluster_name, has_control_plane,has_monitoring,has_alerting,images_storage, monitoring_storage"

	return db.QueryRow("SELECT "+columns+" FROM orders WHERE id=$1", o.ID).
		Scan(&o.PaypalID,
			&o.UserID,
			&o.ClusterName,
			&o.HasControlPlane,
			&o.HasMonitoring,
			&o.HasAlerting,
			&o.ImageStorage,
			&o.MonitoringStorage)
}

func (o *Order) UpdateOrder(db *sql.DB) error {
	columns := "paypal_id=$1, user_id=$2, cluster_name=$3, has_control_plane=$4,has_monitoring=$5,has_alerting=$6,images_storage=$7, monitoring_storage=$8"

	_, err :=
		db.Exec("UPDATE orders SET "+columns+" WHERE id=$9",
			o.PaypalID, o.UserID, o.ClusterName, o.HasControlPlane, o.HasMonitoring, o.HasAlerting, o.ImageStorage, o.MonitoringStorage, o.ID)

	return err
}

func (o *Order) DeleteOrder(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM orders WHERE id=$1", o.ID)

	return err
}

func (o *Order) CreateOrder(db *sql.DB) error {

	err := db.QueryRow(
		"INSERT INTO orders(paypal_id, user_id, cluster_name, has_control_plane, has_monitoring, has_alerting, images_storage, monitoring_storage) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		o.PaypalID, o.UserID, o.ClusterName, o.HasControlPlane, o.HasMonitoring, o.HasAlerting, o.ImageStorage, o.MonitoringStorage).Scan(&o.ID)

	if err != nil {
		return err
	}

	return nil
}

func GetOrders(db *sql.DB, start int, count int, userID ...string) ([]Order, error) {

	queryString := "SELECT id, paypal_id, user_id, cluster_name, has_control_plane, has_monitoring, has_alerting, images_storage, monitoring_storage FROM orders "
	var args []interface{}

	if userID != nil {
		queryString += "WHERE user_id = $1 "
		args = append(args, userID[0])
		queryString += "LIMIT $2 OFFSET $3"
	} else {
		queryString += "LIMIT $1 OFFSET $2"
	}

	args = append(args, count, start)

	rows, err := db.Query(queryString, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := []Order{}

	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.PaypalID, &o.UserID, &o.ClusterName, &o.HasControlPlane, &o.HasMonitoring, &o.HasAlerting, &o.ImageStorage, &o.MonitoringStorage); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}
