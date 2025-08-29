package seeding

import (
	"fmt"
	"grocademy/internal/db/models"
	"grocademy/internal/pkg/string_array"
	"math/rand"
	"reflect"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

type SeederInterface interface {
	Seed(userCount, courseCount, modulePerCourse int)
}

type Seeder struct {
	DB *gorm.DB
}

func NewSeeder(db *gorm.DB) *Seeder {
	seeder := Seeder{DB: db}
	return &seeder
}

func (s *Seeder) Seed(userCount, courseCount, modulePerCourse int) {
	s.CustomGenerator(courseCount)
	s.SeedUser(userCount)
	s.SeedCourse(courseCount)
	s.SeedModule(modulePerCourse, courseCount)
}

func (s *Seeder) SeedUser(userCount int) {
	for range userCount {
		a := models.User{}
		err := faker.FakeData(&a)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%+v\n", a)
		if res := s.DB.Create(&a); res.Error != nil {
			fmt.Println(res.Error)
		}
	}
}

func (s *Seeder) SeedCourse(courseCount int) {
	for range courseCount {
		a := models.Course{}
		err := faker.FakeData(&a)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%+v\n", a)
		if res := s.DB.Create(&a); res.Error != nil {
			fmt.Println(res.Error)
		}
	}
}

func (s *Seeder) SeedModule(modulePerCourse, courseCount int) {
	for range courseCount {
		for range modulePerCourse {
			a := models.Module{}
			err := faker.FakeData(&a)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%+v\n", a)
			if res := s.DB.Create(&a); res.Error != nil {
				fmt.Println(res.Error)
			}
		}
	}
}

func (s *Seeder) CustomGenerator(courseCount int) {
	_ = faker.AddProvider("thumbnail", func(v reflect.Value) (interface{}, error) {
		return "https://res.cloudinary.com/dlybowzgq/image/upload/v1756054842/dafdaf_fsusvf.jpg", nil
	})

	_ = faker.AddProvider("course_id", func(v reflect.Value) (interface{}, error) {
		var course models.Course
		s.DB.Order("RANDOM()").First(&course)
		print("ANJENG ")
		println(course.ID)
		return course.ID, nil
	})

	_ = faker.AddProvider("order", func(v reflect.Value) (interface{}, error) {
		var count int64
		s.DB.Model(&models.Module{}).Count(&count)
		return count + 1, nil
	})

	_ = faker.AddProvider("topics", func(v reflect.Value) (interface{}, error) {
		randomNumber := rand.Intn(5)
		topics := make([]string, randomNumber)
		for i := range randomNumber {
			topics[i] = faker.Word()
		}
		return string_array.StringArray(topics), nil
	})

	_ = faker.AddProvider("pdf_path", func(v reflect.Value) (interface{}, error) {
		return "https://res.cloudinary.com/dlybowzgq/image/upload/v1756046256/modul1.pdf", nil
	})

	_ = faker.AddProvider("video_path", func(v reflect.Value) (interface{}, error) {
		return "https://res.cloudinary.com/dlybowzgq/video/upload/v1756057789/modul2_rcctmc.mp4", nil
	})
}
