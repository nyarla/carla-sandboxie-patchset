.PHONY: docker debug build helpers clean

all: clean helpers

clean:
	cd Carla && git checkout . && git clean -f -d
	rm -rf dist/* && touch dist/.gitkeep
	rm -rf builds/* && rm -rf builds/.done-*

docker:
	docker build -t nyarla/carla .

debug: docker
	docker run --rm -it \
		-v $(shell pwd):/home/builder/src \
		-v $(shell pwd)/builds:/home/builder/builds \
		-u builder nyarla/carla

build: docker
	docker run --rm -it \
		-v $(shell pwd):/home/builder/src \
		-v $(shell pwd)/builds:/home/builder/builds \
		-u builder nyarla/carla /home/builder/src/bin/build.sh

helpers: build
	bash bin/helpers.sh
