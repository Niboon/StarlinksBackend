package api

type Star struct {
	ID		int 	`json:"id"`
	Name	string `json:"name"`
	UserID 	int 	`json:"user_id"`
	Img 	string `json:"img"`
	Link 	string `json:"link"`
}

type User struct {
	ID    	int 		`json:"id"`
	Token  	string 	`json:"token"`
}
