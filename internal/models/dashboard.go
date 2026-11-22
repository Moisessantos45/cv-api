package models

type DashboardMetrics struct {
	Projects ProjectStats `json:"projects"`
	Posts   PostStats   `json:"posts"`
	Videos  VideoStats  `json:"videos"`
}

type ProjectStats struct {
	Total     int64 `json:"total"`
	Draft     int64 `json:"draft"`
	Published int64 `json:"published"`
	Archived  int64 `json:"archived"`
}

type PostStats struct {
	Total     int64 `json:"total"`
	Draft     int64 `json:"draft"`
	Published int64 `json:"published"`
	Archived  int64 `json:"archived"`
	Private   int64 `json:"private"`
	Unlisted  int64 `json:"unlisted"`
}

type VideoStats struct {
	Total int64 `json:"total"`
}
