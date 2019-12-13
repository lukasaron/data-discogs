package model

type Label struct {
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	Images      []Image      `json:"images,omitempty"`
	ContactInfo string       `json:"contactInfo"`
	Profile     string       `json:"profile"`
	DataQuality string       `json:"dataQuality"`
	Urls        []string     `json:"urls"`
	ParentLabel *LabelLabel  `json:"parentLabel,omitempty"`
	SubLabels   []LabelLabel `json:"subLabels"`
}

type LabelLabel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

/*


<?xml version="1.0" encoding="UTF-8"?>
<label>
	<images>
       <image height="600" type="primary" uri="" uri150="" width="600" />
       <image height="338" type="secondary" uri="" uri150="" width="300" />
    </images>
    <id>43</id>
    <name>Axis</name>
    <contactinfo>Axis Records&#xD;
P.O. Box 416600&#xD;
Miami Beach, FL 33141 USA&#xD;
Tel: +1 786 953 4176&#xD;
	</contactinfo>
	<profile>American techno label established in 1991 in Chicago</profile>
	<data_quality>Needs Vote</data_quality>
	<parentLabel id="1175241">Axis Records (10)</parentLabel>
   	<urls>
    	<url>http://www.axisrecords.com</url>
        <url>http://twitter.com/AxisRecords</url>
    </urls>
    <sublabels>
   		<label id="15681">6277</label>
      	<label id="4504">Luxury Records</label>
   	</sublabels>
</label>
*/