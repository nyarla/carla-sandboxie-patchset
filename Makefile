.PHONY: all build proxy finalize submodule patch docker clean

all: clean build finalize

build: patch docker proxy
	docker run --rm -it \
		-v $(shell pwd)/Carla:/home/builder/src \
		-v $(shell pwd)/dist:/home/builder/dist \
		-u builder nyarla/carla \
		-c '(cd src/Carla && make distclean && cd ../../) ; \
				./src/build.sh win64 ; \
			  ./src/build.sh win64 ; \
			  ./src/build.sh win64 ; \
			  rm ./PawPawBuilds/targets/win64/lib/python3.8/site-packages/liblo.pyd ; \
			  ./src/build.sh win64'

proxy:
	cd proxy && (	\
			GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui" -o ../builds/carla-bridge-native.exe carla-bridge.go ; \
			GOOS=windows GOARCH=386 go build -ldflags="-H windowsgui" -o ../builds/carla-bridge-win32.exe carla-bridge.go ; \
			GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui" -o ../builds/carla-discovery-native.exe carla-discovery.go ; \
			GOOS=windows GOARCH=386 go build -ldflags="-H windowsgui" -o ../builds/carla-discovery-win32.exe carla-discovery.go ;)

finalize:
	cp Carla/Carla/Carla-2.3-win64.zip builds/
	cd builds && unzip Carla-2.3-win64.zip
	cd builds/Carla-*/ && bash -c 'for app in Carla Carla.lv2 Carla.vst; do \
		(cd $$app ; \
			mv carla-discovery-win32.exe _carla-discovery-win32.exe ; \
			mv carla-discovery-native.exe _carla-discovery-native.exe ; \
			mv carla-bridge-win32.exe _carla-bridge-win32.exe ; \
			mv carla-bridge-native.exe _carla-bridge-native.exe ; \
			cp ../../carla-*-*.exe . ; \
			cp carla-discovery-native.exe carla-discovery-win64.exe ; \
			cp _carla-discovery-native.exe _carla-discovery-win64.exe ; \
		cd ..) ; \
	done'
	cd builds/Carla-* && cp -R . ../../dist/

submodule:
	git submodule update --init --recursive

patch:
	cd Carla/Carla \
		&& ( patch -p1 -i ../../patches/sandboxie-support.patch; \
				 patch -p1 -i ../../patches/sandboxie-discovery.patch; )

docker:
	docker build -t nyarla/carla .

clean:
	cd Carla/Carla && git restore .
	rm -rf builds/*
