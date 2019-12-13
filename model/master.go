package model

type Master struct {
	Id          string          `json:"id"`
	MainRelease string          `json:"mainRelease"`
	Images      []Image         `json:"images,omitempty"`
	Artists     []ReleaseArtist `json:"artists"`
	Genres      []string        `json:"genres"`
	Styles      []string        `json:"styles"`
	Year        string          `json:"year"`
	Title       string          `json:"title"`
	DataQuality string          `json:"dataQuality"`
	Videos      []Video         `json:"videos"`
}

/*
<?xml version="1.0" encoding="UTF-8"?>
<master id="33228">
   <main_release>341048</main_release>
   <images>
      <image height="602" type="primary" uri="" uri150="" width="600" />
      <image height="600" type="secondary" uri="" uri150="" width="600" />
   </images>
   <artists>
      <artist>
         <id>20691</id>
         <name>Philip Glass</name>
         <anv />
         <join />
         <role />
         <tracks />
      </artist>
   </artists>
   <genres>
      <genre>Electronic</genre>
      <genre>Classical</genre>
   </genres>
   <styles>
      <style>Soundtrack</style>
      <style>Modern Classical</style>
   </styles>
   <year>1988</year>
   <title>Powaqqatsi</title>
   <data_quality>Correct</data_quality>
   <videos>
      <video duration="303" embed="true" src="https://www.youtube.com/watch?v=jRg3agJn1Mg">
         <title>Philip Glass - Powaqqatsi - 01. Serra Pelada</title>
         <description>Philip Glass - Powaqqatsi - 01. Serra Pelada</description>
      </video>
      <video duration="25" embed="true" src="https://www.youtube.com/watch?v=4E4gn9nNGYs">
         <title>Philip Glass - Powaqqatsi - 02. The Title</title>
         <description>Philip Glass - Powaqqatsi - 02. The Title</description>
      </video>
   </videos>
</master>
*/
