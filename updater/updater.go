package updater

import (
	"os"
	"strconv"
	"github.com/go-pg/pg"
	"time"
	"github.com/juju/errors"
	"github.com/sirupsen/logrus"
)

var (
	counter = 0

	countInStockUpdateDelay = 0
	totalProductsUpdateDelay = 0

	db  *pg.DB
	myLog *logrus.Logger
)

func StartUpdater() {
	var err error

	countInStockUpdateDelay, err = strconv.Atoi(os.Getenv("countInStockUpdateDelay"))
	if err != nil {
		panic(err)
	}

	totalProductsUpdateDelay, err = strconv.Atoi(os.Getenv("totalProductsUpdateDelay"))
	if err != nil {
		panic(err)
	}

	db = pg.Connect(&pg.Options{
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Addr: os.Getenv("DB_ADDRESS"),
	})

	myLog = logrus.New()
	myLog.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	myLog.SetReportCaller(false)


	for {
		if counter == totalProductsUpdateDelay {
			counter = 0
			totalProductsUpdate()
			time.Sleep(time.Duration(countInStockUpdateDelay) * time.Second)
		} else {
			counter++
			onlyCountInStockUpdate()
			time.Sleep(time.Duration(countInStockUpdateDelay) * time.Second)
		}

	}
}


func onlyCountInStockUpdate() {
	myLog.Info("only products in stock update begin")
	t := time.Now()

	_, pr, err := getAllProductsFromServer(false)
	if err != nil {
		myLog.Error(errors.Details(errors.Trace(err)))
		return
	}

	myLog.WithFields(logrus.Fields{
		"s": time.Now().Sub(t).Seconds(),
	}).Info("products were successfully parsed")

	var allProductsFromDB []*Product
	err = db.Model(&allProductsFromDB).Select()
	if err != nil {
		myLog.Error(errors.Details(errors.Trace(err)))
		return
	}

	for _, dbProduct := range allProductsFromDB {
		for i :=0; i < len(pr); i++ {
			parsedProduct := pr[i]

			if parsedProduct.Id == dbProduct.Id && parsedProduct.CountInStock == dbProduct.CountInStock {
				pr = pr[:i+copy(pr[i:], pr[i+1:])]
				break
			}
		}
	}

	for _, product := range pr {
		myLog.WithFields(logrus.Fields{
			"name": product.FullName,
			"count in stock": product.CountInStock,
		}).Info("new")
	}

	myLog.WithFields(logrus.Fields{
		"s": time.Now().Sub(t).Seconds(),
		"new products": len(pr),
	}).Info("products were successfully compared")

	for _, product := range pr {
		_, err = db.Model(product).Where("id = ?", product.Id).Update()
		if err != nil {
			myLog.Error(errors.Details(errors.Trace(err)))
			return
		}
	}

	myLog.WithFields(logrus.Fields{
		"s": time.Now().Sub(t).Seconds(),
	}).Info("only products in stock update finished")
}


func totalProductsUpdate() {
	myLog.Info("total products update begin")
	t := time.Now()

	pw, pr, err := getAllProductsFromServer(true)
	if err != nil {
		myLog.Error(errors.Details(errors.Trace(err)))
		return
	}

	myLog.WithFields(logrus.Fields{
		"s": time.Now().Sub(t).Seconds(),
	}).Info("products were successfully parsed")

	tx, err := db.Begin()
	if err != nil {
		myLog.Error(errors.Details(errors.Trace(err)))
		return
	}

	for _, product := range pr {
		product.LastUpdate = t

		exists, err := db.Model(&Product{}).Where("id = ?", product.Id).Exists()
		if err != nil {
			myLog.Error(errors.Details(errors.Trace(err)))
			return
		}

		if exists {
			_, err = tx.Model(product).
				Where("id = ?", product.Id).Update()
		} else {
			err = tx.Insert(product)
		}

		if err != nil {
			myLog.Error(errors.Details(errors.Trace(err)))
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		myLog.Error(errors.Details(errors.Trace(err)))
		return
	}

	myLog.WithFields(logrus.Fields{
		"s": time.Now().Sub(t).Seconds(),
	}).Info("products were successfully inserted")

	for _, wrapper := range pw {
		wrapper.LastUpdate = t

		exists, err := db.Model(&ProductWrapper{}).Where("name = ?", wrapper.Name).Exists()
		if err != nil {
			myLog.Error(errors.Details(errors.Trace(err)))
			return
		}

		myLog.WithFields(logrus.Fields{
			"name": wrapper.Name,
			"options": wrapper.Options,
		}).Info("inserting")

		if exists {
			_, err = db.Model(wrapper).
				Column("options").Column("children").
				Column("last_update").
				Where("name = ?", wrapper.Name).Update()
		} else {
			err = db.Insert(wrapper)
		}

		if err != nil {
			myLog.Error(errors.Details(errors.Trace(err)))
			return
		}
	}

	totalDeleted := 0

	res, err := db.Exec("delete from products where last_update < ? or last_update is null", t)
	if err != nil {
		myLog.Error(errors.Details(errors.Trace(err)))
		return
	}

	totalDeleted += res.RowsAffected()

	res, err = db.Exec("delete from product_wrappers where last_update < ? or last_update is null", t)
	if err != nil {
		myLog.Error(errors.Details(errors.Trace(err)))
		return
	}

	totalDeleted += res.RowsAffected()

	myLog.WithFields(logrus.Fields{
		"s": time.Now().Sub(t).Seconds(),
		"deleted": totalDeleted,
	}).Info("total products update was successfully finished")
}