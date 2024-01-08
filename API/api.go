package api

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	d "docstruction/getconstructionmaterial/Database"
	g "docstruction/getconstructionmaterial/GCalls"
	gpt "docstruction/getconstructionmaterial/GPT"
	w "docstruction/getconstructionmaterial/Websites"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/api/gmail/v1"
)

// How many emails should be sent out for inquiries
const SUPPLIERCONTACTLIMIT = 3

// Currency constants
const (
	CADTOUSD float64 = 0.75
	USTTOCAD float64 = 1.33
)

// Purpose: Run code that will handle a customer inquiry once it is in a database
// Parameters:
// inquiryID string --> ID that is generated to identify the inquiry in a database
// catigorizationTemplate string --> the template that is used to feed gpt a catigorization prompt
// emailTemplate string --> template to feed gpt to get an email written
// Return:
// Error if present
func ProcessCustomerInquiry(p *pgxpool.Pool, inquiryID, catigorizationTemplate, emailTemplate string) (err error) {
	// use the inquiry id to get the row of information in the database
	custInquiry, err := d.ReadCustomerInquiry(p, os.Getenv("CUSTOMER_INQUIRY_TABLE"), inquiryID)
	if err != nil {
		log.Printf("Error in Reading Customer Inquiry Table: %v", err)
		return err
	}

	matForm := g.MaterialFormInfo{
		Email:    custInquiry.Email,
		Material: custInquiry.Material,
		Loc:      custInquiry.Loc,
	}

	srv := g.ConnectToGmailAPI()

	suppInfo, err := ContactSupplierForMaterial(srv, matForm, catigorizationTemplate, emailTemplate)
	if err != nil {
		log.Printf("Error in Contacting Suppliers: %v", err)
		return err
	}

	var emails []string

	for _, s := range suppInfo {
		emails = append(emails, s.Email[0])
		err = d.AddBlankEmailInquiryEntry(p, inquiryID, matForm.Email, matForm.Material, s, true, "emails")
		if err != nil {
			log.Printf("Error in Addting a Email Sent to Table: %v", err)
			fmt.Println("Failed to added email sent into database. ")
		}
	}

	AlertAdmin(srv, matForm, emails)

	return nil

}

func ConcurrentProcessCustomerInquiry(wg *sync.WaitGroup, errStream chan<- error, srv *gmail.Service, p *pgxpool.Pool, inquiryID, catigorizationTemplate, emailTemplate string) {
	defer wg.Done()
	// use the inquiry id to get the row of information in the database
	custInquiry, err := d.ReadCustomerInquiry(p, os.Getenv("CUSTOMER_INQUIRY_TABLE"), inquiryID)
	if err != nil {
		log.Printf("Error in Reading Customer Inquiry Table: %v", err)
		errStream <- err
		return
	}

	matForm := g.MaterialFormInfo{
		Email:    custInquiry.Email,
		Material: custInquiry.Material,
		Loc:      custInquiry.Loc,
	}

	suppInfo, err := ContactSupplierForMaterial(srv, matForm, catigorizationTemplate, emailTemplate)
	if err != nil {
		log.Printf("Error in Contacting Suppliers: %v", err)
		errStream <- err
		return
	}

	var emails []string

	for _, s := range suppInfo {
		emails = append(emails, s.Email[0])
		err = d.AddBlankEmailInquiryEntry(p, inquiryID, matForm.Email, matForm.Material, s, true, "emails")
		if err != nil {
			log.Printf("Error in Addting a Email Sent to Table: %v", err)
			fmt.Println("Failed to added email sent into database. ")
			errStream <- err
			return
		}
	}

	AlertAdmin(srv, matForm, emails)
	fmt.Println("Done Work")
}

// Purpose: Send en email to the admin about th erequest that has come in and what supplier (emails) are used to contact for th erequest
// Parameters:
// srv *gmail.Service --> Service used to access gmail
// matinfo g.MaterialFormINfo --> material infromaiton that is submitted to the website by client
// emailsSentTo []string --> emails list that tells where the inquiry emails are sent to (suppliers)
// Error is any present

