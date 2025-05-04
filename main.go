package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Sale struct {
	InvoiceNo   string
	StockCode   string
	Description string
	Quantity    int
	InvoiceDate time.Time
	UnitPrice   float64
	CustomerID  string
	Country     string
	Total       float64
}

func main() {
	// Открываем CSV файл с данными о продажах
	file, err := os.Open("data.csv")
	if err != nil {
		log.Fatal("Ошибка при открытии файла:", err)
	}
	defer file.Close()

	// Читаем данные из CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Ошибка при чтении CSV:", err)
	}

	// Парсим данные в структуру Sale
	var sales []Sale
	for i, record := range records {
		if i == 0 { // Исправлено: добавлено 0
			continue // Пропускаем заголовок
		}

		// Пропускаем пустые строки
		if len(record) < 8 {
			continue
		}

		quantity, err := strconv.Atoi(record[3])
		if err != nil {
			log.Printf("Ошибка при парсинге количества: %v", err)
			continue
		}

		// Пропускаем возвраты (отрицательное количество)
		if quantity <= 0 {
			continue
		}

		unitPrice, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			log.Printf("Ошибка при парсинге цены: %v", err)
			continue
		}

		// Пропускаем нулевые цены
		if unitPrice <= 0 {
			continue
		}

		// Парсим дату (предполагаем формат "MM/DD/YYYY HH:MM")
		date, err := time.Parse("1/2/2006 15:04", record[4])
		if err != nil {
			log.Printf("Ошибка при парсинге даты %s: %v", record[4], err)
			continue
		}

		total := float64(quantity) * unitPrice

		sale := Sale{
			InvoiceNo:   record[0],
			StockCode:   record[1],
			Description: record[2],
			Quantity:    quantity,
			InvoiceDate: date,
			UnitPrice:   unitPrice,
			CustomerID:  record[6],
			Country:     record[7],
			Total:       total,
		}
		sales = append(sales, sale)
	}

	// Анализ данных
	totalTransactions := len(sales)
	totalRevenue := 0.0
	countryRevenue := make(map[string]float64)
	productSales := make(map[string]int) // Количество продаж по товарам

	for _, sale := range sales {
		totalRevenue += sale.Total
		countryRevenue[sale.Country] += sale.Total
		productSales[sale.Description] += sale.Quantity
	}

	// Выводим результаты анализа
	fmt.Printf("Всего транзакций: %d\n", totalTransactions)
	fmt.Printf("Общий доход: $%.2f\n", totalRevenue)

	// Выводим доход по странам (только ненулевые)
	fmt.Println("\nДоход по странам:")
	for country, revenue := range countryRevenue {
		if revenue > 0 {
			fmt.Printf("- %s: $%.2f\n", country, revenue)
		}
	}

	// Находим самый популярный товар (с ненулевыми продажами)
	popularProduct := ""
	maxQuantity := 0
	for product, quantity := range productSales {
		if quantity > maxQuantity {
			maxQuantity = quantity
			popularProduct = product
		}
	}

	if maxQuantity > 0 {
		fmt.Printf("\nСамый популярный товар: %s (продано %d единиц)\n", popularProduct, maxQuantity)
	} else {
		fmt.Println("\nНет данных о продажах товаров")
	}

	// Дополнительная статистика
	// Находим самого ценного клиента (с ненулевыми тратами)
	clientSpending := make(map[string]float64)
	for _, sale := range sales {
		if sale.CustomerID != "" && sale.Total > 0 {
			clientSpending[sale.CustomerID] += sale.Total
		}
	}

	topClient := ""
	maxSpending := 0.0
	for client, spending := range clientSpending {
		if spending > maxSpending {
			maxSpending = spending
			topClient = client
		}
	}

	if maxSpending > 0 {
		fmt.Printf("\nСамый ценный клиент: %s (потратил $%.2f)\n", topClient, maxSpending)
	} else {
		fmt.Println("\nНет данных о клиентах")
	}
}
