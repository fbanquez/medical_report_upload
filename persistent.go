package main

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Structure added only for the Metropolitan Clinic to get patient information in its CRM.
type crmField struct {
	//
	PatientAddress   string `json:"patient_address"`
	PatientPhone     string `json:"patient_phone"`
	PatientCellphone string `json:"patient_cellphone"`
}

type control struct {
	//
	Active    bool      `json:"active"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedBy string    `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

type complexField struct {
	//
	Mime    string `json:"mime"`
	Content string `json:"content"`
}

type innerReport struct {
	//
	Header    complexField `json:"header"`
	Body      complexField `json:"body"`
	Signature complexField `json:"signature"`
	Footer    complexField `json:"footer"`
}

// HealthReport struct
type HealthReport struct {
	//
	AccessionNumber    string      `json:"accession_number"`
	Institution        string      `json:"institution"`
	PatientId          string      `json:"patient_identification"`
	PatientName        string      `json:"patient_name"`
	PatientGender      string      `json:"patient_gender"`
	PatientBirth       string      `json:"patient_birth"`
	PatientEmail       string      `json:"patient_email"`
	PhysicianName      string      `json:"physician_name"`
	PhysicianSpecialty string      `json:"physician_specialty"`
	ReferInstitution   string      `json:"referred_institution"`
	ReferEmailInst     string      `json:"email_ref_institution"`
	ReferPhys          string      `json:"referred_physician"`
	ReferEmailPhys     string      `json:"email_ref_physician"`
	DateReport         string      `json:"date_report"`
	Attached           []string    `json:"attached"`
	Report             innerReport `json:"report"`
	StudyDate          string      `json:"study_date"`
	StudyName          string      `json:"study_name"`
	StudyType          string      `json:"study_type"`
	Tag1               string      `json:"tag1"`
	Tag2               string      `json:"tag2"`
	Tag3               crmField    `json:"tag3"`
	Control            control     `json:"control"`
}

//var db *sql.DB

// connectDBMedicalReports non-exportable function used to get a connection instance to
// local storage medical reports.
func connectDBMedicalReports() (db *sql.DB, err error) {
	//
	db, err = sql.Open("sqlite3", "./persistent/medicalReports")
	if err != nil {
		Error.Println("There were problems trying to establish a connection to the medical reports local storage . Validate path and repository. ", err)
	}

	return
}

// connectDBMedicalProcedures non-exportable function used to get a connection instance to
// local storage medical reports.
func connectDBMedicalProcedures() (db *sql.DB, err error) {
	//
	db, err = sql.Open("sqlite3", "./persistent/medicalProcedures")
	if err != nil {
		Error.Println("There were problems trying to establish a connection to the medical procedures local storage. Validate path and repository. ", err)
	}

	return
}

// connectDBCodeImages non-exportable function used to get a connection instance to
// local storage coded images.
func connectDBCodeImages() (db *sql.DB, err error) {
	//
	db, err = sql.Open("sqlite3", "./persistent/codeImages")
	if err != nil {
		Error.Println("There were problems trying to establish a connection to the code images local storage. Validate path and repository. ", err)
	}

	return
}

// checkTableReports non-exportable function used to determine wether or not the named table exists
// into Medical Reports local storage.
func checkTableReports(tableName string) (existsTable bool, err error) {
	//
	db, err := connectDBMedicalReports()
	if err != nil {
		Error.Println("Problems validating the existence of the table into local storage. ", err)
	}
	defer db.Close()

	statement := "SELECT CAST(COUNT(*) AS BIT) FROM sqlite_master WHERE type = 'table' AND name = $1"
	stmt, err := db.Prepare(statement)
	if err != nil {
		Error.Println("Problems creating the statement that validates the existence of local storage. ", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(tableName).Scan(&existsTable)
	if err != nil {
		Error.Println("Could not find the tables within local storage. ", err)
	}

	return
}

// checkTableProcedures non-exportable function used to determine wether or not the named table exists
// into Medical Procedures local storage.
func checkTableProcedures(tableName string) (existsTable bool, err error) {
	//
	db, err := connectDBMedicalProcedures()
	if err != nil {
		Error.Println("Problems validating the existence of the table into local storage. ", err)
	}
	defer db.Close()

	statement := "SELECT CAST(COUNT(*) AS BIT) FROM sqlite_master WHERE type = 'table' AND name = $1"
	stmt, err := db.Prepare(statement)
	if err != nil {
		Error.Println("Problems creating the statement that validates the existence of local storage. ", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(tableName).Scan(&existsTable)
	if err != nil {
		Error.Println("Could not find the tables within local storage. ", err)
	}

	return
}

// createTable non-exportable function used to create a table, in this case the medical reports table.
func createTable(table int) (err error) {
	//
	statement := ""
	var db *sql.DB
	if table == 1 {
		db, err = connectDBMedicalReports()
		if err != nil {
			Error.Println("Could not create medical reports table. ", err)
		}
		statement = `CREATE TABLE IF NOT EXISTS medical_reports (
			acs_number TEXT, 
			pat_id TEXT, 
			pat_name TEXT,
			pat_ext_id TEXT,  
			pat_gender TEXT, 
			pat_birth TEXT, 
			pat_email TEXT, 
			pat_address TEXT, 
			pat_phone TEXT, 
			pat_cellphone TEXT, 
			phys_name TEXT, 
			phys_specialty TEXT,
			ref_phys_name TEXT,
			ref_phys_mail TEXT,
			date_report TEXT, 
			sequence INTEGER NOT NULL, 
			content TEXT,
			status INTEGER DEFAULT 0,
			PRIMARY KEY (acs_number, sequence))`
	} else {
		db, err = connectDBMedicalProcedures()
		if err != nil {
			Error.Println("Could not create medical procedure table. ", err)
		}
		statement = `CREATE TABLE IF NOT EXISTS medical_procedures (
		    acs_number TEXT, 
		    pat_id TEXT, 
		    pat_name TEXT, 
		    pat_ext_id TEXT, 
		    pat_gender TEXT, 
		    pat_birth TEXT, 
		    pat_email TEXT,
		    phys_name TEXT,
		    phys_email TEXT,
		    status INTEGER DEFAULT 0,
		    PRIMARY KEY (acs_number, pat_email))`
	}
	defer db.Close()

	stmt, err := db.Prepare(statement)
	if err != nil {
		Error.Println("Problems preparing the statement that create the local storage. ", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()

	return
}

/*
StoreMedicalReportList exportable function that receives data from various medical
reports extracted from the RIS to store them locally.
*/
func StoreMedicalReportList(mrs []MedicalReport) (err error) {
	//
	db, err := connectDBMedicalReports()
	if err != nil {
		Error.Println("Could not insert medical reports into local storage. ", err)
	}
	defer db.Close()

	statement := `INSERT INTO medical_reports (acs_number, pat_id, 
		pat_name, pat_ext_id, pat_gender, pat_birth, pat_email, 
		pat_address, pat_phone, pat_cellphone, phys_name, phys_specialty, 
		ref_phys_name, ref_phys_mail, date_report, sequence, content, status) 
		VALUES ($1, $2, $3, $4, $5, $6,	$7, $8, $9, $10, $11, $12, $13, $14, 
			$15, $16, $17, $18)`

	stmt, err := db.Prepare(statement)
	if err != nil {
		Error.Println("Problems building the statement that insert medical report's data into local storage. ", err)
	}
	defer stmt.Close()

	Info.Println("Saving medical reports within local storage.")

	for _, reg := range mrs {
		//
		//fmt.Println(reg)
		_, err := stmt.Exec(reg.AccessionNumber, reg.PatientId, reg.PatientName,
			reg.PatientExtId, reg.PatientGender, reg.PatientBirth, reg.PatientEmail,
			reg.PatientAddress, reg.PatientPhone, reg.PatientCellphone,
			reg.PhysicianName, reg.PhysicianSpecialty, reg.PhysicianReferred,
			reg.PhysicianRefMail, reg.DateReport, reg.Sequence, reg.Content, Registered)
		if err != nil {
			//Warning.Println("An issue occurred while storing medical reports within local storage. ", err)
			continue
		}
		//fmt.Println(res.RowsAffected())
	}
	return
}

/*
StoreMedicalProcedureList exportable function that receives data from various medical
procedure extracted from the RIS to store them locally.
*/
func StoreMedicalProcedureList(mrs []MedicalReport) (err error) {
	//
	db, err := connectDBMedicalProcedures()
	if err != nil {
		Error.Println("Could not insert medical procedure into local storage. ", err)
	}
	defer db.Close()

	statement := `INSERT INTO medical_procedures (acs_number, pat_id, 
		pat_name, pat_ext_id, pat_gender, pat_birth, pat_email, phys_name, 
		phys_email, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	stmt, err := db.Prepare(statement)
	if err != nil {
		Error.Println("Problems building the statement that insert medical procedure's data into local storage. ", err)
	}
	defer stmt.Close()

	Info.Println("Saving medical procedures within local storage.")

	for _, reg := range mrs {
		//
		//fmt.Println(reg)
		_, err := stmt.Exec(reg.AccessionNumber, reg.PatientId,
			reg.PatientName, reg.PatientExtId, reg.PatientGender, reg.PatientBirth,
			reg.PatientEmail, reg.PhysicianReferred, reg.PhysicianRefMail, Registered)
		if err != nil {
			//Warning.Println("An issue occurred while storing medical procedures within local storage. ", err)
			continue
		}
		//fmt.Println(res.RowsAffected())
	}
	return
}

/*
changeMedicalReportStatus is a non-exportable function used to modify the
report's status during the upload process.
*/
func changeMedicalReportStatus(accession string, status int) (count int64, err error) {
	//
	db, err := connectDBMedicalReports()
	if err != nil {
		Error.Println("Could not modify medical reports status into local storage. ", err)
	}
	defer db.Close()

	statement := "UPDATE medical_reports SET status = $1 WHERE acs_number = $2"
	stmt, err := db.Prepare(statement)
	if err != nil {
		Error.Println("Problems preparing the statement that changes the status of medical reports. ", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(status, accession)
	if err != nil {
		Error.Println("The status update could not be executed in the medical report. ", err)
	}

	count, err = res.RowsAffected()
	if err != nil {
		Error.Println("Problems calculating the records affected by the status update in the medical reports. ", err)
	}

	return
}

/*
changeMedicalProcedureStatus is a non-exportable function used to modify the
medical procedure's status during the upload process.
*/
func changeMedicalProcedureStatus(accession string, status int) (count int64, err error) {
	//
	db, err := connectDBMedicalProcedures()
	if err != nil {
		Error.Println("Could not modify medical procedures status into local storage. ", err)
	}
	defer db.Close()

	statement := "UPDATE medical_procedures SET status = $1 WHERE acs_number = $2"
	stmt, err := db.Prepare(statement)
	if err != nil {
		Error.Println("Problems preparing the statement that changes the status of medical procedure. ", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(status, accession)
	if err != nil {
		Error.Println("The status update could not be executed in the medical procedure. ", err)
	}

	count, err = res.RowsAffected()
	if err != nil {
		Error.Println("Problems calculating the records affected by the status update in the medical procedure. ", err)
	}

	return
}

/*
getNextMedicalReport is a non-exportable function used to obtain from local
storage the data of the next medical report to be processed. Returns a list
of medical reports
*/
func getNextMedicalReport() (mrs []MedicalReport, err error) {
	//
	db, err := connectDBMedicalReports()
	if err != nil {
		Error.Println("Unable to access local storage to get the next medical report. ", err)
	}
	defer db.Close()

	statement := `SELECT acs_number, 
	pat_id, pat_name, pat_ext_id, pat_gender, pat_birth, 
	pat_email, pat_address, pat_phone, pat_cellphone, 
	phys_name, phys_specialty, ref_phys_name, 
	ref_phys_mail, date_report, sequence, content 
	FROM medical_reports
	WHERE acs_number = (SELECT acs_number
		FROM medical_reports WHERE status = 0
		LIMIT 1)
	ORDER BY sequence`
	stmt, err := db.Prepare(statement)
	if err != nil {
		Error.Println("Problems preparing the statement that gets the next report to process. ", err)
	}
	defer stmt.Close()

	Info.Println("Obtaining the next medical report to process.")

	rows, err := stmt.Query()
	if err != nil {
		Error.Println("Problems executing the query that returns the next report to process. ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var reg MedicalReport

		err = rows.Scan(&reg.AccessionNumber, &reg.PatientId, &reg.PatientName, &reg.PatientExtId,
			&reg.PatientGender, &reg.PatientBirth, &reg.PatientEmail, &reg.PatientAddress,
			&reg.PatientPhone, &reg.PatientCellphone, &reg.PhysicianName, &reg.PhysicianSpecialty,
			&reg.PhysicianReferred, &reg.PhysicianRefMail, &reg.DateReport, &reg.Sequence, &reg.Content)
		if err != nil {
			Error.Println("Cannot get report list from local storage. ", err)
		}

		mrs = append(mrs, reg)
	}

	return
}

/*
getNextMedicalProcedure is a non-exportable function used to obtain from local
storage the data of the next medical procedure to be processed. Returns a list
of medical procedure
*/
func getNextMedicalProcedure() (reg MedicalReport, err error) {
	//
	db, err := connectDBMedicalProcedures()
	if err != nil {
		Error.Println("Unable to access local storage to get the next medical procedure. ", err)
	}
	defer db.Close()

	statement := `SELECT acs_number,
	pat_id, pat_name, pat_ext_id, pat_gender, 
	pat_birth, pat_email, phys_name, phys_email
	FROM medical_procedures
	WHERE status = 0
    LIMIT 1`
	stmt, err := db.Prepare(statement)
	if err != nil {
		Error.Println("Problems preparing the statement that gets the next procedure to process. ", err)
	}
	defer stmt.Close()

	Info.Println("Obtaining the next medical procedure to process.")

	rows, err := stmt.Query()
	if err != nil {
		Error.Println("Problems executing the query that returns the next procedure to process. ", err)
	}
	defer rows.Close()

	for rows.Next() {
		//
		err = rows.Scan(&reg.AccessionNumber, &reg.PatientId,
			&reg.PatientName, &reg.PatientExtId, &reg.PatientGender, &reg.PatientBirth,
			&reg.PatientEmail, &reg.PhysicianReferred, &reg.PhysicianRefMail)
		if err != nil {
			Error.Println("Cannot get medical report from local storage. ", err)
		}
	}

	return
}

func getImage(id int) (content string) {
	//
	db, err := connectDBCodeImages()
	if err != nil {
		Error.Println("Unable to access local storage to get the encoded image. ", err)
	}
	defer db.Close()

	statement := `SELECT content
	FROM code_images
	WHERE id = $1`

	stmt, err := db.Prepare(statement)
	if err != nil {
		Error.Println("Problems preparing the statement that gets the encoded image. ", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		Error.Println("Problems executing the query that returns the encoded image. ", err)
	}
	defer rows.Close()

	for rows.Next() {
		//
		err = rows.Scan(&content)
		if err != nil {
			Error.Println("Cannot get encoded image from local storage. ", err)
		}
	}

	return
}

/*
buildJSONObjectReport is a non-exportable function used to build the JSON object
that will be used in the call to the microservice from the several data of
a medical report extracted from the local storage.
*/
func buildJSONObjectReport(mrs []MedicalReport) (jsonObj string, err error) {
	//
	var bodyReport string
	for _, mr := range mrs {
		//
		bodyReport += mr.Content
	}

	dateTime := time.Now() //.Format("01-02-2006 15:04:05")

	reg := HealthReport{
		AccessionNumber:    mrs[0].AccessionNumber,
		Institution:        config.Institution.Id,
		PatientId:          mrs[0].PatientId,
		PatientName:        mrs[0].PatientName,
		PatientGender:      mrs[0].PatientGender,
		PatientBirth:       mrs[0].PatientBirth,
		PatientEmail:       mrs[0].PatientEmail,
		PhysicianName:      mrs[0].PhysicianName,
		PhysicianSpecialty: mrs[0].PhysicianSpecialty,
		ReferInstitution:   "",
		ReferEmailInst:     "",
		ReferPhys:          mrs[0].PhysicianReferred,
		ReferEmailPhys:     "",
		DateReport:         mrs[0].DateReport,
		Attached:           nil,
		Report: innerReport{
			Header: complexField{
				Mime:    "image/png",
				Content: getImage(0),
			},
			Body: complexField{
				Mime:    "text/rtf",
				Content: bodyReport,
			},
			Signature: complexField{
				Mime:    "",
				Content: "",
			},
			Footer: complexField{
				Mime:    "image/png",
				Content: getImage(1),
			},
		},
		StudyDate: "",
		StudyName: "",
		StudyType: "",
		Tag1:      mrs[0].PatientExtId,
		Tag2:      "",
		Tag3: crmField{ // This embedded document is only used for the Metropolitan Clinic CRM.
			PatientAddress:   mrs[0].PatientAddress,
			PatientPhone:     mrs[0].PatientPhone,
			PatientCellphone: mrs[0].PatientCellphone,
		},
		Control: control{
			Active:    true,
			CreatedBy: "1a11aaaa11a1a1111a11a1a1",
			CreatedAt: dateTime,
			UpdatedBy: "1a11aaaa11a1a1111a11a1a1",
			UpdatedAt: dateTime,
		},
	}

	var jsonData []byte
	jsonData, err = json.Marshal(reg)
	if err != nil {
		Error.Println("Cannot parse report data to JSON object. ", err)
	}

	jsonObj = string(jsonData)

	return
}

/*
buildJSONObjectProcedure is a non-exportable function used to build the JSON object
that will be used in the call to the microservice from data of a medical procedure
extracted from the local storage.
*/
func buildJSONObjectProcedure(mp MedicalReport) (jsonObj string, err error) {
	//
	dateTime := time.Now() //.Format("01-02-2006 15:04:05")

	reg := HealthReport{
		AccessionNumber:    mp.AccessionNumber,
		Institution:        config.Institution.Id,
		PatientId:          mp.PatientId,
		PatientName:        mp.PatientName,
		PatientGender:      mp.PatientGender,
		PatientBirth:       mp.PatientBirth,
		PatientEmail:       mp.PatientEmail,
		PhysicianName:      "",
		PhysicianSpecialty: "",
		ReferInstitution:   "",
		ReferEmailInst:     "",
		ReferPhys:          mp.PhysicianReferred,
		ReferEmailPhys:     "",
		DateReport:         "",
		Attached:           nil,
		Report: innerReport{
			Header: complexField{
				Mime:    "",
				Content: "",
			},
			Body: complexField{
				Mime:    "",
				Content: "",
			},
			Signature: complexField{
				Mime:    "",
				Content: "",
			},
			Footer: complexField{
				Mime:    "",
				Content: "",
			},
		},
		StudyDate: "",
		StudyName: "",
		StudyType: "",
		Tag1:      mp.PatientExtId,
		Tag2:      "",
		Tag3: crmField{ // This embedded document is only used for the Metropolitan Clinic CRM.
			PatientAddress:   "",
			PatientPhone:     "",
			PatientCellphone: "",
		},
		Control: control{
			Active:    true,
			CreatedBy: "1a11aaaa11a1a1111a11a1a1",
			CreatedAt: dateTime,
			UpdatedBy: "1a11aaaa11a1a1111a11a1a1",
			UpdatedAt: dateTime,
		},
	}

	var jsonData []byte
	jsonData, err = json.Marshal(reg)
	if err != nil {
		Error.Println("Cannot parse procedure data to JSON object. ", err)
	}

	jsonObj = string(jsonData)

	return
}
