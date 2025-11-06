images:
	docker build -t loopvar:go122 -f go1.22.dockerfile .
	docker build -t loopvar:go121 -f go1.21.dockerfile .

run: images
	docker run loopvar:go121
	docker run loopvar:go122