func AlertAdmin(srv *gmail.Service, matInfo g.MaterialFormInfo, emailsSentTo []string) error {

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }

	adminEmail := os.Getenv("ADMIN_EMAIL")

	subj := fmt.Sprintf("Docstruction Notificaiton: %s", matInfo.Material)

	msg := fmt.Sprintf("Inquiry from: %s\nInquiry material: %s\nInquiry Location: %s\n Emailed Suppliers:\n", matInfo.Email, matInfo.Material, matInfo.Loc)

	for _, email := range emailsSentTo {
		msg = msg + "- " + email + "\n"
	}

	_, err := g.SendEmail(srv, subj, msg, adminEmail)
	if err != nil {
		return err
	}

	return nil
}

// Purpose: Execute logic that takes the material info from form and sends out emails to supplier
// Parameters:
// MatInfo g.MaterialFromInfo --> Struct that carried the information in the form. (material name and user request email)
// catigoorizationTemplate string --> Pathway to the tempalte dues for the gpt promp maker
// emailTemplate string --> Pathway to the template used for the gpt email prompt maker
// loc *mapts.LatLng --> Google maps struct for holding the llat and lng for the place the user is requesting from.
// Return:
// error if any present
func ContactSupplierForMaterial(srv *gmail.Service, matInfo g.MaterialFormInfo, catigorizationTemplate, emailTemplate string) (emailsSentto []g.SupplierEmailInfo, err error) {
	//Call chat gpt to catigorized the item

	fmt.Printf("Inputted material form info:")
	fmt.Println(matInfo)

	fmt.Println(matInfo.Email)
	fmt.Println(matInfo.Loc)
	fmt.Println(matInfo.Material)

	catergory, err := gpt.PromptGPTMaterialCatogorization(catigorizationTemplate, matInfo.Material)
	if err != nil {
		log.Fatalf("Catogirization Error: %v", err)
		time.Sleep(5 * time.Second)
		return []g.SupplierEmailInfo{}, err
	}

	// Search for near by supplies for the category
	c, err := g.GetMapsClient()
	if err != nil {
		log.Fatalf("Map Client Connection Error: %v", err)
		time.Sleep(5 * time.Second)

		return []g.SupplierEmailInfo{}, err
	}

	//Get Lat and lng coordinates
	geometry, _, err := g.GeocodeGeneralLocation(c, matInfo.Loc)
	if err != nil {
		log.Fatalf("Geocoding Converstion Error: %v", err)
		time.Sleep(5 * time.Second)

		return []g.SupplierEmailInfo{}, err
	}

	searchResp, err := g.SearchSuppliers(c, catergory, &geometry.Location)
	if err != nil {
		log.Fatalf("Map Search Supplier Error: %v", err)
		time.Sleep(5 * time.Second)

		return []g.SupplierEmailInfo{}, err
	}

	var supplierInfo []g.SupplierEmailInfo

	for _, supplier := range searchResp.Results {
		supplier, _ := g.GetSupplierInfo(c, supplier)

		supplierInfo = append(supplierInfo, supplier)
	}

	//Get the supplier emails from the info that is found
	var filteredSuppliers []g.SupplierEmailInfo // Assuming SupplierInfo is the type of your slice elements

	counter := 0
	for _, supInfo := range supplierInfo {
		if counter < SUPPLIERCONTACTLIMIT+2 {
			email, err := w.FindSupplierContactEmail(supInfo.Website)
			if err != nil {
				// log.Printf("Supplier Email Get Error: %v", err)  // Log the error, but don't stop the entire process
				continue // Skip this supplier and continue with the next one
			} else {
				supInfo.Email = email
				filteredSuppliers = append(filteredSuppliers, supInfo) // Add to the new slice
			}
		} else {
			break
		}
	}

	fmt.Println("FilteredSuppliers: ", len(filteredSuppliers))

	supplierInfo = nil //Setting to nil so the memory allocatin is lower.

	counter = 0

	var emailsSentTo []g.SupplierEmailInfo

	for _, supInfo := range filteredSuppliers {
		if counter < SUPPLIERCONTACTLIMIT {
			if len(supInfo.Email) != 0 {
				// get the email prompt from chat gpt
				if w.IsValidEmail(supInfo.Email[0]) {
					subj, body, err := gpt.CreateEmailToSupplier(emailTemplate, supInfo.Name, matInfo.Material)
					if err != nil {
						fmt.Printf("GPT Email Create Error: %v\n", err)
						continue //What should happen if this error occurs?
					} else {

						// send the emal to the supplier
						// msgDetails, err = g.SendEmail(srv, subj, body, supInfo.Email[0])
						msgDetails, err := g.SendEmail(srv, subj, body, "harmand1999@gmail.com")
						if err == nil {
							//email sent successfully
							supInfo.Email_ThreadID = msgDetails.ThreadId
							emailsSentTo = append(emailsSentTo, supInfo)
							counter = counter + 1
						} else {
							fmt.Println("Supplier email was not able to send")
							continue
						}
					}
				}
			}
		} else {
			break
		}
	}

	return emailsSentTo, nil
}

