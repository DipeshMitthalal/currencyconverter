package currencyconverter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	. "github.com/clbanning/mxj"

	"appengine"
	"appengine/urlfetch"
)

func init() {
	http.HandleFunc("/convert", ConvertCurrency)
}

//ConvertCurrency function is called when a url with path/currency is requested
func ConvertCurrency(responseWriter http.ResponseWriter, request *http.Request) {
	//Setting app engine context
	appEngineContext := appengine.NewContext(request)
	//Parsing the request header to see if XML content type is requested
	requestType := request.Header.Get("Accept")
	isXMLResponseRequested := strings.Contains(requestType, "xml")

	//Parsing the query paramter named currency from URL
	//Query paramter is case sensitive. it doesnt work if name of paramter  is Currency instead of currency
	baseCurrency := request.FormValue("currency")
	//Parsing the query paramter named amount from URL
	baseAmount := request.FormValue("amount")

	baseAmountInFloat, _ := strconv.ParseFloat(baseAmount, 64)

	if len(baseAmount) == 0 {
		responseWriter.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(responseWriter, "Invalid Amount or Amount not entered")
		return
	}

	if len(baseCurrency) <= 2 {
		responseWriter.Header().Set("Content-Type", "text/plain")
		responseWriter.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(responseWriter, "Invalid Currency Input or Currency is missing")
		return
	}

	if baseAmountInFloat == 0 {
		responseWriter.Header().Set("Content-Type", "text/plain")
		responseWriter.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(responseWriter, "Amount can be greater then zero")
		return
	}
	//convertedAmountinForeignCurrencies is used to map the response from currency convert API
	responseFromCurrencyConverterAPI, convertedAmountinForeignCurrencies := CalculateCurrency(baseCurrency, baseAmountInFloat, appEngineContext)
	//convertedCurrency is used to map the response as per the requirments
	amountinForeignCurrencies := convertedCurrency{baseAmount, baseCurrency, convertedAmountinForeignCurrencies.Rates}
	//FOREX service API response is different from the requirements for shipwallets currenyConvertion Service
	//Hence we map response from forexapi service (convertedAmountinForeignCurrencies) to amountinForeignCurrencies(convertedCurrency)
	//convertedCurrency is used to map the response as per the requirments
	log.Printf("responseFromCurrencyConverterAPI", responseFromCurrencyConverterAPI)

	if responseFromCurrencyConverterAPI.StatusCode != 200 {
		responseWriter.Header().Set("Content-Type", "text/plain")
		responseWriter.WriteHeader(responseFromCurrencyConverterAPI.StatusCode)
		b, _ := ioutil.ReadAll(responseFromCurrencyConverterAPI.Body)
		responseWriter.Write(b)
		return
	}

	amountinForeignCurrenciesInJSON, _ := json.Marshal(amountinForeignCurrencies)
	if isXMLResponseRequested {
		responseWriter.Header().Set("Content-Type", "application/xml")
	} else {
		responseWriter.Header().Set("Content-Type", "application/json")
	}
	if isXMLResponseRequested {
		//Just a wrapper on json.Unmarshal
		//	Converting JSON to XML is a simple as:
		JSONinXML, err := NewMapJson(amountinForeignCurrenciesInJSON)
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}
		XMLBYTE, err := JSONinXML.Xml()
		if err != nil {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write(XMLBYTE)

	} else {
		//If XML respponse is not requested, respond in json
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write(amountinForeignCurrenciesInJSON)
	}

}

//currencyFromAPI is used to map the response from currency convert API
type currencyFromAPI struct {
	Base  string
	Date  string
	Rates map[string]float64
}

//FOREX servic API response is different from the requirements for shipwallets currenyConvertion Service
//Hence we use convertedCurrency is used to map currencyFromAPI to conform to the shipwallet's response requirements
//convertedCurrency is used to map the response as per the requirments
type convertedCurrency struct {
	Amount    string
	Currency  string
	Converted map[string]float64
}

//Function CalculateCurrency is used to convery the base curreny and amount using and base currency and amount is passed as arguments
func CalculateCurrency(baseCurrency string, baseAmount float64, appEngineContext appengine.Context) (response *http.Response, convertedCurrency currencyFromAPI) {

	url := "http://api.fixer.io/latest?base=" + baseCurrency
	appEngineContext.Infof("Requested URL: %v", url)
	appEngineHttpClient := urlfetch.Client(appEngineContext)
	request, err := http.NewRequest("GET", url, nil)
	response, err = appEngineHttpClient.Do(request)
	if err != nil {
		panic(err)
		appEngineContext.Errorf("%v", err)
	}
	defer response.Body.Close()

	appEngineContext.Infof("response Status: %v", response.Status)
	appEngineContext.Infof("response Headers: %v", response.Header)
	body, _ := ioutil.ReadAll(response.Body)

	convertedCurrency = currencyFromAPI{}
	//umarshalling JSON into interface valyes
	err = json.Unmarshal(body, &convertedCurrency)
	if err != nil {
		panic(err)
	}

	for nameOfCurrency, exchangeRate := range convertedCurrency.Rates {
		//fmt.Println("k:", nameOfCurrency, "v:", exchangeRate)
		convertedCurrency.Rates[nameOfCurrency] = Round(baseAmount*exchangeRate, 0.5, 2)
	}

	return
}

//Function to round the float64 to 2 given decimal places
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
