package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/godror/godror"
)

// MedicalReport struct
type MedicalReport struct {
	AccessionNumber    string
	PatientId          string
	PatientName        string
	PatientExtId       string
	PatientGender      string
	PatientBirth       string
	PatientEmail       string
	PatientAddress     string
	PatientPhone       string
	PatientCellphone   string
	PhysicianName      string
	PhysicianSpecialty string
	PhysicianReferred  string
	PhysicianRefMail   string
	DateReport         string
	Sequence           int
	Content            string
}

// GetMedicalReports is a afunction that
func GetMedicalReports() (mrs []MedicalReport, err error) {
	//
	uri := fmt.Sprintf("%s/%s@%s:%s/%s",
		config.Database.User,
		config.Database.Passwd,
		config.Database.Host,
		config.Database.Port,
		config.Database.Db)
	db, err := sql.Open("godror", uri)
	if err != nil {
		Error.Println("Cannot establish a connection with the RIS database. ", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		Error.Println("Unable to communicate with the RIS database. ", err)
		return
	}

	// January 2, 2006
	currentDate := time.Now().Local().Format("2006-01-02")
	Info.Println("Looking for medical reports within the RIS dated ", currentDate)

	sqlQuery := `SELECT UB.TERMIN_UBTID, NVL(PA.PD_MISC_TEXT_1, 'S/C'), 
		NP.VORNAME || ' ' || NP.NAME, NVL(PA.PATIENT_FREMD_ID, 'S/N'), 
		DECODE(NP.GESCHLECHT,'W','F',NP.GESCHLECHT), 
		NP.GEBURTSDATUM, NVL(LOWER(PA.PD_MISC_TEXT_2), 'S/C'),
		NVL(APAT.STRASSE_HAUSNR, 'S/D'), NVL(APAT.TEL_RUFNUMMER, 'S/N'), NVL(APAT.AD_MOBILE_PHONE, 'S/N'),
		NPM.VORNAME || ' ' || NPM.NAME, 'Radiologo',
		NVL(US.UEBERWEISER_KURZNAME, 'S/N'), NVL(LOWER(AD.EMAIL_ADRESSE), 'S/C'),
		B.BEFUND_DATUM, BT.BEFUND_TEXT_SEQUENZ, BT.BEFUND_TEXT_SEQUENZ_TEXT
	FROM NATUERLICHE_PERSONEN NP 
	RIGHT OUTER JOIN PATIENT PA ON NP.PERSONEN_ID = PA.PERSONEN_ID
	LEFT OUTER JOIN ADRESSEN APAT ON APAT.PERSONEN_ID = PA.PERSONEN_ID
	LEFT OUTER JOIN UNTBEH UB ON PA.PERSONEN_ID = UB.PATIENT_PID AND UB.TERMIN_UBTID > 0
	LEFT OUTER JOIN UEBERWEISENDE_STELLE US ON UB.UNTBEH_UEBERWEISER_PID = US.PERSONEN_ID
	LEFT OUTER JOIN ADRESSEN AD ON UB.UNTBEH_UEBERWEISER_PID = AD.PERSONEN_ID
	LEFT OUTER JOIN UBERGEBNIS UE ON UB.TERMIN_UBTID = UE.TERMIN_UBTID AND UE.UBERGEBNIS_MASKE = 'BEFUND'
	LEFT OUTER JOIN BEFUND B ON UE.UBERGEBNIS_UBEID = B.UBERGEBNIS_UBEID AND UE.BEFUND_UBEID = B.BEFUND_UBEID
	LEFT OUTER JOIN BEFUND_TEXT BT ON TO_CHAR(UE.UBERGEBNIS_UBEID) = BT.UBERGEBNIS_UBEID
	LEFT OUTER JOIN NATUERLICHE_PERSONEN NPM ON B.BEFUND_SIGNIERER_PID = NPM.PERSONEN_ID
	WHERE B.BEFUND_DATUM >= TO_DATE(:1,'YYYY-MM-DD HH24:MI:SS')
	ORDER BY UB.TERMIN_UBTID, BT.BEFUND_TEXT_SEQUENZ ASC`

	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		Error.Println("Problems preparing the statement that will extract the data from the RIS. ", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(currentDate)
	if err != nil {
		Error.Println("Problems performing the process that will obtain the medical reports from the RIS. ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var reg MedicalReport

		err = rows.Scan(&reg.AccessionNumber, &reg.PatientId, &reg.PatientName,
			&reg.PatientExtId, &reg.PatientGender, &reg.PatientBirth, &reg.PatientEmail,
			&reg.PatientAddress, &reg.PatientPhone, &reg.PatientCellphone,
			&reg.PhysicianName, &reg.PhysicianSpecialty, &reg.PhysicianReferred,
			&reg.PhysicianRefMail, &reg.DateReport, &reg.Sequence, &reg.Content)
		if err != nil {
			Error.Println("Data cannot be parsed into the medical report structure. ", err)
		}

		mrs = append(mrs, reg)

	}

	return
}

// GetMedicalProcedures is a afunction that
func GetMedicalProcedures() (mrs []MedicalReport, err error) {
	//
	uri := fmt.Sprintf("%s/%s@%s:%s/%s",
		config.Database.User,
		config.Database.Passwd,
		config.Database.Host,
		config.Database.Port,
		config.Database.Db)
	db, err := sql.Open("godror", uri)
	if err != nil {
		Error.Println("Cannot establish a connection with the RIS database. ", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		Error.Println("Unable to communicate with the RIS database. ", err)
		return
	}

	// January 2, 2006
	currentDate := time.Now().Local().Format("2006-01-02")
	Info.Println("Looking for RIS patient with medical procedure at", currentDate)

	sqlQuery := `SELECT UB.TERMIN_UBTID, NVL(PA.PD_MISC_TEXT_1, 'S/C'), 
	                 NP.VORNAME || ' ' || NP.NAME, NVL(PA.PATIENT_FREMD_ID, 'S/N'), 
	                 DECODE(NP.GESCHLECHT,'W','F',NP.GESCHLECHT), 
                     NP.GEBURTSDATUM, NVL(LOWER(PA.PD_MISC_TEXT_2), 'S/C'),
                     NVL(US.UEBERWEISER_KURZNAME, 'S/N'), NVL(LOWER(AD.EMAIL_ADRESSE), 'S/C')
                 FROM NATUERLICHE_PERSONEN NP 
                 RIGHT OUTER JOIN PATIENT PA ON NP.PERSONEN_ID = PA.PERSONEN_ID
                 LEFT OUTER JOIN UNTBEH UB ON PA.PERSONEN_ID = UB.PATIENT_PID
                 LEFT OUTER JOIN UEBERWEISENDE_STELLE US ON UB.UNTBEH_UEBERWEISER_PID = US.PERSONEN_ID
                 LEFT OUTER JOIN ADRESSEN AD ON UB.UNTBEH_UEBERWEISER_PID = AD.PERSONEN_ID
                 WHERE UB.UNTBEH_START >= TO_DATE(:1,'YYYY-MM-DD HH24:MI:SS')`

	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		Error.Println("Problems preparing the statement that will extract the data of RIS patients. ", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(currentDate)
	if err != nil {
		Error.Println("Problems performing the process that will obtain the patients with medical procedures. ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var reg MedicalReport

		err = rows.Scan(&reg.AccessionNumber, &reg.PatientId, &reg.PatientName,
			&reg.PatientExtId, &reg.PatientGender, &reg.PatientBirth,
			&reg.PatientEmail, &reg.PhysicianReferred, &reg.PhysicianRefMail)
		if err != nil {
			Error.Println("Data can't be parsed into the medical report structure when trying to get RIS patients. ", err)
		}

		mrs = append(mrs, reg)

	}

	return
}