// Purpose: Provides a simple function that is called to refresh the push notification service atleast once perday
// Return:
// Error if present
func RefreshPushNotificationWatch() (err error) {
	srv := g.ConnectToGmailAPI()
	err = g.WatchPushNotification(srv)
	if err != nil {
		return err
	}
	fmt.Println("Refreshed Push Notification")
	return nil
}

func RefreshPushNotificationWatchAfterDay() {
	ticker := time.NewTicker(24 * time.Hour)

	go func() {
		for range ticker.C {
			err := RefreshPushNotificationWatch()
			if err != nil {
				log.Fatalf("Error refreshing the push notification alert: %v", err)
			}
		}
	}()

}

// Purpose: go thorugh all unread emails and process them to see what the outcome of the exchange was and update the database accordingly
// Parameters:
// srv *gmail.Service --> pointer to the establish gmail api service
// user string --> user email that you want to check the unread messages of
// Return:
// Error if present
func AddressPushNotification(p *pgxpool.Pool, srv *gmail.Service, user, receiveAnalysisTemplate, emailInquiryTableName, customerInquiryTableName string) (err error) {
	messages, err := g.GetUnreadMessagesData(srv, user)
	if err != nil {
		fmt.Printf("Retrive error for unread Messages: %v\n", err)
		return err
	}

	// implement concourrency tools here
	// Make a different thread read the different unread messages
	// Create the gpt prompt and send to gpt
	// break the gpt response into what is in the email.
	// fill the table as needed

	for _, message := range messages.Messages {
		// need to make a function (concurrent) that runs through the unread emails to do all of the needed tasks
		emailInfo, _, err := g.GetMessage(srv, message, user)
		if err != nil {
			fmt.Printf("Specific Message retrive error: %v\n", err)
			return err
		}

		// Need to analize the the email body in chat gpt to see what I should do next
		presentInfo, err := gpt.PromptGPTReceiveEmailAnalysis(receiveAnalysisTemplate, emailInfo.Body)
		fmt.Println("Present Info:", presentInfo)
		if err != nil {
			fmt.Printf("Incoming Message Analysis error: %v\n", err)
			return err
		}

		sup_thread_id := message.ThreadId

		if presentInfo.Present {
			id_opt := d.IDOption{
				Thread_ID: sup_thread_id,
			}

			// find the threadID and use it to pull out the row of information in the inquiry table
			emailInquiry, err := d.ReadEmailInquiryEntry(p, emailInquiryTableName, id_opt)
			if err != nil {
				fmt.Printf("Error Reading database email inquiry entry: -->\n")
				fmt.Println(err)
				return err
			}

			emailInquiry.Present = presentInfo.Present
			emailInquiry.Price = presentInfo.Price
			emailInquiry.Currency = presentInfo.Currency
			emailInquiry.Data_Sheet = nil
			if presentInfo.Data_Sheet {
				// TODO: get the attachement from the email
				placeholder := []byte("DataSheetPlaceHolder")
				emailInquiry.Data_Sheet = &placeholder
			}

			fmt.Println("Updating the email inquiry entery")

			err = d.UpdateEmailInquiryEntryMaterialPresent(p, emailInquiry.Inquiry_ID, emailInquiryTableName, emailInquiry.Price, emailInquiry.Currency, emailInquiry.Data_Sheet)
			if err != nil {
				fmt.Printf("Error updating the material table: %v\n", err)
				return err
			}

			fmt.Println("Reading customer inquiry entery")

			// If present then the data in the main customer inquiry should be compared and changed if this inquiry is better
			// Read the customer inquiry row into
			custInquiry, err := d.ReadCustomerInquiry(p, customerInquiryTableName, emailInquiry.Inquiry_ID)
			if err != nil {
				fmt.Printf("Error reading Customer inquiry table: %v\n", err)
				return err
			}

			fmt.Println("Comparing the incoming inquiry to what is present")

			// TODO: Implement a currency compare mechanism generally.
			result, err := incomingInquiryBetter(custInquiry, emailInquiry)
			if err != nil {
				fmt.Printf("Error deciding inquiry quality: %v\n", err)
				return err
			}

			var currency string

			fmt.Println("Checking what needs to be updated")

			if result {
				if emailInquiry.Currency != "" {
					currency = emailInquiry.Currency
				} else {
					currency = custInquiry.Currency
				}

				// More competative item has come in therefor update the customer inquiry table
				err = d.UpdateCustomerInquiryMaterial(p, customerInquiryTableName, emailInquiry.Inquiry_ID, sup_thread_id, emailInquiry.Price, currency, emailInquiry.Data_Sheet)
				if err != nil {
					fmt.Printf("Error updating customer inquiry table: %v\n", err)
					return err
				}

			}
		}

		err = g.MarkEmailAsRead(srv, user, message.Id)
		if err != nil {
			fmt.Printf("Error marking email as read: %v\n", err)
			return err
		}

		fmt.Println("Marked the email as read and working ")
	}

	return nil
}

