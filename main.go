package main

import (
	"fmt"
	"encoding/json"
	"github.com/jung-kurt/gofpdf"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
	"bytes"
	"os"
	gomail "gopkg.in/gomail.v2"
	
)

type pitch struct {
	CompanyName        string `json:"companyName"`
	Email              string `json:"email"`
	CustomerSegments   string `json:"customerSegments"`
	ProblemOrNeed      string `json:"problemOrNeed"`
	ProductDescription string `json:"productDescription"`
	Competitors        string `json:"competitors"`
	Differentiation    string `json:"differentiation"`
	BirthDate          string `json:"birthDate"`
	Opportunities      string `json:"opportunities"`
	StageOfDevelopment string `json:"stageOfDevelopment"`
	CurrentNeeds       string `json:"currentNeeds"`
	Stakeholders       string `json:"stakeholders"`
	TakeOffDate		   string `json:"takeOffDate"`
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("public/*"))	
}

type info struct {
	Name string
}

func (i info) sendMail(email string) {


	var tplemail bytes.Buffer
	if err := tpl.Execute(&tplemail, i); err != nil {
		log.Println(err)
	}

	result := tplemail.String()
	m := gomail.NewMessage()
	m.SetHeader("From", "admin@thestartuptools.com")
	m.SetHeader("To", email)
	m.SetAddressHeader("Cc", "thestartuptools@gmail.com", "Dr. Cedric Aimal Edwin")
	m.SetHeader("Subject", "One Minute Elevator Pitch")
	m.SetBody("text/html", result)
	m.Attach("pdf/" + i.Name + "_pitch.pdf")// attach whatever you want

	d := gomail.NewDialer("smtp.gmail.com", 587, "thestartuptools@gmail.com", "Helloworld1")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func submit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	
	// Decode json data from request
	d := json.NewDecoder(r.Body)
	xp := pitch{}
	err := d.Decode(&xp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Pdf Generator Code
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTopMargin(40)
	pdf.SetHeaderFuncMode(func() {
		pdf.Image("logo.png", 10, 6, 30, 0, false, "", 0, "")
		pdf.SetY(5)
		pdf.SetFont("Arial", "B", 15)
		pdf.CellFormat(0, 35, "One Minute Elevator Pitch", "B", 0, "C", false, 0, "")
		pdf.Ln(20)
	}, true)
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(200, 10, "Copyrights 2018",
			"", 0, "C", false, 0, "")
	})
	pdf.AliasNbPages("")
	pdf.AddPage()
	pdf.SetFont("Times", "", 12)
	
	pdf.Ln(20)
	// Company Name
	pdf.CellFormat(0, 10, fmt.Sprintf("Company : %s", xp.CompanyName),
			"", 1, "", false, 0, "")
	// Date
	pdf.CellFormat(0, 10, fmt.Sprintf("Date : %s", xp.CompanyName),"", 1, "", false, 0, "")
	pdf.Ln(5);
	//Insert Image Here

	// For (Identify customers segments that use or will use your offerings) who are 
	// dissatisfied with [or demand] (explain the compelling problem or need you are 
	// attempting to address), our company (insert name from above) provides 
	// (briefly describe your product or service offering in terms of a solution to your 
	// customers’ compelling needs). Unlike competitor products [or services] and alternatives, 
	// such as (name key competitors and alternatives), we offer superior value because 
	// (explain what differentiates your business from competitors and substitutes). 
	// The business was funded (When was the business established?) and plan to peruse 
	// the following opportunity: (list key opportunities you wish to pursue or are already pursing). 
	// We have currently completed the following: (describe the stage of development, including whether 
	// there is a working prototype and whether it has been validated by customers). 
	// We currently need (explain the resources needed, including money, talent, etc.). 
	// We are projecting to deliver the following benefits: (identify key stakeholders, for example, 
	// customers, suppliers, owners, partners, and then describe the tangible business value you 
	// will create) and expect to begin delivering that value by (when do you expect to take off, 
	// provide date). 
	

	//  For (Identify customers segments that use or will use your offerings)

	pdf.MultiCell(0, 10, string(fmt.Sprintf(`For %s who are dissatisfied with %s,our company %s provides %s. Unlike competitor products and alternatives, such as %s, we offer superior value because %s. The business was funded %s and plan to persue the following oppurtunity: %s. We have currently compeleted the following: %s. We currently need %s. We are projecting to deliver the following benefits: %s and expect to begin delivering that value by %s.`, 
									   	xp.CustomerSegments, xp.ProblemOrNeed,
										xp.CompanyName, xp.ProductDescription,
										xp.Competitors, xp.Differentiation,
										xp.BirthDate, xp.Opportunities,
										xp.StageOfDevelopment, xp.CurrentNeeds,
										xp.Stakeholders, xp.TakeOffDate)),
	"", "", false)

	pdf.Ln(15);
	pdf.CellFormat(0, 10, "Signature : Dr.Cedric Aimal Edwin","", 1, "", false, 0, "")
	
	err = pdf.OutputFileAndClose(fmt.Sprintf("pdf/%s_pitch.pdf", xp.CompanyName))
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalln(err)
	}
	
	

	pitchemail := info{xp.CompanyName}
	go pitchemail.sendMail(xp.Email)
	fmt.Printf("One minute pitch send to company: %s at email: %s \n", xp.CompanyName, xp.Email)

	// Saving a record in file for future references
	go fileRecord(xp.CompanyName, xp.Email)

}

func fileRecord(companyName string, email string) {
	f, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte("One minute pitch send to company: " + companyName + " at email:" + email + "\n")); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
func index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    err := tpl.ExecuteTemplate(w, "index.html","")
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalln(err)
	}
}
func main() {

	router := httprouter.New()
    // router.GET("/", Index)
	router.GET("/", index)
	router.ServeFiles("/assets/*filepath", http.Dir("assets"))
	// http.Handle("/", http.FileServer(http.Dir("assets/")))

	router.POST("/submit", submit)

    log.Fatal(http.ListenAndServe(":80", router))
	
}
