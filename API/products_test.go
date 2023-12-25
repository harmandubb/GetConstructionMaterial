package api

// func resetTestDataBase() error {
// 	sqlString := "DELETE FROM products"

// 	p := ConnectToDataBase("mynewdatabase")

// 	tx, err := p.Begin(context.Background())
// 	if err != nil {
// 		return err
// 	}

// 	defer tx.Rollback(context.Background())

// 	_, err = tx.Exec(context.Background(), sqlString)
// 	if err != nil {
// 		return err
// 	}

// 	tx.Commit(context.Background())

// 	return nil

// }

// // func writeFileFromBytes(filePath string, data []byte) error {
// // 	// Write data to filePath using os.WriteFile
// // 	err := os.WriteFile(filePath, data, 0644)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	return nil
// // }

// func TestCheckDatabase(t *testing.T) {
// 	database := "mynewdatabase"
// 	result := CheckDataBase(database)

// 	if result != database {
// 		t.Errorf("Database = %s, but wanted %s", result, database)
// 	}
// }

// // func TestAddProductBasic(t *testing.T) {
// // 	name := "Meta Caulk Collar"
// // 	category := "Firestopping"
// // 	price := 10.01

// // 	resetTestDataBase()

// // 	AddProductBasic(name, category, price)

// // 	//Read the database to see if the action occured
// // 	sqlString := "SELECT * FROM products"
// // 	rows, err := dataBaseRead(sqlString)
// // 	if err != nil {
// // 		log.Fatalln(err)
// // 	}

// // 	got := Product{}
// // 	rows.Next()

// // 	p := reflect.ValueOf(&got).Elem()
// // 	numCols := p.NumField()
// // 	columns := make([]interface{}, numCols)
// // 	for i := 0; i < numCols; i++ {
// // 		field := p.Field(i)
// // 		columns[i] = field.Addr().Interface()
// // 	}
// // 	err = rows.Scan(columns...)
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}

// // 	if got.Name != name {
// // 		t.Errorf("Database name = %s; want %s", got.Name, name)
// // 	}

// // 	if got.Category != category {
// // 		t.Errorf("Database category = %s; want %s", got.Category, category)
// // 	}

// // 	if got.Price != price {
// // 		t.Errorf("Database price = %f; want %f", got.Price, price)
// // 	}
// // }

// func TestAddProductDataSheet(t *testing.T) {
// 	p := ConnectToDataBase("mynewdatabase")
// 	oidVal, err := AddProductDataSheet("Meta Caulk Collar", "./1.pdf", "mynewdatabase", p)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	err = getProductDataSheet(oidVal, "mynewdatabase", "./output.pdf")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	f1, err := os.Open("./1.pdf")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer f1.Close()

// 	f2, err := os.Open("./output.pdf")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer f2.Close()

// 	const chunkSize = 1024

// 	buf1 := make([]byte, chunkSize)
// 	buf2 := make([]byte, chunkSize)

// 	for {
// 		n1, err1 := f1.Read(buf1)
// 		n2, err2 := f2.Read(buf2)

// 		if !bytes.Equal(buf1[:n1], buf2[:n2]) {
// 			t.Error("Files are not the same")
// 		}
// 		// Check for errors
// 		if err1 != nil || err2 != nil {
// 			if err1 != err2 || err1 != io.EOF { // Different errors or not EOF
// 				log.Fatal(err1)
// 			}
// 			if err1 == io.EOF && err2 == io.EOF { // Both files ended together
// 				break
// 			}
// 		}
// 	}

// }

// func TestAddPorductPicture(t *testing.T) {
// 	oidVal, img_w, img_h, err := AddProductPicture("Meta Caulk Collar", "./img1.jpg", "mynewdatabase")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	err = getProductPicture(oidVal, img_w, img_h, "mynewdatabase", "./output")
// 	if err != nil {
// 		t.Error(err)
// 	}

// }

// func TestAddProduct(t *testing.T) {
// 	err := resetTestDataBase()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	dateFormatted, err := time.Parse("DateTime", time.Now().Format("DateTime"))
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	pInfo := ProductInfo{
// 		Date:      dateFormatted,
// 		Name:      "Hilti Fire Stop Collars",
// 		Category:  "Fire Stop",
// 		Price:     10.01,
// 		Currency:  "CAD",
// 		DataSheet: []string{"./Attachment/Technical-information-ASSET-DOC-LOC-1540917.pdf"},
// 	}

// 	err = addProduct(pInfo)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	//Need to read back the information to verify that it was correctly transmitted
// 	product, err := readDataBaseRow("products", pInfo.Name)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if product.Name != pInfo.Name {
// 		t.Errorf("Product Name Error")
// 	}
// 	if product.Created != pInfo.Date {
// 		t.Errorf("Created Error")
// 	}
// 	if product.Category != pInfo.Category {
// 		t.Error("Category Error")
// 	}
// 	if product.Price != pInfo.Price {
// 		t.Error("Price Error")
// 	}
// 	if product.Currency != pInfo.Currency {
// 		t.Error("Currency Error")
// 	}

// 	err = getProductDataSheet(*product.Data_Sheet, "mynewdatabase", "output.pdf")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }
