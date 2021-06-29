package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Enumerate constants used to model the status of medical reports during the
// upload process to the viewmed platform
const (
	Registered = iota
	Processing
	Uploaded
)

func main() {
	//
	Info.Println("Starting MRU Application.")
	// Setting the Get Procedure timer
	timerGetProcedure := time.Duration(config.Timer.GetProcedure)
	tickGetProcedure := time.NewTicker(time.Second * timerGetProcedure)

	// Setting the Up Procedure timer
	timerUpProcedure := time.Duration(config.Timer.UpProcedure)
	tickUpProcedure := time.NewTicker(time.Second * timerUpProcedure)

	// Setting the Get Report timer
	timerGetReport := time.Duration(config.Timer.GetReport)
	tickGetReport := time.NewTicker(time.Second * timerGetReport)

	// Setting the Up Report timer
	timerUpReport := time.Duration(config.Timer.UpReport)
	tickUpReport := time.NewTicker(time.Second * timerUpReport)

	// Launching the goroutines that will take the data and upload them.
	go schedulerGetProcedure(tickGetProcedure)
	go schedulerUpProcedure(tickUpProcedure)
	go schedulerGetReport(tickGetReport)
	go schedulerUpReport(tickUpReport)

	// Get the system signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	tickGetProcedure.Stop()
	tickUpProcedure.Stop()
	tickGetReport.Stop()
	tickUpReport.Stop()
	//
	Warning.Println("Ending MRU Application.")
}

// Task Planner Get Procedure
func schedulerGetProcedure(tick *time.Ticker) {
	//
	for range tick.C {
		task1()
	}

}

//Task Planner Up Procedure
func schedulerUpProcedure(tick *time.Ticker) {
	//
	for range tick.C {
		if config.Service.Active {
			task2()
		}
	}

}

// Task Planner Get Report
func schedulerGetReport(tick *time.Ticker) {
	//
	for range tick.C {
		task3()
	}

}

//Task Planner Up Report
func schedulerUpReport(tick *time.Ticker) {
	//
	for range tick.C {
		if config.Service.Active {
			task4()
		}
	}

}

/*
Get Procedure task, it's responsible for obtaining the medical report data from
the RIS and stores it into the application.
*/
func task1() {
	// Extracting the data from Oracle database
	mp, err := GetMedicalProcedures()
	if err != nil {
		Error.Println("Problems from get medical procedures from Centricity. ", err)
	}
	//fmt.Println(mrs)
	Info.Println(len(mp), "procedures were obtained from the RIS.")

	// Continue the process if there are medical procedure into the RIS to process
	if len(mp) > 0 {
		// Validating destination table
		existe, err := checkTableProcedures("medical_procedures")
		Info.Println("Validating local storage's existence:", existe)
		if !existe {
			//
			Warning.Println("Local medical_reports table does not exist. It will be created.")
			if err = createTable(2); err != nil {
				Error.Println("Problem creating local table of medical_procedures. ", err)
			}
		}

		// Inserting into local database
		err = StoreMedicalProcedureList(mp)
		if err != nil {
			Error.Println("Problem inserting medical reports into local storage. ", err)
		}
	}
}

/*
Up Procedure task is responsible for extracting all the information of the medical
procedure from the local storage, processing the data, constructing the JSON object and
sending it to the viewmed platform using the call to the corresponding microservice.
*/
func task2() {
	//
	// Validating destination table
	existe, err := checkTableProcedures("medical_procedures")
	continuar := true
	Info.Println("Validating medical_procedures local storage's existence:", existe)
	if !existe {
		//
		Warning.Println("Local medical_procedures table does not exist. It will be created.")
		if err = createTable(2); err != nil {
			Error.Println("Problem creating local table of medical_procedures. ", err)
		}
	} else {
		// Get all the data of a medical procedure from local storage
		mp, err := getNextMedicalProcedure()
		if err != nil {
			Error.Println("Problems getting medical report from local storage. ", err)
			continuar = false
		}

		// Continue the process if there are medical procedure into the local storage to process
		if continuar {
			//
			medicalProcAN := mp.AccessionNumber
			Info.Println("Getting medical procedure with accession_number", medicalProcAN)

			// Updating local medical procedure to PROCESSING
			upd, err := changeMedicalProcedureStatus(medicalProcAN, Processing)
			if err != nil {
				Error.Println("The medical procedure", medicalProcAN, "could not be updated to the \"Processing\" status. ", err)
				// TODO retornar cuando hay un error
			}
			Info.Println("The medical procedure updated to \"Processing\":", upd)

			// Building the JSON object with the medical procedure data
			jsonObject, err := buildJSONObjectProcedure(mp)

			// Call to the microservice
			_, err = sendMR(jsonObject)
			if err != nil {
				Error.Println("Problems calling the viewmed microservice.", err)
				// Updating local medical procedure to REGISTERED
				upd, err = changeMedicalProcedureStatus(medicalProcAN, Registered)
				if err != nil {
					Error.Println("The medical procedure", medicalProcAN, "could not be updated to the \"Registered\" status. ", err)
					// TODO retornar cuando hay un error
				}
			} else {

				// Updating local medical procedure to UPLOADED
				upd, err = changeMedicalProcedureStatus(medicalProcAN, Uploaded)
				if err != nil {
					Error.Println("The medical procedure", medicalProcAN, "could not be updated to the \"Uploaded\" status. ", err)
					// TODO retornar cuando hay un error
				}
				Info.Println("The medical procedure updated to \"Uploaded\":", upd)
			}
		}
	}
}

