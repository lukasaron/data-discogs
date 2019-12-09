package model

type Track struct {
	Position string `json:"position"`
	Title    string `json:"title"`
	Duration string `json:"duration"`
}

/*
<tracklist>
  <track>
	 <position>A1</position>
	 <title>A Sea Apart</title>
	 <duration>5:08</duration>
  </track>
</tracklist>
*/
