package dto

type CreateUserRequest struct {
	Firebase_id string `json:"firebase_id"`
	Cases       string `json:"cases"`
}

type UserResponse struct {
	Mongo_id    string `json:"_id"`
	Firebase_id string `json:"firebase_id"`
	Cases       string `json:"cases"`
}

type UpdateUserRequest struct {
	Firebase_id string `json:"firebase_id"`
	Cases       string `json:"cases"`
}

type DeleteUserRequest struct {
	Firebase_id string `json:"firebase_id"`
}
