package order

import (
	"database/sql"
)

type Order struct {
	ID              int    `json:"id"`
	UserID          string `json:"user_id"`
	ClusterName     string `json:"cluster_name"`
	HasControlPlane bool   `json:"has_control_plane"`
	HasMonitoring   bool   `json:"has_monitoring"`
	HasAlerting     bool   `json:"has_alerting"`
	StorageSize     int    `json:"storage_size"`
}

func (o *Order) GetOrder(db *sql.DB) error {
	columns := "user_id, cluster_name, has_control_plane,has_monitoring,has_alerting,storage_size"

	return db.QueryRow("SELECT "+columns+" FROM orders WHERE id=$1", o.ID).
		Scan(&o.UserID, &o.ClusterName, &o.HasControlPlane, &o.HasMonitoring, &o.HasAlerting, &o.StorageSize)
}

func (o *Order) UpdateOrder(db *sql.DB) error {
	columns := "user_id=$1, cluster_name=$2, has_control_plane=$3,has_monitoring=$4,has_alerting=$5,storage_size=$6"

	_, err :=
		db.Exec("UPDATE orders SET "+columns+" WHERE id=$7",
			o.UserID, o.ClusterName, o.HasControlPlane, o.HasMonitoring, o.HasAlerting, o.StorageSize, o.ID)

	return err
}

func (o *Order) DeleteOrder(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM orders WHERE id=$1", o.ID)

	return err
}

func (o *Order) CreateOrder(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO orders(user_id, cluster_name, has_control_plane, has_monitoring, has_alerting, storage_size) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		o.UserID, o.ClusterName, o.HasControlPlane, o.HasMonitoring, o.HasAlerting, o.StorageSize).Scan(&o.ID)

	if err != nil {
		return err
	}

	return nil
}

func GetOrders(db *sql.DB, start, count int) ([]Order, error) {
	rows, err := db.Query(
		"SELECT id, user_id,cluster_name,has_control_plane,has_monitoring,has_alerting,storage_size FROM orders LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := []Order{}

	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.ClusterName, &o.HasControlPlane, &o.HasMonitoring, &o.HasAlerting, &o.StorageSize); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}
