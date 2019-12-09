package model

type Video struct {
	Duration    string `json:"duration"`
	Embed       string `json:"embed"`
	Src         string `json:"src"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

/*
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
*/
