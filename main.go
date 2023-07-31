package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Course struct {
	ID            int    `json:"id"`
	NamaMk        string `json:"namaMk"`
	Jurusan       string `json:"jurusan"`
	Fakultas      string `json:"fakultas"`
	JumlahSks     int    `json:"jumlahSks"`
	SemesterMin   int    `json:"semesterMin"`
	PrediksiNilai string `json:"prediksiNilai"`
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("mysql", "root:Aremaniak1_@tcp(db)/course_scheduler")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected to the database")
}

func setupCorsConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func getAllCourses(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM courses")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		err := rows.Scan(&course.ID, &course.NamaMk, &course.Jurusan, &course.Fakultas, &course.JumlahSks, &course.SemesterMin, &course.PrediksiNilai)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		courses = append(courses, course)
	}

	c.JSON(http.StatusOK, courses)
}

func removeCourse(c *gin.Context) {
	id := c.Param("id")
	courseID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	_, err = db.Exec("DELETE FROM courses WHERE id=?", courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}

func addCourse(c *gin.Context) {
	var course Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
		return
	}

	result, err := db.Exec("INSERT INTO courses (namaMk, jurusan, fakultas, jumlahSks, semesterMin, prediksiNilai) VALUES (?, ?, ?, ?, ?, ?)",
		course.NamaMk, course.Jurusan, course.Fakultas, course.JumlahSks, course.SemesterMin, course.PrediksiNilai)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	courseID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course added successfully", "id": courseID})
}

func addDataJson(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	data, err := ioutil.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var courses []Course
	if err := json.Unmarshal(data, &courses); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()

	for _, course := range courses {
		_, err := tx.Exec("INSERT INTO courses (namaMk, jurusan, fakultas, jumlahSks, semesterMin, prediksiNilai) VALUES (?, ?, ?, ?, ?, ?)",
			course.NamaMk, course.Jurusan, course.Fakultas, course.JumlahSks, course.SemesterMin, course.PrediksiNilai)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data added successfully"})
}

var nilaiKonversi = map[string]float64{
	"A":  4.0,
	"AB": 3.5,
	"B":  3.0,
	"BC": 2.5,
	"C":  2.0,
	"D":  1.0,
}

func getKonversiValue(nilai string) float64 {
	return nilaiKonversi[nilai]
}

func searchCourses(jurusan, fakultas string, semesterPengambilan, minSKS, maxSKS int) []Course {
	courses := getCoursesByJurusanFakultas(jurusan, fakultas)

	// Filter courses based on the semester constraint
	var validCourses []Course
	for _, course := range courses {
		if course.SemesterMin <= semesterPengambilan {
			validCourses = append(validCourses, course)
		}
	}

	n := len(validCourses)

	// Initialize DP table
	dp := make([][]float64, n+1)
	for i := range dp {
		dp[i] = make([]float64, maxSKS+1)
	}

	// Dynamic Programming - knapsack variation
	for i := 1; i <= n; i++ {
		for w := 1; w <= maxSKS; w++ {
			// Check if including the current course is possible
			if validCourses[i-1].JumlahSks <= w {
				// Choose the maximum between including or excluding the current course
				dp[i][w] = max(dp[i-1][w], dp[i-1][w-validCourses[i-1].JumlahSks]+getKonversiValue(validCourses[i-1].PrediksiNilai))
			} else {
				dp[i][w] = dp[i-1][w]
			}
		}
	}

	// Find the combination of courses that maximizes the total predicted value
	var selectedCourses []Course
	w := maxSKS
	for i := n; i > 0 && w > 0; i-- {
		if dp[i][w] != dp[i-1][w] {
			selectedCourses = append(selectedCourses, validCourses[i-1])
			w -= validCourses[i-1].JumlahSks
		}
	}

	return selectedCourses
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func getCoursesByJurusanFakultas(jurusan, fakultas string) []Course {
	// Fetch courses from the database based on jurusan and fakultas
	// (Implement this part based on your database schema)
	// Here's a placeholder implementation:

	// Assuming 'db' is the global database connection object
	rows, err := db.Query("SELECT * FROM courses WHERE jurusan=? AND fakultas=?", jurusan, fakultas)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		err := rows.Scan(&course.ID, &course.NamaMk, &course.Jurusan, &course.Fakultas, &course.JumlahSks, &course.SemesterMin, &course.PrediksiNilai)
		if err != nil {
			panic(err.Error())
		}
		courses = append(courses, course)
	}
	fmt.Println(courses)
	return courses
}

func searchCoursesAPI(c *gin.Context) {
	jurusan := c.Param("jurusan")
	fmt.Println("jurusan " + jurusan)
	fakultas := c.Param("fakultas")
	fmt.Println("fakultas " + fakultas)
	semester, _ := strconv.Atoi(c.Param("semester"))
	fmt.Println("semester " + c.Param("semester"))
	minSKS, _ := strconv.Atoi(c.Param("minSKS"))
	fmt.Println("minSKS " + c.Param("minSKS"))
	maxSKS, _ := strconv.Atoi(c.Param("maxSKS"))
	fmt.Println("maxSKS " + c.Param("maxSKS"))

	// Call the searchCourses function to get the selected courses based on the provided parameters
	selectedCourses := searchCourses(jurusan, fakultas, semester, minSKS, maxSKS)

	// Return the selected courses in the response
	c.JSON(http.StatusOK, selectedCourses)
}

func main() {
	initDB()

	r := gin.Default()

	// Enable CORS
	r.Use(setupCorsConfig())

	r.GET("/api/getAllCourses", getAllCourses)
	r.DELETE("/api/removeCourses/:id", removeCourse)
	r.POST("/api/addCourses", addCourse)
	r.POST("/api/addDataJson", addDataJson)

	// Add the new API endpoint for course search
	r.GET("/api/searchCourses/:jurusan/:fakultas/:semester/:minSKS/:maxSKS", searchCoursesAPI)
	r.Run(":5001")
}
