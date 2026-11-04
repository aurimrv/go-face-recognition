# Go-Face-Recognition

> **Teaching fork.** This repository is a teaching fork of [André Oliveira Silva](https://github.com/andre-ols)'s project [andre-ols/go-face-recognition](https://github.com/andre-ols/go-face-recognition), used in a hands-on class on GenAI-assisted refactoring of source and test code. Only the images were replaced; the source code is unchanged. All credit for the original project goes to its author, and it remains under the original [MIT License](LICENSE).
>
> To set up the development environment (Cursor + Dev Container), see [DEVCONTAINER_SETUP.md](DEVCONTAINER_SETUP.md).

Go-Face-Recognition is an facial recognition system, based on the principles of FaceNet and developed entirely in Go language. It leverages cutting-edge technology and utilizes the go-face library, which is built upon the powerful dlib C++ library for high-performance facial analysis.

## Table of Contents

- [Overview](#overview)
  - [About FaceNet](#about-facenet)
  - [About dlib](#about-dlib)
- [Key Features](#key-features)
- [Usage](#usage)
  - [Dynamic Loading of People](#dynamic-loading-of-people)
  - [Recognition of Faces](#recognition-of-faces)
  - [Output Generation](#output-generation)
  - [Capabilities](#capabilities)
- [Installation and Usage](#installation-and-usage)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

## About FaceNet

Go-Face-Recognition is based on the principles of [FaceNet](https://arxiv.org/abs/1503.03832), a groundbreaking facial recognition system developed by Google. FaceNet employs a deep neural network to directly learn a mapping from facial images to a compact Euclidean space, where distances between embeddings correspond directly to a measure of facial similarity. By leveraging this learned embedding space, tasks such as facial recognition, verification, and clustering become straightforward, as FaceNet embeddings serve as feature vectors that capture essential facial characteristics. This integration enables Go-Face-Recognition to achieve state-of-the-art performance in facial recognition tasks, making it a versatile and powerful tool for various applications.

### About dlib

[dlib](http://dlib.net/) is a modern C++ toolkit containing machine learning algorithms and tools for creating complex software in C++ to solve real-world problems. It is renowned for its robustness, efficiency, and versatility in various applications, including computer vision, machine learning, and artificial intelligence.

## Key Features

- **Dynamic Person Loading:** Go-Face-Recognition dynamically loads people from within the 'persons' directory, enhancing flexibility.

- **Precise Face Recognition:** Utilizing go-face library powered by dlib, the system performs accurate and reliable face recognition, even in complex scenarios.

- **Easy Deployment with Docker:** To streamline dependency management and deployment, the project is encapsulated within a Docker environment, ensuring seamless integration into any development or production environment.

## Usage

### Dynamic Loading of People

This project dynamically loads people from within the `persons/` directory. Each person should have a subfolder with the person's name, containing images of that person to be used in the model. It is ideal to provide more than one image per person to improve classification accuracy. The images provided for each person should contain only one face, which is the face of the person.

### Recognition of Faces

After loading the people, the software reads an image from the `images/` directory. By default, it searches for an image named `unknown.jpg`. It then recognizes the faces in the image based on the provided people. The input image can contain multiple people, and the software attempts to recognize all of them.

### Output Generation

The output of the system is a new image with the faces marked and the name of each identified person. The generated image will be saved in the `images/` directory with the name `result.jpg`.

### Capabilities

The system enables effortless recognition of faces within images, empowering users with a powerful tool for various applications, including security, authentication, access control, and more.

## Installation and Usage

1. Clone this repository:

```bash
git clone https://github.com/aurimrv/go-face-recognition.git
```

2. Navigate to the project directory:

```bash
cd go-face-recognition
```

3. Build the Docker image:

```bash
docker compose up --build -d
```

4. Run the Docker container:

```bash
docker exec go-face-recognition go run cmd/main.go
```

This will execute the project inside the Docker container and demonstrate facial recognition.

## Project Structure

```
go-face-recognition/
├── cmd
│   └── main.go
├── DEVCONTAINER_SETUP.md
├── docker-compose.yml
├── Dockerfile
├── fonts
│   └── Roboto-Regular.ttf
├── go.mod
├── go.sum
├── images
│   ├── Denzel-Washington
│   │   ├── unknown-01.jpg
│   │   ├── unknown-02.jpg
│   │   └── unknown-03.jpg
│   ├── Kate-Winslet
│   │   ├── unknown-01.jpg
│   │   └── unknown-02.jpg
│   ├── Megan-Fox
│   │   └── unknown-01.jpg
│   ├── Tom-Cruise
│   │   ├── unknown-01.avif
│   │   ├── unknown-01.jpg
│   │   ├── unknown-02.avif
│   │   └── unknown-02.jpg
│   └── unknown.jpg
├── internal
│   ├── entity
│   │   ├── drawer.go
│   │   └── person.go
│   └── usecases
│       ├── classify_person.go
│       ├── load_persons.go
│       └── recognize_persons.go
├── LICENSE
├── models
│   ├── dlib_face_recognition_resnet_model_v1.dat
│   ├── dnn_age_predictor_v1.dat
│   ├── dnn_gender_classifier_v1.dat
│   ├── mmod_human_face_detector.dat
│   ├── shape_predictor_5_face_landmarks.dat
│   └── shape_predictor_68_face_landmarks.dat
├── persons
│   ├── Denzel-Washington
│   │   ├── denzel-01.jpg
│   │   ├── denzel-02.jpg
│   │   └── denzel-03.jpg
│   ├── Kate-Winslet
│   │   ├── kate-01.jpg
│   │   ├── kate-02.jpg
│   │   └── kate-03.jpg
│   ├── Megan-Fox
│   │   ├── megan-01.jpg
│   │   ├── megan-02.jpg
│   │   └── megan-03.jpg
│   └── Tom-Cruise
│       ├── tom-01.jpg
│       ├── tom-02.jpg
│       └── tom-03.jpg
└── README.md
```

The `persons/` directory contains one sub-folder per known person (Denzel Washington, Kate Winslet, Megan Fox and Tom Cruise), with images of that person used to build the model. The `images/` directory contains the input images: `unknown.jpg` is the default input, and each person's sub-folder holds extra test images of that person that are not used for training. The output image, `result.jpg`, is generated in `images/` when the program runs. The `fonts/` directory contains the font used to write names on the output image. The `models/` directory contains the pre-trained dlib models for face detection, landmarks, recognition, age and gender. The `internal/` directory contains the core logic of the system, including entities and use cases. The `cmd/` directory contains the main entry point of the system.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvement, please open an issue or submit a pull request on the GitHub repository.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
