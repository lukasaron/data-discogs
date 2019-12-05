package model

type Format struct {
	Name         string
	Quantity     string
	Text         string
	Descriptions []string
}

/*
 <formats>
      <format name="Vinyl" qty="1" text="">
         <descriptions>
            <description>12"</description>
            <description>33 ⅓ RPM</description>
         </descriptions>
      </format>
   </formats>
*/
