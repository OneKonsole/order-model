package order

import (
	"database/sql"
)

type Order struct {
	ID                int    `json:"id"`
	UserID            string `json:"user_id" validate:"required,uuid"`
	ClusterName       string `json:"cluster_name" validate:"required,min=1,max=63,isvalidclustername"`
	HasControlPlane   bool   `json:"has_control_plane"`
	HasMonitoring     bool   `json:"has_monitoring"`
	HasAlerting       bool   `json:"has_alerting"`
	ImageStorage      int    `json:"images_storage" validate:"required"`
	MonitoringStorage int    `json:"monitoring_storage" validate:"required_with=has_monitoring"`
}

func (o *Order) GetOrder(db *sql.DB) error {
	columns := "user_id, cluster_name, has_control_plane,has_monitoring,has_alerting,images_storage, monitoring_storage"

	return db.QueryRow("SELECT "+columns+" FROM orders WHERE id=$1", o.ID).
		Scan(&o.UserID,
			&o.ClusterName,
			&o.HasControlPlane,
			&o.HasMonitoring,
			&o.HasAlerting,
			&o.ImageStorage,
			&o.MonitoringStorage)
}

func (o *Order) UpdateOrder(db *sql.DB) error {
	columns := "user_id=$1, cluster_name=$2, has_control_plane=$3,has_monitoring=$4,has_alerting=$5,images_storage=$6, monitoring_storage=$7"

	_, err :=
		db.Exec("UPDATE orders SET "+columns+" WHERE id=$8",
			o.UserID, o.ClusterName, o.HasControlPlane, o.HasMonitoring, o.HasAlerting, o.ImageStorage, o.MonitoringStorage, o.ID)

	return err
}

func (o *Order) DeleteOrder(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM orders WHERE id=$1", o.ID)

	return err
}

func (o *Order) CreateOrder(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO orders(user_id, cluster_name, has_control_plane, has_monitoring, has_alerting, images_storage, monitoring_storage) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		o.UserID, o.ClusterName, o.HasControlPlane, o.HasMonitoring, o.HasAlerting, o.ImageStorage, o.MonitoringStorage).Scan(&o.ID)

	if err != nil {
		return err
	}

	return nil
}

func GetOrders(db *sql.DB, start, count int) ([]Order, error) {
	rows, err := db.Query(
		"SELECT id, user_id,cluster_name,has_control_plane,has_monitoring,has_alerting,images_storage,monitoring_storage FROM orders LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := []Order{}

	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.ClusterName, &o.HasControlPlane, &o.HasMonitoring, &o.HasAlerting, &o.ImageStorage, &o.MonitoringStorage); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}
