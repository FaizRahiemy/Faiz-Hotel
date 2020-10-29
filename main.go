// main.go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
    "strconv"
    "strings"

    "github.com/gorilla/mux"
)

type rooms struct {
    Id              int    `json:"id"`
    Room_number     string `json:"room_number"`
    Room_type       int `json:"room_type"`
    Price           int `json:"price"`
}

type room_type struct {
    Id              int    `json:"id"`
    Room_type       int `json:"room_type"`
    Price           int `json:"price"`
}

type price_rule struct {
    Id              int    `json:"id"`
    Name            string `json:"name"`
    Price_rule_type int `json:"price_rule_type"`
    Type_id         int `json:"type_id"`
    Room_id         int `json:"room_id"`
    Price_rule      int `json:"price_rule"`
    Price_value     int `json:"price_value"`
    Price_threshold int `json:"price_threshold"`
    Start_date      time.Time `json:"start_date"`
    End_date        time.Time `json:"end_date"`
}

type room_occupied struct {
    Id              int   `json:"id"`
    Room_id         int `json:"room_id"`
    Start_date      time.Time `json:"start_date"`
    End_date        time.Time `json:"end_date"`
    Customer_id         int `json:"customer_id"`
}

type promo struct {
    Id              int   `json:"id"`
    Name            string `json:"name"`
    Promo_type      int `json:"promo_type"`
    Value           int `json:"value"`
    Minimum_nights  int `json:"minimum_nights"`
    Minimum_rooms   int `json:"minimum_rooms"`
    Checkin_day     []int `json:"checkin_day"`
    Booking_day     []int `json:"booking_day"`
    Booking_hour    []int `json:"booking_hour"`
}

type response struct {
    List            string `json:"list"`
    Promo           int `json:"promo"`
    Final           int `json:"final"`
}

var roomss []rooms
var room_types []room_type
var price_rules []price_rule
var room_occupieds []room_occupied
var promos []promo

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to Faiz Hotel!")
    fmt.Println("Endpoint Hit: Home")
}

func returnAllRooms(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: returnAllRooms")
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    
    json.NewEncoder(w).Encode(roomss)
}

func returnSingleRoom(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: returnSingleRoom")
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    
    vars := mux.Vars(r)
    key, err := strconv.Atoi(vars["id"])
    
    if err != nil {
        // handle error
        fmt.Println(err)
    } else {
        for _, room := range roomss {
            if room.Id == key {
                json.NewEncoder(w).Encode(room)
            }
        }
    }
}

func returnPromoPrice(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: returnPromoPrice")
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)

    priceGet, err := r.URL.Query()["price"]
    totals, err2 := r.URL.Query()["total"]
    promoGet, err3 := r.URL.Query()["promo_id"]
    dayGet, err4 := r.URL.Query()["days"]
    
    if !err || len(priceGet[0]) < 1 || !err2 || len(totals[0]) < 1 || !err3 || len(promoGet[0]) < 1 || !err4 || len(dayGet[0]) < 1 {
        log.Println("Url Param is missing")
        return
    }
    
    prices := strings.Split(priceGet[0], ",")
    total, errTotal := strconv.Atoi(totals[0])
    promo, errPromo := strconv.Atoi(promoGet[0])
    days := strings.Split(dayGet[0], ",")
    
    list := ""
    totalPrice := 0
    
    if errTotal != nil || errPromo != nil {
        // handle error
        fmt.Println(err)
    } else {
        prom_type := 0
        value := 0
        for _, prom := range promos {
            if (prom.Id == promo) {
                if (len(prices) >= prom.Minimum_rooms && len(days) >= prom.Minimum_nights) {
                    for _, checkin := range prom.Checkin_day {
                        checkinStr := strconv.Itoa(checkin)
                        if (strings.Contains(dayGet[0], checkinStr)) {
                            prom_type = prom.Promo_type
                            value = prom.Value
                        }
                    }
                }
            }
        }
        for _, price := range prices {
            princeInt, errPrc := strconv.Atoi(price)

            if errPrc != nil {
                // handle error
                fmt.Println(err)
            } else {
                newprice := 0
                if (prom_type == 0) {
                    newprice = princeInt - value
                    total -= value
                } else {
                    newprice = princeInt - (princeInt * value)
                    total -= (princeInt * value)
                }
                if (list != "") {
                    list += ","
                }
                list = list + strconv.Itoa(newprice)
                totalPrice += newprice
            }
        }
                
        res := response{List: list, Promo: total, Final: totalPrice}
        json.NewEncoder(w).Encode(res)
    }
}

func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/", homePage)
    myRouter.HandleFunc("/rooms", returnAllRooms)
    myRouter.HandleFunc("/room/{id}", returnSingleRoom)
    myRouter.HandleFunc("/promo", returnPromoPrice)
    log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
    roomss = []rooms{
        rooms{Id: 1, Room_number: "2001", Room_type: 1, Price: -1},
        rooms{Id: 2, Room_number: "2002", Room_type: 1, Price: -1},
    }
    promos = []promo{
        promo{Id: 1, Name: "Monday Promo", Promo_type: 0, Value: 1000, Minimum_nights: 1, Minimum_rooms: 1, Checkin_day: []int{1}, Booking_day: []int{1}, Booking_hour: []int{10}},
        promo{Id: 2, Name: "Weekend Promo", Promo_type: 1, Value: 10, Minimum_nights: 2, Minimum_rooms: 2, Checkin_day: []int{0,6}, Booking_day: []int{0,1,2,3,4,5,6}, Booking_hour: []int{23}},
    }
    handleRequests()
}