// Purpose: Account for currency difference when comparing the pricing of material. As of right now only the american and canadian currencies will be compared.
// Parameters:
// present Inquiry, incoming Inquiry d.CustomerInquiry --> information that is releated to what is present in the database and what is being given by a new supplier.
// Return:
// Result bool --> true if the incoming inquiry is better
// Error if present

func incomingInquiryBetter(currentInquiry d.CustomerInquiry, incomingInquiry d.EmailInquiries) (bool, error) {
	// 3 Paths
	// 1. Nothing in the present Inquiry --> just accept the incoming
	// 2. Incoming better --> return true after currency check
	// 3. Curent better --> return false

	// 1. Check if nothing is in the present checking the present field is the easiest
	if currentInquiry.Present {
		//Check if there is a currency difference
		if currentInquiry.Currency != incomingInquiry.Currency {
			// Need to perfrom a currency change
			if currentInquiry.Currency == "CAD" && incomingInquiry.Currency == "USD" {
				if currentInquiry.Price*CADTOUSD < incomingInquiry.Price {
					return false, nil
				} else {
					return true, nil
				}
			} else if currentInquiry.Currency == "USD" && incomingInquiry.Currency == "CAD" {
				if currentInquiry.Price < incomingInquiry.Price*CADTOUSD {
					return false, nil
				} else {
					return true, nil
				}
			} else {
				return false, errors.New("currency exchange is not present")
			}
		} else {
			// Don't need to perform a currency change
			if incomingInquiry.Price < currentInquiry.Price {
				return true, nil
			} else {
				return false, nil
			}
		}
	} else {
		return true, nil
	}
}
