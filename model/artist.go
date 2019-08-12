package model

type Artist struct {
	Id             string
	Name           string
	RealName       string
	Images         []*Image
	Profile        string
	DataQuality    string
	NameVariations []string
	Urls           []string
	Aliases        []*ArtistAlias
}

type ArtistAlias struct {
	Id   string
	Name string
}

/*
<?xml version="1.0" encoding="UTF-8"?>
<artist>
   <images>
      <image height="337" type="primary" uri="" uri150="" width="600" />
      <image height="554" type="secondary" uri="" uri150="" width="600" />
      <image height="325" type="secondary" uri="" uri150="" width="500" />
      <image height="300" type="secondary" uri="" uri150="" width="300" />
      <image height="600" type="secondary" uri="" uri150="" width="600" />
      <image height="400" type="secondary" uri="" uri150="" width="600" />
      <image height="376" type="secondary" uri="" uri150="" width="600" />
      <image height="302" type="secondary" uri="" uri150="" width="456" />
      <image height="682" type="secondary" uri="" uri150="" width="600" />
      <image height="177" type="secondary" uri="" uri150="" width="285" />
      <image height="225" type="secondary" uri="" uri150="" width="225" />
      <image height="350" type="secondary" uri="" uri150="" width="600" />
      <image height="225" type="secondary" uri="" uri150="" width="400" />
      <image height="250" type="secondary" uri="" uri150="" width="600" />
      <image height="450" type="secondary" uri="" uri150="" width="600" />
      <image height="285" type="secondary" uri="" uri150="" width="600" />
      <image height="399" type="secondary" uri="" uri150="" width="600" />
      <image height="400" type="secondary" uri="" uri150="" width="600" />
      <image height="322" type="secondary" uri="" uri150="" width="600" />
      <image height="353" type="secondary" uri="" uri150="" width="530" />
      <image height="364" type="secondary" uri="" uri150="" width="405" />
      <image height="398" type="secondary" uri="" uri150="" width="600" />
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
      <name id="140604">Fishtail</name>
      <name id="37287">Jovan Blade</name>
      <name id="87224">Lake Mead Drive</name>
      <name id="60876">Madd Phlavor</name>
      <name id="32119">Plastic Soul Junkies</name>
      <name id="237455">Prolific</name>
      <name id="124">Seven Grand Housing Authority</name>
      <name id="1204622">Telephone (2)</name>
      <name id="119">Terrence Parker</name>
      <name id="13595">The Lost Articles</name>
      <name id="31905">Tia's Daddy</name>
   </aliases>
</artist>
*/