/*
Get Report task, it's responsible for obtaining the medical report data from
the RIS and stores it into the application.
*/
func task3() {
	// Extracting the data from Oracle database
	mrs, err := GetMedicalReports()
	if err != nil {
		Error.Println("Problems from get medical reports from Centricity. ", err)
	}
	//fmt.Println(mrs)
	Info.Println(len(mrs), "reports were obtained from the RIS.")

	// Continue the process if there are medical reports into the RIS to process
	if len(mrs) > 0 {
		// Validating destination table
		existe, err := checkTableReports("medical_reports")
		Info.Println("Validating local storage's existence:", existe)
		if !existe {
			//
			Warning.Println("Local medical_reports table does not exist. It will be created.")
			if err = createTable(1); err != nil {
				Error.Println("Problem creating local table of medical_reports. ", err)
			}
		}

		// Inserting into local database
		err = StoreMedicalReportList(mrs)
		if err != nil {
			Error.Println("Problem inserting medical reports into local storage. ", err)
		}
	}
}

/*
Up Report task is responsible for extracting all the information of the medical
report from the local storage, processing the data, constructing the JSON object and
sending it to the viewmed platform using the call to the corresponding microservice.
*/
func task4() {
	//
	// Validating destination table
	existe, err := checkTableReports("medical_reports")
	Info.Println("Validating medical_reports local storage's existence:", existe)
	if !existe {
		//
		Warning.Println("Local medical_reports table does not exist. It will be created.")
		if err = createTable(1); err != nil {
			Error.Println("Problem creating local table of medical_reports. ", err)
		}
	} else {
		// Get all the data of a medical report from local storage
		mrs, err := getNextMedicalReport()
		if err != nil {
			Error.Println("Problems getting medical report from local storage. ", err)
		}
		Info.Println("The medical report to upload is divided into", len(mrs), "parts.")

		// Continue the process if there are medical reports into the local storage to process
		if len(mrs) > 0 {
			//
			medicalReprtAN := mrs[0].AccessionNumber
			Info.Println("Working with accession_number", medicalReprtAN)

			// Updating local medical report to PROCESSING
			upd, err := changeMedicalReportStatus(medicalReprtAN, Processing)
			if err != nil {
				Error.Println("The medical report", medicalReprtAN, "could not be updated to the \"Processing\" status. ", err)
				// TODO retornar cuando hay un error
			}
			Info.Println("Updated to \"Processing\" this amount of parts:", upd)

			// Building the JSON object with the medical report data
			jsonObject, err := buildJSONObjectReport(mrs)

			// Call to the microservice
			_, err = sendMR(jsonObject)
			if err != nil {
				Error.Println("Problems calling the viewmed microservice.", err)
				// Updating local medical report to REGISTERED
				upd, err = changeMedicalReportStatus(medicalReprtAN, Registered)
				if err != nil {
					Error.Println("The medical report", medicalReprtAN, "could not be updated to the \"Registered\" status. ", err)
					// TODO retornar cuando hay un error
				}
			} else {

				// Updating local medical report to UPLOADED
				upd, err = changeMedicalReportStatus(medicalReprtAN, Uploaded)
				if err != nil {
					Error.Println("The medical report", medicalReprtAN, "could not be updated to the \"Uploaded\" status. ", err)
					// TODO retornar cuando hay un error
				}
				Info.Println("Updated to \"Uploaded\" this amount of parts:", upd)
			}
		}
	}
}
