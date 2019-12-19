package model

type Artist struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	RealName       string   `json:"realName"`
	Images         []Image  `json:"images,omitempty"`
	Profile        string   `json:"profile"`
	DataQuality    string   `json:"dataQuality"`
	NameVariations []string `json:"nameVariations"`
	Urls           []string `json:"urls"`
	Aliases        []Alias  `json:"aliases"`
	Members        []Member `json:"members,omitempty"`
}

type Alias struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Member struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

/*
<?xml version="1.0" encoding="UTF-8"?>
<artist>
   <images>
      <image height="337" type="primary" uri="" uri150="" width="600" />
      <image height="554" type="secondary" uri="" uri150="" width="600" />
   </images>
   <id>132</id>
   <name>Minimum Wage Brothers</name>
   <realname>Terrence Parker</realname>
   <profile />
   <data_quality>Correct</data_quality>
   <urls>
      <url>http://www.terrenceparker.net</url>
      <url>http://www.myspace.com/terrenceparker</url>
   </urls>
   <namevariations>
      <name>Minimum Wage Bros.</name>
      <name>The Minimum Wage Bros.</name>
   </namevariations>
   <aliases>
      <name id="10678">2 Sweat Doctors</name>
      <name id="121">Disco Revisited</name>
   </aliases>
</artist>
*/
