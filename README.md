# currencyconverter
One needs to have google app engine for GO SDK and related packages
GOROOT has to set correctly to place where GO is installed and also packages related to googleapp engine is present

Test 1 for JSON response:
curl -X GET -H "Accept: application/json, */*" 'http://currencyconverter-1106.appspot.com/currency?amount=1&currency=SEK'
{"Amount":"1","Currency":"SEK","Converted":{"AUD":0.17,"BGN":0.21,"BRL":0.48,"CAD":0.16,"CHF":0.12,"CNY":0.76,"CZK":2.87,"DKK":0.79,"EUR":0.11,"GBP":0.08,"HKD":0.93,"HRK":0.81,"HUF":33.03,"IDR":1634.9,"ILS":0.46,"INR":7.82,"JPY":14.38,"KRW":136.39,"MXN":1.99,"MYR":0.51,"NOK":0.98,"NZD":0.18,"PHP":5.59,"PLN":0.45,"RON":0.47,"RUB":7.51,"SGD":0.17,"THB":4.27,"TRY":0.35,"USD":0.12,"ZAR":1.62}}

Test 2 for XML response:
curl -X GET -H "Accept: application/xml, */*" 'http://currencyconverter-1106.appspot.com/currency?amount=1&currency=SEK'
<doc><Amount>1</Amount><Converted><AUD>0.17</AUD><BGN>0.21</BGN><BRL>0.48</BRL><CAD>0.16</CAD><CHF>0.12</CHF><CNY>0.76</CNY><CZK>2.87</CZK><DKK>0.79</DKK><EUR>0.11</EUR><GBP>0.08</GBP><HKD>0.93</HKD><HRK>0.81</HRK><HUF>33.03</HUF><IDR>1634.9</IDR><ILS>0.46</ILS><INR>7.82</INR><JPY>14.38</JPY><KRW>136.39</KRW><MXN>1.99</MXN><MYR>0.51</MYR><NOK>0.98</NOK><NZD>0.18</NZD><PHP>5.59</PHP><PLN>0.45</PLN><RON>0.47</RON><RUB>7.51</RUB><SGD>0.17</SGD><THB>4.27</THB><TRY>0.35</TRY><USD>0.12</USD><ZAR>1.62</ZAR></Converted><Currency>SEK</Currency></doc>
