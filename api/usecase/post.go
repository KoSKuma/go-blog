package usecase

import (
	"time"

	"github.com/KoSKuma/go-blog/api/entity"
	"github.com/KoSKuma/go-blog/api/repository"
)

type PostUsecase struct {
	DBRepo repository.DatabaseRepo
	Logger repository.Logger
}

func (p *PostUsecase) CreatePost(post entity.Post) (string, error) {
	post.PostTime = time.Now()
	id, err := p.DBRepo.InsertOne(post)
	if err != nil {
		return "", err
	}
	err = p.Logger.Log("CREATE POST", post)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (p *PostUsecase) FindPost(id string) (entity.Post, error) {
	post, err := p.DBRepo.FindOne(id)
	if err != nil {
		return post, err
	}
	err = p.Logger.Log("GET POST", post)
	if err != nil {
		return post, err
	}
	return post, nil
}

func (p *PostUsecase) UpdatePost(id string, post entity.PostUpdate) error {
	err := p.Logger.Log("UPDATE POST", post)
	if err != nil {
		return err
	}
	return p.DBRepo.UpdateOne(id, post)
}

func (p *PostUsecase) DeletePost(id string) error {
	err := p.Logger.Log("DELETE POST", id)
	if err != nil {
		return err
	}
	return p.DBRepo.DeleteOne(id)
}

func (p *PostUsecase) FindPosts() ([]entity.Post, error) {
	posts, err := p.DBRepo.FindAll()
	if err != nil {
		return posts, err
	}
	err = p.Logger.Log("GET POSTS", posts)
	if err != nil {
		return posts, err
	}
	return posts, nil
}
