package model

type Label struct {
	Id          string
	Name        string
	Images      []*Image
	ContactInfo string
	Profile     string
	DataQuality string
	Urls        []string
	ParentLabel *LabelLabels
	SubLabels   []*LabelLabels
}

type LabelLabels struct {
	Id   string
	Name string
}

/*


<?xml version="1.0" encoding="UTF-8"?>
<label>
   <images>
      <image height="600" type="primary" uri="" uri150="" width="600" />
      <image height="338" type="secondary" uri="" uri150="" width="300" />
      <image height="169" type="secondary" uri="" uri150="" width="160" />
   </images>
   <id>43</id>
   <name>Axis</name>
   <contactinfo>Axis Records&#xD;
P.O. Box 416600&#xD;
Miami Beach, FL 33141 USA&#xD;
Tel: +1 786 953 4176&#xD;
&#xD;
Early contact address:&#xD;
Axis Records&#xD;
180 North Wabash Avenue &#xD;
Suite 315&#xD;
Chicago, Illinois 60601&#xD;
USA &#xD;
&#xD;
tel/fax: +1 1-312-917-0800&#xD;
info@axisrecords.com / mills@axisrecords.com</contactinfo>
   <profile>American techno label established in 1991 in Chicago, operated by producer and owner, [a=Jeff Mills]. &#xD;
&#xD;
The label releases most of Mills' own work.</profile>
   <data_quality>Needs Vote</data_quality>
   <parentLabel id="1175241">Axis Records (10)</parentLabel>
   <urls>
      <url>http://www.axisrecords.com</url>
      <url>http://twitter.com/AxisRecords</url>
      <url>http://www.youtube.com/user/AxisRecords1</url>
   </urls>
   <sublabels>
      <label id="15681">6277</label>
      <label id="4504">Luxury Records</label>
      <label id="46">Purpose Maker</label>
      <label id="176113">Something In The Sky</label>
      <label id="349357">Taken</label>
      <label id="333">Tomorrow</label>
   </sublabels>
</label>
*/
