package main
import(
	"encoding/json"
    "fmt"
    "net/http"
	"io/ioutil"
    "log"
    "os"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
  _ "github.com/go-sql-driver/mysql"

)

func hello(c echo.Context) error {
    return c.String(http.StatusOK, "Hello Welcome to Our Application!")
}

// func getCats(c echo.Context) error {
//     catName := c.QueryParam("name")
//     catType := c.QueryParam("type")

//     dataType := c.Param("data")

//     if dataType == "string" {
//         return c.String(http.StatusOK, fmt.Sprintf("your cat name is: %s\nand his type is: %s\n", catName, catType))
//     }

//     if dataType == "json" {
//         return c.JSON(http.StatusOK, map[string]string{
//             "name": catName,
//             "type": catType,
//         })
//     }

//     return c.JSON(http.StatusBadRequest, map[string]string{
//         "error": "you need to lets us know if you want json or string data",
//     })
// }

func getData(c echo.Context) error {
	db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/goexample2?charset=utf8&parseTime=True")
	defer db.Close()
 	if err!=nil{
 	log.Println("Connection Failed to Open")
 	} 
	 log.Println("Connection Established")
	 db.SingularTable(true)
	dataObject := []Dataset{}
	db.Find(&dataObject)
	return c.JSON(http.StatusOK,dataObject)
}

type Response struct{
	gorm.Model
	ConformTo string `json:"conformsTo"`
	DescribedBy string `json:"describedBy"`
	Dataset []Dataset `json:"dataset"`
}

type Dataset struct{
	AccessLevel string `json:"accessLevel"`
	Description string `json:"description"`
	Identifier string `json:"identifier"`
	Title string `json:"title"`
	ResponseID uint
}

func main() {
    fmt.Println("Welcome to the server")

	// Below lines of code is required to read the Json Data and parse date into struct
	response, err := http.Get("https://www.ssa.gov/data.json")
    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }
    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
	responseObject := Response{}
	json.Unmarshal([]byte(responseData), &responseObject)
	fmt.Println(responseObject.ConformTo)
	fmt.Println(responseObject.DescribedBy)
	fmt.Println(len(responseObject.Dataset))

	
	// Below lines of code is used to connect to mysql with the help of Gorm. Please change the credentials based on your login credentials.
	db, err := gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/goexample2?charset=utf8&parseTime=True")
	defer db.Close()
 	if err!=nil{
 	log.Println("Connection Failed to Open")
 	} 
 	log.Println("Connection Established")
	db.SingularTable(true)
	db.DropTableIfExists(&Response{})
	db.DropTableIfExists(&Dataset{})
	db.CreateTable(&Response{})
	db.CreateTable(&Dataset{})
	db.Create(&responseObject)
	
	//Implemented Echo Web Framework
    e := echo.New()
	e.GET("/", hello)
	e.GET("/retreiveData",getData)
    e.Start(":8000")
}