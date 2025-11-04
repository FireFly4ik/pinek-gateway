package models

type UploadDeleteFilesResponse struct {
	FileIDs []string `json:"file_ids"`
}

type UploadFileResponse struct {
	FileID  string `json:"file_id"`
	FileURL string `json:"file_url"`
	Message string `json:"message"`
}

type GetFilesResponse struct {
	FileURLs []string `json:"file_urls"`
}

type GetFileResponse struct {
	FileURL string `json:"file_url"`
}
