package structs

type AniData struct {
	Data MediaData `json:"data"`
}
type Title struct {
	Romaji  string `json:"romaji"`
	English string `json:"english"`
}
type CoverImage struct {
	Large string `json:"large"`
	Color string `json:"color"`
}

type NextAiringEpisode struct {
	AiringAt int `json:"airingAt"`
	Episode  int `json:"episode"`
}
type Media struct {
	ID                int               `json:"id"`
	Title             Title             `json:"title"`
	Type              string            `json:"type"`
	Genres            []string          `json:"genres"`
	CoverImage        CoverImage        `json:"coverImage"`
	Status            string            `json:"status"`
	Season            string            `json:"season"`
	SeasonYear        int               `json:"seasonYear"`
	Episodes          int               `json:"episodes"`
	AverageScore      int               `json:"averageScore"`
	MeanScore         int               `json:"meanScore"`
	Format            string            `json:"format"`
	Description       string            `json:"description"`
	NextAiringEpisode NextAiringEpisode `json:"nextAiringEpisode"`
}
type MediaData struct {
	Media Media `json:"Media"`
}

type AniStaffData struct {
	Data StaffData `json:"data"`
}
type DateOfBirth struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}
type Name struct {
	Full string `json:"full"`
}
type Image struct {
	Large string `json:"large"`
}
type Nodes struct {
	ID    int   `json:"id"`
	Title Title `json:"title"`
}
type MediaNodes struct {
	Nodes []Media `json:"nodes"`
}
type CharacterNode struct {
	ID    int        `json:"id"`
	Name  Name       `json:"name"`
	Media MediaNodes `json:"media"`
}
type Characters struct {
	Nodes []CharacterNode `json:"nodes"`
}
type Staff struct {
	ID                 int         `json:"id"`
	Gender             string      `json:"gender"`
	Age                int         `json:"age"`
	PrimaryOccupations []string    `json:"primaryOccupations"`
	DateOfBirth        DateOfBirth `json:"dateOfBirth"`
	Name               Name        `json:"name"`
	Image              Image       `json:"image"`
	Characters         Characters  `json:"characters"`
}
type StaffData struct {
	Staff Staff `json:"Staff"`
}

type AniCharData struct {
	Data CharData `json:"data"`
}

type Character struct {
	ID          int         `json:"id"`
	Gender      string      `json:"gender"`
	Age         string      `json:"age"`
	Name        Name        `json:"name"`
	DateOfBirth DateOfBirth `json:"dateOfBirth"`
	Image       Image       `json:"image"`
	Media       MediaNodes  `json:"media"`
}
type CharData struct {
	Character Character `json:"Character"`
}