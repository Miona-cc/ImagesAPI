package types

type Image struct {
	Id   string `json:"id" bson:"id"`
	Type string `json:"type" bson:"type"`
	Src  string `json:"src" bson:"src"`
}

type Match struct {
	Id   string   `json:"id" bson:"id"`
	Type string   `json:"type" bson:"type"`
	Srcs []string `json:"srcs" bson:"srcs"`
}
