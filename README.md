
# Assestment by Saham Rakyat

Reza Wahyu Ramadhan
[![Built with GO-Echo Framework](https://img.shields.io/badge/built%20with-GoEcho-ff69b4.svg?logo=cookiecutter)](https://github.com/karec/cookiecutter-flask-restful)

## Documentation

Copy the .env-example to .env and set up the credential.

```
cp .env-example .env
```

run main program

```
go run main.go
```

The API documentation is on [Documentation](https://documenter.getpostman.com/view/17785126/2s93m1b5Uv) or access the json file GO API DOCUMENTATION.postman_collection.json on this repository

The cache is expired every 5 minutes and right now there is neither observer nor listener is implemented
## BONUS QUESTION

#### Why clean architecture is good for your application ? if yes, please explain


Yes, I always try my best to code with clean architecture. most of my projects either my corporate project or my freelance project is using boilerplate to initiate the projects. Most of them is using cookiecutter boilerplate. I hardly tried to find the boilerplate for echo framework but I couldnt find it, so I just tried to built this application as close as I can to Laravel boilerplate.

Using clean architecture also provides a solid foundation for building maintainable, testable, and flexible applications.

#### How to scale up your application and when it needs to be

Scaling is needed when I have a serious issue within traffic, workload, or data. And there is 2 different ways to scale the application either vertical or horizontal. Currently I just ever scale my application with vertical scaling by updating the framework that I used for the application.

