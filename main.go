package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	route.PathPrefix("/Public/").Handler(http.StripPrefix("/Public", http.FileServer(http.Dir("./Public"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/project", project).Methods("GET")
	route.HandleFunc("/add-project", addProject).Methods("POST")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/project-detail/{index}", projectDetail).Methods("GET")
	// EDIT PROJECT
	route.HandleFunc("/edit-project/{index}", editproject).Methods("GET")
	// route.HandleFunc("/update-project", updateproject).Methods("POST")
	// DELETE PROJECT
	route.HandleFunc("/delete-project/{index}", DeleteProject).Methods("GET")

	port := "8000"

	fmt.Println("Server sedang berjalan di port " + port)
	http.ListenAndServe("localhost:"+port, route)

}

type Project struct {
	ProjectName  string
	StartDate    string
	EndDate      string
	Duration     string
	Description  string
	Technologies []string
}

var Projects = []Project{
	{
		ProjectName:  "Percobaan Name",
		StartDate:    "01 October 2022",
		EndDate:      "01 november 2022",
		Duration:     "1 Months",
		Description:  "unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo.",
		Technologies: []string{"nodejs", "react-native", "nextjs", "vuejs"},
	},
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contet-type", "text/html; charset-utf-8")
	tmpl, err := template.ParseFiles("views/index.html")

	if err != nil {
		log.Println(err)
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	respndata := map[string]interface{}{
		"Projects": Projects,
	}

	tmpl.Execute(w, respndata)
}
func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contet-type", "text/html; charset-utf-8")
	tmpl, err := template.ParseFiles("views/project.html")

	if err != nil {
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contet-type", "text/html; charset-utf-8")

	tmpl, err := template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("Message: " + err.Error()))
		return
	}
	tmpl.Execute(w, nil)
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contet-type", "text/html; charset-utf-8")
	tmpl, err := template.ParseFiles("views/project-detail.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		var renderDetail = Project{}
		index, _ := strconv.Atoi(mux.Vars(r)["index"])

		for i, data := range Projects {
			if index == i {
				renderDetail = Project{
					ProjectName:  data.ProjectName,
					StartDate:    data.StartDate,
					EndDate:      data.EndDate,
					Duration:     data.Duration,
					Description:  data.Description,
					Technologies: data.Technologies,
				}
			}
		}

		data := map[string]interface{}{
			"renderDetail": renderDetail,
		}
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, data)
	}
}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {
		projectName := r.PostForm.Get("project")
		StartDate := r.PostForm.Get("date-start")
		EndDate := r.PostForm.Get("date-end")
		Description := r.PostForm.Get("description")
		Technologies := r.Form["technologies"]

		var newProject = Project{
			ProjectName:  projectName,
			StartDate:    FormatDate(StartDate),
			EndDate:      FormatDate(EndDate),
			Duration:     GetDuration(StartDate, EndDate),
			Description:  Description,
			Technologies: Technologies,
		}

		Projects = append(Projects, newProject)

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

// UPDATE PROJECT
func editproject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("views/edit-project.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		var updateData = Project{}
		index, _ := strconv.Atoi(mux.Vars(r)["index"])

		for i, data := range Projects {
			if index == i {
				updateData = Project{
					ProjectName:  data.ProjectName,
					StartDate:    ReturnDate(data.StartDate),
					EndDate:      ReturnDate(data.EndDate),
					Description:  data.Description,
					Technologies: data.Technologies,
				}
				Projects = append(Projects[:index], Projects[index+1:]...)
			}
		}
		data := map[string]interface{}{
			"updateData": updateData,
		}
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, data)
	}
}

// func updateproject(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseForm()

// 	if err != nil {
// 		log.Fatal(err)
// 	} else {
// 		projectName := r.PostForm.Get("project")
// 		StartDate := r.PostForm.Get("date-start")
// 		EndDate := r.PostForm.Get("date-end")
// 		Description := r.PostForm.Get("description")
// 		Technologies := r.Form["technologies"]

// 		var newProject = Project{
// 			ProjectName:  projectName,
// 			StartDate:    FormatDate(StartDate),
// 			EndDate:      FormatDate(EndDate),
// 			Duration:     GetDuration(StartDate, EndDate),
// 			Description:  Description,
// 			Technologies: Technologies,
// 		}

// 		Projects = append(Projects, newProject)

// 		http.Redirect(w, r, "/", http.StatusMovedPermanently)
// 	}

// }

// DELETE PROJECT
func DeleteProject(w http.ResponseWriter, r *http.Request) {

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	Projects = append(Projects[:index], Projects[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}

// GET DURATION
func GetDuration(startDate string, endDate string) string {

	layout := "2006-01-02"

	date1, _ := time.Parse(layout, startDate)
	date2, _ := time.Parse(layout, endDate)

	margin := date2.Sub(date1).Hours() / 24
	var duration string

	if margin > 30 {
		if (margin / 30) <= 1 {
			duration = "1 Month"
		} else {
			duration = strconv.Itoa(int(margin)/30) + " Months"
		}
	} else {
		if margin <= 1 {
			duration = "1 Day"
		} else {
			duration = strconv.Itoa(int(margin)) + " Days"
		}
	}

	return duration
}

// CHANGE DATE FORMAT
func FormatDate(InputDate string) string {

	layout := "2006-01-02"
	t, _ := time.Parse(layout, InputDate)

	Formated := t.Format("02 January 2006")

	return Formated
}

// RETURN DATE FORMAT
func ReturnDate(InputDate string) string {

	layout := "02 January 2006"
	t, _ := time.Parse(layout, InputDate)

	Formated := t.Format("2006-01-02")

	return Formated
}
