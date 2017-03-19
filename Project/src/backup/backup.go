package backup

import (
	. "../def"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func WriteToBackup(file string, orders Orders) {
	ordersJson, errEncode := json.Marshal(orders)

	if errEncode != nil {
		fmt.Println("error encoding json: ", errEncode)
	}

	_, errOpen := os.Open(file)
	if errOpen != nil {
		fmt.Println("No file to write to, creating file...")
		_, _ = os.Create(file)
	}

	errWrite := ioutil.WriteFile(file, ordersJson, 0644)
	if errWrite != nil {
		fmt.Println("Error writing to file")
		log.Fatal(errWrite)
	}
}

func ReadFromBackup(file string) Orders {
	var orders Orders
	filename, errOpen := os.Open(file)
	if errOpen != nil {
		fmt.Println("No file to read from, creating file...")
		_, _ = os.Create(file)
		orders = generateEmptyOrders(file)
		WriteToBackup(file, orders)
		return orders
	}

	data := make([]byte, 1024)
	n, errRead := filename.Read(data)
	if errRead != nil {
		fmt.Println("Error reading from file")
		fmt.Println(errRead)
	}

	errDecode := json.Unmarshal(data[:n], &orders)
	if errDecode != nil {
		fmt.Println("Error decoding orders from backup")
	}
	return orders

}

func generateEmptyOrders(file string) Orders {
	var emptyOrders Orders
	for f := 0; f < NumFloors; f++ {
		for b := 0; b < NumTypes; b++ {
			emptyOrders[f][b] = false
		}
	}
	return emptyOrders
}
