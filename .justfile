# just run
run:
    env ENV=dev go run main.go


# podman build image
podbuild:
    podman build -t go-mall:v1 .